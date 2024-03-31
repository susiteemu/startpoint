package managetui

import (
	"fmt"
	"goful/core/client/validator"
	"goful/core/model"
	"goful/core/print"
	"os/exec"
	"time"

	preview "goful/tui/request/preview"
	prompt "goful/tui/request/prompt"
	"os"

	list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
	"github.com/rs/zerolog/log"
)

type ActiveView int

const (
	List ActiveView = iota
	Prompt
	Duplicate
	Preview
	Stopwatch
)

type Mode int

type PostAction struct {
	Type             string
	Payload          interface{}
	AddtionalContext interface{}
}

const (
	Select Mode = iota
	Edit
)

const (
	CreateSimpleRequestLabel  = "Choose a name for your complex request. Make it filename compatible and unique within this workspace. After pressing <enter> program will open your $EDITOR and quit. You will then be able to write the contents of the request."
	CreateComplexRequestLabel = "Choose a name for your request. Make it filename compatible and unique within this workspace. After pressing <enter> program will open your $EDITOR and quit. You will then be able to write the contents of the request."
	RenameRequestLabel        = "Rename your request."
	CopyRequestLabel          = "Choose name for your request."
)

const (
	CreateSimpleRequest  = "CSmplReq"
	CreateComplexRequest = "CCmplxReq"
	EditRequest          = "EReq"
	PrintRequest         = "PReq"
	RenameRequest        = "RnReq"
	CopyRequest          = "CpReq"
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

type Request struct {
	Name   string
	Url    string
	Method string
	Mold   model.RequestMold
}

func (i Request) Title() string {
	return i.Name
}

func (i Request) Description() string {
	validMethod := validator.IsValidMethod(i.Method)

	var methodStyle = lipgloss.NewStyle()

	var color = ""
	if validMethod {
		color = methodColors[i.Method]
		if color == "" {
			color = "#cdd6f4"
		}
		methodStyle = methodStyle.Background(lipgloss.Color(color)).Foreground(lipgloss.Color("#1e1e2e")).PaddingRight(1).PaddingLeft(1)
	} else {
		methodStyle = methodStyle.Foreground(lipgloss.Color("#f38ba8")).Border(lipgloss.Border{Bottom: "^"}, false, false, true, false).BorderForeground(lipgloss.Color("#f38ba8"))

	}

	var urlStyle = lipgloss.NewStyle()

	validUrl := validator.IsValidUrl(i.Url)
	if validUrl {
		urlStyle = urlStyle.Foreground(lipgloss.Color("#b4befe"))
	} else {
		urlStyle = urlStyle.Foreground(lipgloss.Color("#f38ba8")).Border(lipgloss.Border{Bottom: "^"}, false, false, true, false).BorderForeground(lipgloss.Color("#f38ba8"))
	}

	method := i.Method
	if i.Method == "" {
		method = "<method>"
	}

	url := i.Url
	if url == "" {
		url = "<url>"
	}

	return lipgloss.JoinHorizontal(0, methodStyle.Render(method), " ", urlStyle.Render(url))
}
func (i Request) FilterValue() string { return fmt.Sprintf("%s %s %s", i.Name, i.Method, i.Url) }

func updateStatusbar(m *uiModel, msg string) {

	profileText := ""
	if m.mode == Edit {
		m.statusbar.FirstColumnColors.Background = statusbarModeEditBg
		m.statusbar.ThirdColumnColors.Background = statusbarSecondColBg
	} else {
		profileText = "dev" // TODO get profile for realz
		m.statusbar.FirstColumnColors.Background = statusbarModeSelectBg
		m.statusbar.ThirdColumnColors.Background = statusbarThirdColBg
	}

	m.statusbar.SetContent(modeStr(m.mode), msg, profileText, "goful")
}

type uiModel struct {
	mode       Mode
	active     ActiveView
	list       list.Model
	preview    preview.Model
	prompt     prompt.Model
	stopwatch  stopwatch.Model
	statusbar  statusbar.Model
	width      int
	height     int
	postAction PostAction
}

func (m uiModel) Init() tea.Cmd {
	return nil
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusbar.SetSize(msg.Width)
		updateStatusbar(&m, "")
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(m.height - 2)

	case tea.KeyMsg:
		// if we are filtering, it gets all the input
		if m.list.FilterState() == list.Filtering {
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
			if m.active == Preview || m.active == Prompt {
				m.active = List
				return m, nil
			}
			if m.mode == Edit && m.active == List {
				m.mode = Select
				m.list.SetDelegate(newSelectDelegate())
				updateStatusbar(&m, "")
				return m, nil
			}
			return m, nil
		case "i":
			if m.mode == Select && m.active == List {
				m.mode = Edit
				m.list.SetDelegate(newEditModeDelegate())
				updateStatusbar(&m, "")
				return m, nil
			}
		}

	case RunRequestMsg:
		m.active = Stopwatch
		return m, tea.Batch(
			m.stopwatch.Init(),
			doRequest(msg.Request),
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
			promptLabel := CreateSimpleRequestLabel
			if !msg.Simple {
				promptKey = CreateComplexRequest
				promptLabel = CreateComplexRequestLabel
			}
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key: promptKey,
			}, "", promptLabel, checkRequestWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case EditRequestMsg:
		if m.mode == Edit && m.active == List {
			cmd, err := openFileToEditorCmd(msg.Request)
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
		deleted := msg.Request.Mold.DeleteFromFS()
		if deleted {
			index := m.list.Index()
			m.list.RemoveItem(index)
			return m, createStatusMsg(fmt.Sprintf("Deleted %s", msg.Request.Title()))
		} else {
			return m, createStatusMsg(fmt.Sprintf("Failed to delete %s", msg.Request.Title()))
		}
	case EditRequestFinishedMsg:
		oldRequest := msg.Request
		if msg.err == nil {
			newRequest, ok := readRequest(oldRequest.Mold.Root, oldRequest.Mold.Filename)
			if ok {
				setCmd := m.list.SetItem(m.list.Index(), newRequest)
				statusCmd := createStatusMsg(fmt.Sprintf("Edited request %s", oldRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			}
		}
		return m, createStatusMsg(fmt.Sprintf("Failed to edit request %s", oldRequest.Title()))
	case PreviewRequestMsg:
		if m.active == List {
			m.active = Preview
			m.preview.Viewport.Width = m.width
			m.preview.Viewport.Height = m.height - m.preview.VerticalMarginHeight()
			selected := msg.Request
			var formatted string
			var err error
			switch selected.Mold.ContentType {
			case "yaml":
				formatted, err = print.SprintYaml(selected.Mold.Raw())
			case "star":
				formatted, err = print.SprintStar(selected.Mold.Raw())
			}

			if formatted == "" || err != nil {
				formatted = selected.Mold.Raw()
			}
			m.preview.Viewport.SetContent(formatted)
			m.preview.Viewport.YPosition = 0
			return m, nil
		}
	case RenameRequestMsg:
		if m.active == List {
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        RenameRequest,
				Additional: msg.Request,
			}, msg.Request.Name, RenameRequestLabel, checkRequestWithNameDoesNotExist(m), m.width)
			return m, nil
		}
	case CopyRequestMsg:
		if m.active == List {
			m.active = Prompt
			m.prompt = prompt.New(prompt.PromptContext{
				Key:        CopyRequest,
				Additional: msg.Request,
			}, fmt.Sprintf("Copy of %s", msg.Request.Name), CopyRequestLabel, checkRequestWithNameDoesNotExist(m), m.width)
		}
	case prompt.PromptAnsweredMsg:
		m.active = List
		if msg.Context.Key == RenameRequest {
			renamedRequest, ok := renameRequest(msg.Input, msg.Context.Additional.(Request))
			if ok {
				log.Debug().Msgf("Index of renamed item is %d", m.list.Index())
				setCmd := m.list.SetItem(m.list.Index(), renamedRequest)
				statusCmd := createStatusMsg(fmt.Sprintf("Renamed request to %s", renamedRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			} else {
				return m, createStatusMsg("Failed to rename request")
			}
		} else if msg.Context.Key == CopyRequest {
			copiedRequest, ok := copyRequest(msg.Input, msg.Context.Additional.(Request))
			if ok {
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
	case CreateRequestFinishedMsg:
		if msg.err == nil {
			newRequest, ok := readRequest(msg.root, msg.filename)
			if ok {
				setCmd := m.list.InsertItem(m.list.Index()+1, newRequest)
				statusCmd := createStatusMsg(fmt.Sprintf("Created request %s", newRequest.Title()))
				return m, tea.Batch(setCmd, statusCmd)
			}
		}
		return m, createStatusMsg("Failed to create request")

	case StatusMessage:
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
	case Stopwatch:
		m.stopwatch, cmd = m.stopwatch.Update(msg)
	}
	return m, cmd
}

func (m uiModel) View() string {
	switch m.active {
	case List:
		return renderList(m)
	case Prompt:
		return renderPrompt(m)
	case Preview:
		return m.preview.View()
	case Stopwatch:
		return stopwatchStyle.Render("Running request... :: Elapsed time: " + m.stopwatch.View())
	default:
		return renderList(m)
	}
}

func renderList(m uiModel) string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(m.list.View()),
		m.statusbar.View(),
	)
}

func renderPrompt(m uiModel) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.prompt.View())
}

func Start(loadedRequests []model.RequestMold) {
	log.Info().Msgf("Starting up manage TUI with %d loaded requests", len(loadedRequests))

	var requests []list.Item

	for _, v := range loadedRequests {
		r := Request{
			Name:   v.Name(),
			Url:    v.Url(),
			Method: v.Method(),
			Mold:   v,
		}
		requests = append(requests, r)
	}

	var d list.DefaultDelegate
	var modeColor lipgloss.AdaptiveColor
	d = newSelectDelegate()
	modeColor = statusbarModeSelectBg

	requestList := list.New(requests, d, 0, 0)
	requestList.Title = "Requests"
	requestList.Styles.Title = titleStyle

	requestList.Help.Styles.FullKey = helpKeyStyle
	requestList.Help.Styles.FullDesc = helpDescStyle
	requestList.Help.Styles.ShortKey = helpKeyStyle
	requestList.Help.Styles.ShortDesc = helpDescStyle
	requestList.Help.Styles.ShortSeparator = helpSeparatorStyle
	requestList.Help.Styles.FullSeparator = helpSeparatorStyle

	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: statusbarFourthColFg,
			Background: modeColor,
		},
		statusbar.ColorConfig{
			Foreground: statusbarSecondColFg,
			Background: statusbarSecondColBg,
		},
		statusbar.ColorConfig{
			Foreground: statusbarThirdColFg,
			Background: statusbarThirdColBg,
		},
		statusbar.ColorConfig{
			Foreground: statusbarFourthColFg,
			Background: statusbarFourthColBg,
		},
	)
	m := uiModel{list: requestList, active: List, mode: Select, stopwatch: stopwatch.NewWithInterval(time.Millisecond), statusbar: sb}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(uiModel); ok {
		handlePostAction(m)
	}
}
