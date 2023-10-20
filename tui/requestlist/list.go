package listtui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center)

type Request struct {
	Name    string
	Url     string
	Method  string
	Headers map[string]string
	Body    []byte
}

func (i Request) Title() string       { return i.Name }
func (i Request) Description() string { return fmt.Sprintf("%v %v", i.Method, i.Url) }
func (i Request) FilterValue() string { return i.Name }

type Model struct {
	list      list.Model
	Selection Request
	Selected  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(Request)
			if ok {
				m.Selection = i
				m.Selected = true
			}
			return m, tea.Cmd(func() tea.Msg { return RequestSelectedMsg{} })
		}

	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func New(requests []Request, width, height int) Model {
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
	// TODO AdditionalFullHelpKeys

	m := Model{
		list:      requestList,
		Selection: Request{},
		Selected:  false,
	}
	return m
}

func (m Model) View() string {
	return listStyle.Render(m.list.View())
}

type RequestSelectedMsg struct{}
