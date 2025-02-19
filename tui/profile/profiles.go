package profileui

import (
	"fmt"
	"os/exec"

	"github.com/susiteemu/startpoint/core/model"
	"github.com/susiteemu/startpoint/core/print"
	messages "github.com/susiteemu/startpoint/tui/messages"
	"github.com/susiteemu/startpoint/tui/overlay"
	preview "github.com/susiteemu/startpoint/tui/preview"
	prompt "github.com/susiteemu/startpoint/tui/prompt"
	statusbar "github.com/susiteemu/startpoint/tui/statusbar"
	"github.com/susiteemu/startpoint/tui/styles"

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
	CreateProfileLabel = "Choose a name for your profile (for default profile, leave the name blank). Make it filename compatible and unique within this workspace. After choosing \"ok\" your $EDITOR will open and you will be able to write the contents of the profile. Remember to quit your editor window to return back."
	RenameProfileLabel = "Rename your profile"
	CopyProfileLabel   = "Choose a name to your profile"
)

const (
	List ActiveView = iota
	Prompt
	Preview
)

type Mode int

const (
	Normal Mode = iota
	Embedded
)

type Profile struct {
	Name             string
	Variables        int
	PrivateVariables int
	OsEnvVariables   int
	ProfileModel     *model.Profile
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
		{k.Select, k.Cancel},
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
	list          list.Model
	prompt        prompt.Model
	statusbar     statusbar.Model
	preview       preview.Model
	embeddedHelp  help.Model
	active        ActiveView
	mode          Mode
	width         int
	height        int
	widthPercent  float64
	heightPercent float64
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	width := int(float64(m.width) * m.widthPercent)
	m.list.SetWidth(width)
	if m.mode == Normal {
		m.statusbar.SetWidth(width)
		listHeight := calculateListHeight(*m)
		m.list.SetHeight(listHeight)
		updateStatusbar(m, "")
	} else {
		capHeight := len(m.list.Items())*2 + 2
		height := min(capHeight, int(float64(m.height)*m.heightPercent))
		m.list.SetHeight(height)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.mode == Embedded {
			m.SetSize(msg.Width, msg.Height)
		}
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeyEsc.String():
			if m.active == Preview || m.active == Prompt {
				m.active = List
			}
			return m, nil
		}
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
				changeCmd := CreateChangeCmd()
				return m, tea.Batch(setCmd, statusCmd, changeCmd)
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
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Deleted %s", msg.Profile.Name))
				changeCmd := CreateChangeCmd()
				return m, tea.Batch(statusCmd, changeCmd)
			} else {
				return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Profile.Name))
			}
		}
	case EditProfileMsg:
		if m.mode == Normal && m.active == List {
			cmd, err := openFileToEditorCmd(msg.Profile.ProfileModel.Root, msg.Profile.ProfileModel.Filename)
			if err != nil {
				statusCmd := messages.CreateStatusMsg("Failed preparing editor")
				return m, statusCmd
			}
			cb := func(err error) tea.Msg {
				return EditProfileFinishedMsg{
					Profile: msg.Profile,
					err:     err,
				}
			}
			return m, tea.ExecProcess(cmd, cb)
		}
	case EditProfileFinishedMsg:
		oldProfile := msg.Profile
		if msg.err == nil {
			editedProfile, ok := readProfile(oldProfile.ProfileModel.Root, oldProfile.ProfileModel.Filename)
			if ok {
				setCmd := m.list.SetItem(m.list.Index(), editedProfile)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Edited profile %s", oldProfile.Title()))
				changeCmd := CreateChangeCmd()
				return m, tea.Batch(setCmd, statusCmd, changeCmd)
			}
		}
		return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to edit profile %s", oldProfile.Title()))
	case PreviewProfileMsg:
		log.Debug().Msgf("Preview Profile, %v", m)
		if m.active == List {
			m.active = Preview
			selected := msg.Profile.ProfileModel
			formatted, err := print.SprintDotenv(selected.Raw)

			if formatted == "" || err != nil {
				formatted = selected.Raw
			}

			m.preview = preview.New(selected.Filename, formatted, m.width, m.height, 0.8, 0.8)
			return m, nil
		}
	case prompt.PromptAnsweredMsg:
		m.active = List
		if msg.Context.Key == RenameProfile {
			profile := msg.Context.Additional.(Profile)
			renamedProfile, ok := renameProfile(msg.Input, profile)
			if ok {
				setCmd := m.list.SetItem(m.list.Index(), renamedProfile)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Renamed profile to %s", renamedProfile.Title()))
				changeCmd := CreateChangeCmd()
				return m, tea.Batch(setCmd, statusCmd, changeCmd)
			} else {
				return m, messages.CreateStatusMsg("Failed to rename profile")
			}
		} else if msg.Context.Key == CopyProfile {
			profile := msg.Context.Additional.(Profile)
			copiedProfile, ok := copyProfile(msg.Input, profile)
			if ok {
				setCmd := m.list.InsertItem(m.list.Index()+1, copiedProfile)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Copied profile to %s", copiedProfile.Title()))
				changeCmd := CreateChangeCmd()
				return m, tea.Batch(setCmd, statusCmd, changeCmd)
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
	case Preview:
		m.preview, cmd = m.preview.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Prompt:
		return renderPrompt(m)
	case Preview:
		return renderPreview(m)
	default:
		return renderList(m)
	}
}

