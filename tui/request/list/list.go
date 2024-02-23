package listtui

import (
	"fmt"
	"goful/core/client/validator"
	"goful/core/model"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.NoColor{}).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#a6e3a1")).
	Padding(0)

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#f38ba8"))

var listStyle = lipgloss.NewStyle().BorderBackground(lipgloss.Color("#cdd6f4"))
var methodColors = map[string]string{
	"GET":    "#89b4fa",
	"POST":   "#a6e3a1",
	"PUT":    "#f9e2af",
	"DELETE": "#f38ba8",
	"PATCH":  "#94e2d5",
	// TODO etc
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

type Model struct {
	List      list.Model
	Selection Request
	Selected  bool
	width     int
	height    int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := m.List.SelectedItem().(Request)
			if ok {
				if !validator.IsValidUrl(i.Url) || !validator.IsValidMethod(i.Method) {
					statusCmd := m.List.NewStatusMessage(statusMessageStyle.Render("\ue654 Invalid request."))
					return m, tea.Batch(statusCmd)
				} else {
					m.Selection = i
					m.Selected = true
					return m, tea.Cmd(func() tea.Msg { return RequestSelectedMsg{} })
				}
			} else {
				statusCmd := m.List.NewStatusMessage(statusMessageStyle.Render("\ue654 Invalid request."))
				return m, tea.Batch(statusCmd)
			}
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.List.SetSize(msg.Width, msg.Height)
		m.width = h
		m.height = v
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

func (m Model) SelectedItem() Request {
	return m.List.SelectedItem().(Request)
}

func New(requests []Request, width, height int, additionalFullHelpKeys []key.Binding) Model {
	items := []list.Item{}

	for _, v := range requests {
		items = append(items, v)
	}

	d := list.NewDefaultDelegate()
	d.SetHeight(3)

	// Change colors
	titleColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#cdd6f4"}
	descColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#bac2de"}
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(titleColor).BorderLeftForeground(titleColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(descColor).BorderLeftForeground(descColor)

	requestList := list.New(items, d, width, height)
	requestList.Title = "Requests"
	requestList.Styles.Title = titleStyle
	requestList.Help.ShowAll = true
	if additionalFullHelpKeys != nil {
		requestList.AdditionalFullHelpKeys = func() []key.Binding {
			return additionalFullHelpKeys
		}
	}

	m := Model{
		List:      requestList,
		Selection: Request{},
		Selected:  false,
	}
	return m
}

func (m Model) View() string {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	f.Write([]byte(fmt.Sprintf("w:%v, h: %v\n", m.width, m.height)))

	defer f.Close()
	return m.List.View()
}

type RequestSelectedMsg struct{}
