package builder

import (
	"goful/core/model"
	starlarkng "goful/core/scripting/starlark"
	"goful/core/templating/yamlng"
	"reflect"

	"github.com/rs/zerolog/log"
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
		headers := yamlRequest.Headers
		for k, v := range profile.Variables {
			yamlRequest.Url = yamlng.ProcessTemplateVariable(yamlRequest.Url, k, v)
			for headerName, headerValues := range yamlRequest.Headers {
				headers[headerName] = yamlng.ProcessTemplateVariables(headerValues, k, v)
			}
		}
		yamlRequest.Headers = headers
	}

	request := model.Request{
		Url:     yamlRequest.Url,
		Method:  yamlRequest.Method,
		Headers: yamlRequest.Headers,
		Body:    yamlRequest.Body,
	}

	return request, true, nil
}

func buildStarlarkRequest(requestMold *model.RequestMold, previousResponse *model.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMold.Starlark == nil {
		return model.Request{}, false, nil
	}

	res, err := starlarkng.RunStarlarkScript(*requestMold, previousResponse, profile)
	if err != nil {
		log.Error().Err(err).Msg("Running Starlark script resulted to error")
		return model.Request{}, true, err
	}

	headers := make(map[string][]string)
	for k, headerVal := range res["headers"].(map[string]interface{}) {
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

	log.Debug().Msgf("Converted headers %v", headers)

	req := model.Request{
		Url:     res["url"].(string),
		Method:  res["method"].(string),
		Headers: new(model.Headers).FromMap(headers),
		Body:    res["body"],
	}

	log.Debug().Msgf("Built request %v", req)

	return req, true, nil
}
