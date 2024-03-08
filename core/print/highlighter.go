package print

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/spf13/viper"
)

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
	style := styles.Get(viper.GetString("theme.syntax"))
	if style == nil {
		style = styles.Fallback
	}
	return style
}

func resolveFormatter() chroma.Formatter {
	formatter := formatters.Get(viper.GetString("printer.response.formatter"))
	if formatter == nil {
		formatter = formatters.Fallback
	}
	return formatter
}
