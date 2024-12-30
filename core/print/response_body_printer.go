package print

import (
	"bytes"
	"encoding/json"
	"errors"
	"startpoint/core/model"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/rs/zerolog/log"
	"github.com/yosssi/gohtml"
)

func SprintBody(size int64, body []byte, headers model.Headers, pretty bool) (string, error) {
	respBodyStr := ""
	if size > 0 {
		dispatcher := NewBodyFormatter(&JsonContentTypeBodyHandler{}, &XmlContentTypeBodyHandler{}, &HtmlContentTypeBodyHandler{}, &DefaultContentTypeBodyHandler{})

		contentType, err := headers.ContentType()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get content type")
			respBodyStr = string(body)
		} else {
			respBodyStr, err = dispatcher.Format(contentType, body)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to format body")
				respBodyStr = string(body)
			}
		}
		if pretty && len(respBodyStr) > 0 {
			respBodyStr, err = prettyPrintBody(respBodyStr, headers)
			if err != nil {
				return "", err
			}
		}

	}
	return respBodyStr, nil
}

func prettyPrintBody(respBodyStr string, headers model.Headers) (string, error) {
	buf := new(bytes.Buffer)
	lexer := resolveBodyLexer(headers)
	style := resolveStyle()
	formatter := resolveFormatter()
	iterator, err := lexer.Tokenise(nil, respBodyStr)
	if err != nil {
		return "", err
	}
	err = formatter.Format(buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func resolveBodyLexer(headers model.Headers) chroma.Lexer {
	var lexer chroma.Lexer
	contentType, err := headers.ContentType()
	if err != nil {
		lexer = lexers.Fallback
		log.Warn().Err(err).Msgf("Failed to get content type: using fallback lexer %v", lexer)
	} else {
		if contentType == model.CONTENT_TYPE_PLAINTEXT {
			lexer = lexers.Get("plaintext")
		} else {
			lexer = lexers.MatchMimeType(contentType)
		}
		log.Debug().Msgf("Matched mimetype %s with lexer %v", contentType, lexer)
	}

	if lexer == nil {
		lexer = lexers.Fallback
		log.Debug().Msgf("Using fallback lexer %v", lexer)
	}

	lexer = chroma.Coalesce(lexer)
	return lexer
}

type BodyFormatHandler interface {
	Supports(contentType string) bool
	Handle(body []byte) (string, error)
}

type JsonContentTypeBodyHandler struct{}

func (h *JsonContentTypeBodyHandler) Supports(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), model.CONTENT_TYPE_APPLICATION_JSON)
}

func (h *JsonContentTypeBodyHandler) Handle(body []byte) (string, error) {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, body, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to indent json")
		return "", err
	}
	return prettyJson.String(), nil
}

type XmlContentTypeBodyHandler struct{}

func (h *XmlContentTypeBodyHandler) Supports(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), model.CONTENT_TYPE_APPLICATION_XML)
}

func (h *XmlContentTypeBodyHandler) Handle(body []byte) (string, error) {
	gohtml.Condense = true
	return string(gohtml.FormatBytes(body)), nil
}

type HtmlContentTypeBodyHandler struct{}

func (h *HtmlContentTypeBodyHandler) Supports(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), model.CONTENT_TYPE_TEXT_HTML)
}

func (h *HtmlContentTypeBodyHandler) Handle(body []byte) (string, error) {
	return string(gohtml.FormatBytes(body)), nil
}

type DefaultContentTypeBodyHandler struct{}

func (h *DefaultContentTypeBodyHandler) Supports(contentType string) bool {
	return true
}

func (h *DefaultContentTypeBodyHandler) Handle(body []byte) (string, error) {
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
