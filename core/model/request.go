package model

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

type Request struct {
	Url     string
	Method  string
	Headers Headers
	Body    Body
	Output  string
}

type RequestMold struct {
	Yaml        *YamlRequest
	Starlark    *StarlarkRequest
	ContentType string
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
	Output  string  `yaml:"output"`
	Raw     string
}

type StarlarkRequest struct {
	Script string
}

func (r *Request) IsForm() bool {
	contentType, ok := r.Headers["Content-Type"]
	if !ok {
		return false
	}
	if len(contentType) == 0 {
		return false
	}
	return strings.ToLower(contentType[0]) == "application/x-www-form-urlencoded"
}

func (r *Request) BodyAsMap() (map[string]string, bool) {
	asMapInterface, ok := r.Body.(map[string]interface{})
	if !ok {
		return map[string]string{}, false
	}
	asMapString := make(map[string]string)
	for k, v := range asMapInterface {
		asMapString[k] = v.(string)
	}
	return asMapString, true
}

func (r *RequestMold) Name() string {
	if r.Yaml != nil {
		return r.Yaml.Name
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(`(?mU)^.*meta:name:(.*)$`)
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			return strings.TrimSpace(match[1])
		}
	}
	return ""
}

func (r *RequestMold) Url() string {
	var url = ""
	if r.Yaml != nil {
		url = r.Yaml.Url
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(`(?mU)^.*doc:url:(.*)$`)
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			url = strings.TrimSpace(match[1])
		} else {
			pattern = regexp.MustCompile(`(?mU)^\s*url\s*=(.*)$`)
			match = pattern.FindStringSubmatch(r.Starlark.Script)
			if len(match) == 2 {
				url = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(match[1]), "\"", ""), "'", "")
			}
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
		pattern := regexp.MustCompile(`(?mU)^.*doc:method:(.*)$`)
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			method = strings.TrimSpace(match[1])
		} else {
			pattern = regexp.MustCompile(`(?mU)^\s*method\s*=(.*)$`)
			match = pattern.FindStringSubmatch(r.Starlark.Script)
			if len(match) == 2 {
				method = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(match[1]), "\"", ""), "'", "")
			}
		}
	}
	return method
}

func (r *RequestMold) Raw() string {
	if r.Yaml != nil {
		return r.Yaml.Raw
	} else if r.Starlark != nil {
		return r.Starlark.Script
	}
	return ""
}

func (r *RequestMold) PreviousReq() string {
	var prevReq = ""
	if r.Yaml != nil {
		prevReq = r.Yaml.PrevReq
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(`(?mU)^.*meta:prev_req:(.*)$`)
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			prevReq = strings.TrimSpace(match[1])
		}
	}
	return prevReq
}

func (r *RequestMold) DeleteFromFS() bool {
	err := os.Remove(fmt.Sprintf("%s/%s", r.Root, r.Filename))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to remove file %s", r.Filename)
		return false
	}
	return true
}

func (r *RequestMold) Clone() RequestMold {
	copy := RequestMold{
		ContentType: r.ContentType,
		Root:        r.Root,
		Filename:    r.Filename,
	}

	if r.Yaml != nil {
		yamlRequest := YamlRequest{
			Name:    r.Yaml.Name,
			PrevReq: r.Yaml.PrevReq,
			Url:     r.Yaml.Url,
			Method:  r.Yaml.Method,
			Headers: r.Yaml.Headers,
			Body:    r.Yaml.Body,
			Output:  r.Yaml.Output,
			Raw:     r.Yaml.Raw,
		}
		copy.Yaml = &yamlRequest
	} else if r.Starlark != nil {
		starlarkRequest := StarlarkRequest{
			Script: r.Starlark.Script,
		}
		copy.Starlark = &starlarkRequest
	}

	return copy
}
