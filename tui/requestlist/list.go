package listtui

import (
	"fmt"
	"goful/core/model"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().BorderBackground(lipgloss.Color("#cdd6f4"))

type Request struct {
	Name   string
	Url    string
	Method string
	Mold   *model.RequestMold
}

func (i Request) Title() string       { return i.Name }
func (i Request) Description() string { return fmt.Sprintf("%v %v", i.Method, i.Url) }
func (i Request) FilterValue() string { return i.Name }

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
				m.Selection = i
				m.Selected = true
			}
			return m, tea.Cmd(func() tea.Msg { return RequestSelectedMsg{} })
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

func New(requests []Request, width, height int, additionalFullHelpKeys []key.Binding) Model {
	items := []list.Item{}

	for _, v := range requests {
		items = append(items, v)
	}

	d := list.NewDefaultDelegate()

	// Change colors
	titleColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#cdd6f4"}
	descColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#bac2de"}
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(titleColor).BorderLeftForeground(titleColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(descColor).BorderLeftForeground(descColor)

	requestList := list.New(items, d, width, height)
	requestList.Title = "Requests"
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
