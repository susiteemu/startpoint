package requestui

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"startpoint/core/loader"
	"startpoint/core/model"
	"startpoint/core/print"
	"time"

	"os"
	keyprompt "startpoint/tui/keyprompt"
	preview "startpoint/tui/preview"
	profiles "startpoint/tui/profile"
	prompt "startpoint/tui/prompt"
	statusbar "startpoint/tui/statusbar"
	"startpoint/tui/styles"

	"github.com/charmbracelet/bubbles/help"
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

func findRequestMold(r Request, m Model) (*model.RequestMold, error) {
	var requestMold *model.RequestMold
	for _, m := range m.requestMolds {
		if m.Name() == r.Name {
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

	msgItem := statusbar.StatusbarItem{
		Text: msg, BackgroundColor: style.statusbarPrimaryBg, ForegroundColor: style.statusbarPrimaryFg,
	}

	profileItem := statusbar.StatusbarItem{
		Text: profileText, BackgroundColor: profileBg, ForegroundColor: style.statusbarSecondaryFg,
	}

	m.statusbar.SetItem(modeItem, 0)
	m.statusbar.SetItem(msgItem, 1)
	m.statusbar.SetItem(profileItem, 2)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusbar.SetWidth(msg.Width)
		m.list.SetWidth(msg.Width)
		m.help.Width = msg.Width
		listHeight := calculateListHeight(m)
		m.list.SetHeight(listHeight)
		updateStatusbar(&m, "")
	case tea.KeyMsg:
		// if we are filtering, it gets all the input
		if m.active == List && m.list.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case "q":
			if m.active == Preview {
				m.active = List
				return m, nil
			}
			if m.active == List {
				return m, tea.Quit
			}
		case tea.KeyEsc.String():
			if m.active == Preview || m.active == Prompt || m.active == Profiles || m.active == Keyprompt {
				m.active = List
				return m, nil
			}
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
		}

	case ShowKeyprompt:
		if m.mode == Edit && m.active == List {
			m.active = Keyprompt
			m.keyprompt = keyprompt.New(msg.Label, msg.Entries)
		}

	case RunRequestMsg:
		request := msg.Request
		m.active = Stopwatch
		m.list.ResetFilter()
		requestMold, err := findRequestMold(request, m)
		if err != nil {
			// TODO show error if request mold is not found
			return m, createStatusMsg(fmt.Sprintf("Failed to run request %s", request.Title()))
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

	case CreateRequestMsg:
		if m.mode == Edit && m.active == List {
			promptKey := CreateSimpleRequest
			promptLabel := CreateRequestLabel
			if !msg.Simple {
				promptKey = CreateComplexRequest
				promptLabel = CreateRequestLabel
			}
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key: promptKey,
			}, "", promptLabel, checkRequestWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case CreateRequestFinishedMsg:
		if msg.err == nil {
			newRequest, newRequestMold, ok := readRequest(msg.root, msg.filename)
			if ok {
				setCmd := m.list.InsertItem(m.list.Index()+1, newRequest)
				// note: order is not relevant here
				m.requestMolds = append(m.requestMolds, newRequestMold)
				statusCmd := createStatusMsg(fmt.Sprintf("Created request %s", newRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			}
		}
		return m, createStatusMsg("Failed to create request")
	case EditRequestMsg:
		if m.mode == Edit && m.active == List {
			requestMold, err := findRequestMold(msg.Request, m)
			if err != nil {
				// TODO handle err
				return m, createStatusMsg(fmt.Sprintf("Failed to edit request %s", msg.Request.Title()))
			}
			cmd, err := openFileToEditorCmd(requestMold.Root, requestMold.Filename)
			if err != nil {
				statusCmd := createStatusMsg("Failed preparing editor")
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
				log.Error().Err(err).Msgf("Failed to find request mold with %v", msg.Request)
				return m, createStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Request.Title()))
			}
			deleted := requestMold.DeleteFromFS()
			if deleted {
				index := m.list.Index()
				m.list.RemoveItem(index)
				removeIndex := slices.Index(m.requestMolds, requestMold)
				m.requestMolds = slices.Delete(m.requestMolds, removeIndex, removeIndex+1)
				return m, createStatusMsg(fmt.Sprintf("Deleted %s", msg.Request.Title()))
			} else {
				return m, createStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Request.Title()))
			}
		}
	case EditRequestFinishedMsg:
		oldRequest := msg.Request
		if msg.err == nil {
			requestMold, err := findRequestMold(oldRequest, m)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to find request mold with %v", msg.Request)
				return m, createStatusMsg(fmt.Sprintf("Failed to edit request %s", oldRequest.Title()))
			}
			editedRequest, editedRequestMold, ok := readRequest(requestMold.Root, requestMold.Filename)
			if ok {
				setCmd := m.list.SetItem(m.list.Index(), editedRequest)
				index := slices.Index(m.requestMolds, requestMold)
				m.requestMolds = slices.Replace(m.requestMolds, index, index+1, editedRequestMold)
				statusCmd := createStatusMsg(fmt.Sprintf("Edited request %s", oldRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			}
		}
		return m, createStatusMsg(fmt.Sprintf("Failed to edit request %s", oldRequest.Title()))
	case PreviewRequestMsg:
		if m.active == List {
			m.active = Preview
			selected, err := findRequestMold(msg.Request, m)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to find request mold with %v", msg.Request)
				return m, createStatusMsg(fmt.Sprintf("Failed to open %s for preview", msg.Request.Title()))
			}
			var formatted string
			switch selected.ContentType {
			case "yaml":
				formatted, err = print.SprintYaml(selected.Raw())
			case "star":
				formatted, err = print.SprintStar(selected.Raw())
			}

			if formatted == "" || err != nil {
				formatted = selected.Raw()
			}
			// note: cannot give correct height before preview is created
			// and we can calculate vertical margin height
			m.preview = preview.New(selected.Filename, formatted)
			height := m.height - m.preview.VerticalMarginHeight()
			m.preview.SetSize(m.width, height)

			return m, nil
		}
	case RenameRequestMsg:
		if m.mode == Edit && m.active == List {
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
				return m, createStatusMsg("Failed to rename request")
			}
			moldIndex := slices.Index(m.requestMolds, requestMold)
			renamedRequest, renamedRequestMold, ok := renameRequest(msg.Input, request, *requestMold)
			if ok {
				m.requestMolds = slices.Replace(m.requestMolds, moldIndex, moldIndex+1, renamedRequestMold)
				setCmd := m.list.SetItem(m.list.Index(), renamedRequest)
				statusCmd := createStatusMsg(fmt.Sprintf("Renamed request to %s", renamedRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, createStatusMsg("Failed to rename request")
			}
		} else if msg.Context.Key == CopyRequest {
			request := msg.Context.Additional.(Request)
			requestMold, err := findRequestMold(request, m)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to find request mold with %v", request)
				return m, createStatusMsg("Failed to copy request")
			}
			copiedRequest, copiedRequestMold, ok := copyRequest(msg.Input, request, *requestMold)
			if ok {
				// note order is not relevant here
				m.requestMolds = append(m.requestMolds, copiedRequestMold)
				setCmd := m.list.InsertItem(m.list.Index()+1, copiedRequest)
				statusCmd := createStatusMsg(fmt.Sprintf("Copied request to %s", copiedRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, createStatusMsg("Failed to copy request")
			}

		} else if msg.Context.Key == CreateSimpleRequest || msg.Context.Key == CreateComplexRequest {
			var (
				root     string
				filepath string
				cmd      *exec.Cmd
				err      error
			)
			if msg.Context.Key == CreateSimpleRequest {
				root, filepath, cmd, err = createSimpleRequestFileCmd(msg.Input)
			} else {
				root, filepath, cmd, err = createComplexRequestFileCmd(msg.Input)
			}
			if err != nil {
				return m, createStatusMsg("Failed preparing editor")
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
			m.profileui = profiles.NewEmbedded(allProfiles, m.width, m.height)
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
			return m, createStatusMsg("Failed to set profile")
		}

		activeProfile = activedProfile

		updateStatusbar(&m, "")
		return m, nil

	case keyprompt.KeypromptAnsweredMsg:
		m.active = List
		if msg.Key == "y" {
			return m, tea.Cmd(func() tea.Msg {
				return CreateRequestMsg{
					Simple: true,
				}
			})
		} else if msg.Key == "s" {
			return m, tea.Cmd(func() tea.Msg {
				return CreateRequestMsg{
					Simple: false,
				}
			})
		} else {
			return m, nil
		}

	case StatusMessage:
		log.Debug().Msgf("Show status message %s", string(msg))
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
		return m.preview.View()
	case Stopwatch:
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			style.stopwatchStyle.Render("Running request\n\n"+m.stopwatch.View()),
			lipgloss.WithWhitespaceChars("\u28FF"),
			lipgloss.WithWhitespaceForeground(style.whitespaceFg))
	case Profiles:
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.profileui.View())
	default:
		return renderList(m)
	}
}

