package print

import (
	"bytes"
	"errors"
	"startpoint/core/model"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/rs/zerolog/log"
)

type PrintOpts struct {
	PrettyPrint    bool
	PrintHeaders   bool
	PrintBody      bool
	PrintTraceInfo bool
}

func SprintResponse(resp *model.Response, printOpts PrintOpts) (string, error) {
	return sprintResponse(resp, printOpts)
}

func SprintFullResponse(resp *model.Response) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    false,
		PrintHeaders:   true,
		PrintBody:      true,
		PrintTraceInfo: true,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPrettyFullResponse(resp *model.Response) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    true,
		PrintHeaders:   true,
		PrintBody:      true,
		PrintTraceInfo: true,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPlainResponse(resp *model.Response, printHeaders bool, printBody bool) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    false,
		PrintHeaders:   printHeaders,
		PrintBody:      printBody,
		PrintTraceInfo: false,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPrettyResponse(resp *model.Response, printHeaders bool, printBody bool) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    true,
		PrintHeaders:   printHeaders,
		PrintBody:      printBody,
		PrintTraceInfo: false,
	}
	return sprintResponse(resp, printOpts)
}

func sprintResponse(resp *model.Response, printOpts PrintOpts) (string, error) {
	pretty := printOpts.PrettyPrint
	respStr := ""

	if printOpts.PrintTraceInfo {
		traceInfo, _ := SprintTraceInfo(resp.TraceInfo, pretty)
		respStr += traceInfo + "\n"
	}

	if printOpts.PrintHeaders {
		respStatusStr, err := SprintStatus(resp, pretty)
		if err != nil {
			return "", err
		}
		respHeadersStr, err := SprintHeaders(resp, pretty)
		if err != nil {
			return "", err
		}
		respStr += respStatusStr + "\n" + respHeadersStr
	}

	if printOpts.PrintBody {
		respBodyStr, err := SprintBody(resp)
		if err != nil {
			return "", err
		}
		if printOpts.PrintHeaders {
			respStr += "\n"
		}

		if pretty {
			respBodyStr, err = prettyPrintBody(respBodyStr, resp)
			if err != nil {
				return "", err
			}
		}

		respStr += respBodyStr
	}

	return respStr, nil
}

func getContentType(headers map[string]model.HeaderValues) (string, error) {
	for k := range headers {
		if k == "Content-Type" {
			contentType := headers[k][0]
			if strings.Contains(contentType, ";") {
				contentType = strings.Split(contentType, ";")[0]
			}
			return contentType, nil
		}
	}
	return "", errors.New("could not find Content-Type")
}

func prettyPrintBody(respBodyStr string, resp *model.Response) (string, error) {
	buf := new(bytes.Buffer)
	lexer := resolveBodyLexer(resp)
	style := resolveStyle()
	formatter := resolveFormatter()
	iterator, err := lexer.Tokenise(nil, respBodyStr)
	if err != nil {
		return "", err
	}
	err = formatter.Format(buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func resolveBodyLexer(resp *model.Response) chroma.Lexer {
	var lexer chroma.Lexer
	contentType, err := getContentType(resp.Headers)
	log.Debug().Msgf("Content-type %s", contentType)
	if err != nil {
		lexer = lexers.Fallback
	} else {
		lexer = lexers.MatchMimeType(contentType)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)
	return lexer
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
