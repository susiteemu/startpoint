package model

import (
	"testing"
)

func TestStarlarkRequestDocString(t *testing.T) {

	starlarkRequest := RequestMold{
		Starlark: &StarlarkRequest{
			Script: `"""
meta:name: Starlark request
meta:prev_req: Some previous request
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
	if starlarkRequest.Name() != wantedName {
		t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Name(), wantedName)
	}
	/*
		wantedPrevReq := "Some previous request"
		if starlarkRequest. != wantedPrevReq {
			t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Name(), wantedName)
		}
	*/
	wantedUrl := "http://foobar.com"
	if starlarkRequest.Url() != wantedUrl {
		t.Errorf("url is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Url(), wantedUrl)
	}

	wantedMethod := "POST"
	if starlarkRequest.Method() != wantedMethod {
		t.Errorf("method is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Method(), wantedMethod)
	}

}
func TestStarlarkRequestDocStringMissingName(t *testing.T) {

	starlarkRequest := RequestMold{
		Starlark: &StarlarkRequest{
			Script: `"""
meta:name:
meta:prev_req: Some previous request
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
	if starlarkRequest.Name() != wantedName {
		t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Name(), wantedName)
	}
	/*
		wantedPrevReq := "Some previous request"
		if starlarkRequest. != wantedPrevReq {
			t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Name(), wantedName)
		}
	*/
	wantedUrl := "http://foobar.com"
	if starlarkRequest.Url() != wantedUrl {
		t.Errorf("url is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Url(), wantedUrl)
	}

	wantedMethod := "POST"
	if starlarkRequest.Method() != wantedMethod {
		t.Errorf("method is not equal!\ngot\n%v\nwanted\n%v", starlarkRequest.Method(), wantedMethod)
	}

}

func TestYamlRequestNameRename(t *testing.T) {

	yamlRequest := RequestMold{
		Yaml: &YamlRequest{
			Name:   "A GET request",
			Url:    "https://httpbin.org/anything",
			Method: "GET",
		},
		Raw: `name: A GET request
url: 'https://httpbin.org/anything'
method: GET
headers:
  X-Foo-Bar: SomeValue
body: >
  {
    "id": 1,
    "name": "Jane"
  }`,
	}

	ok := yamlRequest.Rename("A changed GET request")
	if !ok {
		t.Error("result is not true")
	}
	wantedName := "A changed GET request"
	if yamlRequest.Name() != wantedName {
		t.Errorf("name is not equal!\ngot\n%v\nwanted\n%v", yamlRequest.Name(), wantedName)
	}

	wantedRaw := `name: A changed GET request
url: 'https://httpbin.org/anything'
method: GET
headers:
  X-Foo-Bar: SomeValue
body: >
  {
    "id": 1,
    "name": "Jane"
  }`

	if wantedRaw != yamlRequest.Raw {
		t.Errorf("raw is not equal!\ngot\n%v\nwanted\n%v", yamlRequest.Raw, wantedRaw)
	}
}
