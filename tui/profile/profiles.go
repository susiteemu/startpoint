package profileui

import (
	"fmt"
	"goful/core/model"
	prompt "goful/tui/prompt"

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

func (i Profile) Title() string       { return i.Name }
func (i Profile) Description() string { return fmt.Sprintf("Vars: %d", i.Variables) }
func (i Profile) FilterValue() string { return i.Name }

type Model struct {
	list         list.Model
	prompt       prompt.Model
	embeddedHelp help.Model
	active       ActiveView
	mode         Mode
	width        int
	height       int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h := msg.Height
		if m.mode == Embedded {
			h = max(0, msg.Height-2)
		}
		m.list.SetSize(msg.Width, h)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "a":
			if m.mode == Normal && m.active != Create {
				m.active = Create
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
		helpView := m.embeddedHelp.View(embeddedKeys)

		views := []string{}
		views = append(views, m.list.View())
		views = append(views, helpView)

		return lipgloss.JoinVertical(lipgloss.Center, views...)

	}
	return m.list.View()
}

func renderCreate(m Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.prompt.View())
}

func NewEmbedded(loadedProfiles []model.Profile, width, height int) Model {
	return newModel(loadedProfiles, true, width, height)
}

func New(loadedProfiles []model.Profile) Model {
	return newModel(loadedProfiles, false, 0, 0)
}

func newModel(loadedProfiles []model.Profile, embedded bool, width, height int) Model {
	var profiles []list.Item

	for _, v := range loadedProfiles {
		r := Profile{
			Name:      v.Name,
			Variables: len(v.Variables), // TODO
		}
		profiles = append(profiles, r)
	}

	d := newNormalDelegate()
	if embedded {
		d = newEmbeddedDelegate()
	}
	profileList := list.New(profiles, d, width, max(0, height-2))
	profileList.Title = "Profiles"
	profileList.Styles.Title = titleStyle

	var embeddedHelp help.Model
	if !embedded {
		profileList.Help.Styles.FullKey = helpKeyStyle
		profileList.Help.Styles.FullDesc = helpDescStyle
		profileList.Help.Styles.ShortKey = helpKeyStyle
		profileList.Help.Styles.ShortDesc = helpDescStyle
		profileList.Help.Styles.ShortSeparator = helpSeparatorStyle
		profileList.Help.Styles.FullSeparator = helpSeparatorStyle
	} else {
		profileList.SetShowHelp(false)
		embeddedHelp = help.New()
		embeddedHelp.Styles.ShortKey = helpKeyStyle
		embeddedHelp.Styles.ShortDesc = helpDescStyle
		embeddedHelp.Styles.FullKey = helpKeyStyle
		embeddedHelp.Styles.FullDesc = helpDescStyle
		embeddedHelp.Styles.ShortSeparator = helpSeparatorStyle
		embeddedHelp.Styles.FullSeparator = helpSeparatorStyle
	}
	profileList.SetShowStatusBar(false)
	profileList.SetFilteringEnabled(false)

	mode := Normal
	if embedded {
		profileList.DisableQuitKeybindings()
		mode = Embedded
	}

	return Model{list: profileList, active: List, mode: mode, width: width, height: height, embeddedHelp: embeddedHelp}
}
