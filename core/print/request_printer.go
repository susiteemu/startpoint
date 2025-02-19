package print

import (
	"errors"
	"fmt"
	"strings"

	"github.com/susiteemu/startpoint/core/model"

	"github.com/rs/zerolog/log"
)

func SprintRequest(request *model.Request, pretty bool) (string, string, error) {
	if request == nil {
		return "", "", errors.New("Request must not be nil!")
	}

	var requestBuilder []string
	methodLabel := request.Method
	urlLabel := request.Url
	requestBuilder = append(requestBuilder, fmt.Sprintf("%s %s", methodLabel, urlLabel))
	headers, _, err := SprintHeaders(request.Headers, false)
	if err != nil {
		return "", "", err
	}
	requestBuilder = append(requestBuilder, headers)

	if request.IsForm() || request.IsMultipartForm() || request.HasBodyAsMap() {
		if request.Body != nil {
			bodyAsMap, ok := request.BodyAsMap()
			log.Debug().Msgf("Body as map %v", bodyAsMap)
			if !ok {
				return "", "", errors.New("cannot convert body to map")
			}
			for k, v := range bodyAsMap {
				requestBuilder = append(requestBuilder, fmt.Sprintf("%v: %v", k, v))
			}
			if len(bodyAsMap) > 0 {
				requestBuilder = append(requestBuilder, "")
			}
		}
	} else {
		if request.Body != nil {
			log.Debug().Msgf("Request body type %T, value=%v, headers=%v", request.Body, request.Body, request.Headers)
			bodyAsStr, ok := request.Body.(string)
			if !ok {
				return "", "", errors.New("cannot convert request body to string")
			}
			bodyBytes := []byte(bodyAsStr)
			printedBody, _, err := SprintBody(int64(len(bodyBytes)), bodyBytes, request.Headers, false)
			if err != nil {
				log.Warn().Err(err).Msg("Error occurred while printing body")
				return "", "", err
			}
			if len(printedBody) > 0 {
				requestBuilder = append(requestBuilder, printedBody)
				requestBuilder = append(requestBuilder, "")
			}
		}
	}

	if len(request.Output) > 0 {
		requestBuilder = append(requestBuilder, fmt.Sprintf("Save response to %s", request.Output))
		requestBuilder = append(requestBuilder, "")
	}

	printed := strings.Join(requestBuilder, "\n")
	prettyPrinted := ""
	if pretty {
		prettyPrinted = SprintFaint(printed)
	}

	return printed, prettyPrinted, nil
}
