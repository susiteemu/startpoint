package selectui

import (
	"fmt"
	"goful/core/client"
	"goful/core/client/builder"
	"goful/core/model"
	"goful/core/print"
	"os"
	"time"

	list "goful/tui/requestlist"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var stopwatchStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center).Bold(true)

type selectModel struct {
	list      list.Model
	stopwatch stopwatch.Model
	resp      string
	width     int
	height    int
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
		case "ctrl+c":
			return m, tea.Quit
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
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m selectModel) View() string {
	if m.list.Selected {
		return stopwatchStyle.Render("Making request to " + m.list.Selection.Url + " :: Elapsed time: " + m.stopwatch.View())
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.list.View())
}

type requestFinishedMsg string

func doRequest(r list.Request) tea.Cmd {
	// TODO handle errors
	return func() tea.Msg {
		req, err := builder.BuildRequest(*r.Mold, model.Profile{})
		if err != nil {
			return requestFinishedMsg(fmt.Sprintf("err: %v", err))
		}
		resp, err := client.DoRequest(req)
		if err != nil {
			return requestFinishedMsg(fmt.Sprintf("err: %v", err))
		}

		printed, err := print.SprintPrettyFullResponse(resp)
		if err != nil {
			return requestFinishedMsg(fmt.Sprintf("err: %v", err))
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
			Mold:   &v,
		}
		requests = append(requests, r)
	}

	m := selectModel{list: list.New(requests, 0, 0, nil), stopwatch: stopwatch.NewWithInterval(time.Millisecond)}

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
