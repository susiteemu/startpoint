package editor

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func OpenFileToEditorCmd(filepath string) (*exec.Cmd, error) {
	if len(filepath) <= 0 {
		return nil, errors.New("filepath must not be empty.")
	}

	editor, args, err := getEditor()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get editor")
		return nil, err
	}
	args = append(args, filepath)
	log.Debug().Msgf("Using %s as editor with args %v", editor, args)
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func getEditor() (string, []string, error) {
	editor := strings.Fields(viper.GetString("editor"))
	if len(editor) == 1 {
		return editor[0], []string{}, nil
	}
	if len(editor) > 1 {
		return editor[0], editor[1:], nil
	}
	// TODO read default editor from configuration
	return "", []string{}, errors.New("Editor is not configured through configuration file or $EDITOR environment variable.")
}
