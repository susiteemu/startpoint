package builder

import (
	"github.com/google/go-cmp/cmp"
	"goful/core/model"
	"testing"
)

func TestBuildRequestYaml(t *testing.T) {

	requestMetadata := model.RequestMetadata{
		Name:       "yaml_request",
		PrevReq:    "",
		Request:    "yaml",
		WorkingDir: "testdata",
	}

	wantedRequest := model.Request{
		Url:    "foobar.com",
		Method: "POST",
		Headers: map[string]model.HeaderValues{
			"X-Foo-Bar": {"SomeValue"},
		},
		Body: []byte("{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}"),
	}

	request, err := BuildRequest(requestMetadata, model.Profile{})
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest) {
		t.Errorf("got %q, wanted %q", request, wantedRequest)
	}
}

func TestBuildRequestYamlWithTemplateVariables(t *testing.T) {

	requestMetadata := model.RequestMetadata{
		Name:       "yaml_request_with_tmpl_vars",
		PrevReq:    "",
		Request:    "yaml",
		WorkingDir: "testdata",
	}

	wantedRequest := model.Request{
		Url:    "prodfoobar.com/api",
		Method: "POST",
		Headers: map[string]model.HeaderValues{
			"X-Foo-Bar":  {"SomeValue"},
			"X-Tmpl-Var": {"Value from template var"},
		},
		Body: []byte("{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}"),
	}

	profile := model.Profile{
		Name: "test",
		Variables: map[string]interface{}{
			"domain":            "prodfoobar.com",
			"header-value-test": "Value from template var",
		},
	}

	request, err := BuildRequest(requestMetadata, profile)
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest) {
		t.Errorf("got %q, wanted %q", request, wantedRequest)
	}
}
