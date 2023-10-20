package printer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
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

	if printBody {
		resp_body_str, err := SprintBody(resp)
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

		dispatcher := NewBodyFormatter(&JsonContentTypeBodyHandler{}, &XmlContentTypeBodyHandler{}, &DefaultContentTypeBodyHandler{})

		contentType, err := getContentType(resp.Header)
		if err != nil {
			resp_body_str = string(resp_body[:])
		} else {
			resp_body_str, _ = dispatcher.Format(contentType, resp_body)
		}
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

type BodyFormatHandler interface {
	Supports(contentType string) bool
	Handle(body []byte) (string, error)
}

type JsonContentTypeBodyHandler struct{}

func (h *JsonContentTypeBodyHandler) Supports(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), "application/json")
}

func (h *JsonContentTypeBodyHandler) Handle(body []byte) (string, error) {
	var pretty_json bytes.Buffer
	err := json.Indent(&pretty_json, body[:], "", "    ")
	if err != nil {
		return "", err
	}
	return string(pretty_json.Bytes()), nil
}

type XmlContentTypeBodyHandler struct{}

func (h *XmlContentTypeBodyHandler) Supports(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), "application/xml")
}

func (h *XmlContentTypeBodyHandler) Handle(body []byte) (string, error) {
	// TODO actual indentation
	return string(body), nil
}

type DefaultContentTypeBodyHandler struct{}

func (h *DefaultContentTypeBodyHandler) Supports(contentType string) bool {
	return true
}

func (h *DefaultContentTypeBodyHandler) Handle(body []byte) (string, error) {
	// TODO actual indentation
	return string(body), nil
}

type BodyFormatter struct {
	handlers []BodyFormatHandler
}

func NewBodyFormatter(handlers ...BodyFormatHandler) *BodyFormatter {
	return &BodyFormatter{handlers: handlers}
}

func (d *BodyFormatter) Format(contentType string, body []byte) (string, error) {
	for _, handler := range d.handlers {
		if handler.Supports(contentType) {
			return handler.Handle(body)
		}
	}
	return "", errors.New("No handler found for the given content-type")
}
