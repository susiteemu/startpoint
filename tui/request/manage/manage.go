package managetui

import (
	"fmt"
	"goful/core/model"
	create "goful/tui/request/create"
	list "goful/tui/request/list"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type ActiveView int

const (
	List ActiveView = iota
	Create
	CreateComplex
	Update
)

var keys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add simple"),
	),
	key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "Add complex"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
		key.WithHelp(tea.KeyEnter.String(), "Edit"),
	),
}

type uiModel struct {
	list          list.Model
	create        create.Model
	createComplex create.Model
	active        ActiveView
	selected      list.Request
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
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "a":
			if m.active == List {
				m.active = Create
				return m, nil
			}
		case "A":
			if m.active == List {
				m.active = CreateComplex
				return m, nil
			}
		}
	case list.RequestSelectedMsg:
		m.active = Update
		m.selected = m.list.Selection
		return m, tea.Quit
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

func Start(loadedRequests []model.RequestMold) {
	var requests []list.Request

	for _, v := range loadedRequests {
		r := list.Request{
			Name:   v.Name(),
			Url:    v.Url(),
			Method: v.Method(),
			Mold:   v,
		}
		requests = append(requests, r)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := uiModel{list: list.New(requests, 0, 0, keys), create: create.New(false), createComplex: create.New(true), active: List}

	p := tea.NewProgram(m, tea.WithAltScreen())

	r, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(uiModel); ok {

		if m.active == Create || m.active == CreateComplex {
			createRequestFile(m, f)
		} else if m.active == Update {
			openRequestFileForUpdate(m, f)
		}

	}
}

func createRequestFile(m uiModel, logFile *os.File) {

	fileName := ""
	content := ""
	createFile := false
	if m.active == Create {
		fileName = fmt.Sprintf("%s.yaml", m.create.Name)
		createFile = true
		// TODO read from a template file
		content = fmt.Sprintf(`name: %s
# Possible request to call _before_ this one
prev_req:
# Request url, may contain template variables in a form of {var}
url:
# HTTP method
method:
# HTTP headers as key-val list, e.g. X-Foo-Bar: SomeValue
headers:
# Request body, e.g.
# {
#    "id": 1,
#    "name": "Jane">
# }
body: >
`, m.create.Name)
	} else if m.active == CreateComplex {
		fileName = fmt.Sprintf("%s.star", m.createComplex.Name)
		// TODO read from template
		content = fmt.Sprintf(`"""
meta:name: %s
meta:prev_req: <call other request before this>
doc:url: <your url for display>
doc:method: <your http method for display>
"""
# insert contents of your script here, for more see https://github.com/google/starlark-go/blob/master/doc/spec.md
# Request url
url = ""
# HTTP method
method = ""
# HTTP headers, e.g. { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
headers = {}
# Request body, e.g. { "id": 1, "people": [ {"name": "Joe"}, {"name": "Jane"}, ] }
body = {}
`, m.createComplex.Name)
		createFile = true
	}

	if !createFile {
		return
	}

	logFile.WriteString(fmt.Sprintf("About to create new request with name %v\n", fileName))
	if len(fileName) > 0 {
		file, err := os.Create("tmp/" + fileName)
		if err == nil {
			defer file.Close()
			// TODO handle err
			file.WriteString(content)
			file.Sync()
			filename := file.Name()
			editor := viper.GetString("editor")
			if editor == "" {
				logFile.WriteString("Editor is not configured through configuration file or $EDITOR environment variable.")
			}

			logFile.WriteString(fmt.Sprintf("Opening file %s\n", filename))
			cmd := exec.Command(editor, filename)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				logFile.WriteString(fmt.Sprintf("Failed to open file with editor: %v\n", err))
			}
			log.Printf("Successfully edited file %v", file.Name())
			fmt.Printf("Saved new request to file %v", file.Name())
		} else {
			logFile.WriteString(fmt.Sprintf("Failed to create file %v\n", err))
		}
	}
}

func openRequestFileForUpdate(m uiModel, logFile *os.File) {
	if m.active == Update && m.selected.Name != "" {

		fileName := fmt.Sprintf("tmp/%s", m.selected.Mold.Filename)
		logFile.WriteString(fmt.Sprintf("About to open request file %v\n", fileName))
		if len(fileName) > 0 {
			// TODO handle err
			editor := viper.GetString("editor")
			if editor == "" {
				logFile.WriteString("Editor is not configured through configuration file or $EDITOR environment variable.")
			}

			cmd := exec.Command(editor, fileName)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				logFile.WriteString(fmt.Sprintf("Failed to open file with editor: %v\n", err))
			}
			logFile.WriteString(fmt.Sprintf("Successfully edited file %v\n", fileName))
		}
	}
}
