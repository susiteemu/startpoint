package managetui

import (
	"fmt"
	"goful/core/model"
	list "goful/tui/profile/list"
	create "goful/tui/request/create"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActiveView int

const (
	List ActiveView = iota
	Create
	Update
)

var keys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add new profile"),
	)}

type uiModel struct {
	list     list.Model
	create   create.Model
	active   ActiveView
	selected list.Profile
	width    int
	height   int
	debug    string
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
	case list.ProfileSelectedMsg:
		return m, tea.Quit
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
		lipgloss.Left,
		lipgloss.Top,
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

func Start(loadedProfiles []model.Profile) {
	var profiles []list.Profile

	for _, v := range loadedProfiles {
		r := list.Profile{
			Name:      v.Name,
			Variables: len(v.Variables), // TODO
		}
		profiles = append(profiles, r)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := uiModel{list: list.New(profiles, 0, 0, keys), create: create.New(false), active: List}

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
