package profileui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.NoColor{}).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#a6e3a1")).
			Padding(0)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f38ba8"))

	helpKeyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	helpDescStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)
	helpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)

	stopwatchStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center).Bold(true)

	profileTitleColor = lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#cdd6f4"}
	profileDescColor  = lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#bac2de"}

	statusbarFirstColBg  = lipgloss.Color("#11111b")
	statusbarFirstColFg  = lipgloss.Color("#FFFDF5")
	statusbarSecondColBg = lipgloss.Color("#89b4fa")
	statusbarSecondColFg = lipgloss.Color("#1e1e2e")

	listStyle = lipgloss.NewStyle().BorderBackground(lipgloss.Color("#cdd6f4"))
)
