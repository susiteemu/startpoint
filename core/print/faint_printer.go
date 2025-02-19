package print

import (
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func SprintFaint(str string) string {
	theme := styles.LoadTheme()
	faintStyle := lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true)
	return faintStyle.Render(str)
}
