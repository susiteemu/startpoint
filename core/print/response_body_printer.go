package print

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func SprintBody(resp *http.Response) (string, error) {
	respBodyStr := ""
	if resp.ContentLength > 0 {
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		dispatcher := NewBodyFormatter(&JsonContentTypeBodyHandler{}, &XmlContentTypeBodyHandler{}, &DefaultContentTypeBodyHandler{})

		contentType, err := getContentType(resp.Header)
		if err != nil {
			respBodyStr = string(respBody[:])
		} else {
			respBodyStr, _ = dispatcher.Format(contentType, respBody)
		}
	}
	return respBodyStr, nil
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
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, body[:], "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJson.Bytes()), nil
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
	return "", errors.New("no handler found for the given content-type")
}
