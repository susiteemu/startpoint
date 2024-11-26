package model

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const CONTENT_TYPE_YAML = "yaml"
const CONTENT_TYPE_STARLARK = "star"
const HEADER_NAME_AUTHORIZATION = "Authorization"
const HEADER_VALUE_BASIC_AUTH = "Basic"
const HEADER_VALUE_BEARER_AUTH = "Bearer"

var (
	starlarkNameFields = []string{
		"meta:name",
	}
	starlarkUrlFields = []string{
		"doc:url",
		"url",
	}
	starlarkMethodFields = []string{
		"doc:method",
		"method",
	}
	starlarkPrevReqFields = []string{
		"prev_req",
	}
	starlarkOutputFields = []string{
		"meta:output",
	}
)

type BasicAuth struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Auth struct {
	Basic  BasicAuth `yaml:"basic_auth"`
	Bearer string    `yaml:"bearer_token"`
}

type Request struct {
	Headers Headers
	Options map[string]interface{}
	Body    Body
	Url     string
	Method  string
	Output  string
}

type RequestMold struct {
	Yaml        *YamlRequest
	Starlark    *StarlarkRequest
	ContentType string
	Root        string
	Name        string
	Filename    string
}

type YamlRequest struct {
	PrevReq string                 `yaml:"prev_req,omitempty"`
	Url     string                 `yaml:"url"`
	Method  string                 `yaml:"method"`
	Headers Headers                `yaml:"headers,omitempty"`
	Body    Body                   `yaml:"body,omitempty"`
	Output  string                 `yaml:"output,omitempty"`
	Options map[string]interface{} `yaml:"options,omitempty"`
	Raw     string                 `yaml:"raw,omitempty"`
	Auth    Auth                   `yaml:"auth,omitempty"`
}

type StarlarkRequest struct {
	Script string
}

func (r *Request) IsForm() bool {
	contentType, ok := r.ContentType()
	if !ok {
		return false
	}
	return strings.ToLower(contentType) == "application/x-www-form-urlencoded"
}

func (r *Request) IsMultipartForm() bool {
	contentType, ok := r.ContentType()
	if !ok {
		return false
	}
	return strings.ToLower(strings.TrimSpace(contentType)) == "multipart/form-data"
}

func (r *Request) ContentType() (string, bool) {
	contentType, ok := r.Headers["Content-Type"]
	if !ok {
		return "", false
	}
	if len(contentType) == 0 {
		return "", false
	}
	return strings.Split(contentType[0], ";")[0], true
}

func (r *Request) BodyAsMap() (map[string]string, bool) {
	asMapInterface, ok := r.Body.(map[string]interface{})
	if !ok {
		return map[string]string{}, false
	}
	asMapString := make(map[string]string)
	for k, v := range asMapInterface {
		asInt, isInt := v.(int)
		if isInt {
			asMapString[k] = strconv.Itoa(asInt)
		} else {
			asMapString[k] = v.(string)
		}
	}
	return asMapString, true
}

func (r *RequestMold) Url() string {
	if r.Yaml != nil {
		return r.Yaml.Url
	} else if r.Starlark != nil {
		return extractValueFromAlternativeFieldNames(r.Starlark.Script, starlarkUrlFields)
	}
	return ""
}

func (r *RequestMold) Method() string {
	if r.Yaml != nil {
		return r.Yaml.Method
	} else if r.Starlark != nil {
		return extractValueFromAlternativeFieldNames(r.Starlark.Script, starlarkMethodFields)
	}
	return ""
}

func (r *RequestMold) Raw() string {
	if r.Yaml != nil {
		if r.Yaml.Raw == "" {
			asYaml, err := yaml.Marshal(r.Yaml)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal YAML request to YAML")
			} else {
				r.Yaml.Raw = string(asYaml)
			}
		}
		return r.Yaml.Raw
	} else if r.Starlark != nil {
		return r.Starlark.Script
	}
	return ""
}

