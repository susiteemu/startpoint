package starlarkng

import (
	"errors"
	"goful/core/model"
	"goful/core/scripting/starlark/goconv"
	"goful/core/scripting/starlark/starlarkconv"
	"net/http"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func RunStarlarkScript(request model.RequestMold, previousResponse *http.Response, profile model.Profile) (map[string]interface{}, error) {

	if request.Starlark == nil {
		return nil, errors.New("starlark request must not be nil")
	}

	profileValues, err := starlarkconv.Convert(profile.Variables)
	if err != nil {
		// TODO handle err
		return nil, err
	}

	predeclared := starlark.StringDict{
		"profile": profileValues,
	}

	thread := &starlark.Thread{Name: "starlark runner thread"}

	// TODO read from config
	fileOptions := syntax.FileOptions{
		Set:               true,
		While:             true,
		TopLevelControl:   true,
		GlobalReassign:    true,
		LoadBindsGlobally: true,
		Recursion:         true,
	}

	starlarkRequest := request.Starlark

	globals, _ := starlark.ExecFileOptions(&fileOptions, thread, request.Name(), starlarkRequest.Script, predeclared)
	values := make(map[string]interface{})
	for _, name := range globals.Keys() {
		starlarkValue := globals[name]
		goValue, err := goconv.ConvertValue(starlarkValue)
		if err != nil {
			return nil, err
		}
		values[name] = goValue
	}

	return values, nil
}
