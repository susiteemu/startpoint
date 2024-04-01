package requestui

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func checkRequestWithNameDoesNotExist(m uiModel) func(s string) error {
	return func(s string) error {
		log.Debug().Msgf("Validating %s against list of existing requests %d", s, len(m.list.Items()))
		for _, item := range m.list.Items() {
			r := item.(Request)
			if r.Name == s {
				return errors.New("Request with the same name already exists.")
			}
		}
		return nil
	}
}
