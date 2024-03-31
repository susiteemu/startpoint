package managetui

import (
	"errors"
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/loader"
	"goful/core/model"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"goful/core/print"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func doRequest(r Request) tea.Cmd {
	// TODO handle errors
	return func() tea.Msg {
		req, err := builder.BuildRequest(r.Mold, model.Profile{})
		if err != nil {
			return RunRequestFinishedMsg(fmt.Sprintf("failed to build request err: %v", err))
		}
		resp, err := client.DoRequest(req)
		if err != nil {
			return RunRequestFinishedMsg(fmt.Sprintf("failed to do request err: %v", err))
		}

		printed, err := print.SprintPrettyFullResponse(resp)
		if err != nil {
			return RunRequestFinishedMsg(fmt.Sprintf("failed to sprint response err: %v", err))
		}
		return RunRequestFinishedMsg(printed)
	}
}

func handlePostAction(m uiModel) {
	switch m.postAction.Type {
	case PrintRequest:
		fmt.Printf("%s\n", m.postAction.Payload.(string))
	}
}

func createSimpleRequestFileCmd(name string) (string, string, *exec.Cmd, error) {
	if len(name) == 0 {
		return "", "", nil, errors.New("name must not be empty")
	}

	filename := fmt.Sprintf("%s.yaml", name)
	// TODO read from a template file
	content := fmt.Sprintf(`name: %s
# Possible request to call _before_ this one
prev_req:
# Request url, may contain template variables in a form of {var}
url:
# HTTP method
method:
# HTTP headers as key-val list, e.g. X-Foo-Bar: SomeValue
headers:
# Request body, e.g.
# {
#    "id": 1,
#    "name": "Jane">
# }
body: >
`, name)

	cmd, err := createFileAndReturnOpenToEditorCmd(filename, content)
	// TODO get root
	return "tmp/", filename, cmd, err
}

func createComplexRequestFileCmd(name string) (string, string, *exec.Cmd, error) {
	if len(name) == 0 {
		return "", "", nil, errors.New("name must not be empty")
	}

	filename := fmt.Sprintf("%s.star", name)
	// TODO read from template
	content := fmt.Sprintf(`"""
meta:name: %s
meta:prev_req: <call other request before this>
doc:url: <your url for display>
doc:method: <your http method for display>
"""
# insert contents of your script here, for more see https://github.com/google/starlark-go/blob/master/doc/spec.md
# Request url
url = ""
# HTTP method
method = ""
# HTTP headers, e.g. { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
headers = {}
# Request body, e.g. { "id": 1, "people": [ {"name": "Joe"}, {"name": "Jane"}, ] }
body = {}
`, name)

	cmd, err := createFileAndReturnOpenToEditorCmd(filename, content)
	// TODO get root
	return "tmp/", filename, cmd, err
}

func createFileAndReturnOpenToEditorCmd(filename string, content string) (*exec.Cmd, error) {
	if len(filename) > 0 {
		// TODO get workdir from configuration
		file, err := os.Create(filepath.Join("tmp", filename))
		if err == nil {
			defer file.Close()
			// todo handle err
			_, err = file.WriteString(content)
			if err != nil {
				log.Error().Err(err).Msg("Failed to write file")
				return nil, err
			}
			file.Sync()
			fileName := file.Name()

			editor, args, err := getEditor()
			log.Debug().Msgf("Using %s as editor", editor)
			if err != nil {
				log.Error().Err(err)
				return nil, err
			}
			args = append(args, fileName)

			cmd := exec.Command(editor, args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd, nil
		} else {
			log.Error().Err(err).Msg("Failed to create file")
			return nil, err
		}
	}
	return nil, errors.New("filename must not be empty")

}

func openFileToEditorCmd(r Request) (*exec.Cmd, error) {
	if r.Mold.Filename != "" {
		fileName := filepath.Join(r.Mold.Root, r.Mold.Filename)

		log.Info().Msgf("About to open request file %v\n", fileName)
		if len(fileName) > 0 {
			editor, args, err := getEditor()
			log.Debug().Msgf("Using %s as editor", editor)
			if err != nil {
				log.Error().Err(err)
				return nil, err
			}
			args = append(args, fileName)

			cmd := exec.Command(editor, args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd, nil
		}
	}
	return nil, errors.New("request mold does not have a filename")
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

func renameRequest(newName string, r Request) (Request, bool) {
	original := Request{
		Name:   r.Name,
		Url:    r.Url,
		Method: r.Method,
		Mold:   r.Mold.Clone(),
	}

	oldPath := filepath.Join(r.Mold.Root, r.Mold.Filename)
	r.Name = newName
	changeMoldName(newName, &r.Mold)

	log.Info().Msgf("Renaming from %s with name %s", oldPath, newName)
	newPath := filepath.Join(r.Mold.Root, r.Mold.Filename)
	log.Debug().Msgf("Renaming file to %s", newPath)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to rename file to %s", newPath)
		return original, false
	}

	file, err := os.Create(newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to open file %s", newPath)
		return original, false
	}
	defer file.Close()
	log.Debug().Msgf("About to write contents %s", r.Mold.Raw())
	_, err = file.WriteString(r.Mold.Raw())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write to file %s", newPath)
		return original, false
	}
	file.Sync()
	return r, true
}

func copyRequest(name string, r Request) (Request, bool) {
	copy := Request{
		Name:   name,
		Url:    r.Url,
		Method: r.Method,
		Mold:   r.Mold.Clone(),
	}

	changeMoldName(name, &copy.Mold)
	path := filepath.Join(copy.Mold.Root, copy.Mold.Filename)
	file, err := os.Create(path)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to open file %s", path)
		return copy, false
	}
	defer file.Close()
	log.Debug().Msgf("About to write contents %s", copy.Mold.Raw())
	_, err = file.WriteString(copy.Mold.Raw())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write to file %s", path)
		return copy, false
	}
	file.Sync()
	return copy, true

}

func changeMoldName(name string, m *model.RequestMold) {
	if m.Yaml != nil {
		m.Filename = fmt.Sprintf("%s.yaml", name)
		m.Yaml.Name = name
		pattern := regexp.MustCompile(`(?mU)^name:(.*)$`)
		nameChanged := pattern.ReplaceAllString(m.Yaml.Raw, fmt.Sprintf("name: %s", name))
		m.Yaml.Raw = nameChanged
	} else if m.Starlark != nil {
		m.Filename = fmt.Sprintf("%s.star", name)
		pattern := regexp.MustCompile(`(?mU)^.*meta:name:(.*)$`)
		nameChanged := pattern.ReplaceAllString(m.Starlark.Script, fmt.Sprintf("meta:name: %s", name))
		m.Starlark.Script = nameChanged
	}
}

func readRequest(root, filename string) (Request, bool) {
	mold, err := loader.ReadRequest(root, filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request")
		return Request{}, false
	}
	request := Request{
		Name:   mold.Name(),
		Method: mold.Method(),
		Url:    mold.Url(),
		Mold:   mold,
	}
	return request, true
}
