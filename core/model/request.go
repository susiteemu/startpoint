package model

import (
	"fmt"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"
)

type Request struct {
	Url     string
	Method  string
	Headers Headers
	Body    Body
}

type RequestMold struct {
	Yaml        *YamlRequest
	Starlark    *StarlarkRequest
	ContentType string
	Raw         string
	Root        string
	Filename    string
}

type YamlRequest struct {
	Name    string
	PrevReq string  `yaml:"prev_req"`
	Url     string  `yaml:"url"`
	Method  string  `yaml:"method"`
	Headers Headers `yaml:"headers"`
	Body    Body    `yaml:"body"`
}

type StarlarkRequest struct {
	Script string
}

func (r *RequestMold) Name() string {
	if r.Yaml != nil {
		return r.Yaml.Name
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(".*meta:name:\\s*(.*)")
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			return match[1]
		}
	}
	return ""
}

func (r *RequestMold) Url() string {
	var url = ""
	if r.Yaml != nil {
		url = r.Yaml.Url
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(".*doc:url:\\s*(.*)")
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			url = match[1]
		}
	}
	if url != "" {
		return url
	}
	return ""

}

func (r *RequestMold) Method() string {
	var method = ""
	if r.Yaml != nil {
		method = r.Yaml.Method
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(".*doc:method:\\s*(.*)")
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			method = match[1]
		}
	}
	if method != "" {
		return method
	}
	return ""

}

func (r *RequestMold) DeleteFromFS() bool {
	err := os.Remove(fmt.Sprintf("%s/%s", r.Root, r.Filename))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to remove file %s", r.Filename)
		return false
	}
	return true
}
