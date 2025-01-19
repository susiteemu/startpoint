package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"startpoint/core/model"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	YAML_EXT = ".yaml"
	YML_EXT  = ".yml"
	STAR_EXT = ".star"
	LUA_EXT  = ".lua"
)

func ReadRequest(root, filename string) (*model.RequestMold, error) {
	path := filepath.Join(root, filename)
	extension := filepath.Ext(filename)
	var request *model.RequestMold

	switch extension {
	case YAML_EXT, YML_EXT:
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			return nil, err
		}
		yamlRequest := &model.YamlRequest{}
		err = yaml.Unmarshal(file, yamlRequest)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to unmarshal file %s", path)
			return nil, err
		}
		yamlRequest.Raw = strings.TrimSuffix(string(file), "\n")
		// TODO: how to filter out yaml files that are not requests?
		if yamlRequest.Url != "" || yamlRequest.Method != "" {
			request = &model.RequestMold{
				Yaml:     yamlRequest,
				Type:     model.CONTENT_TYPE_YAML,
				Root:     root,
				Filename: filename,
				Name:     strings.TrimSuffix(filename, extension),
			}
		}

	case STAR_EXT:
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			return nil, err
		}
		starlarkRequest := &model.ScriptableRequest{
			Script: strings.TrimSuffix(string(file), "\n"),
		}
		request = &model.RequestMold{
			Scriptable: starlarkRequest,
			Type:       model.CONTENT_TYPE_STARLARK,
			Root:       root,
			Filename:   filename,
			Name:       strings.TrimSuffix(filename, extension),
		}

	case LUA_EXT:
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			return nil, err
		}
		luaRequest := &model.ScriptableRequest{
			Script: strings.TrimSuffix(string(file), "\n"),
		}
		request = &model.RequestMold{
			Scriptable: luaRequest,
			Type:       model.CONTENT_TYPE_LUA,
			Root:       root,
			Filename:   filename,
			Name:       strings.TrimSuffix(filename, extension),
		}
	}

	if request == nil {
		log.Error().Msgf("Request could not be read %s", path)
		return nil, fmt.Errorf("failed to read request %s", path)
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

		extension := filepath.Ext(filename)
		if extension == YAML_EXT || extension == YML_EXT || extension == STAR_EXT || extension == LUA_EXT {
			log.Debug().Msgf("Walk crossed a file %s", filename)

			requestMold, err := ReadRequest(root, filename)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to read file %s", filename)
			}
			if requestMold != nil {
				requestSlice = append(requestSlice, requestMold)
			}
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msgf("Error occurred while walking %s", root)
		return nil, err
	}

	return requestSlice, nil
}
