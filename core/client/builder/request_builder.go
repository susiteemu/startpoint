package builder

import (
	b64 "encoding/base64"
	"fmt"
	"startpoint/core/configuration"
	"startpoint/core/model"
	luang "startpoint/core/scripting/lua"
	starlarkng "startpoint/core/scripting/starlark"
	"startpoint/core/templating/templateng"
	"startpoint/core/tools/conv"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var builders = []func(requestMold *model.RequestMold, previousResponse *model.Response, profile model.Profile) (model.Request, bool, error){
	buildYamlRequest,
	buildScriptableRequest,
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

	auth := yamlRequest.Auth
	if auth != (model.Auth{}) {
		if auth.Basic != (model.BasicAuth{}) {
			if auth.Basic.User != "" && auth.Basic.Password != "" {
				userPwd := fmt.Sprintf("%s:%s", auth.Basic.User, auth.Basic.Password)
				base64encoded := b64.StdEncoding.EncodeToString([]byte(userPwd))
				if yamlRequest.Headers == nil {
					yamlRequest.Headers = model.Headers{}
				}
				yamlRequest.Headers[model.HEADER_NAME_AUTHORIZATION] = model.HeaderValues{fmt.Sprintf("%s %s", model.HEADER_VALUE_BASIC_AUTH, base64encoded)}
			}

		} else if auth.Bearer != "" {
			if yamlRequest.Headers == nil {
				yamlRequest.Headers = model.Headers{}
			}
			yamlRequest.Headers[model.HEADER_NAME_AUTHORIZATION] = model.HeaderValues{fmt.Sprintf("%s %s", model.HEADER_VALUE_BEARER_AUTH, auth.Bearer)}
		}
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

func buildScriptableRequest(requestMold *model.RequestMold, previousResponse *model.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMold.Scriptable == nil {
		return model.Request{}, false, nil
	}

	script := requestMold.Scriptable.Script
	if len(profile.Variables) > 0 {
		for k, v := range profile.Variables {
			script, _ = templateng.ProcessTemplateVariable(script, k, v)
		}
	}
	requestMold.Scriptable.Script = script

	var res map[string]interface{}
	var err error
	switch requestMold.Type {
	case model.CONTENT_TYPE_STARLARK:
		res, err = starlarkng.RunStarlarkScript(*requestMold, previousResponse)
	case model.CONTENT_TYPE_LUA:
		res, err = luang.RunLuaScript(*requestMold, previousResponse)
	default:
		return model.Request{}, true, fmt.Errorf("Unsupported script type %s", requestMold.Type)
	}

	if err != nil {
		log.Error().Err(err).Msg("Running script resulted to error")
		return model.Request{}, true, err
	}

	headersResult, has := res["headers"]
	headers := make(map[string][]string)
	if has {
		for k, headerVal := range headersResult.(map[string]interface{}) {
			if asStr, ok := headerVal.(string); ok {
				headers[k] = []string{asStr}
			} else if asInterfaceArr, ok := headerVal.([]interface{}); ok {
				var l []string
				for _, singleHeaderVal := range asInterfaceArr {
					l = append(l, singleHeaderVal.(string))
				}
				headers[k] = l
			}
		}
	}

	authResult, has := res["auth"]
	if has {
		log.Debug().Msgf("Auth %v", authResult)
		if authMap, ok := authResult.(map[string]interface{}); ok {
			basicAuth, has := authMap["basic_auth"]
			if has {
				log.Debug().Msgf("Basic auth %T", basicAuth)
				var username, password interface{}
				var hasUser, hasPwd bool
				basicAuthStringMap, ok := basicAuth.(map[string]interface{})
				if ok {
					username, hasUser = basicAuthStringMap["username"]
					password, hasPwd = basicAuthStringMap["password"]
				} else {
					basicAuthAnyMap, ok := basicAuth.(map[interface{}]interface{})
					if ok {
						username, hasUser = basicAuthAnyMap["username"]
						password, hasPwd = basicAuthAnyMap["password"]
					}
				}
				if hasUser && hasPwd {
					userPwd := fmt.Sprintf("%s:%s", username, password)
					userPwdBytes := []byte(userPwd)
					base64encoded := b64.StdEncoding.EncodeToString(userPwdBytes)
					headers[model.HEADER_NAME_AUTHORIZATION] = []string{fmt.Sprintf("%s %s", model.HEADER_VALUE_BASIC_AUTH, base64encoded)}
				}
			} else {
				bearerToken, has := authResult.(map[string]interface{})["bearer_token"]
				if has {
					headers[model.HEADER_NAME_AUTHORIZATION] = []string{fmt.Sprintf("%s %s", model.HEADER_VALUE_BEARER_AUTH, bearerToken)}
				}
			}
		} else {
			log.Warn().Msgf("Auth %v is in invalid format %T", authResult, authResult)
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

	url, err := conv.AssertAndConvert[string](res, "url")
	if err != nil {
		return model.Request{}, true, err
	}
	method, err := conv.AssertAndConvert[string](res, "method")
	if err != nil {
		return model.Request{}, true, err
	}
	body, err := conv.AssertAndConvert[interface{}](res, "body")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to convert body")
	}

	req := model.Request{
		Url:     url,
		Method:  method,
		Headers: new(model.Headers).FromMap(headers),
		Body:    body,
		Options: options,
		Output:  output,
	}

	log.Debug().Msgf("Built request %v", req)

	return req, true, nil
}
