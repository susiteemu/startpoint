package luang

import (
	"startpoint/core/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Generate tests
//

func TestRunLuaScript(t *testing.T) {

	tests := []struct {
		name             string
		mold             model.RequestMold
		previousResponse model.Response
		expected         map[string]interface{}
	}{
		{
			name: "Run script with basic request mold",
			mold: model.RequestMold{
				Scriptable: &model.ScriptableRequest{
					Script: `
print(prevResponse.headers["X-Custom-Header"][1])
return {
	url = "http://foo.bar",
	method = "POST",
	headers = {
		["X-Custom-Header"] = {"FooBar","Barz"}
	},
	body = {
		id=1,
		name="Jane"
	}
}`,
				},
			},
			previousResponse: model.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Headers: map[string]model.HeaderValues{
					"X-Custom-Header": {"SomeValue"},
				},
			},
			expected: map[string]interface{}{
				"url":    "http://foo.bar",
				"method": "POST",
				"headers": map[string][]string{
					"X-Custom-Header": {"FooBar", "Barz"},
				},
				"body": map[string]interface{}{
					"id":   1,
					"name": "Jane",
				},
				"auth":    map[string]interface{}{},
				"options": map[string]interface{}{},
				"output":  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunLuaScript(tt.mold, &tt.previousResponse)
			assert.Nil(t, err, "did not expect error to happen")
			assert.NotEqual(t, tt.expected, result, "results should match")
		})
	}
}
