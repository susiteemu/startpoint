package builder

import (
	"reflect"
	"startpoint/core/configuration"
	"startpoint/core/model"
	starlarkng "startpoint/core/scripting/starlark"
	"startpoint/core/templating/templateng"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var builders = []func(requestMold *model.RequestMold, previousResponse *model.Response, profile model.Profile) (model.Request, bool, error){
	buildYamlRequest,
	buildStarlarkRequest,
}

func BuildRequest(requestMold *model.RequestMold, profile model.Profile) (model.Request, error) {
	log.Debug().Msgf("Searching suitable builder for %v", requestMold)

	var request model.Request
	for _, builder := range builders {
		result, accept, err := builder(requestMold, nil, profile)
		if err != nil {
			return model.Request{}, err
		}
		if accept {
			request = result
			break
		}
	}
	return request, nil
}

func BuildRequestUsingPreviousResponse(requestMold *model.RequestMold, previousResponse *model.Response, profile model.Profile) (model.Request, error) {
	var request model.Request
	for _, builder := range builders {
		result, accept, err := builder(requestMold, previousResponse, profile)
		if err != nil {
			return model.Request{}, err
		}
		if accept {
			request = result
			break
		}
	}
	return request, nil
}

func buildYamlRequest(requestMold *model.RequestMold, _ *model.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMold.Yaml == nil {
		return model.Request{}, false, nil
	}

	yamlRequest := requestMold.Yaml

	if len(profile.Variables) > 0 {

		rawYaml := requestMold.Yaml.Raw
		for k, v := range profile.Variables {
			rawYaml, _ = templateng.ProcessTemplateVariable(rawYaml, k, v)
		}
		log.Debug().Msgf("Processed raw into %s", rawYaml)

		yamlRequest = &model.YamlRequest{}
		err := yaml.Unmarshal([]byte(rawYaml), yamlRequest)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to unmarshal yaml %s", rawYaml)
			return model.Request{}, false, err
		}
		yamlRequest.Raw = rawYaml

		log.Debug().Msgf("Processed into yaml request %v", yamlRequest)
	}

	options := make(map[string]interface{})
	if len(yamlRequest.Options) > 0 {
		configuration.Flatten("", yamlRequest.Options, options)
	}

	request := model.Request{
		Url:     yamlRequest.Url,
		Method:  yamlRequest.Method,
		Headers: yamlRequest.Headers,
		Body:    yamlRequest.Body,
		Options: options,
		Output:  yamlRequest.Output,
	}

	return request, true, nil
}

func buildStarlarkRequest(requestMold *model.RequestMold, previousResponse *model.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMold.Starlark == nil {
		return model.Request{}, false, nil
	}

	script := requestMold.Starlark.Script
	if len(profile.Variables) > 0 {
		for k, v := range profile.Variables {
			script, _ = templateng.ProcessTemplateVariable(script, k, v)
		}
	}
	requestMold.Starlark.Script = script
	res, err := starlarkng.RunStarlarkScript(*requestMold, previousResponse)
	if err != nil {
		log.Error().Err(err).Msg("Running Starlark script resulted to error")
		return model.Request{}, true, err
	}

	headersResult, has := res["headers"]
	headers := make(map[string][]string)
	if has {
		for k, headerVal := range headersResult.(map[string]interface{}) {
			t := reflect.TypeOf(headerVal)
			if t.String() == "string" {
				headers[k] = []string{headerVal.(string)}
			} else if t.String() == "[]interface {}" {
				var l []string
				for _, singleHeaderVal := range headerVal.([]interface{}) {
					l = append(l, singleHeaderVal.(string))
				}
				headers[k] = l
			}
		}
	}

	log.Debug().Msgf("Converted headers %v", headers)

	optionsResult, has := res["options"]
	options := make(map[string]interface{})
	if has {
		configuration.Flatten("", optionsResult.(map[string]interface{}), options)
	}

	outputResult, has := res["output"]
	output := ""
	if has {
		output = outputResult.(string)
	}
	if len(output) == 0 {
		output = requestMold.Output()
	}

	// FIXME: add checks
	req := model.Request{
		Url:     res["url"].(string),
		Method:  res["method"].(string),
		Headers: new(model.Headers).FromMap(headers),
		Body:    res["body"],
		Options: options,
		Output:  output,
	}

	log.Debug().Msgf("Built request %v", req)

	return req, true, nil
}
