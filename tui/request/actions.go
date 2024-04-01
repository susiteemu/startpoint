package requestui

import (
	"errors"
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/editor"
	"goful/core/loader"
	"goful/core/model"
	"goful/core/writer"
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

func handlePostAction(m Model) {
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

	cmd, err := createFileAndReturnOpenToEditorCmd("tmp", filename, content)
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

func openFileToEditorCmd(r Request) (*exec.Cmd, error) {
	if len(r.Mold.Filename) <= 0 {
		return nil, errors.New("request mold does not have a filename")
	}
	path := filepath.Join(r.Mold.Root, r.Mold.Filename)
	log.Info().Msgf("About to open request file %v\n", path)
	return editor.OpenFileToEditorCmd(path)
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

	log.Info().Msgf("Renaming from %s to %s", oldPath, newName)
	newPath := filepath.Join(r.Mold.Root, r.Mold.Filename)
	err := writer.RenameFile(oldPath, newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to rename file to %s", newPath)
		return original, false
	}

	_, err = writer.WriteFile(newPath, r.Mold.Raw())
	if err != nil {
		// FIXME in this case, rename file back to what it was?
		log.Error().Err(err).Msgf("Failed to write to file %s", newPath)
		return original, false
	}
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
	_, err := writer.WriteFile(path, copy.Mold.Raw())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write to file %s", path)
		return copy, false
	}
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
