package writer

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

func WriteFile(path, contents string) (string, error) {
	if len(path) <= 0 {
		return "", errors.New("path must not be empty.")
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString(contents)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write file")
		file.Close()
		return "", err
	}
	err = file.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close the file")
		return "", err
	}
	return file.Name(), nil
}

func RenameFile(oldPath, newPath string) error {
	if len(oldPath) <= 0 {
		return errors.New("old path must not be empty.")
	}

	if len(newPath) <= 0 {
		return errors.New("new path must not be empty.")
	}

	log.Debug().Msgf("Renaming file from %s to %s", oldPath, newPath)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to rename file from %s to %s", oldPath, newPath)
		return err
	}
	return nil
}
