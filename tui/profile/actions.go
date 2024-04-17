package profileui

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"startpoint/core/editor"
	"startpoint/core/writer"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func createProfileCmd(name string) (string, string, *exec.Cmd, error) {
	filename := fmt.Sprintf("%s.env", name)
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
