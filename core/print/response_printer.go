package print

import (
	"github.com/susiteemu/startpoint/core/model"
	"strings"
)

type PrintOpts struct {
	PrettyPrint    bool
	PrintHeaders   bool
	PrintBody      bool
	PrintTraceInfo bool
	PrintRequest   bool
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
		PrintRequest:   true,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPrettyFullResponse(resp *model.Response) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    true,
		PrintHeaders:   true,
		PrintBody:      true,
		PrintTraceInfo: true,
		PrintRequest:   true,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPlainResponse(resp *model.Response, printHeaders bool, printBody bool) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    false,
		PrintHeaders:   printHeaders,
		PrintBody:      printBody,
		PrintTraceInfo: false,
		PrintRequest:   false,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPrettyResponse(resp *model.Response, printHeaders bool, printBody bool) (string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    true,
		PrintHeaders:   printHeaders,
		PrintBody:      printBody,
		PrintTraceInfo: false,
		PrintRequest:   false,
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

	if printOpts.PrintRequest {
		request, err := SprintRequest(&resp.Request, pretty)
		if err != nil {
			return "", err
		}
		if len(request) > 0 {
			responseBuilder = append(responseBuilder, request)
		}
	}

	if printOpts.PrintHeaders {
		respStatusStr, err := SprintStatus(resp, pretty)
		if err != nil {
			return "", err
		}
		responseBuilder = append(responseBuilder, respStatusStr)
		respHeadersStr, err := SprintHeaders(resp.Headers, pretty)
		if err != nil {
			return "", err
		}
		if len(respHeadersStr) > 0 {
			responseBuilder = append(responseBuilder, respHeadersStr)
		}
	}

	if printOpts.PrintBody {
		respBodyStr, err := SprintBody(resp.Size, resp.Body, resp.Headers, pretty)
		if err != nil {
			return "", err
		}

		if len(respBodyStr) > 0 {
			responseBuilder = append(responseBuilder, respBodyStr)
		}
	}

	return strings.Join(responseBuilder, "\n"), nil
}
