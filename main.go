package main

import (
	logrus "github.com/sirupsen/logrus"
	"goful-cli/printer"
	"goful-cli/request"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	var url = "https://httpbin.org/anything"
	var headers = map[string]string{"X-Foo": "bar", "X-Bar": "foo"}
	var body = []byte("{\"foo\":\"hello\"}")

	logrus.Debugf("Request url=%s, headers=%v, body=%v", url, headers, body)

	resp, err := request.DoRequest(url, "POST", headers, body)

	if err != nil {
		logrus.Errorf("Error occurred while performing request %v", err)
	}

	printer.PrintResponse(resp)
}
