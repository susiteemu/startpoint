package printer

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

func PrintResponse(resp *http.Response) {
	fmt.Printf("%v %v", resp.Proto, resp.Status)
	fmt.Println("")

	resp_headers := resp.Header

	// sort header names
	resp_header_names := make([]string, 0, len(resp_headers))
	for k := range resp_headers {
		resp_header_names = append(resp_header_names, k)
	}
	sort.Strings(resp_header_names)

	for _, resp_header := range resp_header_names {
		resp_header_vals := strings.Join(resp_headers.Values(resp_header), ", ")
		fmt.Printf("%v: %v", resp_header, resp_header_vals)
		fmt.Println("")
	}

	if resp.ContentLength > 0 {
		defer resp.Body.Close()
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Error occurred while reading body", err)
		}

		bodyStr := string(resp_body[:])
		fmt.Print(bodyStr)
	}
}
