package loader

import (
	"goful/core/model"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadRequestMetadata(root string) ([]model.RequestMetadata, error) {
	var metadataSlice []model.RequestMetadata
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
		if extension != ".yaml" && extension != ".yml" {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		metadata := &model.RequestMetadata{}
		err = yaml.Unmarshal(file, metadata)
		if err != nil {
			return err
		}
		if metadata.Request != "" {
			var nameWithoutExt = filename[0 : len(filename)-len(extension)]
			metadata.Name = nameWithoutExt
			metadata.WorkingDir = root
			metadataSlice = append(metadataSlice, *metadata)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return metadataSlice, nil
}
