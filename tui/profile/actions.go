package profileui

import (
	"errors"
	"fmt"
	"goful/core/editor"
	"goful/core/writer"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func createProfileCmd(name string) (string, string, *exec.Cmd, error) {
	filename := fmt.Sprintf("%s.env", name)
	content := ""
	cmd, err := createFileAndReturnOpenToEditorCmd("tmp", filename, content)
	// TODO get root
	return "tmp/", filename, cmd, err
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