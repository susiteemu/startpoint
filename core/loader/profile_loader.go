package loader

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
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
	profileName := ""
	if filename == ".env" {
		profileName = "default"
	} else if filename == ".env.local" {
		profileName = "default.local"
	} else {
		splits := strings.Split(filename, ".")
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

type groupedProfile struct {
	public  *model.Profile
	private *model.Profile
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

	groupedProfiles := []groupedProfile{}
	for _, v := range profileSlice {
		if v.IsPrivateProfile() {
			hasPublicProfile := false
			for _, vv := range profileSlice {
				if !vv.IsPrivateProfile() && strings.TrimSuffix(v.Name, ".local") == vv.Name {
					hasPublicProfile = true
					break
				}
			}
			if !hasPublicProfile {
				groupedProfiles = append(groupedProfiles, groupedProfile{
					public:  nil,
					private: v,
				})
			}
		} else {
			var privateProfile *model.Profile
			for _, vv := range profileSlice {
				if vv.IsPrivateProfile() && strings.TrimSuffix(vv.Name, ".local") == v.Name {
					privateProfile = vv
					break
				}
			}
			groupedProfiles = append(groupedProfiles, groupedProfile{
				public:  v,
				private: privateProfile,
			})

		}
	}
	slices.SortFunc(groupedProfiles, func(a, b groupedProfile) int {
		nameA := ""
		nameB := ""
		if a.public != nil {
			nameA = a.public.Name
		} else if a.private != nil {
			nameA = a.private.Name
		}
		if b.public != nil {
			nameB = b.public.Name
		} else if b.private != nil {
			nameB = b.private.Name
		}
		// make comparison case insensitive
		nameA = strings.ToLower(nameA)
		nameB = strings.ToLower(nameB)
		if nameA < nameB {
			return -1
		}
		if nameA > nameB {
			return 1
		}
		return 0
	})

	profileSlice = []*model.Profile{}
	for _, v := range groupedProfiles {
		if v.public != nil {
			v.public.HasPublicProfile = true
			v.public.HasPrivateProfile = v.private != nil
			profileSlice = append(profileSlice, v.public)
		}
		if v.private != nil {
			v.private.HasPrivateProfile = true
			v.private.HasPublicProfile = v.public != nil
			profileSlice = append(profileSlice, v.private)
		}
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
