package loader

import (
	"goful/core/model"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadRequests(t *testing.T) {
	requests, err := ReadRequests("testdata")

	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if len(requests) != 2 {
		t.Errorf("got %d, wanted %d", len(requests), 2)
		return
	}

	var wantedRequests []model.RequestMold

	script := `"""
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
`

	starlarkRequest := model.RequestMold{
		Starlark: &model.StarlarkRequest{
			Script: script,
		},
		ContentType: "star",
		Filename:    "starlark_request.star",
	}

	wantedRequests = append(wantedRequests, starlarkRequest)

	yamlRequest := model.RequestMold{
		Yaml: &model.YamlRequest{
			Name:    "yaml_request",
			PrevReq: "",
			Url:     "foobar.com",
			Method:  "POST",
			Headers: map[string]model.HeaderValues{
				"X-Foo-Bar": {"SomeValue"},
			},
			Body: "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}\n",
			Raw: `name: yaml_request
prev_req:
url: foobar.com
method: POST
headers:
  X-Foo-Bar: SomeValue
body: >
  {
    "id": 1,
    "name": "Jane"
  }
`,
		},
		ContentType: "yaml",
		Filename:    "yaml_request.yaml",
	}

	wantedRequests = append(wantedRequests, yamlRequest)

	for i := 0; i < len(requests); i++ {

		request := requests[i]
		wantedRequest := wantedRequests[i]
		if request.Yaml != nil {
			r := request.Yaml
			w := wantedRequest.Yaml
			if !cmp.Equal(r, w) {
				t.Errorf("structs are not equal!\ngot\n%v\nwanted\n%v", r, w)
			}
		}
		if request.Starlark != nil {
			r := request.Starlark
			w := wantedRequest.Starlark
			if !cmp.Equal(r, w) {
				t.Errorf("structs are not equal!\ngot\n%v\nwanted\n%v", r, w)
			}
		}
		if !cmp.Equal(request.Raw, wantedRequest.Raw) {
			t.Errorf("structs are not equal!\ngot\n%v\nwanted\n%v", request, wantedRequest)
		}
		if !cmp.Equal(request.ContentType, wantedRequest.ContentType) {
			t.Errorf("structs are not equal!\ngot\n%v\nwanted\n%v", request, wantedRequest)
		}

		if !cmp.Equal(request.Filename, wantedRequest.Filename) {
			t.Errorf("structs are not equal!\ngot\n%v\nwanted\n%v", request, wantedRequest)
		}
	}

}

func TestReadRequestsWithInvalidRoot(t *testing.T) {

	_, err := ReadRequests("non_existent")

	if err == nil {
		t.Errorf("did expect error")
		return
	}

}
