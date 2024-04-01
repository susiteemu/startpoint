package profilemgmtui

import (
	"fmt"
	"goful/core/model"
	profiles "goful/tui/profile"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	profiles profiles.Model
	width    int
	height   int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.profiles, cmd = m.profiles.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.profiles.View()
}

func Start(loadedProfiles []model.Profile) {
	m := Model{
		profiles: profiles.New(loadedProfiles),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
