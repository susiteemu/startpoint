package print

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func SprintHeaders(resp *http.Response) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}

	respHeaders := resp.Header
	respHeadersStr := ""
	// sort header names
	respHeaderNames := sortHeaderNames(respHeaders)
	for _, respHeader := range respHeaderNames {
		respHeaderVals := strings.Join(respHeaders.Values(respHeader), ", ")
		respHeadersStr += fmt.Sprintf("%v: %v", respHeader, respHeaderVals)
		respHeadersStr += fmt.Sprintln("")
	}

	return respHeadersStr, nil
}

func sortHeaderNames(headers http.Header) []string {
	sortedHeaders := make([]string, 0, len(headers))
	for k := range headers {
		sortedHeaders = append(sortedHeaders, k)
	}
	sort.Strings(sortedHeaders)
	return sortedHeaders
}
