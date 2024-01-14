package builder

import (
	"goful/core/loader"
	"goful/core/model"
	starlarkng "goful/core/scripting/starlark"
	"goful/core/templating/yamlng"
	"net/http"
	"reflect"
)

var builders = []func(requestMetadata model.RequestMetadata, previousResponse *http.Response, profile model.Profile) (model.Request, bool, error){
	buildYamlRequest,
	buildStarlarkRequest,
}

func BuildRequest(requestMetadata model.RequestMetadata, profile model.Profile) (model.Request, error) {
	var request model.Request
	for _, builder := range builders {
		result, accept, err := builder(requestMetadata, nil, profile)
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
func BuildRequestUsingPreviousResponse(requestMetadata model.RequestMetadata, previousResponse *http.Response, profile model.Profile) (model.Request, error) {
	return model.Request{}, nil
}

func buildYamlRequest(requestMetadata model.RequestMetadata, _ *http.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMetadata.Request != "yaml" && requestMetadata.Request != "yml" {
		return model.Request{}, false, nil
	}
	request, err := loader.ReadYamlRequest(requestMetadata)
	if err != nil {
		return model.Request{}, true, err
	}

	if len(profile.Variables) > 0 {
		headers := request.Headers
		for k, v := range profile.Variables {
			request.Url = yamlng.ProcessTemplateVariable(request.Url, k, v)
			for headerName, headerValues := range request.Headers {
				headers[headerName] = yamlng.ProcessTemplateVariables(headerValues, k, v)
			}
		}
		request.Headers = headers
	}

	return request, true, err
}

func buildStarlarkRequest(requestMetadata model.RequestMetadata, previousResponse *http.Response, profile model.Profile) (model.Request, bool, error) {
	if requestMetadata.Request != "star" {
		return model.Request{}, false, nil
	}

	res, err := starlarkng.RunStarlarkScript(requestMetadata, previousResponse, profile)
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
