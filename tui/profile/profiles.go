package managetui

import (
	"fmt"
	"goful/core/model"
	create "goful/tui/request/prompt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActiveView int

const (
	List ActiveView = iota
	Create
	Update
)

type Profile struct {
	Name      string
	Variables int
}

func (i Profile) Title() string       { return i.Name }
func (i Profile) Description() string { return fmt.Sprintf("Vars: %d", i.Variables) }
func (i Profile) FilterValue() string { return i.Name }

type uiModel struct {
	list   list.Model
	create create.Model
	active ActiveView
	width  int
	height int
}

func (m uiModel) Init() tea.Cmd {
	return nil
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)

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
	case ProfileSelectedMsg:
		return m, tea.Quit
	case create.PromptAnsweredMsg:
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
	var profiles []list.Item

	for _, v := range loadedProfiles {
		r := Profile{
			Name:      v.Name,
			Variables: len(v.Variables), // TODO
		}
		profiles = append(profiles, r)
	}

	d := newSelectDelegate()

	profileList := list.New(profiles, d, 0, 0)
	profileList.Title = "Profiles"
	profileList.Styles.Title = titleStyle

	profileList.Help.Styles.FullKey = helpKeyStyle
	profileList.Help.Styles.FullDesc = helpDescStyle
	profileList.Help.Styles.ShortKey = helpKeyStyle
	profileList.Help.Styles.ShortDesc = helpDescStyle
	profileList.Help.Styles.ShortSeparator = helpSeparatorStyle
	profileList.Help.Styles.FullSeparator = helpSeparatorStyle

	profileList.SetShowStatusBar(false)
	profileList.SetFilteringEnabled(false)

	m := uiModel{list: profileList, active: List}

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
