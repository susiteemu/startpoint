package requestui

import (
	"errors"
	"fmt"
	"github.com/susiteemu/startpoint/core/model"
	"github.com/susiteemu/startpoint/core/print"
	"os/exec"
	"slices"
	"strings"
	"time"

	keyprompt "github.com/susiteemu/startpoint/tui/keyprompt"
	messages "github.com/susiteemu/startpoint/tui/messages"
	"github.com/susiteemu/startpoint/tui/overlay"
	preview "github.com/susiteemu/startpoint/tui/preview"
	profiles "github.com/susiteemu/startpoint/tui/profile"
	prompt "github.com/susiteemu/startpoint/tui/prompt"
	statusbar "github.com/susiteemu/startpoint/tui/statusbar"
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type PostAction struct {
	Type             string
	Payload          interface{}
	AddtionalContext interface{}
}

// hackish solution for bubbletea not supporting passing our own model into list rendering functions
var (
	activeProfile            *model.Profile
	allProfiles              []*model.Profile
	processTemplateVariables bool
)

func modeStr(mode Mode) string {
	switch mode {
	case Select:
		return "SELECT"
	case Edit:
		return "EDIT"
	}
	return ""
}

func indexOfByName(name string, m Model) int {
	index := -1
	for i, item := range m.list.Items() {
		itemRequest := item.(Request)
		if itemRequest.Name == name {
			index = i
			break
		}
	}
	return index
}

func indexOfNew(name string, m Model) int {
	index := len(m.list.Items()) + 1
	for i, item := range m.list.Items() {
		itemRequest := item.(Request)
		if strings.Compare(itemRequest.Name, name) > 0 {
			index = i
			break
		}
	}
	return index
}

func findRequestMold(r Request, m Model) (*model.RequestMold, error) {
	var requestMold *model.RequestMold
	for _, m := range m.requestMolds {
		if m.Name == r.Name {
			requestMold = m
			break
		}
	}
	if requestMold == nil {
		return nil, errors.New(fmt.Sprintf("could not find corresponding request mold for %s", r.Name))
	}
	return requestMold, nil
}

func updateStatusbar(m *Model, msg string) {

	var modeBg lipgloss.Color
	var profileBg lipgloss.Color
	var profileText string
	if m.mode == Edit {
		modeBg = style.statusbarModeEditBg
		profileBg = style.statusbarPrimaryBg
	} else {
		if activeProfile != nil {
			profileText = activeProfile.Name
		}
		if profileText == "" {
			profileText = "<no profile>"
		}
		modeBg = style.statusbarModeSelectBg
		profileBg = style.statusbarThirdColBg
	}

	modeItem := statusbar.StatusbarItem{
		Text: modeStr(m.mode), BackgroundColor: modeBg, ForegroundColor: style.statusbarSecondaryFg,
	}

	profileItem := statusbar.StatusbarItem{
		Text: profileText, BackgroundColor: profileBg, ForegroundColor: style.statusbarSecondaryFg,
	}

	m.statusbar.SetItem(modeItem, 0)
	m.statusbar.ChangeText(msg, 1)
	m.statusbar.SetItem(profileItem, 2)
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.statusbar.SetWidth(w)
	m.topbar.SetWidth(w)
	m.list.SetWidth(w)
	m.help.Width = w
	listHeight := calculateListHeight(*m)
	m.list.SetHeight(listHeight)
	updateStatusbar(m, "")
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if we are filtering, it gets all the input
		if m.active == List && m.list.SettingFilter() {
			if msg.String() == tea.KeyEsc.String() {
				// hide help after cancelling filtering
				m.help.ShowAll = false
				listHeight := calculateListHeight(m)
				m.list.SetHeight(listHeight)
			}
			break
		}

		switch keypress := msg.String(); keypress {
		case tea.KeyEsc.String():
			if m.active == Preview || m.active == Prompt || m.active == Profiles || m.active == Keyprompt {
				m.active = List
				return m, nil
			}
			if m.active == List && !m.list.IsFiltered() {
				if m.mode == Edit && m.active == List {
					m.mode = Select
					processTemplateVariables = true
					m.list.SetDelegate(newSelectDelegate())
					listHeight := calculateListHeight(m)
					m.list.SetHeight(listHeight)
					updateStatusbar(&m, "")
					return m, nil
				}
				return m, nil
			}
		case "i":
			if m.mode == Select && m.active == List {
				m.mode = Edit
				processTemplateVariables = false
				m.list.SetDelegate(newEditModeDelegate())
				listHeight := calculateListHeight(m)
				m.list.SetHeight(listHeight)
				updateStatusbar(&m, "")
				return m, nil
			}
		case "?":
			if m.active == List {
				m.help.ShowAll = !m.help.ShowAll
				listHeight := calculateListHeight(m)
				m.list.SetHeight(listHeight)
				return m, nil
			}
		case "/":
			if m.active == List {
				// show help when filtering:
				// filter gets all key input so pressing ? would not work
				// NOTE: important thing is not to return anything here
				// this way the message will propagate to list bubble
				m.help.ShowAll = true
				listHeight := calculateListHeight(m)
				m.list.SetHeight(listHeight)
			}
		}

	case ShowKeyprompt:
		if m.mode == Edit && m.active == List {
			m.active = Keyprompt
			m.keyprompt = keyprompt.New(msg.Label, msg.Entries, msg.Type, msg.Payload, m.width)
		}

	case RunRequestMsg:
		request := msg.Request
		m.active = Stopwatch
		m.list.ResetFilter()
		requestMold, err := findRequestMold(request, m)
		if err != nil {
			log.Error().Msgf("Could not find request mold with request %v", request)
			return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to run request %s", request.Title()))
		}
		return m, tea.Batch(
			m.stopwatch.Init(),
			doRequest(requestMold, m.requestMolds, activeProfile),
		)
	case RunRequestFinishedMsg:
		m.postAction = PostAction{
			Type:    PrintRequest,
			Payload: string(msg),
		}
		return m, tea.Quit

	case RunRequestFinishedWithFailureMsg:
		m.postAction = PostAction{
			Type:    PrintFailedRequest,
			Payload: string(msg),
		}
		return m, tea.Quit
	case CreateRequestMsg:
		if m.mode == Edit && m.active == List {
			var promptKey, promptLabel string
			switch msg.Type {
			case model.CONTENT_TYPE_YAML:
				promptKey = CreateSimpleRequest
				promptLabel = CreateRequestLabel
			case model.CONTENT_TYPE_STARLARK, model.CONTENT_TYPE_LUA:
				promptKey = CreateComplexRequest
				promptLabel = CreateRequestLabel
			}
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        promptKey,
				Additional: msg.Type,
			}, "", promptLabel, checkRequestWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case CreateRequestFinishedMsg:
		if msg.err == nil {
			newRequest, newRequestMold, ok := readRequest(msg.root, msg.filename)
			if ok {
				m.list.ResetFilter()
				m.help.ShowAll = false
				listHeight := calculateListHeight(m)
				m.list.SetHeight(listHeight)
				index := indexOfNew(newRequest.Name, m)
				setCmd := m.list.InsertItem(index, newRequest)
				m.list.Select(index)
				// NOTE: order is not relevant here
				m.requestMolds = append(m.requestMolds, newRequestMold)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Created request %s", newRequest.Title()))
				return m, tea.Sequence(setCmd, statusCmd)
			}
		}
		return m, messages.CreateStatusMsg("Failed to create request")
	case EditRequestMsg:
		if m.mode == Edit && m.active == List {
			log.Debug().Msgf("Starting to edit request %v", msg.Request)
			requestMold, err := findRequestMold(msg.Request, m)
			if err != nil {
				log.Error().Msgf("Could not find request mold with request %v", msg.Request)
				return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to edit request %s", msg.Request.Title()))
			}
			cmd, err := openFileToEditorCmd(requestMold.Root, requestMold.Filename)
			if err != nil {
				statusCmd := messages.CreateStatusMsg("Failed preparing editor")
				return m, statusCmd
			}
			cb := func(err error) tea.Msg {
				return EditRequestFinishedMsg{
					Request: msg.Request,
					err:     err,
				}
			}
			return m, tea.ExecProcess(cmd, cb)
		}
	case DeleteRequestMsg:
		if m.mode == Edit && m.active == List {
			requestMold, err := findRequestMold(msg.Request, m)
			if err != nil {
				log.Error().Err(err).Msgf("Could not find request mold with request %v", msg.Request)
				return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Request.Title()))
			}

			label := fmt.Sprintf("You are about to delete '%s'.\n\n", requestMold.Name)
			label += "Are you sure you want to delete this request?"
			if isUsedAsPrevReq(requestMold.Name, m.requestMolds) {
				label += " It is used by other requests."
			}

			var keys []keyprompt.KeypromptEntry
			keys = append(keys, keyprompt.KeypromptEntry{
				Text: "yes", Key: "y",
			})
			keys = append(keys, keyprompt.KeypromptEntry{
				Text: "no", Key: "n",
			})
			return m, tea.Cmd(func() tea.Msg {
				return ShowKeyprompt{
					Label:   label,
					Entries: keys,
					Type:    DeleteRequest,
					Payload: msg.Request,
				}
			})

		}
	case DeleteRequestConfirmedMsg:
		requestMold, err := findRequestMold(msg.Request, m)
		if err != nil {
			log.Error().Err(err).Msgf("Could not find request mold with request %v", msg.Request)
			return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Request.Title()))
		}
		deleted := requestMold.DeleteFromFS()
		if deleted {
			// NOTE: use index of selected item with delete: for some reason it works whereas with edit/rename it does not
			index := m.list.Index()
			log.Debug().Msgf("Index of request to delete is %d", index)
			m.list.RemoveItem(index)
			removeIndex := slices.Index(m.requestMolds, requestMold)
			m.requestMolds = slices.Delete(m.requestMolds, removeIndex, removeIndex+1)
			return m, messages.CreateStatusMsg(fmt.Sprintf("Deleted %s", msg.Request.Title()))
		} else {
			return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Request.Title()))
		}
	case EditRequestFinishedMsg:
		oldRequest := msg.Request
		oldName := msg.Request.Name
		if msg.err == nil {
			log.Debug().Msgf("Finishing editing request %v", oldRequest)
			requestMold, err := findRequestMold(oldRequest, m)
			if err != nil {
				log.Error().Err(err).Msgf("Could not find request mold with request %v", msg.Request)
				return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to edit request %s", oldRequest.Title()))
			}
			log.Debug().Msgf("Found corresponding request mold %v", requestMold)
			editedRequest, editedRequestMold, ok := readRequest(requestMold.Root, requestMold.Filename)
			if ok {
				index := indexOfByName(oldName, m)
				log.Debug().Msgf("Index of request is %v", index)
				setCmd := m.list.SetItem(index, editedRequest)
				index = slices.Index(m.requestMolds, requestMold)
				m.requestMolds = slices.Replace(m.requestMolds, index, index+1, editedRequestMold)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Edited request %s", oldRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			}
		}
		return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to edit request %s", oldRequest.Title()))
	case PreviewRequestMsg:
		if m.active == List {
			m.active = Preview
			selected, err := findRequestMold(msg.Request, m)
			if err != nil {
				log.Error().Err(err).Msgf("Could not find request mold with request %v", msg.Request)
				return m, messages.CreateStatusMsg(fmt.Sprintf("Failed to open %s for preview", msg.Request.Title()))
			}
			var formatted string
			formatted, err = print.SprintRequestMold(selected)
			if formatted == "" || err != nil {
				formatted = selected.Raw()
			}
			w := int(float64(m.width) * 0.8)
			h := int(float64(m.height) * 0.8)
			m.preview = preview.New(selected.Filename, formatted, w, h)

			return m, nil
		}
	case RenameRequestMsg:
		if m.mode == Edit && m.active == List {
			log.Debug().Msgf("Starting to rename request %v", msg.Request.Name)
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        RenameRequest,
				Additional: msg.Request,
			}, msg.Request.Name, RenameRequestLabel, checkRequestWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case CopyRequestMsg:
		if m.mode == Edit && m.active == List {
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        CopyRequest,
				Additional: msg.Request,
			}, fmt.Sprintf("Copy of %s", msg.Request.Name), CopyRequestLabel, checkRequestWithNameDoesNotExist(m), m.width)
		}
	case prompt.PromptAnsweredMsg:
		m.active = List
		if msg.Context.Key == RenameRequest {
			request := msg.Context.Additional.(Request)
			requestMold, err := findRequestMold(request, m)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to find request mold with %v", request)
				return m, messages.CreateStatusMsg("Failed to rename request")
			}
			oldName := requestMold.Name
			moldIndex := slices.Index(m.requestMolds, requestMold)
			renamedRequest, renamedRequestMold, ok := renameRequest(msg.Input, request, *requestMold)
			if ok {
				m.requestMolds = slices.Replace(m.requestMolds, moldIndex, moldIndex+1, renamedRequestMold)
				_, ok = changePrevReq(oldName, renamedRequestMold.Name, m.requestMolds)
				if !ok {
					return m, messages.CreateStatusMsg("Failed to refactor dependent requests")
				}

				log.Debug().Msgf("Length of items %d", len(m.list.Items()))
				index := indexOfByName(oldName, m)
				log.Debug().Msgf("Going to replace item in index %d", index)
				setCmd := m.list.SetItem(index, renamedRequest)

				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Renamed request to %s", renamedRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, messages.CreateStatusMsg("Failed to rename request")
			}
		} else if msg.Context.Key == CopyRequest {
			request := msg.Context.Additional.(Request)
			requestMold, err := findRequestMold(request, m)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to find request mold with %v", request)
				return m, messages.CreateStatusMsg("Failed to copy request")
			}
			copiedRequest, copiedRequestMold, ok := copyRequest(msg.Input, request, *requestMold)
			if ok {
				// note order is not relevant here
				m.requestMolds = append(m.requestMolds, copiedRequestMold)
				index := indexOfByName(request.Name, m)
				setCmd := m.list.InsertItem(index+1, copiedRequest)
				statusCmd := messages.CreateStatusMsg(fmt.Sprintf("Copied request to %s", copiedRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, messages.CreateStatusMsg("Failed to copy request")
			}

		} else if msg.Context.Key == CreateSimpleRequest || msg.Context.Key == CreateComplexRequest {
			var (
				root     string
				filepath string
				cmd      *exec.Cmd
				err      error
			)
			requestType := msg.Context.Additional.(string)
			root, filepath, cmd, err = createRequestFileCmd(msg.Input, requestType)
			if err != nil {
				return m, messages.CreateStatusMsg("Failed preparing editor")
			}
			cb := func(err error) tea.Msg {
				return CreateRequestFinishedMsg{
					root:     root,
					filename: filepath,
					err:      err,
				}
			}
			return m, tea.ExecProcess(cmd, cb)
		}
	case ActivateProfile:
		if m.mode == Select && m.active == List {
			m.active = Profiles
			m.profileui = profiles.NewEmbedded(allProfiles, m.width, m.height, 0.5, 0.8)
		}

	case profiles.ProfileSelectedMsg:
		m.active = List
		log.Debug().Msgf("Selected profile %v", msg.Profile.Name)
		var activedProfile *model.Profile
		for _, p := range allProfiles {
			if p.Name == msg.Profile.Name {
				activedProfile = p
				break
			}
		}
		log.Debug().Msgf("Matched with profile %v", activedProfile)
		if activedProfile == nil {
			return m, messages.CreateStatusMsg("Failed to set profile")
		}

		activeProfile = activedProfile

		updateStatusbar(&m, "")
		return m, nil

	case keyprompt.KeypromptAnsweredMsg:
		m.active = List
		if msg.Type == CreateRequest {
			switch msg.Key {
			case "y":
				return m, tea.Cmd(func() tea.Msg {
					return CreateRequestMsg{
						Type: model.CONTENT_TYPE_YAML,
					}
				})
			case "s":
				return m, tea.Cmd(func() tea.Msg {
					return CreateRequestMsg{
						Type: model.CONTENT_TYPE_STARLARK,
					}
				})
			case "l":
				return m, tea.Cmd(func() tea.Msg {
					return CreateRequestMsg{
						Type: model.CONTENT_TYPE_LUA,
					}
				})
			}
		} else if msg.Type == DeleteRequest {
			if msg.Key == "y" {
				return m, tea.Cmd(func() tea.Msg {
					return DeleteRequestConfirmedMsg{
						Request: msg.Payload.(Request),
					}
				})
			}

		}

		return m, nil

	case messages.StatusMessage:
		updateStatusbar(&m, string(msg))
		return m, nil
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.active {
	case List:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case Prompt:
		m.prompt, cmd = m.prompt.Update(msg)
		cmds = append(cmds, cmd)
	case Keyprompt:
		m.keyprompt, cmd = m.keyprompt.Update(msg)
		cmds = append(cmds, cmd)
	case Preview:
		m.preview, cmd = m.preview.Update(msg)
		cmds = append(cmds, cmd)
	case Profiles:
		m.profileui, cmd = m.profileui.Update(msg)
		cmds = append(cmds, cmd)
	case Stopwatch:
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Prompt:
		return renderPrompt(m)
	case Keyprompt:
		return renderKeyprompt(m)
	case Preview:
		return renderPreview(m)
	case Stopwatch:
		return renderStopwatch(m)
	case Profiles:
		return renderProfiles(m)
	default:
		return renderList(m)
	}
}

