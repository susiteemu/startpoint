package loader

import (
	"io/fs"
	"os"
	"path/filepath"
	"startpoint/core/model"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func ReadRequest(root, filename string) (*model.RequestMold, error) {
	path := filepath.Join(root, filename)
	extension := filepath.Ext(filename)
	var request *model.RequestMold

	switch {
	case extension == ".yaml" || extension == ".yml":
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			return nil, err
		}
		yamlRequest := &model.YamlRequest{}
		err = yaml.Unmarshal(file, yamlRequest)
		if err != nil {
			return nil, err
		}
		yamlRequest.Raw = strings.TrimSuffix(string(file), "\n")
		// TODO: how to filter out yaml files that are not requests?
		if yamlRequest.Url != "" || yamlRequest.Method != "" {
			request = &model.RequestMold{
				Yaml:        yamlRequest,
				ContentType: "yaml",
				Root:        root,
				Filename:    filename,
				Name:        strings.TrimSuffix(filename, extension),
			}
		}

	case extension == ".star":
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			return nil, err
		}
		starlarkRequest := &model.StarlarkRequest{
			Script: strings.TrimSuffix(string(file), "\n"),
		}
		request = &model.RequestMold{
			Starlark:    starlarkRequest,
			ContentType: "star",
			Root:        root,
			Filename:    filename,
			Name:        strings.TrimSuffix(filename, extension),
		}

	}

	return request, nil
}

func ReadRequests(root string) ([]*model.RequestMold, error) {
	var requestSlice []*model.RequestMold
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// this prevents walking into subdirectories
		if info.IsDir() && path != root {
			return fs.SkipDir
		}

		filename := info.Name()
		log.Debug().Msgf("Walk crossed a file %s", filename)

		requestMold, err := ReadRequest(root, filename)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read file %s", filename)
		}
		if requestMold != nil {
			requestSlice = append(requestSlice, requestMold)
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msgf("Error occurred while walking %s", root)
		return nil, err
	}

	return requestSlice, nil
}
