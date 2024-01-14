package selectui

import (
	"fmt"
	"goful/core/client"
	"goful/core/model"
	"goful/core/print"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	list "goful/tui/requestlist"
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
		return stopwatchStyle.Render("Making request to " + m.list.Selection.Name + " :: Elapsed time: " + m.stopwatch.View())
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
		resp, _ := client.DoRequest(model.Request{
			Url:     r.Url,
			Method:  r.Method,
			Headers: new(model.Headers).FromMap(r.Headers),
			Body:    r.Body,
		})
		printed, _ := print.SprintPrettyFullResponse(resp)
		return requestFinishedMsg(printed)
	}
}

func Start() {
	requests := []list.Request{
		{Name: "Raspberry Pi’s", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Raspberry Pi’s\"}")},
		{Name: "Nutella", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Nutella\"}")},
		{Name: "Bitter melon", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Bitter melon\"}")},
		{Name: "Nice socks", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Nice socks\"}")},
		{Name: "Eight hours of sleep", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Eight hours of sleep\"}")},
		{Name: "Cats", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Cats\"}")},
		{Name: "Plantasia, the album", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Plantasia, the album\"}")},
		{Name: "Pour over coffee", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Pour over coffee\"}")},
		{Name: "VR", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"VR\"}")},
		{Name: "Noguchi Lamps", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Noguchi Lamps\"}")},
		{Name: "Linux", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Linux\"}")},
		{Name: "Business school", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Business school\"}")},
		{Name: "Pottery", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Pottery\"}")},
		{Name: "Shampoo", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Shampoo\"}")},
		{Name: "Table tennis", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Table tennis\"}")},
		{Name: "Milk crates", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Milk crates\"}")},
		{Name: "Afternoon tea", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Afternoon tea\"}")},
		{Name: "Stickers", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Stickers\"}")},
		{Name: "20° Weather", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"20° Weather\"}")},
		{Name: "Warm light", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Warm light\"}")},
		{Name: "The vernal equinox", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"The vernal equinox\"}")},
		{Name: "Gaffer’s tape", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Gaffer’s tape\"}")},
		{Name: "Terrycloth", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string][]string{"X-Foo": {"bar"}, "X-Bar": {"foo"}, "Content-Type": {"application/json"}}, Body: []byte("{\"foo\":\"Terrycloth\"}")},
	}

	m := selectModel{list: list.New(requests, 0, 0), stopwatch: stopwatch.NewWithInterval(time.Millisecond)}

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
