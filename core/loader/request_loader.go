package loader

import (
	"goful/core/model"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYamlRequest(metadata model.RequestMetadata) (model.Request, error) {
	path := metadata.ToRequestPath()
	file, err := os.ReadFile(path)
	if err != nil {
		return model.Request{}, err
	}
	request := &model.Request{}
	err = yaml.Unmarshal(file, request)
	if err != nil {
		return model.Request{}, err
	}
	return *request, nil
}
