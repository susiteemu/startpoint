package print

import (
	"errors"
	"fmt"
	"goful/core/model"
	"sort"
	"strings"
)

func SprintHeaders(resp *model.Response) (string, error) {
	if resp == nil {
		return "", errors.New("response must not be nil")
	}

	respHeaders := resp.Headers
	respHeadersStr := ""
	// sort header names
	respHeaderNames := sortHeaderNames(resp.Headers)
	for _, respHeader := range respHeaderNames {
		respHeaderVals := strings.Join(respHeaders[respHeader], ", ")
		respHeadersStr += fmt.Sprintf("%v: %v", respHeader, respHeaderVals)
		respHeadersStr += fmt.Sprintln("")
	}

	return respHeadersStr, nil
}

func sortHeaderNames(headers map[string]model.HeaderValues) []string {
	sortedHeaders := make([]string, 0, len(headers))
	for k := range headers {
		sortedHeaders = append(sortedHeaders, k)
	}
	sort.Strings(sortedHeaders)
	return sortedHeaders
}
