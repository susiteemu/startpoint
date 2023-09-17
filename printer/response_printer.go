package printer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

func SprintResponse(resp *http.Response) (string, error) {
	return sprintResponse(resp, false)
}

func SprintPrettyResponse(resp *http.Response) (string, error) {
	return sprintResponse(resp, true)
}

func sprintResponse(resp *http.Response, pretty bool) (string, error) {
	resp_str := ""

	resp_status_str, err := SprintStatus(resp)
	if err != nil {
		return "", err
	}
	resp_str += resp_status_str + "\n"

	resp_headers_str, err := SprintHeaders(resp)
	if err != nil {
		return "", err
	}
	resp_str += resp_headers_str

	resp_body_str, err := SprintBody(resp)
	if err != nil {
		return "", err
	}

	if pretty {
		buf := new(bytes.Buffer)
		// TODO content-type
		err = quick.Highlight(buf, resp_body_str, "json", "terminal16m", "catppuccin-mocha")
		if err != nil {
			return "", err
		}
		resp_str += buf.String()
	} else {
		resp_str += resp_body_str
	}

	return resp_str, nil
}

func SprintStatus(resp *http.Response) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}
	return fmt.Sprintf("%v %v", resp.Proto, resp.Status), nil
}

func SprintHeaders(resp *http.Response) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}

	resp_headers := resp.Header
	resp_headers_str := ""
	// sort header names
	resp_header_names := sortHeaderNames(resp_headers)
	for _, resp_header := range resp_header_names {
		resp_header_vals := strings.Join(resp_headers.Values(resp_header), ", ")
		resp_headers_str += fmt.Sprintf("%v: %v", resp_header, resp_header_vals)
		resp_headers_str += fmt.Sprintln("")
	}

	return resp_headers_str, nil
}

func SprintBody(resp *http.Response) (string, error) {
	resp_body_str := ""
	if resp.ContentLength > 0 {
		defer resp.Body.Close()
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		bodyStr := string(resp_body[:])
		resp_body_str += fmt.Sprint(bodyStr)
	}
	return resp_body_str, nil
}

func sortHeaderNames(headers http.Header) []string {
	sorted_headers := make([]string, 0, len(headers))
	for k := range headers {
		sorted_headers = append(sorted_headers, k)
	}
	sort.Strings(sorted_headers)
	return sorted_headers
}
