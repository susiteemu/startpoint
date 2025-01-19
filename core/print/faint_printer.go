package print

import (
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func SprintFaint(str string) string {
	theme := styles.GetTheme()
	faintStyle := lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true)
	return faintStyle.Render(str)
}
