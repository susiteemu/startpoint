package model

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

const CONTENT_TYPE_YAML = "yaml"
const CONTENT_TYPE_STARLARK = "star"
const CONTENT_TYPE_LUA = "lua"
const HEADER_NAME_AUTHORIZATION = "Authorization"
const HEADER_VALUE_BASIC_AUTH = "Basic"
const HEADER_VALUE_BEARER_AUTH = "Bearer"
const HEADER_NAME_CONTENT_TYPE = "Content-Type"
const CONTENT_TYPE_PLAINTEXT = "text/plain"
const CONTENT_TYPE_APPLICATION_JSON = "application/json"
const CONTENT_TYPE_APPLICATION_XML = "application/xml"
const CONTENT_TYPE_TEXT_HTML = "text/html"
const CONTENT_TYPE_FORM_URLENCODED = "application/x-www-form-urlencoded"
const CONTENT_TYPE_MULTIPART_FORM = "multipart/form-data"

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
func (headerValues HeaderValues) MarshalYAML() (interface{}, error) {
	return strings.Join(headerValues, ","), nil
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

func (headers *Headers) ContentType() (string, error) {
	for k, v := range *headers {
		if k == HEADER_NAME_CONTENT_TYPE {
			contentType := v[0]
			if strings.Contains(contentType, ";") {
				contentType = strings.Split(contentType, ";")[0]
			}
			return contentType, nil
		}
	}
	return "", errors.New("could not find content type")
}
