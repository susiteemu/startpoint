package runner

import (
	"errors"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/model"
	"time"

	"github.com/rs/zerolog/log"
)

func RunRequestChain(reqs []*model.RequestMold, profile *model.Profile, interimResultCb func(took time.Duration, statusCode int)) ([]*model.Response, error) {

	if reqs == nil {
		return nil, errors.New("Requests must not be nil")
	}
	if profile == nil {
		// if passed nil profile we create empty one
		profile = &model.Profile{}
	}

	log.Debug().Msgf("About to run request chain of length %d with profile %v", len(reqs), profile)

	var responses []*model.Response

	var prevResponse *model.Response
	for _, r := range reqs {
		log.Debug().Msgf("Building request %v", r)
		var request model.Request
		var err error
		if prevResponse != nil {
			request, err = builder.BuildRequestUsingPreviousResponse(r, *prevResponse, *profile)
		} else {
			request, err = builder.BuildRequest(r, *profile)
		}

		if err != nil {
			log.Error().Err(err).Msgf("Building request failed with %v", r)
			return responses, err
		}

		response, err := client.DoRequest(request)
		if err != nil {
			log.Error().Err(err).Msgf("Request failed with %v", request)
			return responses, err
		}

		interimResultCb(response.Time, response.StatusCode)

		responses = append(responses, response)
		prevResponse = response
	}

	return responses, nil

}
