package managetui

import (
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/model"
	"os"
	"os/exec"

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

func createRequestFile(m uiModel) {

	fileName := ""
	content := ""
	createFile := false
	if m.active == Create {
		fileName = fmt.Sprintf("%s.yaml", m.create.Name)
		createFile = true
		// TODO read from a template file
		content = fmt.Sprintf(`name: %s
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
`, m.create.Name)
	} else if m.active == CreateComplex {
		fileName = fmt.Sprintf("%s.star", m.createComplex.Name)
		// TODO read from template
		content = fmt.Sprintf(`"""
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
`, m.createComplex.Name)
		createFile = true
	}

	if !createFile {
		return
	}

	log.Info().Msgf("About to create new request with name %v", fileName)
	if len(fileName) > 0 {
		file, err := os.Create("tmp/" + fileName)
		if err == nil {
			defer file.Close()
			// TODO handle err
			file.WriteString(content)
			file.Sync()
			filename := file.Name()
			editor := viper.GetString("editor")
			if editor == "" {
				log.Error().Msg("Editor is not configured through configuration file or $EDITOR environment variable.")
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

func openRequestFileForUpdate(m uiModel) {
	if m.active == Update && m.selected.Name != "" {

		fileName := fmt.Sprintf("tmp/%s", m.selected.Mold.Filename)

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
