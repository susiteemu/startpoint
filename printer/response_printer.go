package printer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	//	"github.com/alecthomas/chroma/v2/quick"
	"github.com/alecthomas/chroma/v2/styles"
)

func SprintFullResponse(resp *http.Response) (string, error) {
	return sprintResponse(resp, false, true, true)
}

func SprintPrettyFullResponse(resp *http.Response) (string, error) {
	return sprintResponse(resp, true, true, true)
}

func SprintResponse(resp *http.Response, printHeaders bool, printBody bool) (string, error) {
	return sprintResponse(resp, false, printHeaders, printBody)
}

func SprintPrettyResponse(resp *http.Response, printHeaders bool, printBody bool) (string, error) {
	return sprintResponse(resp, true, printHeaders, printBody)
}

func sprintResponse(resp *http.Response, pretty bool, printHeaders bool, printBody bool) (string, error) {
	resp_str := ""

	if printHeaders {
		resp_status_str, err := SprintStatus(resp)
		if err != nil {
			return "", err
		}
		resp_headers_str, err := SprintHeaders(resp)
		if err != nil {
			return "", err
		}
		resp_str += resp_status_str + "\n" + resp_headers_str
	}

	var resp_body_str string
	if printBody {
		var err error
		resp_body_str, err = SprintBody(resp)
		if err != nil {
			return "", err
		}
		if printHeaders {
			resp_str += "\n"
		}
		resp_str += resp_body_str
	}

	if pretty {
		buf := new(bytes.Buffer)

		lexer := resolveLexer(resp, printHeaders, printBody)
		style := resolveStyle()
		formatter := resolveFormatter()
		iterator, err := lexer.Tokenise(nil, resp_str)
		if err != nil {
			return "", err
		}
		err = formatter.Format(buf, style, iterator)
		if err != nil {
			return "", err
		}

		resp_str = buf.String()
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

func getContentType(headers http.Header) (string, error) {
	for k := range headers {
		if k == "Content-Type" {
			return headers.Get(k), nil
		}
	}
	return "", errors.New("Could not find Content-Type!")
}

func resolveLexer(resp *http.Response, printHeaders bool, printBody bool) chroma.Lexer {
	var lexer chroma.Lexer
	if printHeaders || (!printBody) {
		lexer = lexers.Get("http")
	} else if printBody {
		contentType, err := getContentType(resp.Header)
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

func resolveStyle() *chroma.Style {
	style := styles.Get("catppuccin-mocha")
	if style == nil {
		style = styles.Fallback
	}
	return style
}

func resolveFormatter() chroma.Formatter {
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	return formatter
}
