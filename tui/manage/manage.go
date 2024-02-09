package managetui

import (
	"fmt"
	"goful/core/model"
	create "goful/tui/requestcreate"
	list "goful/tui/requestlist"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
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

type uiModel struct {
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

func (m uiModel) Init() tea.Cmd {
	return nil
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m uiModel) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Create:
		return renderCreate(m)
	default:
		return renderList(m)
	}
}

func renderList(m uiModel) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.list.View())
}

func renderCreate(m uiModel) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.create.View())
}

func Start(loadedRequests []model.RequestMold) {
	var requests []list.Request

	for _, v := range loadedRequests {
		r := list.Request{
			Name:   v.Name(),
			Url:    v.Url(),
			Method: v.Method(),
		}
		requests = append(requests, r)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := uiModel{list: list.New(requests, 0, 0), create: create.New(), active: List}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(uiModel); ok {
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
