package loader

import (
	"startpoint/core/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadProfiles(t *testing.T) {

	profiles, err := ReadProfiles("testdata")

	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	assert.Equal(t, 2, len(profiles))

	wantedProfiles := []*model.Profile{
		{
			Name: "default",
			Variables: map[string]string{
				"domain": "foobar.com",
				"foo":    "bar",
			},
			Root:              "testdata",
			Filename:          ".env",
			HasPublicProfile:  true,
			HasPrivateProfile: false,
			Raw: `domain=foobar.com
foo=bar`,
		},
		{
			Name: "production",
			Variables: map[string]string{
				"domain": "foobarprod.com",
				"foo":    "bar2",
			},
			Root:              "testdata",
			Filename:          ".env.production",
			HasPublicProfile:  true,
			HasPrivateProfile: false,
			Raw: `domain=foobarprod.com
foo=bar2`,
		},
	}

	assert.Equal(t, profiles, wantedProfiles)

}

func TestGetProfileValues(t *testing.T) {
	profiles := []*model.Profile{
		{
			Name:     "default",
			Filename: ".env",
			Variables: map[string]string{
				"domain": "foobar.com",
				"foo":    "bar",
				"bar":    "foo",
			},
		},
		{
			Name:     "production",
			Filename: ".env.production",
			Variables: map[string]string{
				"domain": "foobarprod.com",
				"foo":    "bar2",
				"bar2":   "foo2",
			},
		},
		{
			Name:     "production.local",
			Filename: ".env.production.local",
			Variables: map[string]string{
				"secret":       "very secret",
				"another_var":  "foobar",
				"var-in-var":   "{another_var}",
				"var-in-var2":  "{another_var2}",
				"another_var2": "BAR_{another_var}",
				"var-in-var3":  "{another_var}_{bar}_{foo}",
			},
		},
	}

	profileValues := GetProfileValues(profiles[1], profiles, []string{})

	wantedProfileValues := map[string]string{
		"domain":       "foobarprod.com",
		"foo":          "bar2",
		"bar":          "foo",
		"bar2":         "foo2",
		"secret":       "very secret",
		"var-in-var":   "foobar",
		"another_var":  "foobar",
		"var-in-var2":  "BAR_foobar",
		"another_var2": "BAR_foobar",
		"var-in-var3":  "foobar_foo_bar2",
	}

	assert.Equal(t, len(wantedProfileValues), len(profileValues))

	for k, got := range profileValues {
		wanted := wantedProfileValues[k]
		assert.Equal(t, wanted, got)
	}

}
