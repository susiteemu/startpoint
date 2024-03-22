package selectui

import (
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/model"
	"goful/core/print"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	list "goful/tui/request/list"
	preview "goful/tui/request/preview"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActiveView int

const (
	List ActiveView = iota
	Preview
)

var keys = []key.Binding{
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Preview"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
		key.WithHelp(tea.KeyEnter.String(), "Run request"),
	),
}

var stopwatchStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center).Bold(true)

type selectModel struct {
	list      list.Model
	preview   preview.Model
	stopwatch stopwatch.Model
	resp      string
	width     int
	height    int
	active    ActiveView
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			if m.active == Preview {
				m.active = List
				return m, nil
			}
			return m, tea.Quit
		case "p":
			if m.active != Preview {
				m.active = Preview
				m.preview.Viewport.Width = m.width
				m.preview.Viewport.Height = m.height - m.preview.VerticalMarginHeight()
				selected := m.list.SelectedItem()
				var formatted string
				var err error
				switch selected.Mold.ContentType {
				case "yaml":
					formatted, err = print.SprintYaml(selected.Mold.Raw)
				case "star":
					formatted, err = print.SprintStar(selected.Mold.Raw)
				}

				if formatted == "" || err != nil {
					formatted = selected.Mold.Raw
				}
				m.preview.Viewport.SetContent(formatted)
				m.preview.Viewport.YPosition = 0
				return m, nil
			}
		}

	case list.RequestSelectedMsg:
		selected := m.list.Selected
		if selected {
			return m, tea.Batch(
				m.stopwatch.Init(),
				doRequest(m.list.Selection),
			)
		}
	case requestFinishedMsg:
		m.resp = string(msg)
		return m, tea.Quit
	}

	var cmd tea.Cmd
	if m.list.Selected {
		m.stopwatch, cmd = m.stopwatch.Update(msg)
	} else if m.active == Preview {
		m.preview, cmd = m.preview.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m selectModel) View() string {
	if m.list.Selected {
		return stopwatchStyle.Render("Making request to " + m.list.Selection.Url + " :: Elapsed time: " + m.stopwatch.View())
	} else if m.active == Preview {
		return m.preview.View()
	}

	return m.list.View()
}

type requestFinishedMsg string

func doRequest(r list.Request) tea.Cmd {
	// TODO handle errors
	return func() tea.Msg {
		req, err := builder.BuildRequest(r.Mold, model.Profile{})
		if err != nil {
			return requestFinishedMsg(fmt.Sprintf("failed to build request err: %v", err))
		}
		resp, err := client.DoRequest(req)
		if err != nil {
			return requestFinishedMsg(fmt.Sprintf("failed to do request err: %v", err))
		}

		printed, err := print.SprintPrettyFullResponse(resp)
		if err != nil {
			return requestFinishedMsg(fmt.Sprintf("failed to sprint response err: %v", err))
		}
		return requestFinishedMsg(printed)
	}
}

func Start(loadedRequests []model.RequestMold) {
	var requests []list.Request

	for _, v := range loadedRequests {
		r := list.Request{
			Name:   v.Name(),
			Url:    v.Url(),
			Method: v.Method(),
			Mold:   v,
		}
		requests = append(requests, r)
	}

	m := selectModel{list: list.New(requests, true, 0, 0, keys), stopwatch: stopwatch.NewWithInterval(time.Millisecond), preview: preview.New()}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(selectModel); ok && m.resp != "" {
		fmt.Printf("%s", m.resp)
	}
}
