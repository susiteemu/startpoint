package previewui

import (
	"fmt"
	"github.com/susiteemu/startpoint/core/ansi"
	"github.com/susiteemu/startpoint/tui/styles"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Model struct {
	title    string
	content  string
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
		m.Viewport.SetContent(renderLines(m.content, m.Viewport.Width-4))
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().BorderForeground(styles.GetTheme().BorderFgColor).Border(lipgloss.RoundedBorder(), true, true).Render(contentStyle.Render(m.Viewport.View()))
}

func (m *Model) SetSize(width int, height int) {
	m.Viewport.Width = width
	m.Viewport.Height = height
}

func renderLines(content string, width int) string {
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
		if lipgloss.Width(line) >= width {
			line = wordwrap.String(line, width)
			colorState, _ := ansi.ParseANSI(line, width)
			line = strings.ReplaceAll(line, "\n", "\u001b[0m\n"+colorState.State)
		}

		linesWithLineNrs = append(linesWithLineNrs, line)
	}

	return strings.Join(linesWithLineNrs, "\n")
}

func New(title, content string, w, h int) Model {

	v := viewport.New(w, h)
	v.SetContent(renderLines(content, w-4))
	v.Style = v.Style.Padding(0, 0)

	m := Model{
		title:    title,
		content:  content,
		Viewport: v,
	}
	return m
}
