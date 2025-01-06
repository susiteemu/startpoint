package previewui

import (
	"fmt"
	"startpoint/tui/styles"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title    string
	Viewport viewport.Model
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
		m.Viewport.Width = int(float64(msg.Width) * 0.8)
		m.Viewport.Height = int(float64(msg.Height) * 0.8)
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true).Render(contentStyle.Render(m.Viewport.View()))
}

func (m *Model) SetSize(width int, height int) {
	m.Viewport.Width = width
	m.Viewport.Height = height
}

func New(title, content string) Model {

	theme := styles.GetTheme()

	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	digits := len(strconv.Itoa(len(lines)))
	lineNrFmt := "%" + fmt.Sprintf("%d", digits) + "d"
	linesWithLineNrs := []string{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		lineNr := i + 1
		lineNrSection := lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true).Render(fmt.Sprintf(lineNrFmt, lineNr))
		line = fmt.Sprintf("%s  %s", lineNrSection, line)
		linesWithLineNrs = append(linesWithLineNrs, line)
	}

	v := viewport.New(0, 0)
	v.SetContent(strings.Join(linesWithLineNrs, "\n"))
	v.Style = v.Style.Padding(0, 0)

	m := Model{
		title:    title,
		Viewport: v,
	}
	return m
}
