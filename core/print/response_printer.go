package print

import (
	"errors"
	"startpoint/core/model"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
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
	var responseBuilder []string

	if printOpts.PrintTraceInfo {
		traceInfo, _ := SprintTraceInfo(resp.TraceInfo, pretty)
		responseBuilder = append(responseBuilder, traceInfo)
	}

	if printOpts.PrintHeaders {
		respStatusStr, err := SprintStatus(resp, pretty)
		if err != nil {
			return "", err
		}
		responseBuilder = append(responseBuilder, respStatusStr)
		respHeadersStr, err := SprintHeaders(resp, pretty)
		if err != nil {
			return "", err
		}
		if len(respHeadersStr) > 0 {
			responseBuilder = append(responseBuilder, respHeadersStr)
		}
	}

	if printOpts.PrintBody {
		respBodyStr, err := SprintBody(resp, pretty)
		if err != nil {
			return "", err
		}

		if len(respBodyStr) > 0 {
			responseBuilder = append(responseBuilder, respBodyStr)
		}
	}

	return strings.Join(responseBuilder, "\n"), nil
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
