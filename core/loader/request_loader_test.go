package loader

import (
	"github.com/google/go-cmp/cmp"
	"goful/core/model"
	"testing"
)

func TestReadYamlRequest(t *testing.T) {
	metadata := model.RequestMetadata{
		Name:       "yaml_request",
		PrevReq:    "",
		Request:    "yaml",
		WorkingDir: "testdata",
	}

	request, err := ReadYamlRequest(metadata)
	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	wantedRequest := model.Request{
		Url:    "foobar.com",
		Method: "POST",
		Headers: map[string]model.HeaderValues{
			"X-Foo-Bar": {"SomeValue"},
		},
		Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
	}

	if !cmp.Equal(request, wantedRequest) {
		t.Errorf("got %q, wanted %q", request, wantedRequest)
	}
}

func TestReadYamlRequest_FileDoesNotExist(t *testing.T) {
	metadata := model.RequestMetadata{
		Name:       "yaml_request_with_no_request_file",
		PrevReq:    "",
		Request:    "yaml",
		WorkingDir: "testdata",
	}

	_, err := ReadYamlRequest(metadata)
	if err == nil {
		t.Errorf("did expect error")
		return
	}
}
