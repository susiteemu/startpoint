package builder

import (
	"github.com/google/go-cmp/cmp"
	"goful/core/model"
	"math/big"
	"testing"
)

func TestBuildRequestYaml(t *testing.T) {

	requestMold := model.RequestMold{
		Yaml: &model.YamlRequest{
			Name:   "yaml_request",
			Url:    "http://foobar.com",
			Method: "POST",
			Headers: model.Headers{
				"X-Foo-Bar": {"SomeValue"},
			},
			Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
		},
	}

	wantedRequest := model.Request{
		Url:    "http://foobar.com",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo-Bar": {"SomeValue"},
		},
		Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
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
		Yaml: &model.YamlRequest{
			Name:   "yaml_request",
			Url:    "http://{domain}/api",
			Method: "POST",
			Headers: model.Headers{
				"X-Foo-Bar":  {"SomeValue"},
				"X-Tmpl-Var": {"{header-value-test}"},
			},
			Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
		},
	}

	wantedRequest := model.Request{
		Url:    "http://prodfoobar.com/api",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo-Bar":  {"SomeValue"},
			"X-Tmpl-Var": {"Value from template var"},
		},
		Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
	}

	profile := model.Profile{
		Name: "test",
		Variables: map[string]string{
			"domain":            "prodfoobar.com",
			"header-value-test": "Value from template var",
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
		Starlark: &model.StarlarkRequest{
			Script: `
"""
meta:name: Starlark request
meta:prev_req: Some previous request
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
