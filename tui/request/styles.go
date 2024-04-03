package requestui

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

	stopwatchStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center).Bold(true)

	requestTitleColor = lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#cdd6f4"}
	requestDescColor  = lipgloss.AdaptiveColor{Light: "#cdd6f4", Dark: "#bac2de"}

	statusbarModeSelectBg = lipgloss.Color("#f9e2af")
	statusbarModeEditBg   = lipgloss.Color("#a6e3a1")
	statusbarFirstColFg   = lipgloss.Color("#1e1e2e")
	statusbarSecondColBg  = lipgloss.Color("#11111b")
	statusbarSecondColFg  = lipgloss.Color("#FFFDF5")
	statusbarThirdColBg   = lipgloss.Color("#94e2d5")
	statusbarThirdColFg   = lipgloss.Color("#1e1e2e")
	statusbarFourthColBg  = lipgloss.Color("#89b4fa")
	statusbarFourthColFg  = lipgloss.Color("#1e1e2e")

	listStyle    = lipgloss.NewStyle().BorderBackground(lipgloss.Color("#cdd6f4"))
	methodColors = map[string]string{
		"GET":    "#89b4fa",
		"POST":   "#a6e3a1",
		"PUT":    "#f9e2af",
		"DELETE": "#f38ba8",
		"PATCH":  "#94e2d5",
		// TODO etc
	}
)
