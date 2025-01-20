package builder

import (
	"github.com/susiteemu/startpoint/core/model"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildYamlRequests(t *testing.T) {

	tests := []struct {
		name     string
		mold     model.RequestMold
		profile  model.Profile
		expected model.Request
	}{

		{
			name: "Test with basic request",
			mold: model.RequestMold{
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
					Auth: model.Auth{
						Basic: model.BasicAuth{
							User:     "jane",
							Password: "doe",
						},
					},
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo-Bar":     {"SomeValue"},
					"Authorization": model.HeaderValues{"Basic amFuZTpkb2U="},
				},
				Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
				Output:  "",
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with basic authentication",
			mold: model.RequestMold{
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
					Auth: model.Auth{
						Basic: model.BasicAuth{
							User:     "jane",
							Password: "doe",
						},
					},
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo-Bar":     {"SomeValue"},
					"Authorization": model.HeaderValues{"Basic amFuZTpkb2U="},
				},
				Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
				Output:  "",
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with bearer token authentication",
			mold: model.RequestMold{
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
					Auth: model.Auth{
						Bearer: "some-token",
					},
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo-Bar":     {"SomeValue"},
					"Authorization": model.HeaderValues{"Bearer some-token"},
				},
				Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
				Output:  "",
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with template variables",
			mold: model.RequestMold{
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
			},
			profile: model.Profile{
				Name: "test",
				Variables: map[string]string{
					"domain":            "prodfoobar.com",
					"header-value-test": "Value from template var",
					"name":              "Jane",
				},
			},
			expected: model.Request{
				Url:    "http://prodfoobar.com/api",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo-Bar":  {"SomeValue"},
					"X-Tmpl-Var": {"Value from template var"},
				},
				Body:    "{\n  \"id\": 1,\n  \"name\": \"Jane\"\n}",
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with form request and template variables",
			mold: model.RequestMold{
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
			},
			profile: model.Profile{
				Name: "test",
				Variables: map[string]string{
					"domain":            "prodfoobar.com",
					"header-value-test": "Value from template var",
					"name":              "Jane",
					"password":          "Secret",
				},
			},
			expected: model.Request{
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := BuildRequest(&tt.mold, tt.profile)
			assert.Nil(t, err, "did not expect error to happen")
			assert.Equal(t, tt.expected, request, "requests should match")
		})
	}

}

func TestBuildStarlarkRequests(t *testing.T) {

	tests := []struct {
		name     string
		mold     model.RequestMold
		profile  model.Profile
		expected model.Request
	}{

		{
			name: "Test with basic request",
			mold: model.RequestMold{
				Name: "Starlark request",
				Type: "star",
				Scriptable: &model.ScriptableRequest{
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
			},
			profile: model.Profile{},
			expected: model.Request{
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
			},
		},

		{
			name: "Test with basic authentication",
			mold: model.RequestMold{
				Name: "Starlark request",
				Type: "star",
				Scriptable: &model.ScriptableRequest{
					Script: `
"""
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
"""
url = "http://foobar.com"
headers = { "X-Foo": "Bar", "X-Foos": [ "Bar1", "Bar2" ] }
method = "POST"
auth = {
    "basic_auth": {
        "username": "jane",
		"password": "doe"
	}
}
body = {
    "id": 1,
    "amount": 1.2001,
    "name": "Jane"
}`,
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo":         {"Bar"},
					"X-Foos":        {"Bar1", "Bar2"},
					"Authorization": model.HeaderValues{"Basic amFuZTpkb2U="},
				},
				Body: map[string]interface{}{
					"id":     big.NewInt(1),
					"amount": 1.2001,
					"name":   "Jane",
				},
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with bearer token authentication",
			mold: model.RequestMold{
				Name: "Starlark request",
				Type: "star",
				Scriptable: &model.ScriptableRequest{
					Script: `
"""
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
"""
url = "http://foobar.com"
headers = { "X-Foo": "Bar", "X-Foos": [ "Bar1", "Bar2" ] }
method = "POST"
auth = {
    "bearer_token": "some-token"
}
body = {
    "id": 1,
    "amount": 1.2001,
    "name": "Jane"
}`,
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo":         {"Bar"},
					"X-Foos":        {"Bar1", "Bar2"},
					"Authorization": model.HeaderValues{"Bearer some-token"},
				},
				Body: map[string]interface{}{
					"id":     big.NewInt(1),
					"amount": 1.2001,
					"name":   "Jane",
				},
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with template variables",
			mold: model.RequestMold{
				Name: "Starlark request",
				Type: "star",
				Scriptable: &model.ScriptableRequest{
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
			},
			profile: model.Profile{
				Name: "some-profile",
				Variables: map[string]string{
					"domain":            "foobar.com",
					"header-value-test": "Value from template var",
				},
			},
			expected: model.Request{
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := BuildRequest(&tt.mold, tt.profile)
			assert.Nil(t, err, "did not expect error to happen")
			assert.Equal(t, tt.expected, request, "requests should match")
		})
	}

}

func TestBuildLuaRequests(t *testing.T) {

	tests := []struct {
		name     string
		mold     model.RequestMold
		profile  model.Profile
		expected model.Request
	}{

		{
			name: "Test with basic request",
			mold: model.RequestMold{
				Name: "Lua request",
				Type: "lua",
				Scriptable: &model.ScriptableRequest{
					Script: `
--[[
prev_req: Some previous request
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
}`,
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo":  {"Bar"},
					"X-Foos": {"Bar1", "Bar2"},
				},
				Body: map[interface{}]interface{}{
					"id":     float64(1),
					"amount": 1.2001,
					"name":   "Jane",
				},
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with basic authentication",
			mold: model.RequestMold{
				Name: "Lua request",
				Type: "lua",
				Scriptable: &model.ScriptableRequest{
					Script: `
--[[
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
]]--
return {
	url = "http://foobar.com",
	headers = { ["X-Foo"]="Bar", ["X-Foos"]={ "Bar1", "Bar2" } },
	method = "POST",
	auth = {
		basic_auth={
			username="jane",
			password="doe"
		}
	},
	body = {
		id=1,
		amount=1.2001,
		name="Jane"
	}
}`,
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo":         {"Bar"},
					"X-Foos":        {"Bar1", "Bar2"},
					"Authorization": model.HeaderValues{"Basic amFuZTpkb2U="},
				},
				Body: map[interface{}]interface{}{
					"id":     float64(1),
					"amount": 1.2001,
					"name":   "Jane",
				},
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with bearer token authentication",
			mold: model.RequestMold{
				Name: "Lua request",
				Type: "lua",
				Scriptable: &model.ScriptableRequest{
					Script: `
--[[
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
]]--
return {
	url = "http://foobar.com",
	headers = { ["X-Foo"] = "Bar", ["X-Foos"] = { "Bar1", "Bar2" } },
	method = "POST",
	auth = {
		bearer_token="some-token"
	},
	body = {
		id=1,
		amount=1.2001,
		name="Jane"
	}
}`,
				},
			},
			profile: model.Profile{},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo":         {"Bar"},
					"X-Foos":        {"Bar1", "Bar2"},
					"Authorization": model.HeaderValues{"Bearer some-token"},
				},
				Body: map[interface{}]interface{}{
					"id":     float64(1),
					"amount": 1.2001,
					"name":   "Jane",
				},
				Options: make(map[string]interface{}),
			},
		},

		{
			name: "Test with template variables",
			mold: model.RequestMold{
				Name: "Lua request",
				Type: "lua",
				Scriptable: &model.ScriptableRequest{
					Script: `
--[[
prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
]]--
return {
	url = "http://{domain}",
	headers = { ["X-Foo"] = "Bar", ["X-Foos"] = { "Bar1", "Bar2" }, ["X-Tmpl-Var"] = "{header-value-test}" },
	method = "POST",
	body = {
		id = 1,
		amount = 1.2001,
		name = "Jane"
	}
}`,
				},
			},
			profile: model.Profile{
				Name: "some-profile",
				Variables: map[string]string{
					"domain":            "foobar.com",
					"header-value-test": "Value from template var",
				},
			},
			expected: model.Request{
				Url:    "http://foobar.com",
				Method: "POST",
				Headers: model.Headers{
					"X-Foo":      {"Bar"},
					"X-Foos":     {"Bar1", "Bar2"},
					"X-Tmpl-Var": {"Value from template var"},
				},
				Body: map[interface{}]interface{}{
					"id":     float64(1),
					"amount": 1.2001,
					"name":   "Jane",
				},
				Options: make(map[string]interface{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := BuildRequest(&tt.mold, tt.profile)
			assert.Nil(t, err, "did not expect error to happen")
			assert.EqualValues(t, tt.expected, request, "requests should match")
		})
	}

}
