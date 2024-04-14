package loader

import (
	"fmt"
	"goful/core/model"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadProfiles(t *testing.T) {

	profiles, err := ReadProfiles("testdata")

	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if len(profiles) != 2 {
		t.Errorf("got %d, wanted %d", len(profiles), 2)
		return
	}

	wantedProfiles := []model.Profile{
		{
			Name: "default",
			Variables: map[string]string{
				"domain": "foobar.com",
				"foo":    "bar",
			},
		},
		{
			Name: "production",
			Variables: map[string]string{
				"domain": "foobarprod.com",
				"foo":    "bar2",
			},
		},
	}

	for _, w := range wantedProfiles {
		found := false
		for _, p := range profiles {
			if cmp.Equal(*p, w) {
				found = true
				break
			}
		}
		if !found {
			t.Error(fmt.Sprintf("wanted %v but not found in %v", w, profiles))
		}
	}

}

func TestGetProfileValues(t *testing.T) {
	profiles := []*model.Profile{
		{
			Name: "default",
			Variables: map[string]string{
				"domain": "foobar.com",
				"foo":    "bar",
				"bar":    "foo",
			},
		},
		{
			Name: "production",
			Variables: map[string]string{
				"domain": "foobarprod.com",
				"foo":    "bar2",
				"bar2":   "foo2",
			},
		},
	}

	profileValues := GetProfileValues(profiles[1], profiles)

	wantedProfileValues := map[string]string{
		"domain": "foobarprod.com",
		"foo":    "bar2",
		"bar":    "foo",
		"bar2":   "foo2",
	}

	if len(profileValues) != len(wantedProfileValues) {
		t.Errorf("lengths do not match: got %d, wanted %d", len(profileValues), len(wantedProfileValues))
		return
	}

	for k, got := range profileValues {
		wanted := wantedProfileValues[k]
		if got != wanted {
			t.Errorf("got %s, wanted %s", got, wanted)
		}
	}

}
