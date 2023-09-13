package printer

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

func PrintResponse(resp *http.Response) (string, error) {
	resp_str := fmt.Sprintf("%v %v", resp.Proto, resp.Status)
	resp_str += fmt.Sprintln("")

	resp_headers := resp.Header

	// sort header names
	resp_header_names := make([]string, 0, len(resp_headers))
	for k := range resp_headers {
		resp_header_names = append(resp_header_names, k)
	}
	sort.Strings(resp_header_names)

	for _, resp_header := range resp_header_names {
		resp_header_vals := strings.Join(resp_headers.Values(resp_header), ", ")
		resp_str += fmt.Sprintf("%v: %v", resp_header, resp_header_vals)
		resp_str += fmt.Sprintln("")
	}

	if resp.ContentLength > 0 {
		defer resp.Body.Close()
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Error occurred while reading body", err)
			return "", err
		}

		bodyStr := string(resp_body[:])
		resp_str += fmt.Sprint(bodyStr)

		return resp_str, nil
	}
	return resp_str, nil
}
