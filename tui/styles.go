package tui

import (
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	helpPaneStyle lipgloss.Style
	helpKeyStyle  lipgloss.Style
	helpDescStyle lipgloss.Style
}

var style *Styles

func InitStyle(theme *styles.Theme, commonStyles *styles.CommonStyle) {

	style = &Styles{
		helpPaneStyle: commonStyles.HelpPaneStyle,
		helpKeyStyle:  commonStyles.HelpKeyStyle,
		helpDescStyle: commonStyles.HelpDescStyle,
	}
}
