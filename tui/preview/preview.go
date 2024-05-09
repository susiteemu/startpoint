package previewui

import (
	"fmt"
	statusbar "startpoint/tui/statusbar"
	"startpoint/tui/styles"
	"strconv"
	"strings"

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

	statusbarItems := []statusbar.StatusbarItem{
		{Text: "PREVIEW", BackgroundColor: theme.StatusbarModePrimaryBgColor, ForegroundColor: theme.StatusbarSecondaryFgColor},
		{Text: title, BackgroundColor: theme.StatusbarPrimaryBgColor, ForegroundColor: theme.StatusbarPrimaryFgColor},
		{Text: fmt.Sprintf("%3.f%%", v.ScrollPercent()*100), BackgroundColor: theme.StatusbarThirdColBgColor, ForegroundColor: theme.StatusbarSecondaryFgColor},
		{Text: "startpoint", BackgroundColor: theme.StatusbarFourthColBgColor, ForegroundColor: theme.StatusbarSecondaryFgColor},
	}
	sb := statusbar.New(statusbarItems, 1, 0)

	m := Model{
		title:     title,
		Viewport:  v,
		Statusbar: sb,
	}
	return m
}
