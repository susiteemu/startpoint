package validator

import (
	"net/http"
	"net/url"
	"slices"
)

var ValidMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

var DefaultBodilessMethod = http.MethodGet
var DefaultBodifulMethod = http.MethodPost

func IsValidMethod(rawMethod string) bool {
	if rawMethod == "" {
		return false
	}

	return slices.Contains(ValidMethods, rawMethod)
}

func IsValidUrl(rawUrl string) bool {
	_, err := url.ParseRequestURI(rawUrl)
	return err == nil
}
