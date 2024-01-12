package starlarkng

import (
	"fmt"
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
		fmt.Printf(">>> ERRRRRRR")
		// TODO handle err
		return nil, err
	}

	fmt.Printf(">>> GOING FORWARD %v\n", metadata)

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
	fmt.Printf(">>> path %v\n", metadata.ToRequestPath())
	globals, _ := starlark.ExecFileOptions(&fileOptions, thread, metadata.ToRequestPath(), nil, predeclared)
	fmt.Printf(">>> GOT VALUES %v\n", globals.Keys())
	values := make(map[string]interface{})
	for _, name := range globals.Keys() {
		starlarkValue := globals[name]
		goValue, err := goconv.ConvertValue(starlarkValue)
		fmt.Printf(">>> name %s, value %v", name, goValue)
		if err != nil {
			return nil, err
		}
		values[name] = goValue
	}

	return values, nil
}
