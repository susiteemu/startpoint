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
