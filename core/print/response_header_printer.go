package print

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/susiteemu/startpoint/core/model"
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func SprintHeaders(headers model.Headers, pretty bool) (string, string, error) {
	if headers == nil {
		return "", "", errors.New("headers must not be nil")
	}

	theme := styles.LoadTheme()
	headerStyle := lipgloss.NewStyle().Foreground(theme.ResponseHeaderFgColor)

	var respHeaderNamesBuilder, prettyRespHeaderNamesBuilder []string
	// sort header names
	respHeaderNames := sortHeaderNames(headers)
	for _, respHeader := range respHeaderNames {
		header := respHeader
		headerValues := strings.Join(headers[respHeader], ", ")
		if pretty {
			prettyHeader := headerStyle.Render(header)
			prettyRespHeaderNamesBuilder = append(prettyRespHeaderNamesBuilder, fmt.Sprintf("%v: %v", prettyHeader, headerValues))
		}
		respHeaderNamesBuilder = append(respHeaderNamesBuilder, fmt.Sprintf("%v: %v", header, headerValues))
	}

	return strings.Join(respHeaderNamesBuilder, "\n"), strings.Join(prettyRespHeaderNamesBuilder, "\n"), nil
}

func sortHeaderNames(headers model.Headers) []string {
	sortedHeaders := make([]string, 0, len(headers))
	for k := range headers {
		sortedHeaders = append(sortedHeaders, k)
	}
	sort.Strings(sortedHeaders)
	return sortedHeaders
}
