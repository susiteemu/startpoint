package model

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Profile struct {
	Name      string
	Variables map[string]string
	Raw       string
	Root      string
	Filename  string
}

func (p *Profile) DeleteFromFS() bool {
	err := os.Remove(filepath.Join(p.Root, p.Filename))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to remove file %s", p.Filename)
		return false
	}
	return true
}
