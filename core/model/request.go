package model

type Request struct {
	Url     string  `yaml:"url"`
	Method  string  `yaml:"method"`
	Headers Headers `yaml:"headers"`
	Body    Body    `yaml:"body"`
}

type RequestMetadata struct {
	Name       string
	PrevReq    string `yaml:"prev_req"`
	Request    string `yaml:"request"`
	WorkingDir string
}
