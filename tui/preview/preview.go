package previewui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/susiteemu/startpoint/core/ansi"
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

const RENDER_LINE_MARGIN = 4

type Model struct {
	title    string
	content  string
	Viewport viewport.Model
	wPercent float64
	hPercent float64
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
		m.Viewport.Width = int(float64(msg.Width) * m.wPercent)
		m.Viewport.Height = int(float64(msg.Height) * m.hPercent)
		m.Viewport.SetContent(renderLines(m.content, m.Viewport.Width-RENDER_LINE_MARGIN))
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().BorderForeground(styles.LoadTheme().BorderFgColor).Border(lipgloss.RoundedBorder(), true, true).Render(contentStyle.Render(m.Viewport.View()))
}

func (m *Model) SetSize(width int, height int) {
	m.Viewport.Width = width
	m.Viewport.Height = height
}

func renderLines(content string, width int) string {
	theme := styles.LoadTheme()

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
			line = wrap.String(line, width)
			colorState, _ := ansi.ParseANSI(line, width)
			line = strings.ReplaceAll(line, "\n", "\u001b[0m\n"+colorState.State)
		}

		linesWithLineNrs = append(linesWithLineNrs, line)
	}

	return strings.Join(linesWithLineNrs, "\n")
}

func New(title, content string, w, h int, wPercent, hPercent float64) Model {

	width := int(float64(w) * wPercent)
	height := int(float64(h) * hPercent)
	v := viewport.New(width, height)
	v.SetContent(renderLines(content, width-RENDER_LINE_MARGIN))
	v.Style = v.Style.Padding(0, 0)

	m := Model{
		title:    title,
		content:  content,
		Viewport: v,
		wPercent: wPercent,
		hPercent: hPercent,
	}
	return m
}
