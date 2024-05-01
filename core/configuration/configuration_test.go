package configuration

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func TestBuildRequestYaml(t *testing.T) {

	src := make(map[string]interface{})
	dest := make(map[string]interface{})

	src["foo"] = "bar"
	src["foos"] = []string{"one", "two"}
	src["nested"] = map[string]interface{}{
		"deeper": []map[string]string{
			{"key1a": "value1a", "key1b": "value1b"},
			{"key2a": "value2a", "key2b": "value2b"},
		},
	}

	Flatten("", src, dest)

	key := "foo"
	wanted := "bar"
	if dest[key] != wanted {
		t.Errorf("with key %s got %q wanted %q\n", key, dest[key], wanted)
	}

	key = "foos"
	got, ok := dest[key].([]string)
	if !ok {
		t.Errorf("with key %s got %v which is not a slice\n", key, got)
	} else {
		wanted := []string{"one", "two"}
		if !slices.Equal(got, wanted) {
			t.Errorf("with key %s got %q wanted %q\n", key, got, wanted)
		}
	}

	key = "nested.deeper"
	gotNested, ok := dest[key].([]map[string]string)
	if !ok {
		t.Errorf("with key %s got %v which is not a []map[string]string\n", key, got)
	} else {
		fmt.Printf("Got %v\n", gotNested)
		wanted := []map[string]string{
			{"key1a": "value1a", "key1b": "value1b"},
			{"key2a": "value2a", "key2b": "value2b"},
		}
		if len(gotNested) != len(wanted) {
			t.Errorf("with key %s results have different lengths: %d vs %d", key, len(gotNested), len(wanted))
		} else {
			for i := 0; i < len(gotNested); i++ {
				g := gotNested[i]
				w := wanted[i]
				if !reflect.DeepEqual(g, w) {
					t.Errorf("with key %s got %q wanted %q\n", key, got, wanted)
				}
			}
		}
	}
}
