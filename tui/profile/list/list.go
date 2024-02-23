package listtui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().BorderBackground(lipgloss.Color("#cdd6f4"))

type Profile struct {
	Name      string
	Variables int
}

func (i Profile) Title() string       { return i.Name }
func (i Profile) Description() string { return fmt.Sprintf("Vars: %d", i.Variables) }
func (i Profile) FilterValue() string { return i.Name }

type Model struct {
	List      list.Model
	Selection Profile
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
			i, ok := m.List.SelectedItem().(Profile)
			if ok {
				m.Selection = i
				m.Selected = true
			}
			return m, tea.Cmd(func() tea.Msg { return ProfileSelectedMsg{} })
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

func New(profiles []Profile, width, height int, additionalFullHelpKeys []key.Binding) Model {
	items := []list.Item{}

	for _, v := range profiles {
		items = append(items, v)
	}

	d := list.NewDefaultDelegate()

	// Change colors
	titleColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#cdd6f4"}
	descColor := lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#bac2de"}
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(titleColor).BorderLeftForeground(titleColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(descColor).BorderLeftForeground(descColor)

	profileList := list.New(items, d, width, height)
	profileList.Title = "Profiles"
	profileList.Help.ShowAll = true
	if additionalFullHelpKeys != nil {
		profileList.AdditionalFullHelpKeys = func() []key.Binding {
			return additionalFullHelpKeys
		}
	}

	m := Model{
		List:      profileList,
		Selection: Profile{},
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

type ProfileSelectedMsg struct{}
