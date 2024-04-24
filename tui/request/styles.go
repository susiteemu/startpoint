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
	listStatusbarFg    lipgloss.Color
	listFilterPromptFg lipgloss.Color
	listFilterCursorFg lipgloss.Color

	helpPaneStyle lipgloss.Style
	helpKeyStyle  lipgloss.Style
	helpDescStyle lipgloss.Style

	whitespaceFg lipgloss.Color

	httpMethodTextColor    lipgloss.Color
	httpMethodDefaultColor lipgloss.Color
	httpMethodColors       map[string]lipgloss.Color

	urlFg                         lipgloss.Color
	urlBg                         lipgloss.Color
	urlTemplatedSectionFg         lipgloss.Color
	urlTemplatedSectionBg         lipgloss.Color
	urlUnfilledTemplatedSectionFg lipgloss.Color
	urlUnfilledTemplatedSectionBg lipgloss.Color
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
		listStatusbarFg:    theme.TextFgColor,
		listFilterPromptFg: theme.TextFgColor,
		listFilterCursorFg: theme.TextFgColor,

		httpMethodTextColor:    theme.HttpMethodTextFgColor,
		httpMethodDefaultColor: theme.HttpMethodDefaultBgColor,
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

		urlFg:                         theme.UrlFgColor,
		urlBg:                         theme.UrlBgColor,
		urlTemplatedSectionFg:         theme.UrlTemplatedSectionFgColor,
		urlTemplatedSectionBg:         theme.UrlTemplatedSectionBgColor,
		urlUnfilledTemplatedSectionFg: theme.UrlUnfilledTemplatedSectionFgColor,
		urlUnfilledTemplatedSectionBg: theme.UrlUnfilledTemplatedSectionBgColor,
	}
}