func (r *RequestMold) PreviousReq() string {
	if r.Yaml != nil {
		return r.Yaml.PrevReq
	} else if r.Starlark != nil {
		return extractValueFromAlternativeFieldNames(r.Starlark.Script, starlarkPrevReqFields)
	}
	return ""
}

/*
* NOTE: assumes previous request is set before; if previous request is not set, this can't add it to the raw versions
 */
func (r *RequestMold) ChangePreviousReq(prevReq string) {
	if r.Yaml != nil {
		r.Yaml.PrevReq = prevReq
		pattern := regexp.MustCompile(`(?mU)^prev_req:(.*)$`)
		changed := pattern.ReplaceAllString(r.Yaml.Raw, fmt.Sprintf("prev_req: \"%s\"", prevReq))
		r.Yaml.Raw = changed
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(`(?mU)^prev_req:(.*)$`)
		changed := pattern.ReplaceAllString(r.Starlark.Script, fmt.Sprintf("prev_req: %s", prevReq))
		r.Starlark.Script = changed
	}
}

func (r *RequestMold) Output() string {
	if r.Yaml != nil {
		return r.Yaml.Output
	} else if r.Starlark != nil {
		return extractValueFromAlternativeFieldNames(r.Starlark.Script, starlarkOutputFields)
	}
	return ""
}

func (r *RequestMold) DeleteFromFS() bool {
	err := os.Remove(filepath.Join(r.Root, r.Filename))
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
		Name:        r.Name,
	}

	if r.Yaml != nil {
		yamlRequest := YamlRequest{
			PrevReq: r.Yaml.PrevReq,
			Url:     r.Yaml.Url,
			Method:  r.Yaml.Method,
			Headers: r.Yaml.Headers,
			Body:    r.Yaml.Body,
			Output:  r.Yaml.Output,
			Raw:     r.Yaml.Raw,
			Auth:    r.Yaml.Auth,
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

func extractValueFromAlternativeFieldNames(str string, fields []string) string {
	for _, field := range fields {
		match, has := extractValueFromField(field, str)
		if has {
			return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(match), "\"", ""), "'", "")
		}
	}
	return ""
}

func extractValueFromField(field string, str string) (string, bool) {
	const (
		INITIAL                    = -1
		START_MATCHING_FIELD       = 0
		START_DETECTING_ASSIGNMENT = 1
		ASSIGNMENT_DETECTED        = 2
		START_CAPTURING            = 3
	)
	var (
		fieldRunes       = []rune(field)
		fieldMatchIdxMax = len(fieldRunes)
		capture          = []rune{}
		match            = false
	)

	str = strings.ReplaceAll(str, "\r\n", "\n")
	lines := strings.Split(str, "\n")
	for lineNr, line := range lines {
		if len(line) < len(field) {
			continue
		}
		state := INITIAL
		fieldMatchIdxPos := 0
		capture = []rune{}
		breakLoop := false
		for idx, c := range line {
			if unicode.IsSpace(c) && state == INITIAL {
				continue
			}
			if !unicode.IsSpace(c) && state == INITIAL {
				state = START_MATCHING_FIELD
			}
			switch state {
			case START_MATCHING_FIELD:
				if c == fieldRunes[fieldMatchIdxPos] {
					fieldMatchIdxPos++
				} else {
					fieldMatchIdxPos = 0
				}
				if fieldMatchIdxPos == fieldMatchIdxMax {
					state = START_DETECTING_ASSIGNMENT
				}
			case START_DETECTING_ASSIGNMENT:
				if !unicode.IsSpace(c) && c != '=' && c != ':' {
					log.Warn().Msgf("Encountered illegal character at position %d on line %d: %c", idx, lineNr, c)
					breakLoop = true
				} else if c == '=' || c == ':' {
					state = ASSIGNMENT_DETECTED
				}
			case ASSIGNMENT_DETECTED:
				if !unicode.IsSpace(c) {
					state = START_CAPTURING
				}
			}

			if state == START_CAPTURING {
				capture = append(capture, c)
			}

			if breakLoop {
				break
			}
		}
		if state == START_CAPTURING {
			match = true
			break
		}
	}
	return string(capture), match
}
