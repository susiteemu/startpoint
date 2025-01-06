package profilemgmtui

import (
	"fmt"
	"os"
	"startpoint/core/model"
	profiles "startpoint/tui/profile"
	"startpoint/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
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

func Start(loadedProfiles []*model.Profile, workspace string) {
	m := Model{
		profiles: profiles.New(loadedProfiles, workspace),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	theme := styles.GetTheme()

	originalBackground := termenv.DefaultOutput().BackgroundColor()
	termenv.DefaultOutput().SetBackgroundColor(termenv.RGBColor(theme.BgColor))
	_, err := p.Run()
	termenv.DefaultOutput().SetBackgroundColor(originalBackground)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
