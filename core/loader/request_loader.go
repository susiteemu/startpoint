package loader

import (
	"goful/core/model"
	"io/fs"
	"os"
	"path/filepath"

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
		yamlRequest.Raw = string(file)
		if yamlRequest.Name != "" {
			request = &model.RequestMold{
				Yaml:        yamlRequest,
				ContentType: "yaml",
				Root:        root,
				Filename:    filename,
			}
		}

	case extension == ".star":
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			return nil, err
		}
		starlarkRequest := &model.StarlarkRequest{
			Script: string(file),
		}
		request = &model.RequestMold{
			Starlark:    starlarkRequest,
			ContentType: "star",
			Root:        root,
			Filename:    filename,
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
		var extension = filepath.Ext(filename)
		// TODO use ReadRequest function
		switch {
		case extension == ".yaml" || extension == ".yml":

			file, err := os.ReadFile(path)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to read %s", path)
				return nil
			}
			yamlRequest := &model.YamlRequest{}
			err = yaml.Unmarshal(file, yamlRequest)
			if err != nil {
				return nil
			}
			yamlRequest.Raw = string(file)
			if yamlRequest.Name != "" {
				request := model.RequestMold{
					Yaml:        yamlRequest,
					ContentType: "yaml",
					Root:        root,
					Filename:    filename,
				}
				requestSlice = append(requestSlice, &request)
			}

		case extension == ".star":

			file, err := os.ReadFile(path)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to read %s", path)
				return nil
			}
			starlarkRequest := &model.StarlarkRequest{
				Script: string(file),
			}
			request := model.RequestMold{
				Starlark:    starlarkRequest,
				ContentType: "star",
				Root:        root,
				Filename:    filename,
			}
			requestSlice = append(requestSlice, &request)

		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msgf("Error occurred while walking %s", root)
		return nil, err
	}

	return requestSlice, nil
}
