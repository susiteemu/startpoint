package client

import (
	"encoding/json"
	"goful/core/model"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func DoRequest(request model.Request) (*model.Response, error) {
	//client := resty.New()

	requestHeaders := request.Headers.ToMap()
	log.Debug().Msgf("Request %v -- %v -- %v", request.Url, request.Body, request.Method)
	// TODO enable trace?
	// TODO handle body vs formdata, also check if []byte can be string before casting
	//

	client := resty.New().R().SetHeaders(requestHeaders)
	if requestHeaders["Content-Type"] == "application/x-www-form-urlencoded" {
		log.Debug().Msg("Request is form!")
		var bodyAsMap map[string]string
		_ = json.Unmarshal([]byte(request.Body.(string)), &bodyAsMap)
		client = client.SetFormData(bodyAsMap)
	} else {
		client = client.SetBody(request.Body)
	}
	resp, err := client.Execute(request.Method, request.Url)
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
		Time:       resp.Time(),
	}

	// TODO add trace flag to configuration and add tracing information to return object
	/*
		ti := resp.Request.TraceInfo()
		fmt.Println("  DNSLookup     :", ti.DNSLookup)
		fmt.Println("  ConnTime      :", ti.ConnTime)
		fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
		fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
		fmt.Println("  ServerTime    :", ti.ServerTime)
		fmt.Println("  ResponseTime  :", ti.ResponseTime)
		fmt.Println("  TotalTime     :", ti.TotalTime)
		fmt.Println("  IsConnReused  :", ti.IsConnReused)
		fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
		fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
		fmt.Println("  RequestAttempt:", ti.RequestAttempt)
		fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
	*/
	return &r, nil
}
