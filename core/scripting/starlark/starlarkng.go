package starlarkng

import (
	"errors"
	"startpoint/core/model"
	"startpoint/core/scripting/starlark/goconv"
	"startpoint/core/scripting/starlark/starlarkconv"

	"github.com/rs/zerolog/log"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func RunStarlarkScript(request model.RequestMold, previousResponse *model.Response) (map[string]interface{}, error) {

	log.Info().Msgf("Running Starlark script with request %v, previousResponse %v", request, previousResponse)

	if request.Scriptable == nil {
		log.Error().Msg("Starlark request is nil, aborting")
		return nil, errors.New("starlark request must not be nil")
	}

	previousResponseStarlark := starlark.Dict{}
	if previousResponse != nil {
		prevResponseHeaders, err := starlarkconv.Convert(previousResponse.HeadersAsMapString())
		if err != nil {
			return nil, err
		}
		previousResponseStarlark.SetKey(starlark.String("headers"), prevResponseHeaders)

		// convert body to map if possible
		bodyAsMap, err := previousResponse.BodyAsMap()
		if err == nil {
			prevResponseBody, err := starlarkconv.Convert(bodyAsMap)
			if err != nil {
				return nil, err
			}
			previousResponseStarlark.SetKey(starlark.String("body"), prevResponseBody)
			log.Debug().Msgf("previousResponseBody %v", prevResponseBody)
		} else {
			log.Warn().Err(err).Msgf("Could not convert body to map. Setting body as string to previous response.")
			prevResponseBody, err := starlarkconv.Convert(string(previousResponse.Body))
			if err != nil {
				return nil, err
			}
			previousResponseStarlark.SetKey(starlark.String("body"), prevResponseBody)
			log.Debug().Msgf("previousResponseBody %v", prevResponseBody)
		}

	}

	predeclared := starlark.StringDict{
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

	starlarkRequest := request.Scriptable

	globals, err := starlark.ExecFileOptions(&fileOptions, thread, request.Name, starlarkRequest.Script, predeclared)
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
