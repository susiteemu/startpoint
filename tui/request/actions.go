package requestui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	requestchain "startpoint/core/chaining"
	"startpoint/core/client/runner"
	"startpoint/core/configuration"
	"startpoint/core/editor"
	"startpoint/core/loader"
	"startpoint/core/model"
	"startpoint/core/writer"
	"strings"
	"time"

	"startpoint/core/print"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func doRequest(r *model.RequestMold, all []*model.RequestMold, profile *model.Profile) tea.Cmd {
	// TODO: handle errors
	return func() tea.Msg {

		chainedRequests := requestchain.ResolveRequestChain(r, all)

		log.Debug().Msgf("Resolved %d chained requests", len(chainedRequests))

		responses, err := runner.RunRequestChain(chainedRequests, profile, interimResult)
		if err != nil {
			return RunRequestFinishedWithFailureMsg(fmt.Sprintf("Error occurred:\n%v", err))
		}
		var printedResponses []string
		for _, resp := range responses {
			var config *configuration.Configuration = configuration.NewWithRequestOptions(resp.Options)
			printResponse := config.GetBoolWithDefault("print", true)
			if !printResponse {
				log.Info().Msgf("Skipping printing response")
				continue
			}
			printOpts := print.PrintOpts{
				PrettyPrint:    config.GetBoolWithDefault("printer.pretty", true),
				PrintBody:      true,
				PrintHeaders:   true,
				PrintTraceInfo: config.GetBool("httpClient.enableTraceInfo"),
			}
			log.Debug().Msgf("Printing with opts %v", printOpts)
			printed, err := print.SprintResponse(resp, printOpts)
			if err != nil {
				return RunRequestFinishedWithFailureMsg(fmt.Sprintf("Error occurred: %v", err))
			}
			printedResponses = append(printedResponses, printed)
		}

		return RunRequestFinishedMsg(strings.Join(printedResponses, "\n"))
	}

}

func interimResult(took time.Duration, statusCode int) {
	log.Debug().Msgf("Request with statuscode %d took %s", statusCode, took.String())
}

func handlePostAction(m Model) {
	switch m.postAction.Type {
	case PrintFailedRequest:
		fmt.Fprintf(os.Stderr, m.postAction.Payload.(string)+"\n")
		os.Exit(1)
	case PrintRequest:
		fmt.Printf(m.postAction.Payload.(string) + "\n")
	}
}

func createSimpleRequestFileCmd(name string) (string, string, *exec.Cmd, error) {
	if len(name) == 0 {
		return "", "", nil, errors.New("name must not be empty")
	}

	filename := fmt.Sprintf("%s.yaml", name)
	// TODO read from a template file
	content := `# Possible request to call _before_ this one
prev_req:
# Request url, may contain template variables in a form of {var}
url:
# HTTP method
method: GET
# HTTP headers as key-val list, e.g. X-Foo-Bar: SomeValue
headers:
# Request body, e.g.
# {
#    "id": 1,
#    "name": "Jane">
# }
body: >
`

	workspace := viper.GetString("workspace")
	cmd, err := createFileAndReturnOpenToEditorCmd(workspace, filename, content)
	return workspace, filename, cmd, err
}

func createComplexRequestFileCmd(name string) (string, string, *exec.Cmd, error) {
	if len(name) == 0 {
		return "", "", nil, errors.New("name must not be empty")
	}

	filename := fmt.Sprintf("%s.star", name)
	// TODO read from template
	content := `"""
prev_req: <call other request before this>
doc:url: <your url for display>
doc:method: GET
"""
# insert contents of your script here, for more see https://github.com/google/starlark-go/blob/master/doc/spec.md
# Request url
url = ""
# HTTP method
method = "GET"
# HTTP headers, e.g. { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
headers = {}
# Request body, e.g. { "id": 1, "people": [ {"name": "Joe"}, {"name": "Jane"}, ] }
body = {}
`

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
		return nil, errors.New("request mold does not have a filename")
	}
	path := filepath.Join(root, filename)
	log.Info().Msgf("About to open request file %v\n", path)
	return editor.OpenFileToEditorCmd(path)
}

func renameRequest(newName string, r Request, mold model.RequestMold) (Request, *model.RequestMold, bool) {
	original := Request{
		Name:   r.Name,
		Url:    r.Url,
		Method: r.Method,
	}
	originalMold := mold.Clone()

	oldPath := filepath.Join(mold.Root, mold.Filename)
	r.Name = newName
	changeMoldName(newName, &mold)

	log.Info().Msgf("Renaming from %s to %s", oldPath, newName)
	newPath := filepath.Join(mold.Root, mold.Filename)
	err := writer.RenameFile(oldPath, newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to rename file to %s", newPath)
		return original, &originalMold, false
	}

	_, err = writer.WriteFile(newPath, mold.Raw())
	if err != nil {
		// FIXME in this case, rename file back to what it was?
		log.Error().Err(err).Msgf("Failed to write to file %s", newPath)
		return original, &originalMold, false
	}
	return r, &mold, true
}

func copyRequest(name string, r Request, mold model.RequestMold) (Request, *model.RequestMold, bool) {
	copy := Request{
		Name:   name,
		Url:    r.Url,
		Method: r.Method,
	}
	copyMold := mold.Clone()

	changeMoldName(name, &copyMold)

	path := filepath.Join(copyMold.Root, copyMold.Filename)
	_, err := writer.WriteFile(path, copyMold.Raw())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write to file %s", path)
		return r, &mold, false
	}
	return copy, &copyMold, true
}

func changeMoldName(name string, m *model.RequestMold) {
	if m.Yaml != nil {
		m.Filename = fmt.Sprintf("%s.yaml", name)
		m.Name = name
	} else if m.Starlark != nil {
		m.Filename = fmt.Sprintf("%s.star", name)
		m.Name = name
	}
}

func readRequest(root, filename string) (Request, *model.RequestMold, bool) {
	mold, err := loader.ReadRequest(root, filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request")
		return Request{}, nil, false
	}
	request := Request{
		Name:   mold.Name,
		Method: mold.Method(),
		Url:    mold.Url(),
	}
	return request, mold, true
}
