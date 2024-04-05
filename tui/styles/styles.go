package styles

import "github.com/charmbracelet/lipgloss"

var (
	HelpPaneStyle      = lipgloss.NewStyle().Padding(1)
	HelpKeyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	HelpDescStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)
	HelpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)
)
