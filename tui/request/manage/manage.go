package managetui

import (
	"fmt"
	"goful/core/client/validator"
	"goful/core/model"
	"goful/core/print"
	"time"

	preview "goful/tui/request/preview"
	prompt "goful/tui/request/prompt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
)

type ActiveView int

const (
	List ActiveView = iota
	Create
	Update
	Duplicate
	Preview
	Stopwatch
)

type Mode int

const (
	Select Mode = iota
	Edit
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
	list      list.Model
	prompt    prompt.Model
	mode      Mode
	active    ActiveView
	preview   preview.Model
	stopwatch stopwatch.Model
	statusbar statusbar.Model
	selected  Request
	response  string
	width     int
	height    int
	debug     string
	help      help.Model
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
			if m.active == Preview {
				m.active = List
				return m, nil
			}
			if m.mode == Edit && m.active == List {
				m.mode = Select
				m.list.SetDelegate(newSelectDelegate())
				updateStatusbar(&m, "")
				return m, nil
			}
			return m, tea.Quit
		case "a":
			if m.mode == Edit && m.active == List {
				m.active = Create
				m.prompt = prompt.New(false)
				return m, nil
			}
		case "A":
			if m.mode == Edit && m.active == List {
				m.active = Create
				m.prompt = prompt.New(true)
				return m, nil
			}
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
	case RequestFinishedMsg:
		m.response = string(msg)
		return m, tea.Quit
	case EditRequestMsg:
		m.active = Update
		m.selected = msg.Request
		return m, tea.Quit
	case PreviewRequestMsg:
		if m.active != Preview {
			m.active = Preview
			m.preview.Viewport.Width = m.width
			m.preview.Viewport.Height = m.height - m.preview.VerticalMarginHeight()
			selected := msg.Request
			var formatted string
			var err error
			switch selected.Mold.ContentType {
			case "yaml":
				formatted, err = print.SprintYaml(selected.Mold.Raw)
			case "star":
				formatted, err = print.SprintStar(selected.Mold.Raw)
			}

			if formatted == "" || err != nil {
				formatted = selected.Mold.Raw
			}
			m.preview.Viewport.SetContent(formatted)
			m.preview.Viewport.YPosition = 0
			return m, nil
		}
	case prompt.CreateMsg:
		return m, tea.Quit
	case StatusMessage:
		updateStatusbar(&m, string(msg))
		return m, nil
	}

	var cmd tea.Cmd
	switch m.active {
	case List:
		m.list, cmd = m.list.Update(msg)
	case Create:
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
	case Create:
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

		if m.active == Create {
			createRequestFile(m)
		} else if m.active == Update {
			openRequestFileForUpdate(m)
		} else if m.response != "" {
			fmt.Printf("%s", m.response)
		}

	}
}
