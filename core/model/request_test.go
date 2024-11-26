package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStarlarkRequestDocString(t *testing.T) {

	starlarkRequest := RequestMold{
		Name: "Starlark request",
		Starlark: &StarlarkRequest{
			Script: `"""
prev_req: Some previous request
meta:output: ./output.txt
doc:url: http://foobar.com
doc:method: POST
"""
url = "http://foobar.com"
method = "POST"
headers = { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
body = { "id": 1474, "prev": prev, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
`},
	}

	wantedName := "Starlark request"
	assert.Equal(t, wantedName, starlarkRequest.Name)

	wantedPrevReq := "Some previous request"
	assert.Equal(t, wantedPrevReq, starlarkRequest.PreviousReq())

	wantedUrl := "http://foobar.com"
	assert.Equal(t, wantedUrl, starlarkRequest.Url())

	wantedMethod := "POST"
	assert.Equal(t, wantedMethod, starlarkRequest.Method())

	wantedOutput := "./output.txt"
	assert.Equal(t, wantedOutput, starlarkRequest.Output())
}

func TestStarlarkRequestDocStringMissingParts(t *testing.T) {

	starlarkRequest := RequestMold{
		Starlark: &StarlarkRequest{
			Script: `"""
doc:url: http://foobar.com
doc:method: POST
"""
#  prev_req: Some previous request
url = "http://foobar.com"
method = "POST"
headers = { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
body = { "id": 1474, "prev": prev, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
`},
	}

	wantedName := ""
	assert.Equal(t, wantedName, starlarkRequest.Name)

	wantedPrevReq := "Some previous request"
	assert.Equal(t, wantedPrevReq, starlarkRequest.PreviousReq())

	wantedUrl := "http://foobar.com"
	assert.Equal(t, wantedUrl, starlarkRequest.Url())

	wantedMethod := "POST"
	assert.Equal(t, wantedMethod, starlarkRequest.Method())

	wantedOutput := ""
	assert.Equal(t, wantedOutput, starlarkRequest.Output())
}

func TestStarlarkRequestParseValuesFromActualCode(t *testing.T) {

	starlarkRequest := RequestMold{
		Starlark: &StarlarkRequest{
			Script: `"""
prev_req: Some previous request
"""
url = "http://foobar.com"
method = "POST"
headers = { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
body = { "id": 1474, "prev": prev, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
`},
	}

	wantedUrl := "http://foobar.com"
	assert.Equal(t, wantedUrl, starlarkRequest.Url())

	wantedMethod := "POST"
	assert.Equal(t, wantedMethod, starlarkRequest.Method())
}

func TestChangePreviousRequest(t *testing.T) {

	starlarkRequest := RequestMold{
		Starlark: &StarlarkRequest{
			Script: `"""
prev_req: Some previous request
"""
url = "http://foobar.com"
method = "POST"
headers = { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
body = { "id": 1474, "prev": prev, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
`},
	}

	wantedPrevReq := "Some previous request"
	assert.Equal(t, wantedPrevReq, starlarkRequest.PreviousReq())

	starlarkRequest.ChangePreviousReq("Some other previous request")

	wantedPrevReq = "Some other previous request"
	assert.Equal(t, wantedPrevReq, starlarkRequest.PreviousReq())

	yamlRequest := RequestMold{
		Yaml: &YamlRequest{
			PrevReq: "Some previous request",
			Raw: `
prev_req: Some previous request,
url: http://foobar.com
method: POST
`,
		},
	}

	wantedPrevReq = "Some previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())

	yamlRequest.ChangePreviousReq("Some other previous request")

	wantedPrevReq = "Some other previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())

	yamlRequest = RequestMold{
		Yaml: &YamlRequest{
			PrevReq: "Some previous request",
			Raw: `
prev_req: "Some previous request",
url: http://foobar.com
method: POST
`,
		},
	}
	wantedPrevReq = "Some previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())

	yamlRequest.ChangePreviousReq("Some other previous request")

	wantedPrevReq = "Some other previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())
}
