package model

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

type Body interface{}

type FormData interface{}
type HeaderValues []string

type Headers map[string]HeaderValues

/*func (body *Body) UnmarshalYAML(node *yaml.Node) error {
	value := node.Value
	ba := []byte(value)
	*body = ba
	return nil
}*/

func (headerValues *HeaderValues) UnmarshalYAML(node *yaml.Node) error {
	value := node.Value
	sl := strings.Split(value, ",")
	*headerValues = sl
	return nil
}

func (headerValues *HeaderValues) ToString() string {
	return strings.Join(*headerValues, ",")
}
func (headers *Headers) FromMap(m map[string][]string) Headers {
	responseHeaders := make(map[string]HeaderValues)
	for k, v := range m {
		responseHeaders[k] = v
	}
	*headers = responseHeaders
	return *headers
}

func (headers *Headers) ToMap() map[string]string {
	headerMap := make(map[string]string)
	for k, v := range *headers {
		headerMap[k] = v.ToString()
	}
	return headerMap
}

func (metadata *RequestMetadata) ToRequestPath() string {
	return filepath.Join(metadata.WorkingDir, fmt.Sprint(metadata.Name, "_r.", metadata.Request))
}
