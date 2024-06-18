package openapi

import (
	"fmt"
	"os"
	"startpoint/core/client/validator"

	"github.com/go-resty/resty/v2"
	"github.com/pb33f/libopenapi"
)

func ReadSpec(path string) {

	// load an OpenAPI 3 specification from bytes
	specBytes, err := loadSpec(path)
	if err != nil {
		panic(err)
	}

	// create a new document from specification bytes
	document, err := libopenapi.NewDocument(specBytes)

	// if anything went wrong, an error is thrown
	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	// TODO figure out v2/v3 OR ask it from user
	v3Model, errors := document.BuildV3Model()

	// if anything went wrong when building the v3 model, a slice of errors will be returned
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported",
			len(errors)))
	}

	// get a count of the number of paths and schemas.
	paths := v3Model.Model.Paths.PathItems.Len()
	schemas := v3Model.Model.Components.Schemas.Len()

	// print the number of paths and schemas in the document
	fmt.Printf("There are %d paths and %d schemas in the document\n", paths, schemas)

	for pathPairs := v3Model.Model.Paths.PathItems.First(); pathPairs != nil; pathPairs = pathPairs.Next() {
		pathItem := pathPairs.Value()
		if pathItem.Get != nil {
			lowOp := pathItem.Get.GoLow()
			fmt.Printf(">>> %s, %s\n", lowOp.Summary.Value, lowOp.OperationId.Value)
		}
	}
}

func loadSpec(path string) ([]byte, error) {

	if validator.IsValidUrl(path) {
		r := resty.New().R()
		resp, err := r.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	} else {
		file, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

}
