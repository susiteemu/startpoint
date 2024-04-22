package requestui

import (
	"startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	stopwatchStyle        lipgloss.Style
	statusbarModeSelectBg lipgloss.Color
	statusbarModeEditBg   lipgloss.Color
	statusbarPrimaryBg    lipgloss.Color
	statusbarThirdColBg   lipgloss.Color
	statusbarFourthColBg  lipgloss.Color
	statusbarPrimaryFg    lipgloss.Color
	statusbarSecondaryFg  lipgloss.Color

	listTitleStyle     lipgloss.Style
	listItemTitleColor lipgloss.Color
	listItemDescColor  lipgloss.Color

	helpPaneStyle lipgloss.Style
	helpKeyStyle  lipgloss.Style
	helpDescStyle lipgloss.Style

	whitespaceFg lipgloss.Color

	httpMethodColors map[string]lipgloss.Color
}

var style *Styles

func InitStyle(theme *styles.Theme, commonStyles *styles.CommonStyle) {
	style = &Styles{
		stopwatchStyle:        lipgloss.NewStyle().BorderForeground(theme.BorderFgColor).Border(lipgloss.RoundedBorder()).Align(lipgloss.Center).Bold(true).Padding(2, 5),
		statusbarModeSelectBg: theme.StatusbarModePrimaryBgColor,
		statusbarModeEditBg:   theme.StatusbarModeSecondaryBgColor,
		statusbarPrimaryBg:    theme.StatusbarPrimaryBgColor,
		statusbarThirdColBg:   theme.StatusbarThirdColBgColor,
		statusbarFourthColBg:  theme.StatusbarFourthColBgColor,
		statusbarPrimaryFg:    theme.StatusbarPrimaryFgColor,
		statusbarSecondaryFg:  theme.StatusbarSecondaryFgColor,

		listTitleStyle:     lipgloss.NewStyle().Foreground(theme.TitleFgColor).Background(theme.TitleBgColor).Padding(0, 1).MarginTop(1),
		listItemTitleColor: theme.TextFgColor,
		listItemDescColor:  theme.SubtextFgColor,
		httpMethodColors: map[string]lipgloss.Color{
			"GET":    theme.HttpMethodGetBgColor,
			"POST":   theme.HttpMethodPostBgColor,
			"PUT":    theme.HttpMethodPutBgColor,
			"DELETE": theme.HttpMethodDeleteBgColor,
			"PATCH":  theme.HttpMethodPatchBgColor,
			// TODO etc
		},
		helpPaneStyle: commonStyles.HelpPaneStyle.Copy(),
		helpKeyStyle:  commonStyles.HelpKeyStyle.Copy(),
		helpDescStyle: commonStyles.HelpDescStyle.Copy(),

		whitespaceFg: theme.WhitespaceFgColor,
	}
}
