package print

import (
	"errors"
	"fmt"
	"github.com/susiteemu/startpoint/core/model"
	"github.com/susiteemu/startpoint/tui/styles"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SprintHeaders(headers model.Headers, pretty bool) (string, error) {
	if headers == nil {
		return "", errors.New("headers must not be nil")
	}

	theme := styles.GetTheme()
	headerStyle := lipgloss.NewStyle().Foreground(theme.ResponseHeaderFgColor)

	respHeadersStr := ""
	// sort header names
	respHeaderNames := sortHeaderNames(headers)
	for _, respHeader := range respHeaderNames {
		header := respHeader
		headerValues := strings.Join(headers[respHeader], ", ")
		if pretty {
			header = headerStyle.Render(header)
		}
		respHeadersStr += fmt.Sprintf("%v: %v", header, headerValues)
		respHeadersStr += fmt.Sprintln("")
	}

	return respHeadersStr, nil
}

func sortHeaderNames(headers model.Headers) []string {
	sortedHeaders := make([]string, 0, len(headers))
	for k := range headers {
		sortedHeaders = append(sortedHeaders, k)
	}
	sort.Strings(sortedHeaders)
	return sortedHeaders
}
