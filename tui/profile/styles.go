package profileui

import (
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	statusbarFirstColBg  lipgloss.Color
	statusbarSecondColBg lipgloss.Color
	statusbarFirstColFg  lipgloss.Color
	statusbarSecondColFg lipgloss.Color

	listTitleStyle     lipgloss.Style
	listItemTitleColor lipgloss.Color
	listItemDescColor  lipgloss.Color

	helpPaneStyle lipgloss.Style
	helpKeyStyle  lipgloss.Style
	helpDescStyle lipgloss.Style

	httpMethodColors map[string]lipgloss.Color
}

var style *Styles

func InitStyle(theme *styles.Theme, commonStyles *styles.CommonStyle) {

	style = &Styles{
		statusbarFirstColBg:  theme.StatusbarPrimaryBgColor,
		statusbarSecondColBg: theme.StatusbarFourthColBgColor,
		statusbarFirstColFg:  theme.StatusbarPrimaryFgColor,
		statusbarSecondColFg: theme.StatusbarSecondaryFgColor,

		listTitleStyle:     lipgloss.NewStyle().Foreground(theme.TitleFgColor).Background(theme.TitleBgColor).Padding(0, 1).MarginTop(1),
		listItemTitleColor: theme.TextFgColor,
		listItemDescColor:  theme.SubtextFgColor,
		helpPaneStyle:      commonStyles.HelpPaneStyle.Copy(),
		helpKeyStyle:       commonStyles.HelpKeyStyle.Copy(),
		helpDescStyle:      commonStyles.HelpDescStyle.Copy(),
	}
}
