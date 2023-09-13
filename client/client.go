package client

import (
	"bytes"
	"net/http"
)

func DoRequest(url string, method string, headers map[string]string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, bodyReader)
	for key, val := range headers {
		req.Header.Add(key, val)
	}
	// foo bar
	if err != nil {
		// TODO log error
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		// TODO log error
		return nil, err
	}

	return resp, nil

}
