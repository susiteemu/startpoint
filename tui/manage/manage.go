package managetui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	create "goful/tui/requestcreate"
	list "goful/tui/requestlist"
	"log"
	"os"
	"os/exec"
)

type ActiveView int

const (
	List ActiveView = iota
	Create
	Update
)

type keyMap struct {
	Add  key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Quit}, // first column
	}
}

var keys = keyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	list     list.Model
	create   create.Model
	active   ActiveView
	selected list.Request
	width    int
	height   int
	debug    string
	keys     keyMap
	help     help.Model
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
		case "a":
			if m.active != Create {
				m.active = Create
				return m, nil
			}
		}
	case list.RequestSelectedMsg:
		m.active = Update
	case create.CreateMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	switch m.active {
	case List:
		m.list, cmd = m.list.Update(msg)
	case Create:
		m.create, cmd = m.create.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Create:
		return renderCreate(m)
	default:
		return renderList(m)
	}
}

func renderList(m model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.list.View())
}

func renderCreate(m model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.create.View())
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

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := model{list: list.New(requests, 0, 0), create: create.New(), active: List}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(model); ok {
		name := m.create.Name
		log.Printf("About to create new request with name %v", name)
		if len(name) > 0 {
			file, err := os.Create("tmp/" + name)
			if err == nil {
				filename := file.Name()
				editor := viper.GetString("editor")
				if editor == "" {
					log.Fatal("Editor is not configured through configuration file or $EDITOR environment variable.")
				}

				cmd := exec.Command(editor, filename)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Run()
				if err != nil {
					log.Printf("Failed to open file with editor: %v", err)
				}
				log.Printf("Successfully edited file %v", file.Name())
				fmt.Printf("Saved new request to file %v", file.Name())
			} else {
				log.Printf("Failed to create file %v", err)
			}
		}
	}
}
