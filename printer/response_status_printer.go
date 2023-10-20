package printer

import (
	"errors"
	"fmt"
	"net/http"
)

func SprintStatus(resp *http.Response) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}
	return fmt.Sprintf("%v %v", resp.Proto, resp.Status), nil
}
