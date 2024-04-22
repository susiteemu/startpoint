package styles

import (
	"startpoint/core/configuration"

	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type Theme struct {
	BgColor                       lipgloss.Color
	TextFgColor                   lipgloss.Color
	SubtextFgColor                lipgloss.Color
	TitleFgColor                  lipgloss.Color
	TitleBgColor                  lipgloss.Color
	BorderFgColor                 lipgloss.Color
	StatusbarPrimaryBgColor       lipgloss.Color
	StatusbarPrimaryFgColor       lipgloss.Color
	StatusbarSecondaryFgColor     lipgloss.Color
	StatusbarModePrimaryBgColor   lipgloss.Color
	StatusbarModeSecondaryBgColor lipgloss.Color
	StatusbarThirdColBgColor      lipgloss.Color
	StatusbarFourthColBgColor     lipgloss.Color
	HttpMethodTextFgColor         lipgloss.Color
	HttpMethodGetBgColor          lipgloss.Color
	HttpMethodPostBgColor         lipgloss.Color
	HttpMethodPutBgColor          lipgloss.Color
	HttpMethodDeleteBgColor       lipgloss.Color
	HttpMethodPatchBgColor        lipgloss.Color
	WhitespaceFgColor             lipgloss.Color
}

var theme *Theme

func GetTheme() *Theme {
	log.Debug().Msgf("Get theme")
	if theme == nil {
		getColor := configuration.GetStringOrDefault
		theme = &Theme{
			TextFgColor:                   lipgloss.Color(getColor("theme.primaryTextFgColor")),
			SubtextFgColor:                lipgloss.Color(getColor("theme.secondaryTextFgColor")),
			TitleFgColor:                  lipgloss.Color(getColor("theme.titleFgColor")),
			TitleBgColor:                  lipgloss.Color(getColor("theme.titleBgColor")),
			BorderFgColor:                 lipgloss.Color(getColor("theme.borderFgColor")),
			StatusbarPrimaryBgColor:       lipgloss.Color(getColor("theme.statusbar.primaryBgColor")),
			StatusbarPrimaryFgColor:       lipgloss.Color(getColor("theme.statusbar.primaryFgColor")),
			StatusbarSecondaryFgColor:     lipgloss.Color(getColor("theme.statusbar.secondaryFgColor")),
			StatusbarModePrimaryBgColor:   lipgloss.Color(getColor("theme.statusbar.modePrimaryBgColor")),
			StatusbarModeSecondaryBgColor: lipgloss.Color(getColor("theme.statusbar.modeSecondaryBgColor")),
			StatusbarThirdColBgColor:      lipgloss.Color(getColor("theme.statusbar.thirdColBgColor")),
			StatusbarFourthColBgColor:     lipgloss.Color(getColor("theme.statusbar.fourthColBgColor")),
			HttpMethodTextFgColor:         lipgloss.Color(getColor("theme.httpMethods.textFgColor")),
			HttpMethodGetBgColor:          lipgloss.Color(getColor("theme.httpMethods.getBgColor")),
			HttpMethodPostBgColor:         lipgloss.Color(getColor("theme.httpMethods.postBgColor")),
			HttpMethodPutBgColor:          lipgloss.Color(getColor("theme.httpMethods.putBgColor")),
			HttpMethodDeleteBgColor:       lipgloss.Color(getColor("theme.httpMethods.deleteBgColor")),
			HttpMethodPatchBgColor:        lipgloss.Color(getColor("theme.httpMethods.patchBgColor")),
			WhitespaceFgColor:             lipgloss.Color(getColor("theme.whitespaceFgColor")),
		}
	}

	return theme
}

type CommonStyle struct {
	HelpPaneStyle      lipgloss.Style
	HelpKeyStyle       lipgloss.Style
	HelpDescStyle      lipgloss.Style
	HelpSeparatorStyle lipgloss.Style
}

var commonStyle *CommonStyle

func GetCommonStyles(theme *Theme) *CommonStyle {
	if commonStyle == nil {
		commonStyle = &CommonStyle{
			HelpPaneStyle:      lipgloss.NewStyle().Padding(1),
			HelpKeyStyle:       lipgloss.NewStyle().Foreground(theme.TextFgColor),
			HelpDescStyle:      lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true),
			HelpSeparatorStyle: lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true),
		}
	}
	return commonStyle
}
