package print

import (
	"bytes"
	"errors"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"goful/core/model"
)

func SprintFullResponse(resp *model.Response) (string, error) {
	return sprintResponse(resp, false, true, true)
}

func SprintPrettyFullResponse(resp *model.Response) (string, error) {
	return sprintResponse(resp, true, true, true)
}

func SprintPlainResponse(resp *model.Response, printHeaders bool, printBody bool) (string, error) {
	return sprintResponse(resp, false, printHeaders, printBody)
}

func SprintPrettyResponse(resp *model.Response, printHeaders bool, printBody bool) (string, error) {
	return sprintResponse(resp, true, printHeaders, printBody)
}

func sprintResponse(resp *model.Response, pretty bool, printHeaders bool, printBody bool) (string, error) {
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
		// FIXME tokenise different parts separately? Now body content-type dictates whole response formatting style
		lexer := resolveResponseLexer(resp, printHeaders, printBody)
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

func getContentType(headers map[string]model.HeaderValues) (string, error) {
	for k := range headers {
		if k == "Content-Type" {
			return headers[k][0], nil
		}
	}
	return "", errors.New("could not find Content-Type")
}

func resolveResponseLexer(resp *model.Response, printHeaders bool, printBody bool) chroma.Lexer {
	var lexer chroma.Lexer
	if printHeaders || (!printBody) {
		lexer = lexers.Get("http")
	} else if printBody {
		contentType, err := getContentType(resp.Headers)
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
