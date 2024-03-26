package managetui

import (
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/model"
	"os"
	"os/exec"
	"path/filepath"

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
		openRequestFileForUpdate(m.postAction.Payload.(Request))
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

	createFile(filename, content)
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

	createFile(filename, content)
}

func createFile(filename string, content string) {

	if len(filename) > 0 {
		// TODO get workdir from configuration
		file, err := os.Create(filepath.Join("tmp", filename))
		if err == nil {
			defer file.Close()
			// todo handle err
			file.WriteString(content)
			file.Sync()
			filename := file.Name()
			editor := viper.GetString("editor")
			if editor == "" {
				log.Error().Msg("editor is not configured through configuration file or $editor environment variable.")
			}

			log.Info().Msgf("opening file %s\n", filename)
			cmd := exec.Command(editor, filename)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				log.Error().Err(err).Msg("failed to open file with editor")
			}
			log.Printf("successfully edited file %v", file.Name())
			fmt.Printf("saved new request to file %v", file.Name())
		} else {
			log.Error().Err(err).Msg("failed to create file")
		}
	}

}

func openRequestFileForUpdate(r Request) {
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
	renamed := Request{
		Name:   r.Name,
		Url:    r.Url,
		Method: r.Method,
		Mold:   r.Mold,
	}
	oldPath := filepath.Join(r.Mold.Root, r.Mold.Filename)
	log.Info().Msgf("Renaming from %s with newName %s", oldPath, newName)
	ok := r.Mold.Rename(newName)
	log.Debug().Msgf("Renaming data succeed? %v", ok)
	if ok {
		renamed.Name = newName
		renamed.Mold = r.Mold
		newPath := filepath.Join(renamed.Mold.Root, renamed.Mold.Filename)
		log.Debug().Msgf("Renaming file to %s", newPath)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to rename file to %s", newPath)
			return renamed, false
		}

		file, err := os.Create(newPath)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to open file %s", newPath)
			return renamed, false
		}
		defer file.Close()
		log.Debug().Msgf("About to write contents %s", renamed.Mold.Raw)
		_, err = file.WriteString(renamed.Mold.Raw)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to write to file %s", newPath)
			return renamed, false
		}
	}
	return renamed, ok
}
