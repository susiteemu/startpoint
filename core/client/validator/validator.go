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
	u, err := url.Parse(rawUrl)
	return err == nil && u.Scheme != "" && u.Host != ""
}
