package requestui

import (
	"fmt"
	"startpoint/core/model"
	"startpoint/core/print"
	"startpoint/core/templating/templateng"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/lipgloss"
	keyprompt "startpoint/tui/keyprompt"
	preview "startpoint/tui/preview"
	profiles "startpoint/tui/profile"
	prompt "startpoint/tui/prompt"
	statusbar "startpoint/tui/statusbar"
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

	var color = style.httpMethodColors[i.Method]
	if color == "" {
		color = style.httpMethodDefaultColor
	}
	methodStyle = methodStyle.Background(lipgloss.Color(color)).Foreground(style.httpMethodTextColor).PaddingRight(1).PaddingLeft(1)

	var urlStyle = lipgloss.NewStyle()
	url := i.Url
	if url == "" {
		url = "<url>"
	}
	urlStyle = urlStyle.Foreground(style.urlFg).Background(style.urlBg)
	if activeProfile != nil && processTemplateVariables {
		for k, v := range activeProfile.Variables {
			processedUrl, match := templateng.ProcessTemplateVariable(url, k, v)
			if match {
				url = processedUrl
			}
		}
		url = print.HighlightWithRegex(url, `{[^{}]*}`, style.urlFg, style.urlBg, style.urlUnfilledTemplatedSectionFg, style.urlUnfilledTemplatedSectionBg)
	} else {
		url = print.HighlightWithRegex(url, `{[^{}]*}`, style.urlFg, style.urlBg, style.urlTemplatedSectionFg, style.urlTemplatedSectionBg)
	}

	method := i.Method
	if i.Method == "" {
		method = "<method>"
	}

	return lipgloss.JoinHorizontal(0, methodStyle.Render(method), " ", url)
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
