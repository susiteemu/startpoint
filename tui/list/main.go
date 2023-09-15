package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"goful-cli/client"
	"goful-cli/printer"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center)
var stopwatchStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center).Bold(true)

type request struct {
	name    string
	url     string
	method  string
	headers map[string]string
	body    []byte
}

func (i request) Title() string       { return i.name }
func (i request) Description() string { return i.url }
func (i request) FilterValue() string { return i.name }

type model struct {
	list      list.Model
	stopwatch stopwatch.Model
	selection request
	selected  bool
	resp      string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(request)
			if ok {
				m.selection = i
				m.selected = true
				return m, tea.Batch(
					m.stopwatch.Init(),
					doRequest(m.selection),
				)
			}
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case requestFinishedMsg:
		m.resp = string(msg)
		return m, tea.Quit
	}

	var cmd tea.Cmd
	if m.selected {
		m.stopwatch, cmd = m.stopwatch.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.selected {
		return stopwatchStyle.Render("Making request to " + m.selection.name + " :: Elapsed time: " + m.stopwatch.View())
	}

	return listStyle.Render(m.list.View())
}

type requestFinishedMsg string

func doRequest(r request) tea.Cmd {
	// TODO handle errors
	return func() tea.Msg {
		resp, _ := client.DoRequest(r.url, r.method, r.headers, r.body)
		printed, _ := printer.PrintResponse(resp)
		return requestFinishedMsg(printed)
	}
}

func Start() {
	items := []list.Item{
		request{name: "Raspberry Pi’s", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Raspberry Pi’s\"}")},
		request{name: "Nutella", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Nutella\"}")},
		request{name: "Bitter melon", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Bitter melon\"}")},
		request{name: "Nice socks", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Nice socks\"}")},
		request{name: "Eight hours of sleep", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Eight hours of sleep\"}")},
		request{name: "Cats", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Cats\"}")},
		request{name: "Plantasia, the album", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Plantasia, the album\"}")},
		request{name: "Pour over coffee", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Pour over coffee\"}")},
		request{name: "VR", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"VR\"}")},
		request{name: "Noguchi Lamps", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Noguchi Lamps\"}")},
		request{name: "Linux", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Linux\"}")},
		request{name: "Business school", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Business school\"}")},
		request{name: "Pottery", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Pottery\"}")},
		request{name: "Shampoo", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Shampoo\"}")},
		request{name: "Table tennis", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Table tennis\"}")},
		request{name: "Milk crates", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Milk crates\"}")},
		request{name: "Afternoon tea", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Afternoon tea\"}")},
		request{name: "Stickers", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Stickers\"}")},
		request{name: "20° Weather", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"20° Weather\"}")},
		request{name: "Warm light", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Warm light\"}")},
		request{name: "The vernal equinox", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"The vernal equinox\"}")},
		request{name: "Gaffer’s tape", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Gaffer’s tape\"}")},
		request{name: "Terrycloth", url: "https://httpbin.org/anything", method: "POST", headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, body: []byte("{\"foo\":\"Terrycloth\"}")},
	}

	// Create a new default delegate
	d := list.NewDefaultDelegate()

	// Change colors
	titleColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#cdd6f4"}
	descColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#bac2de"}
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(titleColor).BorderLeftForeground(titleColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(descColor).BorderLeftForeground(descColor)

	m := model{list: list.New(items, d, 0, 0), stopwatch: stopwatch.NewWithInterval(time.Millisecond)}
	m.list.Title = "Requests"

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(model); ok && m.resp != "" {
		fmt.Printf("%s", m.resp)
	}
}