func (m *Model) GetHelpKeys() help.KeyMap {
	if m.mode == Embedded {
		return embeddedKeyMap{}
	}
	return m.list
}
func renderList(m Model) string {
	if m.mode == Embedded {
		m.embeddedHelp.Width = m.list.Width()
		helpView := lipgloss.NewStyle().Render(m.embeddedHelp.View(embeddedKeys))

		views := []string{}
		views = append(views, m.list.View())
		views = append(views, helpView)

		return lipgloss.JoinVertical(lipgloss.Center, views...)
	}
	var views []string
	listHeight := calculateListHeight(m)
	views = append(views, lipgloss.NewStyle().Height(listHeight).Padding(1, 0, 0, 0).Render(m.list.View()))
	views = append(views, m.statusbar.View())

	joined := lipgloss.JoinVertical(
		lipgloss.Top,
		views...,
	)

	return joined
}

func calculateListHeight(m Model) int {
	listHeight := m.height - statusbar.Height - 1 // -1 for top padding
	return listHeight
}

func renderPreview(m Model) string {
	w := m.width
	h := m.height
	return renderModalAtCenter(renderList(m), m.preview.View(), w, h)
}

func renderPrompt(m Model) string {
	w := m.width
	h := m.height
	return renderModalAtCenter(renderList(m), m.prompt.View(), w, h)
}

func renderModalAtCenter(bg string, modal string, w, h int) string {
	x := (w / 2) - (lipgloss.Width(modal) / 2)
	y := (h / 2) - (lipgloss.Height(modal) / 2)
	return overlay.PlaceOverlay(x, y, modal, bg)
}

func NewEmbedded(loadedProfiles []*model.Profile, winWidth, winHeight int, wPercent, hPercent float64) Model {
	return newModel(loadedProfiles, true, winWidth, winHeight, wPercent, hPercent)
}

func New(loadedProfiles []*model.Profile) Model {
	return newModel(loadedProfiles, false, 0, 0, 1.0, 1.0)
}

func newModel(loadedProfiles []*model.Profile, embedded bool, winWidth, winHeight int, wPercent, hPercent float64) Model {
	var profiles []list.Item

	theme := styles.LoadTheme()
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

	d := newNormalDelegate()
	if embedded {
		d = newEmbeddedDelegate()
	}

	width := int(float64(winWidth) * wPercent)
	height := int(float64(winHeight) * hPercent)
	if embedded {
		capHeight := len(loadedProfiles)*2 + 2
		height = min(capHeight, int(float64(winHeight)*hPercent))
	}

	profileList := list.New(profiles, d, width, height)
	profileList.SetShowHelp(false)
	profileList.SetShowTitle(false)

	var sb statusbar.Model

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
	var embeddedHelp help.Model
	if embedded {
		profileList.DisableQuitKeybindings()
		mode = Embedded
		embeddedHelp = help.New()
		embeddedHelp.Styles.ShortKey = style.helpKeyStyle
		embeddedHelp.Styles.ShortDesc = style.helpDescStyle
		embeddedHelp.Styles.FullKey = style.helpKeyStyle
		embeddedHelp.Styles.FullDesc = style.helpDescStyle
		embeddedHelp.ShortSeparator = "  "
	}

	return Model{
		list:          profileList,
		active:        List,
		embeddedHelp:  embeddedHelp,
		mode:          mode,
		width:         winWidth,
		height:        winHeight,
		statusbar:     sb,
		widthPercent:  wPercent,
		heightPercent: hPercent,
	}
}
