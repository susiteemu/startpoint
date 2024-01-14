package starlarkng

import (
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"goful/core/model"
	"goful/core/scripting/starlark/goconv"
	"goful/core/scripting/starlark/starlarkconv"
	"net/http"
)

func RunStarlarkScript(metadata model.RequestMetadata, previousResponse *http.Response, profile model.Profile) (map[string]interface{}, error) {

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
	globals, _ := starlark.ExecFileOptions(&fileOptions, thread, metadata.ToRequestPath(), nil, predeclared)
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
