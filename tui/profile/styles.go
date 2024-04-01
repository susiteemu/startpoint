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

	statusbarModeSelectBg = lipgloss.AdaptiveColor{Light: "#f9e2af", Dark: "#f9e2af"}
	statusbarModeEditBg   = lipgloss.AdaptiveColor{Light: "#a6e3a1", Dark: "#a6e3a1"}
	statusbarFirstColFg   = lipgloss.AdaptiveColor{Light: "#1e1e2e", Dark: "#1e1e2e"}
	statusbarSecondColBg  = lipgloss.AdaptiveColor{Light: "#11111b", Dark: "#11111b"}
	statusbarSecondColFg  = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
	statusbarThirdColBg   = lipgloss.AdaptiveColor{Light: "#94e2d5", Dark: "#94e2d5"}
	statusbarThirdColFg   = lipgloss.AdaptiveColor{Light: "#1e1e2e", Dark: "#1e1e2e"}
	statusbarFourthColBg  = lipgloss.AdaptiveColor{Light: "#89b4fa", Dark: "#89b4fa"}
	statusbarFourthColFg  = lipgloss.AdaptiveColor{Light: "#1e1e2e", Dark: "#1e1e2e"}

	listStyle = lipgloss.NewStyle().BorderBackground(lipgloss.Color("#cdd6f4"))
)
