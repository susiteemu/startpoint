package model

import (
	"testing"
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
	if starlarkRequest.Name != wantedName {
		t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Name, wantedName)
	}

	wantedPrevReq := "Some previous request"
	if starlarkRequest.PreviousReq() != wantedPrevReq {
		t.Errorf("prev_req is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.PreviousReq(), wantedPrevReq)
	}

	wantedUrl := "http://foobar.com"
	if starlarkRequest.Url() != wantedUrl {
		t.Errorf("url is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Url(), wantedUrl)
	}

	wantedMethod := "POST"
	if starlarkRequest.Method() != wantedMethod {
		t.Errorf("method is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Method(), wantedMethod)
	}

	wantedOutput := "./output.txt"
	if starlarkRequest.Output() != wantedOutput {
		t.Errorf("output is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Output(), wantedOutput)
	}
}

func TestStarlarkRequestDocStringMissingParts(t *testing.T) {

	starlarkRequest := RequestMold{
		Starlark: &StarlarkRequest{
			Script: `"""
prev_req: Some previous request
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

	wantedName := ""
	if starlarkRequest.Name != wantedName {
		t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Name, wantedName)
	}

	wantedPrevReq := "Some previous request"
	if starlarkRequest.PreviousReq() != wantedPrevReq {
		t.Errorf("prev req is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.PreviousReq(), wantedPrevReq)
	}

	wantedUrl := "http://foobar.com"
	if starlarkRequest.Url() != wantedUrl {
		t.Errorf("url is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Url(), wantedUrl)
	}

	wantedMethod := "POST"
	if starlarkRequest.Method() != wantedMethod {
		t.Errorf("method is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Method(), wantedMethod)
	}

	wantedOutput := ""
	if starlarkRequest.Output() != wantedOutput {
		t.Errorf("output is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Output(), wantedOutput)
	}

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
	if starlarkRequest.Url() != wantedUrl {
		t.Errorf("url is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Url(), wantedUrl)
	}

	wantedMethod := "POST"
	if starlarkRequest.Method() != wantedMethod {
		t.Errorf("method is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Method(), wantedMethod)
	}

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

	if starlarkRequest.PreviousReq() != "Some previous request" {
		t.Errorf("previous request is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.PreviousReq(), "Some previous request")
		return
	}

	starlarkRequest.ChangePreviousReq("Some other previous request")

	if starlarkRequest.PreviousReq() != "Some other previous request" {
		t.Errorf("previous request is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.PreviousReq(), "Some other previous request")
		return
	}

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
	if yamlRequest.PreviousReq() != "Some previous request" {
		t.Errorf("previous request is not equal!\ngot\n%v\nwanted\n%v", yamlRequest.PreviousReq(), "Some previous request")
		return
	}

	yamlRequest.ChangePreviousReq("Some other previous request")

	if yamlRequest.PreviousReq() != "Some other previous request" {
		t.Errorf("previous request is not equal!\ngot\n%v\nwanted\n%v", yamlRequest.PreviousReq(), "Some other previous request")
		return
	}

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
	if yamlRequest.PreviousReq() != "Some previous request" {
		t.Errorf("previous request is not equal!\ngot\n%v\nwanted\n%v", yamlRequest.PreviousReq(), "Some previous request")
		return
	}

	yamlRequest.ChangePreviousReq("Some other previous request")

	if yamlRequest.PreviousReq() != "Some other previous request" {
		t.Errorf("previous request is not equal!\ngot\n%v\nwanted\n%v", yamlRequest.PreviousReq(), "Some other previous request")
		return
	}

}
