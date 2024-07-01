package templateng

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverTemplateVariables(t *testing.T) {
	s := "{easy}"
	vars := DiscoverTemplateVariables(s)
	expected := []string{"easy"}
	assert.Equal(t, expected, vars)

	s = "{first}{second}"
	vars = DiscoverTemplateVariables(s)
	expected = []string{"first", "second"}
	assert.Equal(t, expected, vars)

	s = "foobar{first}sometext {second}barfoo"
	vars = DiscoverTemplateVariables(s)
	expected = []string{"first", "second"}
	assert.Equal(t, expected, vars)
}
