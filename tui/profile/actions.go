package profileui

import (
	"errors"
	"fmt"
	"github.com/susiteemu/startpoint/core/editor"
	"github.com/susiteemu/startpoint/core/loader"
	"github.com/susiteemu/startpoint/core/writer"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func createProfileFileCmd(name string) (string, string, *exec.Cmd, error) {
	filename := ""
	if len(strings.TrimSpace(name)) == 0 || name == "default" {
		filename = ".env"
	} else {
		filename = fmt.Sprintf(".env.%s", name)
	}
	content := ""
	workspace := viper.GetString("workspace")
	cmd, err := createFileAndReturnOpenToEditorCmd(workspace, filename, content)
	return workspace, filename, cmd, err
}

func createFileAndReturnOpenToEditorCmd(root, filename, content string) (*exec.Cmd, error) {
	if len(root) <= 0 {
		return nil, errors.New("root must not be empty")
	}
	if len(filename) <= 0 {
		return nil, errors.New("filename must not be empty")
	}

	path, err := writer.WriteFile(filepath.Join(root, filename), content)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return nil, err
	}

	return editor.OpenFileToEditorCmd(path)
}

func openFileToEditorCmd(root, filename string) (*exec.Cmd, error) {
	if len(filename) <= 0 {
		return nil, errors.New("profile does not have a filename")
	}
	path := filepath.Join(root, filename)
	log.Info().Msgf("About to open profile file %v\n", path)
	return editor.OpenFileToEditorCmd(path)
}

func readProfile(root, filename string) (Profile, bool) {
	profile, err := loader.ReadProfile(root, filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read profile")
		return Profile{}, false
	}
	return Profile{
		Name:         profile.Name,
		Variables:    len(profile.Variables),
		ProfileModel: profile,
	}, true
}

func renameProfile(name string, profile Profile) (Profile, bool) {
	oldPath := filepath.Join(profile.ProfileModel.Root, profile.ProfileModel.Filename)
	newName := fmt.Sprintf(".env.%s", name)
	newPath := filepath.Join(profile.ProfileModel.Root, newName)
	log.Info().Msgf("Renaming from %s to %s", oldPath, newPath)
	err := writer.RenameFile(oldPath, newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to rename file to %s", newPath)
		return profile, false
	}
	return readProfile(profile.ProfileModel.Root, newName)
}

func copyProfile(name string, profile Profile) (Profile, bool) {
	if len(name) == 0 {
		log.Error().Msg("Can't copy profile to empty name")
		return Profile{}, false
	}
	filename := fmt.Sprintf(".env.%s", name)
	path := filepath.Join(profile.ProfileModel.Root, filename)
	_, err := writer.WriteFile(path, profile.ProfileModel.Raw)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write file %s", path)
		return Profile{}, false
	}
	return readProfile(profile.ProfileModel.Root, filename)
}
