package printer

import (
	"bytes"
	"errors"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"net/http"
)

func SprintFullResponse(resp *http.Response) (string, error) {
	return sprintResponse(resp, false, true, true)
}

func SprintPrettyFullResponse(resp *http.Response) (string, error) {
	return sprintResponse(resp, true, true, true)
}

func SprintResponse(resp *http.Response, printHeaders bool, printBody bool) (string, error) {
	return sprintResponse(resp, false, printHeaders, printBody)
}

func SprintPrettyResponse(resp *http.Response, printHeaders bool, printBody bool) (string, error) {
	return sprintResponse(resp, true, printHeaders, printBody)
}

func sprintResponse(resp *http.Response, pretty bool, printHeaders bool, printBody bool) (string, error) {
	respStr := ""

	if printHeaders {
		respStatusStr, err := SprintStatus(resp)
		if err != nil {
			return "", err
		}
		respHeadersStr, err := SprintHeaders(resp)
		if err != nil {
			return "", err
		}
		respStr += respStatusStr + "\n" + respHeadersStr
	}

	if printBody {
		respBodyStr, err := SprintBody(resp)
		if err != nil {
			return "", err
		}
		if printHeaders {
			respStr += "\n"
		}
		respStr += respBodyStr
	}

	if pretty {
		buf := new(bytes.Buffer)

		lexer := resolveLexer(resp, printHeaders, printBody)
		style := resolveStyle()
		formatter := resolveFormatter()
		iterator, err := lexer.Tokenise(nil, respStr)
		if err != nil {
			return "", err
		}
		err = formatter.Format(buf, style, iterator)
		if err != nil {
			return "", err
		}

		respStr = buf.String()
	}

	return respStr, nil
}

func getContentType(headers http.Header) (string, error) {
	for k := range headers {
		if k == "Content-Type" {
			return headers.Get(k), nil
		}
	}
	return "", errors.New("Could not find Content-Type!")
}

func resolveLexer(resp *http.Response, printHeaders bool, printBody bool) chroma.Lexer {
	var lexer chroma.Lexer
	if printHeaders || (!printBody) {
		lexer = lexers.Get("http")
	} else if printBody {
		contentType, err := getContentType(resp.Header)
		if err != nil {
			lexer = lexers.Fallback
		} else {
			lexer = lexers.MatchMimeType(contentType)
		}
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)
	return lexer
}

func resolveStyle() *chroma.Style {
	style := styles.Get("catppuccin-mocha")
	if style == nil {
		style = styles.Fallback
	}
	return style
}

func resolveFormatter() chroma.Formatter {
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	return formatter
}
