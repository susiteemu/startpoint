package styles

import "github.com/charmbracelet/lipgloss"

var (
	HelpKeyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	HelpDescStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)
	HelpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)
)
