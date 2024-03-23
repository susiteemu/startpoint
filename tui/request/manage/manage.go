package managetui

import (
	"fmt"
	"goful/core/client/validator"
	"goful/core/model"
	"goful/core/print"
	"time"

	create "goful/tui/request/create"
	preview "goful/tui/request/preview"
	"os"

	"github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActiveView int

const (
	List ActiveView = iota
	Create
	CreateComplex
	Update
	Preview
	Stopwatch
)

type Mode int

const (
	Select Mode = iota
	Edit
)

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

type uiModel struct {
	list          list.Model
	create        create.Model
	createComplex create.Model
	mode          Mode
	active        ActiveView
	preview       preview.Model
	stopwatch     stopwatch.Model
	selected      Request
	response      string
	width         int
	height        int
	debug         string
	help          help.Model
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
				return m, nil
			}
			return m, tea.Quit
		case "a":
			if m.mode == Edit && m.active == List {
				m.active = Create
				return m, nil
			}
		case "A":
			if m.mode == Edit && m.active == List {
				m.active = CreateComplex
				return m, nil
			}
		case "i":
			if m.mode == Select && m.active == List {
				m.mode = Edit
				m.list.SetDelegate(newEditModeDelegate())
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
	case create.CreateMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	switch m.active {
	case List:
		m.list, cmd = m.list.Update(msg)
	case Create:
		m.create, cmd = m.create.Update(msg)
	case CreateComplex:
		m.createComplex, cmd = m.createComplex.Update(msg)
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
		return renderCreate(m)
	case CreateComplex:
		return renderCreateComplex(m)
	case Preview:
		return m.preview.View()
	case Stopwatch:
		return stopwatchStyle.Render("Running request... :: Elapsed time: " + m.stopwatch.View())
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

func renderCreateComplex(m uiModel) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.createComplex.View())
}

func Start(loadedRequests []model.RequestMold, mode Mode) {
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
	if mode == Select {
		d = newSelectDelegate()
	} else {
		d = newEditModeDelegate()
	}

	requestList := list.New(requests, d, 0, 0)
	requestList.Title = "Requests"
	requestList.Styles.Title = titleStyle

	m := uiModel{list: requestList, create: create.New(false), createComplex: create.New(true), active: List, mode: mode, stopwatch: stopwatch.NewWithInterval(time.Millisecond)}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(uiModel); ok {

		if m.active == Create || m.active == CreateComplex {
			createRequestFile(m)
		} else if m.active == Update {
			openRequestFileForUpdate(m)
		} else if m.response != "" {
			fmt.Printf("%s", m.response)
		}

	}
}
