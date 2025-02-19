package resultsui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/susiteemu/startpoint/core/ansi"
	"github.com/susiteemu/startpoint/core/writer"
	messages "github.com/susiteemu/startpoint/tui/messages"
	"github.com/susiteemu/startpoint/tui/styles"
)

const TOP_NAVIGATION_HEIGHT = 1
const RENDER_LINE_MARGIN = 4

var commonStyles *styles.CommonStyle

type RunResult struct {
	RequestName  string
	RunAt        time.Time
	Results      string
	PlainResults string
}

type Model struct {
	results   []RunResult
	activeIdx int
	Viewport  viewport.Model
	width     int
	height    int
	wPercent  float64
	hPercent  float64
	keyMap    keyMap
}

type keyMap struct {
	Next      key.Binding
	Close     key.Binding
	Copy      key.Binding
	Export    key.Binding
	CloseHelp key.Binding
}

func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keyMap.Next, m.keyMap.Close, m.keyMap.Copy, m.keyMap.Export, m.keyMap.CloseHelp}
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keyMap.Next, m.keyMap.Copy, m.keyMap.Export},
		{m.keyMap.Close, m.keyMap.CloseHelp},
	}
}

var embeddedKeys = keyMap{
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next result"),
	),
	Close: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp(tea.KeyEsc.String(), "close results"),
	),
	Copy: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp("c", "copy to clipboard"),
	),
	Export: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp("w", "write to file"),
	),
	CloseHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "close help"),
	),
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
		m.Viewport.Height = int(float64(msg.Height)*m.hPercent) - TOP_NAVIGATION_HEIGHT
		content := getActiveContent(m)
		m.Viewport.SetContent(renderLines(content, m.Viewport.Width-RENDER_LINE_MARGIN))
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {
		case "n":
			m.activeIdx += 1
			if m.activeIdx >= len(m.results) {
				m.activeIdx = 0
			}
			content := getActiveContent(m)
			m.Viewport.SetContent(renderLines(content, m.Viewport.Width-RENDER_LINE_MARGIN))
			return m, nil
		case "c":
			rawResults := m.results[m.activeIdx].PlainResults
			err := clipboard.WriteAll(rawResults)
			if err != nil {
				log.Error().Err(err).Msg("Failed to copy to clipboard")
				return m, messages.CreateStatusMsg("Failed to copy results to clipboard")
			}
			return m, messages.CreateStatusMsg("Copied results to clipboard")
		case "w":
			results := m.results[m.activeIdx]
			rawResults := results.PlainResults
			workspace := viper.GetString("workspace")
			name := fmt.Sprintf("%s_%s.out", results.RequestName, results.RunAt.Format(time.ANSIC))
			path := filepath.Join(workspace, name)
			_, err := writer.WriteFile(path, rawResults)
			if err != nil {
				log.Error().Err(err).Msg("Failed to write results to file")
				return m, messages.CreateStatusMsg("Failed to write results to file")
			}
			return m, messages.CreateStatusMsg(fmt.Sprintf("Wrote results to \"%s\"", path))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var views []string

	prevResult := lipgloss.NewStyle().Faint(m.activeIdx <= 0).Render("❮")
	nextResult := lipgloss.NewStyle().Faint(m.activeIdx+1 >= len(m.results)).Render("❯")
	activeRun := fmt.Sprintf("[%d/%d] %s", m.activeIdx+1, len(m.results), m.results[m.activeIdx].RunAt.Format(time.Stamp))

	views = append(views, lipgloss.NewStyle().Width(m.Viewport.Width).Align(lipgloss.Center).Render(fmt.Sprintf("%s %s %s", prevResult, activeRun, nextResult)))
	views = append(views, contentStyle.Render(m.Viewport.View()))
	joined := lipgloss.JoinVertical(
		lipgloss.Top,
		views...,
	)

	joined = lipgloss.NewStyle().BorderForeground(styles.LoadTheme().BorderFgColor).Border(lipgloss.RoundedBorder(), true, true).Render(joined)
	return joined
}

func renderLines(content string, width int) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	wrappedLines := []string{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if lipgloss.Width(line) >= width {
			line = wrap.String(line, width)
			colorState, _ := ansi.ParseANSI(line, width)
			line = strings.ReplaceAll(line, "\n", "\u001b[0m\n"+colorState.State)
		}
		wrappedLines = append(wrappedLines, line)
	}

	return strings.Join(wrappedLines, "\n")
}

func getActiveContent(m Model) string {
	return m.results[m.activeIdx].Results
}

func New(results []RunResult, activeIdx, w, h int, wPercent, hPercent float64) Model {
	theme := styles.LoadTheme()
	commonStyles = styles.GetCommonStyles(theme)
	content := results[activeIdx].Results

	width := int(float64(w) * wPercent)
	height := int(float64(h) * hPercent)

	v := viewport.New(width, height-TOP_NAVIGATION_HEIGHT)
	v.SetContent(renderLines(content, width-RENDER_LINE_MARGIN))
	v.Style = v.Style.Padding(0, 0)

	m := Model{
		results:   results,
		activeIdx: activeIdx,
		Viewport:  v,
		wPercent:  wPercent,
		hPercent:  hPercent,
		width:     w,
		height:    h,
		keyMap:    embeddedKeys,
	}
	return m
}
