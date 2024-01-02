package loader

import (
	"goful/core/model"
	"slices"
	"testing"
)

func TestReadRequestMetadata(t *testing.T) {

	metadata, err := ReadRequestMetadata("testdata")

	if err != nil {
		t.Errorf("did not expect error %v", err)
		return
	}

	if len(metadata) != 2 {
		t.Errorf("got %d, wanted %d", len(metadata), 2)
		return
	}

	wantedMetadata := []model.RequestMetadata{
		{
			Name:       "starlark_request",
			PrevReq:    "yaml_request",
			Request:    "star",
			WorkingDir: "testdata",
		},
		{
			Name:       "yaml_request",
			PrevReq:    "",
			Request:    "yaml",
			WorkingDir: "testdata",
		},
	}

	if !slices.Equal(metadata, wantedMetadata) {
		t.Errorf("slices got %q, wanted %q", metadata, wantedMetadata)
		return
	}
}

func TestReadRequestMetadataWithInvalidRoot(t *testing.T) {

	_, err := ReadRequestMetadata("non_existent")

	if err == nil {
		t.Errorf("did expect error")
		return
	}

}
