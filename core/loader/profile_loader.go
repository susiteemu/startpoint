package loader

import (
	"goful/core/model"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func ReadProfiles(root string) ([]*model.Profile, error) {
	var profileSlice []*model.Profile
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() != root {
			return fs.SkipDir
		}

		filename := info.Name()

		if strings.HasPrefix(filename, ".env") {

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			envFile, err := godotenv.Parse(file)
			defer file.Close()
			if err != nil {
				return err
			}
			profileName := "default"
			splits := strings.Split(filename, ".")
			if len(splits) > 2 {
				profileName = splits[len(splits)-1]
			}

			profile := model.Profile{
				Name:      profileName,
				Variables: envFile,
			}

			profileSlice = append(profileSlice, &profile)

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return profileSlice, nil
}