func renderList(m Model) string {
	var views []string
	listHeight := calculateListHeight(m)
	views = append(views, lipgloss.NewStyle().Height(listHeight).Render(m.list.View()))
	views = append(views, m.statusbar.View())

	joined := lipgloss.JoinVertical(
		lipgloss.Top,
		views...,
	)

	if m.help.ShowAll {
		helpModal := style.helpPaneStyle.Render(m.help.View(m.list))
		// position at the bottom
		x := (m.width / 2) - (lipgloss.Width(helpModal) / 2)
		y := m.height - lipgloss.Height(helpModal) - 1
		joined = overlay.PlaceOverlay(x, y, helpModal, joined)
	}

	return joined
}

func calculateListHeight(m Model) int {
	listHeight := m.height - statusbar.Height
	return listHeight
}

func renderModal(bg string, modal string, w, h int) string {
	x := (w / 2) - (lipgloss.Width(modal) / 2)
	y := (h / 2) - (lipgloss.Height(modal) / 2)
	return overlay.PlaceOverlay(x, y, modal, bg)
}

func renderStopwatch(m Model) string {
	w := m.width
	h := m.height
	return renderModal(renderList(m), style.stopwatchStyle.Render("Running request\n\n"+m.stopwatch.View()), w, h)
}

func renderPreview(m Model) string {
	w := m.width
	h := m.height
	return renderModal(renderList(m), m.preview.View(), w, h)
}

