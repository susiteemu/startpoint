package client

import (
	"github.com/go-resty/resty/v2"
	"goful/core/model"
)

func DoRequest(request model.Request) (*model.Response, error) {
	client := resty.New()

	requestHeaders := request.Headers.ToMap()
	// TODO enable trace?
	// TODO handle body vs formdata, also check if []byte can be string before casting
	resp, err := client.R().SetHeaders(requestHeaders).SetBody(request.Body).Execute(request.Method, request.Url)
	if err != nil {
		return &model.Response{}, err
	}
	r := model.Response{
		Headers:    new(model.Headers).FromMap(resp.Header()),
		Body:       resp.Body(),
		Status:     resp.Status(),
		StatusCode: resp.StatusCode(),
		Proto:      resp.Proto(),
		Size:       resp.Size(),
		ReceivedAt: resp.ReceivedAt(),
	}

	return &r, nil
}
