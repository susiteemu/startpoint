package print

import (
	"fmt"
	"regexp"
	"startpoint/core/configuration"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

var config *configuration.Configuration = configuration.New()

func resolveLexer(contentType string) chroma.Lexer {
	var lexer chroma.Lexer
	lexer = lexers.MatchMimeType(contentType)

	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)
	return lexer
}

func resolveStyle() *chroma.Style {
	style := styles.Get(config.GetStringOrDefault("theme.syntax"))
	if style == nil {
		style = styles.Fallback
	}
	return style
}

func resolveFormatter() chroma.Formatter {
	formatter := formatters.Get(config.GetStringOrDefault("printer.response.formatter"))
	if formatter == nil {
		formatter = formatters.Fallback
	}
	return formatter
}

func HighlightWithRegex(text string, pattern string, baseFg lipgloss.Color, baseBg lipgloss.Color, highlightFg lipgloss.Color, highlightBg lipgloss.Color) string {
	var coloredText string
	baseStyle := lipgloss.NewStyle().Foreground(baseFg).Background(baseBg)
	highlightStyle := lipgloss.NewStyle().Foreground(highlightFg).Background(highlightBg)
	regexpPattern := regexp.MustCompile(pattern)
	matches := regexpPattern.FindAllStringIndex(text, -1)
	if len(matches) > 0 {
		cursor := 0
		for _, group := range matches {
			if len(group) < 2 {
				log.Error().Msgf("Expected to have two items, instead got %v", group)
				continue
			}
			startIndex := group[0]
			endIndex := group[1]

			before := text[cursor:startIndex]
			matched := text[startIndex:endIndex]
			coloredText += fmt.Sprintf("%s%s", baseStyle.Render(before), highlightStyle.Render(matched))
			cursor = endIndex
		}
		if cursor < len(text) {
			coloredText += baseStyle.Render(text[cursor:])
		}
	} else {
		coloredText = baseStyle.Render(text)
	}
	return coloredText
}
