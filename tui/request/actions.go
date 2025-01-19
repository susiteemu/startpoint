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
		responseCount := len(responses)
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
				PrintRequest:   config.GetBoolWithDefault("printRequest", false),
			}
			log.Debug().Msgf("Printing with opts %v", printOpts)
			printed, err := print.SprintResponse(resp, printOpts)
			if err != nil {
				return RunRequestFinishedWithFailureMsg(fmt.Sprintf("Error occurred: %v", err))
			}
			// FIXME: instead of responseCount, this should check if there are > 1 responses that would be printed
			if responseCount > 1 {
				printedResponses = append(printedResponses, print.SprintFaint(fmt.Sprintf(`#
# %s
#`, resp.RequestName)))

			}
			printedResponses = append(printedResponses, printed)
			printedResponses = append(printedResponses, "")
		}

		return RunRequestFinishedMsg(strings.Join(printedResponses, "\n"))
	}

}

func interimResult(took time.Duration, statusCode int) {
	log.Debug().Msgf("Request with statuscode %d took %s", statusCode, took.String())
}

func (m Model) HandlePostAction() {
	switch m.postAction.Type {
	case PrintFailedRequest:
		fmt.Fprint(os.Stderr, m.postAction.Payload.(string)+"\n")
		os.Exit(1)
	case PrintRequest:
		fmt.Print(m.postAction.Payload.(string) + "\n")
	}
}

func createRequestFileCmd(name string, requestType string) (string, string, *exec.Cmd, error) {
	if len(name) == 0 {
		return "", "", nil, errors.New("name must not be empty")
	}
	var filename, content string

	switch requestType {
	case model.CONTENT_TYPE_YAML:
		filename = fmt.Sprintf("%s.yaml", name)
		content = YamlTemplate
	case model.CONTENT_TYPE_STARLARK:
		filename = fmt.Sprintf("%s.star", name)
		content = StarlarkTemplate
	case model.CONTENT_TYPE_LUA:
		filename = fmt.Sprintf("%s.lua", name)
		content = LuaTemplate
	default:
		return "", "", nil, errors.New("unsupported request type")
	}

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

// TODO: refactor from bool to error
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

// TODO: refactor from bool to error
func changePrevReq(oldPrevReq string, newPrevReq string, molds []*model.RequestMold) ([]*model.RequestMold, bool) {
	for _, mold := range molds {
		if mold.PreviousReq() == oldPrevReq {
			mold.ChangePreviousReq(newPrevReq)
			path := filepath.Join(mold.Root, mold.Filename)
			_, err := writer.WriteFile(path, mold.Raw())
			if err != nil {
				log.Error().Err(err).Msgf("Failed to write to file %s", path)
				return molds, false
			}
		}
	}
	return molds, true
}

func isUsedAsPrevReq(name string, molds []*model.RequestMold) bool {
	for _, mold := range molds {
		if mold.PreviousReq() == name {
			return true
		}
	}
	return false
}

// TODO: refactor from bool to error
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
	m.Filename = fmt.Sprintf("%s.%s", name, m.Type)
	m.Name = name
}

// TODO: refactor from bool to error
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

func RefreshProfiles(loadedProfiles []*model.Profile) {
	envVars := os.Environ()
	allProfiles = []*model.Profile{}
	for _, p := range loadedProfiles {
		if p.IsPrivateProfile() && p.HasPublicProfile {
			continue
		}
		profile := &model.Profile{
			Name:      p.Name,
			Variables: loader.GetProfileValues(p, loadedProfiles, envVars),
		}
		if profile.Name == "default" {
			activeProfile = profile
		}
		allProfiles = append(allProfiles, profile)
	}

}
