package model

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

type Body []byte
type HeaderValues []string

func (pk *Body) UnmarshalYAML(node *yaml.Node) error {
	value := node.Value
	ba := []byte(value)
	*pk = ba
	return nil
}

func (pk *HeaderValues) UnmarshalYAML(node *yaml.Node) error {
	value := node.Value
	sl := strings.Split(value, ",")
	*pk = sl
	return nil
}

func (metadata *RequestMetadata) ToRequestPath() string {
	return filepath.Join(metadata.WorkingDir, fmt.Sprint(metadata.Name, "_r.", metadata.Request))
}

type Request struct {
	Url     string                  `yaml:"url"`
	Method  string                  `yaml:"method"`
	Headers map[string]HeaderValues `yaml:"headers"`
	Body    Body                    `yaml:"body"`
}

type RequestMetadata struct {
	Name       string
	PrevReq    string `yaml:"prev_req"`
	Request    string `yaml:"request"`
	WorkingDir string
}
