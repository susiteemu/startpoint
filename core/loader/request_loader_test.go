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

	starlarkRequest := model.RequestMold{
		Starlark: &model.StarlarkRequest{
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
		},
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

	}

}

func TestReadRequestsWithInvalidRoot(t *testing.T) {

	_, err := ReadRequests("non_existent")

	if err == nil {
		t.Errorf("did expect error")
		return
	}

}
