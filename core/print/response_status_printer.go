package print

import (
	"errors"
	"fmt"
	"goful/core/model"
)

func SprintStatus(resp *model.Response) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}
	return fmt.Sprintf("%v %v", resp.Proto, resp.Status), nil
}
