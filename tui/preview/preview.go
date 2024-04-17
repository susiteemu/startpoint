package previewui

import (
	"fmt"
	statusbar "startpoint/tui/statusbar"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title     string
	Viewport  viewport.Model
	Statusbar statusbar.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		footerHeight := lipgloss.Height(m.Statusbar.View())

		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - footerHeight
		m.Statusbar.SetWidth(msg.Width)
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.Statusbar.ChangeText(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100), 2)
	return lipgloss.JoinVertical(lipgloss.Left, contentStyle.Render(m.Viewport.View()), m.Statusbar.View())
}

func (m Model) VerticalMarginHeight() int {
	footerHeight := lipgloss.Height(m.Statusbar.View())
	return footerHeight
}

func (m *Model) SetSize(width int, height int) {
	m.Viewport.Width = width
	m.Viewport.Height = height
	m.Statusbar.SetWidth(width)
}

func New(title, content string) Model {
	v := viewport.New(0, 0)
	v.SetContent(content)

	statusbarItems := []statusbar.StatusbarItem{
		{Text: "PREVIEW", BackgroundColor: statusbarFirstColBg, ForegroundColor: statusbarFirstColFg},
		{Text: title, BackgroundColor: statusbarSecondColBg, ForegroundColor: statusbarSecondColFg},
		{Text: fmt.Sprintf("%3.f%%", v.ScrollPercent()*100), BackgroundColor: statusbarThirdColBg, ForegroundColor: statusbarThirdColFg},
		{Text: "startpoint", BackgroundColor: statusbarFourthColBg, ForegroundColor: statusbarFourthColFg},
	}
	sb := statusbar.New(statusbarItems, 1, 0)

	m := Model{
		title:     title,
		Viewport:  v,
		Statusbar: sb,
	}
	return m
}
