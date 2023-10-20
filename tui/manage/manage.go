package managetui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	edit "goful/tui/requestedit"
	list "goful/tui/requestlist"
)

type model struct {
	list     list.Model
	edit     edit.Model
	selected list.Request
	width    int
	height   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.selected = m.list.Selection
		}
	}

	var cmd tea.Cmd
	if m.list.Selected {
		m.edit, cmd = m.edit.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.list.Selected {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Top,
			m.edit.View())
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.list.View())
}

func Start() {
	requests := []list.Request{
		{Name: "Raspberry Pi’s", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Raspberry Pi’s\"}")},
		{Name: "Nutella", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Nutella\"}")},
		{Name: "Bitter melon", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Bitter melon\"}")},
		{Name: "Nice socks", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Nice socks\"}")},
		{Name: "Eight hours of sleep", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Eight hours of sleep\"}")},
		{Name: "Cats", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Cats\"}")},
		{Name: "Plantasia, the album", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Plantasia, the album\"}")},
		{Name: "Pour over coffee", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Pour over coffee\"}")},
		{Name: "VR", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"VR\"}")},
		{Name: "Noguchi Lamps", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Noguchi Lamps\"}")},
		{Name: "Linux", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Linux\"}")},
		{Name: "Business school", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Business school\"}")},
		{Name: "Pottery", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Pottery\"}")},
		{Name: "Shampoo", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Shampoo\"}")},
		{Name: "Table tennis", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Table tennis\"}")},
		{Name: "Milk crates", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Milk crates\"}")},
		{Name: "Afternoon tea", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Afternoon tea\"}")},
		{Name: "Stickers", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Stickers\"}")},
		{Name: "20° Weather", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"20° Weather\"}")},
		{Name: "Warm light", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Warm light\"}")},
		{Name: "The vernal equinox", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"The vernal equinox\"}")},
		{Name: "Gaffer’s tape", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Gaffer’s tape\"}")},
		{Name: "Terrycloth", Url: "https://httpbin.org/anything", Method: "POST", Headers: map[string]string{"X-Foo": "bar", "X-Bar": "foo"}, Body: []byte("{\"foo\":\"Terrycloth\"}")},
	}

	m := model{list: list.New(requests, 0, 0), edit: edit.New()}

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
