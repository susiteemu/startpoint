package builder

import (
	"github.com/google/go-cmp/cmp"
	"goful/core/model"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildRequestYaml(t *testing.T) {

	requestMetadata := model.RequestMetadata{
		Name:       "yaml_request",
		PrevReq:    "",
		Request:    "yaml",
		WorkingDir: filepath.Join(currentDir(), "testdata"),
	}

	wantedRequest := model.Request{
		Url:    "http://foobar.com",
		Method: "POST",
		Headers: model.Headers{
			"X-Foo-Bar": {"SomeValue"},
		},
		Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
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
		WorkingDir: filepath.Join(currentDir(), "testdata"),
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

func TestBuildStarlarkRequest(t *testing.T) {
	requestMetadata := model.RequestMetadata{
		Name:       "starlark_request",
		PrevReq:    "",
		Request:    "star",
		WorkingDir: filepath.Join(currentDir(), "testdata"),
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

	request, err := BuildRequest(requestMetadata, model.Profile{})
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if !cmp.Equal(request, wantedRequest, cmp.AllowUnexported(big.Int{})) {
		t.Errorf("got %q, wanted %q", request, wantedRequest)
	}
}

func currentDir() string {
	wd, _ := os.Getwd()
	return wd
}