func renderList(m Model) string {
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
		m.prompt.View(),
		lipgloss.WithWhitespaceChars("\u28FF"),
		lipgloss.WithWhitespaceForeground(style.whitespaceFg))
}

func renderKeyprompt(m Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.keyprompt.View(),
		lipgloss.WithWhitespaceChars("\u28FF"),
		lipgloss.WithWhitespaceForeground(style.whitespaceFg))
}

func Start(loadedRequests []*model.RequestMold, loadedProfiles []*model.Profile) {
	log.Info().Msgf("Starting up manage TUI with %d loaded requests and %d profiles", len(loadedRequests), len(loadedProfiles))

	theme := styles.GetTheme()
	InitStyle(theme, styles.GetCommonStyles(theme))

	var requests []list.Item

	for _, v := range loadedRequests {
		r := Request{
			Name:   v.Name(),
			Url:    v.Url(),
			Method: v.Method(),
		}
		requests = append(requests, r)
	}

	for _, p := range loadedProfiles {
		profile := &model.Profile{
			Name:      p.Name,
			Variables: loader.GetProfileValues(p, loadedProfiles),
		}
		if profile.Name == "default" {
			activeProfile = profile
		}
		allProfiles = append(allProfiles, profile)
	}

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
	requestList.Title = "Requests"
	requestList.Styles.Title = style.listTitleStyle

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
	m := Model{
		list:         requestList,
		active:       List,
		mode:         mode,
		stopwatch:    stopwatch.NewWithInterval(time.Millisecond * 100),
		statusbar:    sb,
		help:         help,
		requestMolds: loadedRequests,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(Model); ok {
		handlePostAction(m)
	}
}
