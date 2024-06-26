package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func (p *Profile) AsDotEnv() string {
	if p == nil {
		return ""
	}

	if p.Variables == nil {
		return ""
	}

	contents := []string{}
	for k, v := range p.Variables {
		contents = append(contents, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(contents, "\n")
}
