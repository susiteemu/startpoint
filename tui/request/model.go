package requestui

import (
	"fmt"
	"goful/core/model"
	"goful/core/templating/templateng"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/lipgloss"
	keyprompt "goful/tui/keyprompt"
	preview "goful/tui/preview"
	profiles "goful/tui/profile"
	prompt "goful/tui/prompt"
	statusbar "goful/tui/statusbar"
)

type ActiveView int
type Mode int

type Request struct {
	Name    string
	Url     string
	Method  string
	Profile *model.Profile
}

func (i Request) Title() string {
	return i.Name
}

func (i Request) Description() string {
	var methodStyle = lipgloss.NewStyle()

	var color = methodColors[i.Method]
	if color == "" {
		color = "#cdd6f4"
	}
	methodStyle = methodStyle.Background(lipgloss.Color(color)).Foreground(lipgloss.Color("#1e1e2e")).PaddingRight(1).PaddingLeft(1)

	var urlStyle = lipgloss.NewStyle()
	url := i.Url
	if activeProfile != nil && processTemplateVariables {
		for k, v := range activeProfile.Variables {
			processedUrl, match := templateng.ProcessTemplateVariable(url, k, v)
			if match {
				url = processedUrl
			}
		}
	}

	urlStyle = urlStyle.Foreground(lipgloss.Color("#b4befe"))
	method := i.Method
	if i.Method == "" {
		method = "<method>"
	}

	if url == "" {
		url = "<url>"
	}

	return lipgloss.JoinHorizontal(0, methodStyle.Render(method), " ", urlStyle.Render(url))
}
func (i Request) FilterValue() string { return fmt.Sprintf("%s %s %s", i.Name, i.Method, i.Url) }

type Model struct {
	mode         Mode
	active       ActiveView
	list         list.Model
	preview      preview.Model
	prompt       prompt.Model
	keyprompt    keyprompt.Model
	stopwatch    stopwatch.Model
	statusbar    statusbar.Model
	profileui    profiles.Model
	help         help.Model
	width        int
	height       int
	postAction   PostAction
	requestMolds []*model.RequestMold
}
