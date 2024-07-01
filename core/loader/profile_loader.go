package loader

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"startpoint/core/model"
	"startpoint/core/templating/templateng"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func ReadProfile(root, filename string) (*model.Profile, error) {
	path := filepath.Join(root, filename)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	envFile, err := godotenv.Parse(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	profileName := "default"
	splits := strings.Split(filename, ".")
	if len(splits) > 2 {
		profileName = strings.Join(splits[2:], ".")
	}

	profile := model.Profile{
		Name:      profileName,
		Variables: envFile,
		Raw:       strings.TrimSuffix(string(file), "\n"),
		Root:      root,
		Filename:  filename,
	}
	return &profile, nil
}

func ReadProfiles(root string) ([]*model.Profile, error) {
	var profileSlice []*model.Profile
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// this prevents walking into subdirectories
		if info.IsDir() && path != root {
			return fs.SkipDir
		}

		filename := info.Name()

		if strings.HasPrefix(filename, ".env") {
			profile, err := ReadProfile(root, filename)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to read profile with root %s and filename %s", root, filename)
			}
			if profile != nil {
				profileSlice = append(profileSlice, profile)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return profileSlice, nil
}

// FIXME: quick and dirty implementation for overriding profile values
func GetProfileValues(currentProfile *model.Profile, profiles []*model.Profile, environmentVars []string) map[string]string {
	profileMap := make(map[string]string)
	if currentProfile == nil || profiles == nil {
		return profileMap
	}
	// take base from default profile
	for _, profile := range profiles {
		if profile.Name == "default" {
			for k, v := range profile.Variables {
				profileMap[k] = v
			}
		}
	}
	// override values with selected profile
	if currentProfile.Name != "default" {
		for k, v := range currentProfile.Variables {
			profileMap[k] = v
		}
	}
	// override values from .local
	for _, profile := range profiles {
		if profile.Name == fmt.Sprintf("%s.local", currentProfile.Name) {
			for k, v := range profile.Variables {
				profileMap[k] = v
			}
		}
	}

	// override and extend with what comes from environment
	for _, e := range environmentVars {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			profileMap[pair[0]] = pair[1]
		}
	}

	// process possible templated variables in profile
	for k, v := range profileMap {
		val := templateng.ProcessTemplateVariableRecursively(v, profileMap)
		profileMap[k] = val
	}

	return profileMap
}
