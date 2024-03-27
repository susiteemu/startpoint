package managetui

import (
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/model"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

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
			return RequestFinishedMsg(fmt.Sprintf("failed to build request err: %v", err))
		}
		resp, err := client.DoRequest(req)
		if err != nil {
			return RequestFinishedMsg(fmt.Sprintf("failed to do request err: %v", err))
		}

		printed, err := print.SprintPrettyFullResponse(resp)
		if err != nil {
			return RequestFinishedMsg(fmt.Sprintf("failed to sprint response err: %v", err))
		}
		return RequestFinishedMsg(printed)
	}
}

func handlePostAction(m uiModel) {
	switch m.postAction.Type {
	case CreateSimpleRequest:
		createSimpleRequestFile(m.postAction.Payload.(string))
	case CreateComplexRequest:
		createComplexRequestFile(m.postAction.Payload.(string))
	case EditRequest:
		openFileToEditor(m.postAction.Payload.(Request))
	case PrintRequest:
		fmt.Printf("%s\n", m.postAction.Payload.(string))
	}
}

func createSimpleRequestFile(name string) {
	if len(name) == 0 {
		return
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

	createFileAndOpenToEditor(filename, content)
}

func createComplexRequestFile(name string) {
	if len(name) == 0 {
		return
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

	createFileAndOpenToEditor(filename, content)
}

func createFileAndOpenToEditor(filename string, content string) {

	if len(filename) > 0 {
		// TODO get workdir from configuration
		file, err := os.Create(filepath.Join("tmp", filename))
		if err == nil {
			defer file.Close()
			// todo handle err
			_, err = file.WriteString(content)
			if err != nil {
				log.Error().Err(err).Msg("Failed to write file")
			}
			file.Sync()
			filename := file.Name()
			editor := viper.GetString("editor")
			if editor == "" {
				log.Error().Msg("Editor is not configured through configuration file or $editor environment variable.")
			}

			log.Info().Msgf("Opening file %s\n", filename)
			cmd := exec.Command(editor, filename)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				log.Error().Err(err).Msg("Failed to open file with editor")
			}
			log.Printf("Successfully edited file %v", file.Name())
			fmt.Printf("Saved new request to file %v", file.Name())
		} else {
			log.Error().Err(err).Msg("Failed to create file")
		}
	}

}

func openFileToEditor(r Request) {
	if r.Mold.Filename != "" {
		fileName := filepath.Join(r.Mold.Root, r.Mold.Filename)

		log.Info().Msgf("About to open request file %v\n", fileName)
		if len(fileName) > 0 {
			// TODO handle err
			editor := viper.GetString("editor")
			if editor == "" {
				log.Error().Msg("Editor is not configured through configuration file or $EDITOR environment variable.")
			}

			cmd := exec.Command(editor, fileName)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				log.Error().Err(err).Msg("Failed to open file with editor")
			}
			log.Info().Msgf("Successfully edited file %v", fileName)
		}
	}
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
