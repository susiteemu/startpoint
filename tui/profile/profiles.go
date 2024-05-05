package profileui

import (
	"fmt"
	"os/exec"
	"startpoint/core/model"
	messages "startpoint/tui/messages"
	prompt "startpoint/tui/prompt"
	statusbar "startpoint/tui/statusbar"
	"startpoint/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type ActiveView int

const (
	CreateProfile = "CreateProfile"
	EditProfile   = "EditProfile"
	RenameProfile = "RenameProfile"
	DeleteProfile = "DeleteProfile"
	CopyProfile   = "CopyProfile"
)

const (
	CreateProfileLabel = "Choose a name for your profile. Make it filename compatible and unique within this workspace. After choosing \"ok\" your $EDITOR will open and you will be able to write the contents of the profile. Remember to quit your editor window to return back."
	RenameProfileLabel = "Rename your profile"
	CopyProfileLabel   = "Choose a name to your profile"
)

const (
	List ActiveView = iota
	Prompt
	Update
)

type Mode int

const (
	Normal Mode = iota
	Embedded
)

type Profile struct {
	Name         string
	Variables    int
	ProfileModel *model.Profile
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
		Text: msg, BackgroundColor: style.statusbarFirstColBg, ForegroundColor: style.statusbarFirstColFg,
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
		case tea.KeyEsc.String():
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
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
	case CreateProfileMsg:
		if m.mode == Normal && m.active == List {
			log.Debug().Msg("Creating profile")
			promptKey := CreateProfile
			promptLabel := CreateProfileLabel
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key: promptKey,
			}, "", promptLabel, checkProfileWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case CreateProfileFinishedMsg:
		if msg.err == nil {
			newProfile, ok := readProfile(msg.root, msg.filename)
			if ok {
				setCmd := m.list.InsertItem(m.list.Index()+1, newProfile)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Created profile %s", newProfile.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			}
		}
		return m, messages.CreateStatusMsg("Failed to create profile")
	case RenameProfileMsg:
		if m.mode == Normal && m.active == List {
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        RenameProfile,
				Additional: msg.Profile,
			}, msg.Profile.Name, RenameProfileLabel, checkProfileWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case CopyProfileMsg:
		if m.mode == Normal && m.active == List {
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        CopyProfile,
				Additional: msg.Profile,
			}, msg.Profile.Name, CopyProfileLabel, checkProfileWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case DeleteProfileMsg:
		if m.mode == Normal && m.active == List {
			deleted := msg.Profile.ProfileModel.DeleteFromFS()
			if deleted {
				index := m.list.Index()
				m.list.RemoveItem(index)
				return m, messages.CreateStatusMsg(fmt.Sprintf("Deleted %s", msg.Profile.Name))
			} else {
				return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Profile.Name))
			}
		}
	case prompt.PromptAnsweredMsg:
		m.active = List
		if msg.Context.Key == RenameProfile {
			profile := msg.Context.Additional.(Profile)
			renamedProfile, ok := renameProfile(msg.Input, profile)
			if ok {
				setCmd := m.list.SetItem(m.list.Index(), renamedProfile)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Renamed profile to %s", renamedProfile.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, messages.CreateStatusMsg("Failed to rename profile")
			}
		} else if msg.Context.Key == CopyProfile {
			profile := msg.Context.Additional.(Profile)
			copiedProfile, ok := copyProfile(msg.Input, profile)
			if ok {
				setCmd := m.list.InsertItem(m.list.Index()+1, copiedProfile)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Copied profile to %s", copiedProfile.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, messages.CreateStatusMsg("Failed to copy profile")
			}

		} else if msg.Context.Key == CreateProfile {
			var (
				root     string
				filepath string
				cmd      *exec.Cmd
				err      error
			)
			root, filepath, cmd, err = createProfileFileCmd(msg.Input)
			if err != nil {
				return m, messages.CreateStatusMsg("Failed preparing editor")
			}
			cb := func(err error) tea.Msg {
				return CreateProfileFinishedMsg{
					root:     root,
					filename: filepath,
					err:      err,
				}
			}
			return m, tea.ExecProcess(cmd, cb)
		}
	case messages.StatusMessage:
		updateStatusbar(&m, string(msg))
		return m, nil
	}
	var cmd tea.Cmd
	switch m.active {
	case List:
		m.list, cmd = m.list.Update(msg)
	case Prompt:
		m.prompt, cmd = m.prompt.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Prompt:
		return renderPrompt(m)
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
		views = append(views, style.helpPaneStyle.Render(m.help.View(m.list)))
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
		helpHeight := lipgloss.Height(style.helpPaneStyle.Render(m.help.View(m.list)))
		listHeight -= helpHeight
	}
	return listHeight
}

func renderPrompt(m Model) string {
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

	theme := styles.GetTheme()
	commonStyles := styles.GetCommonStyles(theme)
	InitStyle(theme, commonStyles)

	for _, v := range loadedProfiles {
		r := Profile{
			Name:         v.Name,
			Variables:    len(v.Variables),
			ProfileModel: v,
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
	profileList.Styles.Title = style.listTitleStyle
	profileList.SetShowHelp(false)

	var sb statusbar.Model
	help := help.New()
	help.Styles.ShortKey = style.helpKeyStyle
	help.Styles.ShortDesc = style.helpDescStyle
	help.Styles.FullKey = style.helpKeyStyle
	help.Styles.FullDesc = style.helpDescStyle
	help.ShortSeparator = "  "

	if !embedded {
		statusbarItems := []statusbar.StatusbarItem{
			{Text: "", BackgroundColor: style.statusbarFirstColBg, ForegroundColor: style.statusbarFirstColFg},
			{Text: "? Help", BackgroundColor: style.statusbarSecondColBg, ForegroundColor: style.statusbarSecondColFg},
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
