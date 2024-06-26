package print

import (
	"errors"
	"fmt"
	"sort"
	"startpoint/core/model"
	"startpoint/tui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SprintHeaders(resp *model.Response, pretty bool) (string, error) {
	if resp == nil {
		return "", errors.New("response must not be nil")
	}

	theme := styles.GetTheme()
	headerStyle := lipgloss.NewStyle().Foreground(theme.ResponseHeaderFgColor)

	respHeaders := resp.Headers
	respHeadersStr := ""
	// sort header names
	respHeaderNames := sortHeaderNames(resp.Headers)
	for _, respHeader := range respHeaderNames {
		header := respHeader
		headerValues := strings.Join(respHeaders[respHeader], ", ")
		if pretty {
			header = headerStyle.Render(header)
		}
		respHeadersStr += fmt.Sprintf("%v: %v", header, headerValues)
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
