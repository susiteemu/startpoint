package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"goful-cli/client"
	"goful-cli/printer"
)

type model struct {
	resp      string
	stopwatch stopwatch.Model
	keymap    keymap
	help      help.Model
	quitting  bool
}

type keymap struct {
	quit key.Binding
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.stopwatch.Init(),
		doRequest,
	)
}

func (m model) View() string {
	s := m.stopwatch.View() + "\n"
	if !m.quitting {
		s = "Elapsed: " + s
		s += m.helpView()
		return s
	}
	return "Request took: " + s + "\n" + m.resp
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.quit,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		}
	case processFinishedMsg:
		m.resp = string(msg)
		m.quitting = true
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

type processFinishedMsg string

func doRequest() tea.Msg {
	// TODO handle errors
	var url = "https://httpbin.org/anything"
	var headers = map[string]string{"X-Foo": "bar", "X-Bar": "foo"}
	var body = []byte("{\"foo\":\"hello\"}")
	resp, _ := client.DoRequest(url, "POST", headers, body)

	printed, _ := printer.PrintResponse(resp)
	return processFinishedMsg(printed)
}

func main() {
	m := model{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
