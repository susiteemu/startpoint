package styles

import (
	"embed"
	"fmt"
	"os"

	"github.com/susiteemu/startpoint/core/configuration"

	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

//go:embed themes/*.yaml
var embeddedThemes embed.FS

type Theme struct {
	BgColor                            lipgloss.Color
	BgColorStr                         string
	TextFgColor                        lipgloss.Color
	TextFgColorStr                     string
	SubtextFgColor                     lipgloss.Color
	TitleFgColor                       lipgloss.Color
	TitleBgColor                       lipgloss.Color
	CursorFgColor                      lipgloss.Color
	CursorBgColor                      lipgloss.Color
	BorderFgColor                      lipgloss.Color
	StatusbarPrimaryBgColor            lipgloss.Color
	StatusbarPrimaryFgColor            lipgloss.Color
	StatusbarSecondaryFgColor          lipgloss.Color
	StatusbarModePrimaryBgColor        lipgloss.Color
	StatusbarModeSecondaryBgColor      lipgloss.Color
	StatusbarThirdColBgColor           lipgloss.Color
	StatusbarFourthColBgColor          lipgloss.Color
	HttpMethodTextFgColor              lipgloss.Color
	HttpMethodDefaultBgColor           lipgloss.Color
	HttpMethodGetBgColor               lipgloss.Color
	HttpMethodPostBgColor              lipgloss.Color
	HttpMethodPutBgColor               lipgloss.Color
	HttpMethodDeleteBgColor            lipgloss.Color
	HttpMethodPatchBgColor             lipgloss.Color
	UrlFgColor                         lipgloss.Color
	UrlBgColor                         lipgloss.Color
	UrlTemplatedSectionFgColor         lipgloss.Color
	UrlTemplatedSectionBgColor         lipgloss.Color
	UrlUnfilledTemplatedSectionFgColor lipgloss.Color
	UrlUnfilledTemplatedSectionBgColor lipgloss.Color
	WhitespaceFgColor                  lipgloss.Color
	ErrorFgColor                       lipgloss.Color

	ResponseStatus200FgColor lipgloss.Color
	ResponseStatus300FgColor lipgloss.Color
	ResponseStatus400FgColor lipgloss.Color
	ResponseStatus500FgColor lipgloss.Color
	ResponseProtoFgColor     lipgloss.Color
	ResponseHeaderFgColor    lipgloss.Color
}

var config *configuration.Configuration = configuration.New()

var theme *Theme

func loadAdditionalViper(path, themeName string) (*viper.Viper, error) {
	additionalViper := viper.New()
	additionalViper.AddConfigPath(path)
	additionalViper.SetConfigType("yaml")
	additionalViper.SetConfigName(themeName)

	err := additionalViper.ReadInConfig()
	return additionalViper, err
}

func loadThemeByThemeName(themeName string) {
	workspaceViper, err := loadAdditionalViper(config.GetStringOrDefault("workspace"), themeName)
	if err == nil {
		viper.MergeConfigMap(workspaceViper.AllSettings())
	} else {
		home, _ := os.UserHomeDir()
		homeViper, err := loadAdditionalViper(home, themeName)
		if err == nil {
			viper.MergeConfigMap(homeViper.AllSettings())
		} else {
			embeddedTheme, err := embeddedThemes.ReadFile(fmt.Sprintf("themes/%s.yaml", themeName))
			if err != nil {
				log.Error().Msgf("Failed to load embedded theme with name %s", themeName)
			} else {
				themeMap := &map[string]any{}
				err = yaml.Unmarshal(embeddedTheme, themeMap)
				if err != nil {
					log.Error().Msgf("Failed to load unmarshal embedded theme %s", themeName)
				} else {
					viper.MergeConfigMap(*themeMap)
				}
			}
		}
	}
}

func LoadTheme() *Theme {
	if theme == nil {
		getColor := config.GetStringOrDefault
		themeName, found := config.GetString("themeName")
		// Theme is either configured as a pointer to external theme file (with themeName) or directly in config file.
		// Here we check if there is a themeName defined and load the theme from the file.
		if found {
			loadThemeByThemeName(themeName)
		}

		theme = &Theme{
			BgColor:                            lipgloss.Color(getColor("theme.bgColor")),
			BgColorStr:                         getColor("theme.bgColor"),
			TextFgColor:                        lipgloss.Color(getColor("theme.primaryTextFgColor")),
			TextFgColorStr:                     getColor("theme.primaryTextFgColor"),
			SubtextFgColor:                     lipgloss.Color(getColor("theme.secondaryTextFgColor")),
			TitleFgColor:                       lipgloss.Color(getColor("theme.titleFgColor")),
			TitleBgColor:                       lipgloss.Color(getColor("theme.titleBgColor")),
			CursorFgColor:                      lipgloss.Color(getColor("theme.cursorFgColor")),
			CursorBgColor:                      lipgloss.Color(getColor("theme.cursorBgColor")),
			BorderFgColor:                      lipgloss.Color(getColor("theme.borderFgColor")),
			StatusbarPrimaryBgColor:            lipgloss.Color(getColor("theme.statusbar.primaryBgColor")),
			StatusbarPrimaryFgColor:            lipgloss.Color(getColor("theme.statusbar.primaryFgColor")),
			StatusbarSecondaryFgColor:          lipgloss.Color(getColor("theme.statusbar.secondaryFgColor")),
			StatusbarModePrimaryBgColor:        lipgloss.Color(getColor("theme.statusbar.modePrimaryBgColor")),
			StatusbarModeSecondaryBgColor:      lipgloss.Color(getColor("theme.statusbar.modeSecondaryBgColor")),
			StatusbarThirdColBgColor:           lipgloss.Color(getColor("theme.statusbar.thirdColBgColor")),
			StatusbarFourthColBgColor:          lipgloss.Color(getColor("theme.statusbar.fourthColBgColor")),
			HttpMethodTextFgColor:              lipgloss.Color(getColor("theme.httpMethods.textFgColor")),
			HttpMethodDefaultBgColor:           lipgloss.Color(getColor("theme.httpMethods.defaultBgColor")),
			HttpMethodGetBgColor:               lipgloss.Color(getColor("theme.httpMethods.getBgColor")),
			HttpMethodPostBgColor:              lipgloss.Color(getColor("theme.httpMethods.postBgColor")),
			HttpMethodPutBgColor:               lipgloss.Color(getColor("theme.httpMethods.putBgColor")),
			HttpMethodDeleteBgColor:            lipgloss.Color(getColor("theme.httpMethods.deleteBgColor")),
			HttpMethodPatchBgColor:             lipgloss.Color(getColor("theme.httpMethods.patchBgColor")),
			UrlFgColor:                         lipgloss.Color(getColor("theme.urlFgColor")),
			UrlBgColor:                         lipgloss.Color(getColor("theme.urlBgColor")),
			UrlTemplatedSectionFgColor:         lipgloss.Color(getColor("theme.urlTemplatedSectionFgColor")),
			UrlTemplatedSectionBgColor:         lipgloss.Color(getColor("theme.urlTemplatedSectionBgColor")),
			UrlUnfilledTemplatedSectionFgColor: lipgloss.Color(getColor("theme.urlUnfilledTemplatedSectionFgColor")),
			UrlUnfilledTemplatedSectionBgColor: lipgloss.Color(getColor("theme.urlUnfilledTemplatedSectionBgColor")),
			WhitespaceFgColor:                  lipgloss.Color(getColor("theme.whitespaceFgColor")),
			ErrorFgColor:                       lipgloss.Color(getColor("theme.errorFgColor")),
			ResponseStatus200FgColor:           lipgloss.Color(getColor("theme.response.status200FgColor")),
			ResponseStatus300FgColor:           lipgloss.Color(getColor("theme.response.status300FgColor")),
			ResponseStatus400FgColor:           lipgloss.Color(getColor("theme.response.status400FgColor")),
			ResponseStatus500FgColor:           lipgloss.Color(getColor("theme.response.status500FgColor")),
			ResponseProtoFgColor:               lipgloss.Color(getColor("theme.response.protoFgColor")),
			ResponseHeaderFgColor:              lipgloss.Color(getColor("theme.response.headerFgColor")),
		}
	}

	return theme
}

func (t *Theme) HttpMethodBgColor(method string) lipgloss.Color {
	switch method {
	case "GET":
		return theme.HttpMethodGetBgColor
	case "POST":
		return theme.HttpMethodPostBgColor
	case "PUT":
		return theme.HttpMethodPutBgColor
	case "DELETE":
		return theme.HttpMethodDeleteBgColor
	case "PATCH":
		return theme.HttpMethodPatchBgColor
	default:
		return theme.HttpMethodDefaultBgColor
	}
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
			HelpPaneStyle:      lipgloss.NewStyle().Foreground(theme.TextFgColor).BorderForeground(theme.BorderFgColor).Padding(1).Border(lipgloss.RoundedBorder()),
			HelpKeyStyle:       lipgloss.NewStyle().Foreground(theme.TextFgColor),
			HelpDescStyle:      lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true),
			HelpSeparatorStyle: lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true),
		}
	}
	return commonStyle
}
