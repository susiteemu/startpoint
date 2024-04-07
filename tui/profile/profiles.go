package profileui

import (
	"fmt"
	"goful/core/model"
	prompt "goful/tui/prompt"
	statusbar "goful/tui/statusbar"
	"goful/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

type Mode int

const (
	Normal Mode = iota
	Embedded
)

type Profile struct {
	Name      string
	Variables int
}

/*
* In embedded mode we disable help from list bubble and instead show our own:
* could not find a reasonable way to remove key bindings from list's help and
* in embedded mode we only really want to see select/cancel keys
*
 */
func (k embeddedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Cancel}
}

func (k embeddedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Cancel, k.Cancel},
	}
}

func updateStatusbar(m *Model, msg string) {
	msgItem := statusbar.StatusbarItem{
		Text: msg, BackgroundColor: statusbarFirstColBg, ForegroundColor: statusbarFirstColFg,
	}
	m.statusbar.SetItem(msgItem, 0)
}

func (i Profile) Title() string       { return i.Name }
func (i Profile) Description() string { return fmt.Sprintf("Vars: %d", i.Variables) }
func (i Profile) FilterValue() string { return i.Name }

type Model struct {
	list      list.Model
	prompt    prompt.Model
	help      help.Model
	statusbar statusbar.Model
	active    ActiveView
	mode      Mode
	width     int
	height    int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		if m.mode == Normal {
			m.statusbar.SetWidth(msg.Width)
			m.help.Width = msg.Width
			listHeight := calculateListHeight(m)
			m.list.SetHeight(listHeight)
			updateStatusbar(&m, "")
		} else {
			m.list.SetHeight(m.height - 2)
		}

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "a":
			if m.mode == Normal && m.active != Create {
				m.active = Create
				return m, nil
			}
		case "?":
			if m.active == List {
				m.help.ShowAll = !m.help.ShowAll
				listHeight := calculateListHeight(m)
				m.list.SetHeight(listHeight)

				return m, nil
			}
		}
	case ProfileSelectedMsg:
		return m, tea.Quit
	case prompt.PromptAnsweredMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	switch m.active {
	case List:
		m.list, cmd = m.list.Update(msg)
	case Create:
		m.prompt, cmd = m.prompt.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Create:
		return renderCreate(m)
	default:
		return renderList(m)
	}
}

func renderList(m Model) string {
	if m.mode == Embedded {
		helpView := m.help.View(embeddedKeys)

		views := []string{}
		views = append(views, m.list.View())
		views = append(views, helpView)

		return lipgloss.JoinVertical(lipgloss.Center, views...)

	}
	var views []string
	if m.help.ShowAll {
		listHeight := calculateListHeight(m)
		views = append(views, lipgloss.NewStyle().Height(listHeight).Render(m.list.View()))
		views = append(views, m.statusbar.View())
		views = append(views, styles.HelpPaneStyle.Render(m.help.View(m.list)))
	} else {
		listHeight := calculateListHeight(m)
		views = append(views, lipgloss.NewStyle().Height(listHeight).Render(m.list.View()))
		views = append(views, m.statusbar.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		views...,
	)
}

func calculateListHeight(m Model) int {
	listHeight := m.height - statusbar.Height
	if m.help.ShowAll {
		helpHeight := lipgloss.Height(styles.HelpPaneStyle.Render(m.help.View(m.list)))
		listHeight -= helpHeight
	}
	return listHeight
}

func renderCreate(m Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.prompt.View())
}

func NewEmbedded(loadedProfiles []*model.Profile, width, height int) Model {
	return newModel(loadedProfiles, true, width, height)
}

func New(loadedProfiles []*model.Profile) Model {
	return newModel(loadedProfiles, false, 0, 0)
}

func newModel(loadedProfiles []*model.Profile, embedded bool, width, height int) Model {
	var profiles []list.Item

	for _, v := range loadedProfiles {
		r := Profile{
			Name:      v.Name,
			Variables: len(v.Variables), // TODO
		}
		profiles = append(profiles, r)
	}

	title := "Profiles"
	d := newNormalDelegate()
	if embedded {
		title = "Select Profile"
		d = newEmbeddedDelegate()
	}
	profileList := list.New(profiles, d, width, max(0, height-2))
	profileList.Title = title
	profileList.Styles.Title = titleStyle
	profileList.SetShowHelp(false)

	var sb statusbar.Model
	help := help.New()
	help.Styles.ShortKey = helpKeyStyle
	help.Styles.ShortDesc = helpDescStyle
	help.Styles.FullKey = helpKeyStyle
	help.Styles.FullDesc = helpDescStyle
	help.ShortSeparator = "  "

	if !embedded {
		statusbarItems := []statusbar.StatusbarItem{
			{Text: "", BackgroundColor: statusbarFirstColBg, ForegroundColor: statusbarFirstColFg},
			{Text: "? Help", BackgroundColor: statusbarSecondColBg, ForegroundColor: statusbarSecondColFg},
		}

		sb = statusbar.New(statusbarItems, 0, 0)
	}

	profileList.SetShowStatusBar(false)
	profileList.SetFilteringEnabled(false)

	mode := Normal
	if embedded {
		profileList.DisableQuitKeybindings()
		mode = Embedded
	}

	return Model{list: profileList, active: List, mode: mode, width: width, height: height, help: help, statusbar: sb}
}
