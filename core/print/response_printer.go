package print

import (
	"strings"

	"github.com/susiteemu/startpoint/core/model"
)

type PrintOpts struct {
	PrettyPrint    bool
	PrintHeaders   bool
	PrintBody      bool
	PrintTraceInfo bool
	PrintRequest   bool
}

func SprintResponse(resp *model.Response, printOpts PrintOpts) (string, string, error) {
	return sprintResponse(resp, printOpts)
}

func SprintFullResponse(resp *model.Response) (string, string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    false,
		PrintHeaders:   true,
		PrintBody:      true,
		PrintTraceInfo: true,
		PrintRequest:   true,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPrettyFullResponse(resp *model.Response) (string, string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    true,
		PrintHeaders:   true,
		PrintBody:      true,
		PrintTraceInfo: true,
		PrintRequest:   true,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPlainResponse(resp *model.Response, printHeaders bool, printBody bool) (string, string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    false,
		PrintHeaders:   printHeaders,
		PrintBody:      printBody,
		PrintTraceInfo: false,
		PrintRequest:   false,
	}
	return sprintResponse(resp, printOpts)
}

func SprintPrettyResponse(resp *model.Response, printHeaders bool, printBody bool) (string, string, error) {
	printOpts := PrintOpts{
		PrettyPrint:    true,
		PrintHeaders:   printHeaders,
		PrintBody:      printBody,
		PrintTraceInfo: false,
		PrintRequest:   false,
	}
	return sprintResponse(resp, printOpts)
}

func sprintResponse(resp *model.Response, printOpts PrintOpts) (string, string, error) {
	pretty := printOpts.PrettyPrint
	var responseBuilder, prettyResponseBuilder []string

	if printOpts.PrintTraceInfo {
		traceInfo, prettyTraceInfo, _ := SprintTraceInfo(resp.TraceInfo, pretty)
		responseBuilder = append(responseBuilder, traceInfo)
		if printOpts.PrettyPrint {
			prettyResponseBuilder = append(prettyResponseBuilder, prettyTraceInfo)
		}
	}

	if printOpts.PrintRequest {
		request, prettyRequest, err := SprintRequest(&resp.Request, pretty)
		if err != nil {
			return "", "", err
		}
		if len(request) > 0 {
			responseBuilder = append(responseBuilder, request)
		}
		if len(prettyRequest) > 0 && printOpts.PrettyPrint {
			prettyResponseBuilder = append(prettyResponseBuilder, prettyRequest)
		}
	}

	if printOpts.PrintHeaders {
		respStatusStr, prettyRespStatusStr, err := SprintStatus(resp, pretty)
		if err != nil {
			return "", "", err
		}
		responseBuilder = append(responseBuilder, respStatusStr)
		if printOpts.PrettyPrint {
			prettyResponseBuilder = append(prettyResponseBuilder, prettyRespStatusStr)
		}
		respHeadersStr, prettyRespHeadersStr, err := SprintHeaders(resp.Headers, pretty)
		if err != nil {
			return "", "", err
		}
		if len(respHeadersStr) > 0 {
			responseBuilder = append(responseBuilder, respHeadersStr)
		}
		if len(prettyRespHeadersStr) > 0 && printOpts.PrettyPrint {
			prettyResponseBuilder = append(prettyResponseBuilder, prettyRespHeadersStr)
		}
	}

	if printOpts.PrintBody {
		respBodyStr, prettyRespBodyStr, err := SprintBody(resp.Size, resp.Body, resp.Headers, pretty)
		if err != nil {
			return "", "", err
		}

		if len(respBodyStr) > 0 {
			responseBuilder = append(responseBuilder, respBodyStr)
		}
		if len(prettyRespBodyStr) > 0 && printOpts.PrettyPrint {
			prettyResponseBuilder = append(prettyResponseBuilder, prettyRespBodyStr)
		}
	}

	return strings.Join(responseBuilder, "\n"), strings.Join(prettyResponseBuilder, "\n"), nil
}
