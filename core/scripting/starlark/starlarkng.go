package starlarkng

import (
	"encoding/json"
	"errors"
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

	// convert model.HeaderValues into []string to pass it to starlarkconv
	headers := make(map[string][]string)
	for k, v := range previousResponse.Headers {
		headers[k] = v
	}

	prevResponseHeaders, err := starlarkconv.Convert(headers)
	if err != nil {
		return nil, err
	}

	// convert body to map if possible
	// TODO check content-type: is application/json?
	var bodyAsMap map[string]interface{}
	err = json.Unmarshal(previousResponse.Body, &bodyAsMap)
	if err != nil {
		// TODO handle err
	}

	prevResponseBody, err := starlarkconv.Convert(bodyAsMap)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("previousResponseBody %v", prevResponseBody)

	previousResponseStarlark := starlark.Dict{}
	previousResponseStarlark.SetKey(starlark.String("body"), prevResponseBody)
	previousResponseStarlark.SetKey(starlark.String("headers"), prevResponseHeaders)

	predeclared := starlark.StringDict{
		"profile":      profileValues,
		"prevResponse": &previousResponseStarlark,
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
		log.Error().Err(err).Msg("Failed to exec starlark script")
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
