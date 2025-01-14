package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type Profile struct {
	Name              string
	Variables         map[string]string
	Raw               string
	Root              string
	Filename          string
	HasPublicProfile  bool
	HasPrivateProfile bool
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

func (p *Profile) IsPrivateProfile() bool {
	if p == nil {
		return false
	}
	return strings.HasSuffix(p.Filename, ".local")
}

func (p *Profile) IsDefaultProfile() bool {
	if p == nil {
		return false
	}
	return p.Name == "default"
}
func (p *Profile) IsDefaultPrivateProfile() bool {
	if p == nil {
		return false
	}
	return p.IsPrivateProfile() && p.Name == "default.local"
}
