package profileui

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func checkProfileWithNameDoesNotExist(m Model) func(s string) error {
	return func(s string) error {
		log.Debug().Msgf("Validating %s against list of existing profiles %d", s, len(m.list.Items()))
		for _, item := range m.list.Items() {
			r := item.(Profile)
			if r.Name == s {
				return errors.New("Profile with the same name already exists.")
			}
		}
		return nil
	}
}
