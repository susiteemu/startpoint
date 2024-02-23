package loader

import (
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

	for i := 0; i < len(profiles); i++ {
		p := profiles[i]
		w := wantedProfiles[i]
		if !cmp.Equal(p, w) {
			t.Errorf("structs are not equal!\ngot\n%v\nwanted\n%v", p, w)
		}
	}

}
