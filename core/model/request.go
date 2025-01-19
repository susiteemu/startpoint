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
	Yaml       *YamlRequest
	Scriptable *ScriptableRequest
	Type       string
	Root       string
	Name       string
	Filename   string
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

type ScriptableRequest struct {
	Script string
}

func (r *Request) IsForm() bool {
	contentType, ok := r.ContentType()
	if !ok {
		return false
	}
	return strings.ToLower(contentType) == CONTENT_TYPE_FORM_URLENCODED
}

func (r *Request) IsMultipartForm() bool {
	contentType, ok := r.ContentType()
	if !ok {
		return false
	}
	return strings.ToLower(strings.TrimSpace(contentType)) == CONTENT_TYPE_MULTIPART_FORM
}

func (r *Request) HasBodyAsMap() bool {
	if r.Body == nil {
		return false
	}
	_, yes := r.Body.(map[string]interface{})
	if !yes {
		_, yes := r.Body.(map[string][]string)
		if !yes {
			_, yes := r.Body.(map[string]string)
			return yes
		}
	}
	return true
}

func (r *Request) ContentType() (string, bool) {
	contentType, ok := r.Headers[HEADER_NAME_CONTENT_TYPE]
	if !ok {
		return "", false
	}
	if len(contentType) == 0 {
		return "", false
	}
	return strings.Split(contentType[0], ";")[0], true
}

func (r *Request) BodyAsMap() (map[string]string, bool) {
	asMapString := make(map[string]string)
	asMapInterface, ok := r.Body.(map[string]interface{})
	if ok {
		for k, v := range asMapInterface {
			asInt, isInt := v.(int)
			if isInt {
				asMapString[k] = strconv.Itoa(asInt)
			} else {
				asArr, isArr := v.([]interface{})
				if isArr {
					asStrArr := []string{}
					for _, i := range asArr {
						asInt, isInt := i.(int)
						if isInt {
							asStrArr = append(asStrArr, strconv.Itoa(asInt))
						} else {
							asStrArr = append(asStrArr, i.(string))
						}
					}
					asMapString[k] = strings.Join(asStrArr, ", ")
				} else {
					asMapString[k] = v.(string)
				}
			}
		}
	} else {
		asMapStringArr, ok := r.Body.(map[string][]string)
		if !ok {
			return map[string]string{}, false
		}
		for k, v := range asMapStringArr {
			asMapString[k] = strings.Join(v, ", ")
		}
	}
	return asMapString, true
}

func (r *RequestMold) Url() string {
	if r.Yaml != nil {
		return r.Yaml.Url
	} else if r.Scriptable != nil {
		switch r.Type {
		case CONTENT_TYPE_STARLARK, CONTENT_TYPE_LUA:
			return extractValueFromAlternativeFieldNames(r.Scriptable.Script, starlarkUrlFields)
		}
	}
	return ""
}

func (r *RequestMold) Method() string {
	if r.Yaml != nil {
		return r.Yaml.Method
	} else if r.Scriptable != nil {
		switch r.Type {
		case CONTENT_TYPE_STARLARK, CONTENT_TYPE_LUA:
			return extractValueFromAlternativeFieldNames(r.Scriptable.Script, starlarkMethodFields)
		}
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
	} else if r.Scriptable != nil {
		return r.Scriptable.Script
	}
	return ""
}

func (r *RequestMold) PreviousReq() string {
	if r.Yaml != nil {
		return r.Yaml.PrevReq
	} else if r.Scriptable != nil {
		switch r.Type {
		case CONTENT_TYPE_STARLARK, CONTENT_TYPE_LUA:
			return extractValueFromAlternativeFieldNames(r.Scriptable.Script, starlarkPrevReqFields)
		}
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
	} else if r.Scriptable != nil {
		switch r.Type {
		case CONTENT_TYPE_STARLARK, CONTENT_TYPE_LUA:
			pattern := regexp.MustCompile(`(?mU)^prev_req:(.*)$`)
			changed := pattern.ReplaceAllString(r.Scriptable.Script, fmt.Sprintf("prev_req: %s", prevReq))
			r.Scriptable.Script = changed
		}
	}
}

func (r *RequestMold) Output() string {
	if r.Yaml != nil {
		return r.Yaml.Output
	} else if r.Scriptable != nil {
		switch r.Type {
		case CONTENT_TYPE_STARLARK, CONTENT_TYPE_LUA:
			return extractValueFromAlternativeFieldNames(r.Scriptable.Script, starlarkOutputFields)
		}
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
		Type:     r.Type,
		Root:     r.Root,
		Filename: r.Filename,
		Name:     r.Name,
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
	} else if r.Scriptable != nil {
		copy.Scriptable = &ScriptableRequest{
			Script: r.Scriptable.Script,
		}
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
		INITIAL = iota
		START_MATCHING_FIELD
		START_DETECTING_ASSIGNMENT
		ASSIGNMENT_DETECTED
		ASSIGNMENT_START_DOUBLE_QUOTES_DETECTED
		ASSIGNMENT_START_SINGLE_QUOTES_DETECTED
		START_CAPTURING
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
		assignmentHasDoubleQuotes := false
		assignmentHasSingleQuotes := false
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
				if !unicode.IsSpace(c) && c != '"' && c != '\'' {
					state = START_CAPTURING
				} else if c == '"' {
					state = ASSIGNMENT_START_DOUBLE_QUOTES_DETECTED
				} else if c == '\'' {
					state = ASSIGNMENT_START_SINGLE_QUOTES_DETECTED
				}
			case ASSIGNMENT_START_DOUBLE_QUOTES_DETECTED:
				if c != '"' {
					assignmentHasDoubleQuotes = true
					state = START_CAPTURING
				}
			case ASSIGNMENT_START_SINGLE_QUOTES_DETECTED:
				if c != '\'' {
					assignmentHasSingleQuotes = true
					state = START_CAPTURING
				}
			case START_CAPTURING:
				if assignmentHasDoubleQuotes && c == '"' {
					breakLoop = true
				} else if assignmentHasSingleQuotes && c == '\'' {
					breakLoop = true
				}
			}

			if breakLoop {
				break
			}

			if state == START_CAPTURING {
				capture = append(capture, c)
			}

		}
		if state == START_CAPTURING {
			match = true
			break
		}
	}
	return string(capture), match
}
