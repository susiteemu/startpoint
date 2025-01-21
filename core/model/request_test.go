package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractValuesFromScriptableRequest(t *testing.T) {

	tests := []struct {
		name            string
		mold            RequestMold
		expectedPrevReq string
		expectedUrl     string
		expectedOutput  string
		expectedMethod  string
	}{
		{
			name: "Starlark request with all attributes inside doc string",
			mold: RequestMold{
				Type: "star",
				Scriptable: &ScriptableRequest{
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
			},
			expectedPrevReq: "Some previous request",
			expectedUrl:     "http://foobar.com",
			expectedOutput:  "./output.txt",
			expectedMethod:  "POST",
		},
		{
			name: "Starlark request with missing parts",
			mold: RequestMold{
				Type: "star",
				Scriptable: &ScriptableRequest{
					Script: `"""
doc:url: http://foobar.com
doc:method: POST
"""
#  prev_req: Some previous request
url = "http://foobar.com"
method = "POST"
headers = { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
body = { "id": 1474, "prev": prev, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
`},
			},
			expectedPrevReq: "Some previous request",
			expectedUrl:     "http://foobar.com",
			expectedOutput:  "",
			expectedMethod:  "POST",
		},
		{
			name: "Starlark request with values inside actual code",
			mold: RequestMold{
				Type: "star",
				Scriptable: &ScriptableRequest{
					Script: `prev_req = "Some previous request"
url = "http://foobar.com"
method = "POST"
headers = { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
body = { "id": 1474, "prev": prev, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
output = "./output.txt"
`},
			},
			expectedPrevReq: "Some previous request",
			expectedUrl:     "http://foobar.com",
			expectedOutput:  "./output.txt",
			expectedMethod:  "POST",
		},
		{
			name: "Lua request with all attributes inside comment block",
			mold: RequestMold{
				Type: "lua",
				Scriptable: &ScriptableRequest{
					Script: `--[[
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
meta:output: ./output.txt
]]--
return {
	url = "http://foobar.com",
	headers = { ["X-Foo"] = "Bar", ["X-Foos"] = { "Bar1", "Bar2" } },
	method = "POST",
	body = {
		id=1,
		amount=1.2001,
		name="Jane"
	}
}`},
			},
			expectedPrevReq: "Some previous request",
			expectedUrl:     "http://foobar.com",
			expectedOutput:  "./output.txt",
			expectedMethod:  "POST",
		},
		{
			name: "Lua request with missing parts",
			mold: RequestMold{
				Type: "lua",
				Scriptable: &ScriptableRequest{
					Script: `--[[
doc:url: http://foobar.com
doc:method: POST
]]--
return {
	url = "http://foobar.com",
	headers = { ["X-Foo"] = "Bar", ["X-Foos"] = { "Bar1", "Bar2" } },
	method = "POST",
	body = {
		id=1,
		amount=1.2001,
		name="Jane"
	}
}`},
			},
			expectedPrevReq: "",
			expectedUrl:     "http://foobar.com",
			expectedOutput:  "",
			expectedMethod:  "POST",
		},
		{
			name: "Lua request with attributes inside actual code",
			mold: RequestMold{
				Type: "lua",
				Scriptable: &ScriptableRequest{
					Script: `
return {
	url = "http://foobar.com",
	headers = { ["X-Foo"] = "Bar", ["X-Foos"] = { "Bar1", "Bar2" } },
	method = "POST",
	body = {
		id=1,
		amount=1.2001,
		name="Jane"
	},
	output = "./output.txt",
	prev_req = "Some previous request"
}`},
			},
			expectedPrevReq: "Some previous request",
			expectedUrl:     "http://foobar.com",
			expectedOutput:  "./output.txt",
			expectedMethod:  "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedPrevReq, tt.mold.PreviousReq())
			assert.Equal(t, tt.expectedUrl, tt.mold.Url())
			assert.Equal(t, tt.expectedMethod, tt.mold.Method())
			assert.Equal(t, tt.expectedOutput, tt.mold.Output())
		})
	}
}

func TestChangePreviousRequest(t *testing.T) {

	starlarkRequest := RequestMold{
		Type: "star",
		Scriptable: &ScriptableRequest{
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

	wantedPrevReq := "Some previous request"
	assert.Equal(t, wantedPrevReq, starlarkRequest.PreviousReq())

	starlarkRequest.ChangePreviousReq("Some other previous request")

	wantedPrevReq = "Some other previous request"
	assert.Equal(t, wantedPrevReq, starlarkRequest.PreviousReq())

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

	wantedPrevReq = "Some previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())

	yamlRequest.ChangePreviousReq("Some other previous request")

	wantedPrevReq = "Some other previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())

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
	wantedPrevReq = "Some previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())

	yamlRequest.ChangePreviousReq("Some other previous request")

	wantedPrevReq = "Some other previous request"
	assert.Equal(t, wantedPrevReq, yamlRequest.PreviousReq())
}
