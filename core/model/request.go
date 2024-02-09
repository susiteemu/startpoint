package model

import "regexp"

type Request struct {
	Url     string
	Method  string
	Headers Headers
	Body    Body
}

type RequestMold struct {
	Yaml     *YamlRequest
	Starlark *StarlarkRequest
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
	if r.Yaml != nil {
		return r.Yaml.Url
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(".*doc:url:\\s*(.*)")
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			return match[1]
		}
	}
	return ""

}

func (r *RequestMold) Method() string {
	if r.Yaml != nil {
		return r.Yaml.Method
	} else if r.Starlark != nil {
		pattern := regexp.MustCompile(".*doc:method:\\s*(.*)")
		match := pattern.FindStringSubmatch(r.Starlark.Script)
		if len(match) == 2 {
			return match[1]
		}
	}
	return ""

}
