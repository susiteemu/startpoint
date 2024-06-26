package builder

import (
	"github.com/google/go-cmp/cmp"
	"math/big"
	"startpoint/core/model"
	"testing"
)

func TestBuildRequestYaml(t *testing.T) {

	requestMold := model.RequestMold{
		Name: "yaml_request",
		Yaml: &model.YamlRequest{
			Url:    "http://foobar.com",
			Method: "POST",
			Headers: model.Headers{
				"X-Foo-Bar": {"SomeValue"},
			},
			Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
			Output:  "",
			Options: make(map[string]interface{}),
		},
	}

	wantedRequest := model.Request{
		Url:    "http://foobar.com",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo-Bar": {"SomeValue"},
		},
		Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
		Output:  "",
		Options: make(map[string]interface{}),
	}

	request, err := BuildRequest(&requestMold, model.Profile{})
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest) {
		t.Errorf("got\n%q\nwanted\n%q\n", request, wantedRequest)
	}
}

func TestBuildRequestYamlWithTemplateVariables(t *testing.T) {
	requestMold := model.RequestMold{
		Name: "yaml_request",
		Yaml: &model.YamlRequest{
			Url:    "http://{domain}/api",
			Method: "POST",
			Headers: model.Headers{
				"X-Foo-Bar":  {"SomeValue"},
				"X-Tmpl-Var": {"{header-value-test}"},
			},
			Body:    "{\n  \"id\": 1,\n  \"name\": \"{name}\"\n}",
			Options: make(map[string]interface{}),
			Raw: `url: http://{domain}/api
method: POST
headers:
  X-Foo-Bar: SomeValue
  X-Tmpl-Var: {header-value-test}
body: >
  {
    "id": 1,
    "name": "{name}"
  }`,
		},
	}

	wantedRequest := model.Request{
		Url:    "http://prodfoobar.com/api",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo-Bar":  {"SomeValue"},
			"X-Tmpl-Var": {"Value from template var"},
		},
		Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
		Options: make(map[string]interface{}),
	}

	profile := model.Profile{
		Name: "test",
		Variables: map[string]string{
			"domain":            "prodfoobar.com",
			"header-value-test": "Value from template var",
			"name":              "Jane",
		},
	}

	request, err := BuildRequest(&requestMold, profile)
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest) {
		t.Errorf("got\n%q\nwanted\n%q\n", request, wantedRequest)
	}
}

func TestBuildFormRequestYamlWithTemplateVariables(t *testing.T) {
	requestMold := model.RequestMold{
		Name: "yaml_request",
		Yaml: &model.YamlRequest{
			Url:    "http://{domain}/api",
			Method: "POST",
			Headers: model.Headers{
				"X-Foo-Bar":  {"SomeValue"},
				"X-Tmpl-Var": {"{header-value-test}"},
			},
			Body: map[string]interface{}{
				"username": "{ name }",
				"password": "{password}",
			},
			Options: make(map[string]interface{}),
			Raw: `url: http://{domain}/api
method: POST
headers:
  X-Foo-Bar: SomeValue
  X-Tmpl-Var: {header-value-test}
body:
  username: "{name}"
  password: "{password}"`,
		},
	}

	wantedRequest := model.Request{
		Url:    "http://prodfoobar.com/api",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo-Bar":  {"SomeValue"},
			"X-Tmpl-Var": {"Value from template var"},
		},
		Body: map[string]interface{}{
			"username": "Jane",
			"password": "Secret",
		},
		Options: make(map[string]interface{}),
	}

	profile := model.Profile{
		Name: "test",
		Variables: map[string]string{
			"domain":            "prodfoobar.com",
			"header-value-test": "Value from template var",
			"name":              "Jane",
			"password":          "Secret",
		},
	}

	request, err := BuildRequest(&requestMold, profile)
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest) {
		t.Errorf("got\n%q\nwanted\n%q\n", request, wantedRequest)
	}
}

func TestBuildStarlarkRequest(t *testing.T) {
	requestMold := model.RequestMold{
		Name: "Starlark request",
		Starlark: &model.StarlarkRequest{
			Script: `
"""
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
"""
url = "http://foobar.com"
headers = { "X-Foo": "Bar", "X-Foos": [ "Bar1", "Bar2" ] }
method = "POST"
body = {
    "id": 1,
    "amount": 1.2001,
    "name": "Jane"
}`,
		},
	}

	wantedRequest := model.Request{
		Url:    "http://foobar.com",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo":  {"Bar"},
			"X-Foos": {"Bar1", "Bar2"},
		},
		Body: map[string]interface{}{
			"id":     big.NewInt(1),
			"amount": 1.2001,
			"name":   "Jane",
		},
		Options: make(map[string]interface{}),
	}

	request, err := BuildRequest(&requestMold, model.Profile{})
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest, cmp.AllowUnexported(big.Int{})) {
		t.Errorf("got\n%q\nwanted\n%q\n", request, wantedRequest)
	}
}
func TestBuildStarlarkRequestWithTemplateVariables(t *testing.T) {
	requestMold := model.RequestMold{
		Name: "Starlark request",
		Starlark: &model.StarlarkRequest{
			Script: `
"""
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
"""
url = "http://{domain}"
headers = { "X-Foo": "Bar", "X-Foos": [ "Bar1", "Bar2" ], "X-Tmpl-Var": "{header-value-test}" }
method = "POST"
body = {
    "id": 1,
    "amount": 1.2001,
    "name": "Jane"
}`,
		},
	}

	wantedRequest := model.Request{
		Url:    "http://foobar.com",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo":      {"Bar"},
			"X-Foos":     {"Bar1", "Bar2"},
			"X-Tmpl-Var": {"Value from template var"},
		},
		Body: map[string]interface{}{
			"id":     big.NewInt(1),
			"amount": 1.2001,
			"name":   "Jane",
		},
		Options: make(map[string]interface{}),
	}

	profile := model.Profile{
		Name: "some-profile",
		Variables: map[string]string{
			"domain":            "foobar.com",
			"header-value-test": "Value from template var",
		},
	}
	request, err := BuildRequest(&requestMold, profile)
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest, cmp.AllowUnexported(big.Int{})) {
		t.Errorf("got\n%q\nwanted\n%q\n", request, wantedRequest)
	}
}
