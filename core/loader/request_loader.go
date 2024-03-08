package loader

import (
	"goful/core/model"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadRequests(root string) ([]model.RequestMold, error) {
	var requestSlice []model.RequestMold
	maxDepth := 0
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.Count(path, string(os.PathSeparator)) > maxDepth {
			return fs.SkipDir
		}

		filename := info.Name()

		var extension = filepath.Ext(filename)

		switch {
		case extension == ".yaml" || extension == ".yml":

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			yamlRequest := &model.YamlRequest{}
			err = yaml.Unmarshal(file, yamlRequest)
			if err != nil {
				return err
			}
			if yamlRequest.Name != "" {
				request := model.RequestMold{
					Yaml:        yamlRequest,
					Raw:         string(file),
					ContentType: "yaml",
					Filename:    filename,
				}
				requestSlice = append(requestSlice, request)
			}

		case extension == ".star":

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			starlarkRequest := &model.StarlarkRequest{
				Script: string(file),
			}
			request := model.RequestMold{
				Starlark:    starlarkRequest,
				Raw:         string(file),
				ContentType: "star",
				Filename:    filename,
			}
			requestSlice = append(requestSlice, request)

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return requestSlice, nil
}