func renderProfiles(m Model) string {
	w := m.width
	h := m.height
	log.Debug().Msgf("ProfileUi width: %d", lipgloss.Width(m.profileui.View()))
	return renderModal(renderList(m), style.profilesStyle.Render(m.profileui.View()), w, h)
}

func renderPrompt(m Model) string {
	w := m.width
	h := m.height
	return renderModal(renderList(m), m.prompt.View(), w, h)
}

func renderKeyprompt(m Model) string {
	w := m.width
	h := m.height
	return renderModal(renderList(m), m.keyprompt.View(), w, h)
}

func New(loadedRequests []*model.RequestMold, loadedProfiles []*model.Profile) Model {
	log.Info().Msgf("Starting up manage TUI with %d loaded requests and %d profiles", len(loadedRequests), len(loadedProfiles))

	theme := styles.GetTheme()
	InitStyle(theme, styles.GetCommonStyles(theme))

	var requests []list.Item

	for _, v := range loadedRequests {
		r := Request{
			Name:   v.Name,
			Url:    v.Url(),
			Method: v.Method(),
		}
		requests = append(requests, r)
	}

	RefreshProfiles(loadedProfiles)

	mode := Select
	if len(requests) == 0 {
		mode = Edit
	}

	var d list.DefaultDelegate
	var modeColor lipgloss.Color

	switch mode {
	case Select:
		d = newSelectDelegate()
		modeColor = style.statusbarModeSelectBg
		processTemplateVariables = true
	case Edit:
		d = newEditModeDelegate()
		modeColor = style.statusbarModeEditBg
		processTemplateVariables = false
	}

	requestList := list.New(requests, d, 0, 0)
	requestList.Styles.Title = style.listTitleStyle
	requestList.Styles.StatusBar = lipgloss.NewStyle().Foreground(style.listStatusbarFg).Padding(0, 1, 1, 1)
	requestList.Styles.StatusBarFilterCount = requestList.Styles.StatusBar.Copy().UnsetPadding().Faint(true)
	requestList.Styles.StatusEmpty = requestList.Styles.StatusBar.Copy().UnsetPadding()
	requestList.Styles.NoItems = requestList.Styles.StatusBar.Copy()
	requestList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(style.listFilterPromptFg).Padding(1, 0, 0, 0)
	requestList.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(style.listFilterCursorFg)

	// NOTE: removing few default keybindings so that pressing our own keys (e.g. 'd') would not have any side-effects
	requestList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("left", "h", "pgup", "b"),
		key.WithHelp("←/h/pgup", "prev page"),
	)
	requestList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("right", "l", "pgdown", "f"),
		key.WithHelp("→/l/pgdn", "next page"),
	)

	requestList.SetShowStatusBar(false)
	requestList.SetShowTitle(false)
	requestList.SetShowHelp(false)

	statusbarItems := []statusbar.StatusbarItem{
		{Text: modeStr(mode), BackgroundColor: modeColor, ForegroundColor: style.statusbarSecondaryFg},
		{Text: "", BackgroundColor: style.statusbarPrimaryBg, ForegroundColor: style.statusbarPrimaryFg},
		{Text: "", BackgroundColor: style.statusbarThirdColBg, ForegroundColor: style.statusbarSecondaryFg},
		{Text: "? Help", BackgroundColor: style.statusbarFourthColBg, ForegroundColor: style.statusbarSecondaryFg},
	}

	help := help.New()
	help.Styles.ShortKey = style.helpKeyStyle
	help.Styles.ShortDesc = style.helpDescStyle
	help.Styles.FullKey = style.helpKeyStyle
	help.Styles.FullDesc = style.helpDescStyle

	sb := statusbar.New(statusbarItems, 1, 0)
	return Model{
		list:         requestList,
		active:       List,
		mode:         mode,
		stopwatch:    stopwatch.NewWithInterval(time.Millisecond * 100),
		statusbar:    sb,
		help:         help,
		requestMolds: loadedRequests,
	}

}
