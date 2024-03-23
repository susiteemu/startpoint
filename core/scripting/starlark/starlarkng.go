package starlarkng

import (
	"errors"
	"fmt"
	"goful/core/model"
	"goful/core/scripting/starlark/goconv"
	"goful/core/scripting/starlark/starlarkconv"

	"github.com/rs/zerolog/log"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func RunStarlarkScript(request model.RequestMold, previousResponse model.Response, profile model.Profile) (map[string]interface{}, error) {

	log.Info().Msgf("Running Starlark script with request %v, previousResponse %v, profile %v", request, previousResponse, profile)

	if request.Starlark == nil {
		log.Error().Msg("Starlark request is nil, aborting")
		return nil, errors.New("starlark request must not be nil")
	}

	profileValues, err := starlarkconv.Convert(profile.Variables)
	if err != nil {
		return nil, err
	}

	var previousResponseValues starlark.Value
	previousResponseValues, err = starlarkconv.Convert(previousResponse)
	if err != nil {
		return nil, err
	}

	predeclared := starlark.StringDict{
		"profile":          profileValues,
		"previousResponse": previousResponseValues,
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

	globals, err := starlark.ExecFileOptions(&fileOptions, thread, request.Name(), starlarkRequest.Script, predeclared)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	log.Debug().Msgf("Run Starlark script and got result %v", globals)

	values := make(map[string]interface{})
	for _, name := range globals.Keys() {
		starlarkValue := globals[name]
		goValue, err := goconv.ConvertValue(starlarkValue)
		if err != nil {
			return nil, err
		}
		values[name] = goValue
	}

	log.Debug().Msgf("Starlark result converted to Golang values %v", values)

	return values, nil
}
