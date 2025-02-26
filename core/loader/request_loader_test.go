package loader

import (
	"github.com/susiteemu/startpoint/core/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequests(t *testing.T) {
	requests, err := ReadRequests("testdata")

	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	assert.Equal(t, 3, len(requests))

	var wantedRequests []model.RequestMold

	script := `"""
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
] }`

	starlarkRequest := model.RequestMold{
		Scriptable: &model.ScriptableRequest{
			Script: script,
		},
		Type:     "star",
		Root:     "testdata",
		Filename: "starlark_request.star",
		Name:     "starlark_request",
	}

	wantedRequests = append(wantedRequests, starlarkRequest)

	yamlRequest := model.RequestMold{
		Yaml: &model.YamlRequest{
			PrevReq: "",
			Url:     "foobar.com",
			Method:  "POST",
			Headers: map[string]model.HeaderValues{
				"X-Foo-Bar": {"SomeValue"},
			},
			Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}\n",
			Raw: `prev_req:
url: foobar.com
method: POST
headers:
  X-Foo-Bar: SomeValue
body: >
  {
    "id": 1,
    "name": "Jane"
  }`,
		},
		Root:     "testdata",
		Type:     "yaml",
		Filename: "yaml_request.yaml",
		Name:     "yaml_request",
	}

	yamlRequestWithBasicAuth := model.RequestMold{
		Yaml: &model.YamlRequest{
			PrevReq: "",
			Url:     "foobar.com",
			Method:  "POST",
			Headers: map[string]model.HeaderValues{
				"X-Foo-Bar": {"SomeValue"},
			},
			Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}\n",
			Raw: `prev_req:
url: foobar.com
method: POST
headers:
  X-Foo-Bar: SomeValue
auth:
  basic:
    username: user
    password: pw
body: >
  {
    "id": 1,
    "name": "Jane"
  }`,
		},
		Root:     "testdata",
		Type:     "yaml",
		Filename: "yaml_request_with_basic_auth.yaml",
		Name:     "yaml_request_with_basic_auth",
	}

	wantedRequests = append(wantedRequests, yamlRequest)
	wantedRequests = append(wantedRequests, yamlRequestWithBasicAuth)

	for i := 0; i < len(requests); i++ {

		request := requests[i]
		wantedRequest := wantedRequests[i]

		if request.Yaml != nil {
			r := request.Yaml
			w := wantedRequest.Yaml
			assert.Equal(t, w, r)
		}
		if request.Scriptable != nil {
			r := request.Scriptable
			w := wantedRequest.Scriptable
			assert.Equal(t, w, r)
		}
		assert.Equal(t, wantedRequest.Type, request.Type)
		assert.Equal(t, wantedRequest.Root, request.Root)
		assert.Equal(t, wantedRequest.Filename, request.Filename)
		assert.Equal(t, wantedRequest.Name, request.Name)

	}

}

func TestReadRequestsWithInvalidRoot(t *testing.T) {

	_, err := ReadRequests("non_existent")

	if err == nil {
		t.Errorf("did expect error")
		return
	}

}
