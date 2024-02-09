package builder

import (
	"goful/core/model"
	starlarkng "goful/core/scripting/starlark"
	"goful/core/templating/yamlng"
	"net/http"
	"reflect"
)

var builders = []func(requestMold model.RequestMold, previousResponse *http.Response, profile model.Profile) (model.Request, bool, error){
	buildYamlRequest,
	buildStarlarkRequest,
}

func BuildRequest(requestMold model.RequestMold, profile model.Profile) (model.Request, error) {
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

/*
		test := make(map[string]bool)
		test["one"] = true
		test["two"] = false
		r, _, _ := starlarkconv.ConvertDict(test)

		predeclared := starlark.StringDict{
			"prev":    starlark.String("hello"),
			"profile": r,
		}

		thread := &starlark.Thread{Name: "my thread"}
		fileOptions := syntax.FileOptions{
			Set:               true,
			While:             true,
			TopLevelControl:   true,
			GlobalReassign:    true,
			LoadBindsGlobally: true,
			Recursion:         true,
		}
		globals, _ := starlark.ExecFileOptions(&fileOptions, thread, "starlark_request_r.star", nil, predeclared)
		fmt.Println("\nGlobals:")
		for _, name := range globals.Keys() {
			v := globals[name]
			//fmt.Printf("%s (%s) = %s\n", name, v.Type(), v.String())
			goValue, err := goconv.ConvertValue(v)
			if err != nil {
				fmt.Printf("error %v", err)
			} else {
				fmt.Printf("converted %v\n", goValue)
			}
		}

		return model.Request{}, nil
	}
*/
func BuildRequestUsingPreviousResponse(requestMold model.RequestMold, previousResponse *http.Response, profile model.Profile) (model.Request, error) {
	return model.Request{}, nil
}

func buildYamlRequest(requestMold model.RequestMold, _ *http.Response, profile model.Profile) (model.Request, bool, error) {
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

	// TODO map values from YamlRequest to Request
	return request, true, nil
}

func buildStarlarkRequest(requestMold model.RequestMold, previousResponse *http.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMold.Starlark == nil {
		return model.Request{}, false, nil
	}

	res, err := starlarkng.RunStarlarkScript(requestMold, previousResponse, profile)
	if err != nil {
		return model.Request{}, true, err
	}

	headers := make(map[string][]string)
	for k, v := range res["headers"].(map[string]interface{}) {
		t := reflect.TypeOf(v)
		if t.String() == "string" {
			headers[k] = []string{v.(string)}
		} else if t.String() == "[]interface {}" {
			vv := v.([]interface{})
			var l []string
			for _, vvv := range vv {
				l = append(l, vvv.(string))
			}
			headers[k] = l
		}
	}

	req := model.Request{
		Url:     res["url"].(string),
		Method:  res["method"].(string),
		Headers: new(model.Headers).FromMap(headers),
		Body:    res["body"],
	}

	return req, true, nil
}
