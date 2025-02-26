package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/susiteemu/startpoint/core/configuration"
	"github.com/susiteemu/startpoint/core/model"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type restyZeroLogger struct {
	logger *zerolog.Logger
}

func newLogger(zlogger *zerolog.Logger) *restyZeroLogger {
	return &restyZeroLogger{
		logger: zlogger,
	}
}
func (l *restyZeroLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}

func (l *restyZeroLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

func (l *restyZeroLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func DoRequest(request model.Request) (*model.Response, error) {
	requestHeaders := request.Headers.ToMap()
	log.Debug().Msgf("Request %v -- %v -- %v -- %v -- %v", request.Url, request.Body, request.Method, request.Headers, request.Options)

	// NOTE: creating new client for each request to simplify configuring it
	// (no need to reset to default values after request)
	var client = resty.New().SetLogger(newLogger(&log.Logger))

	config := configuration.NewWithRequestOptions(request.Options)

	debug := config.GetBool("debug")
	if debug {
		client.SetDebug(debug)
	} else {
		debug := config.GetBool("httpClient.debug")
		client.SetDebug(debug)
	}

	timeoutSeconds, set := config.GetInt("httpClient.timeoutSeconds")
	if set && timeoutSeconds >= 0 {
		client.SetTimeout(time.Duration(timeoutSeconds * int(time.Second)))
	}

	proxy, set := config.GetString("httpClient.proxyUrl")
	if set && len(proxy) > 0 {
		client = client.SetProxy(proxy)
	}

	insecure := config.GetBoolWithDefault("httpClient.insecure", false)
	if insecure {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		rootCertificates, _ := config.GetStringSlice("httpClient.rootCertificates")
		if len(rootCertificates) > 0 {
			for _, cert := range rootCertificates {
				client.SetRootCertificate(cert)
			}
		}
		clientCertificates, set := config.GetSliceMapString("httpClient.clientCertificates")
		if set && len(clientCertificates) > 0 {
			var certs []tls.Certificate
			for _, pair := range clientCertificates {
				// NOTE: viper makes all keys lowercase which is a bit annoying because
				// configuration coming from request does not have lowercase keys: we now must handle both cases
				certFile, found := pair["certfile"]
				if !found {
					certFile = pair["certFile"]
				}
				keyFile, found := pair["keyfile"]
				if !found {
					keyFile = pair["keyFile"]
				}
				cert, err := tls.LoadX509KeyPair(certFile, keyFile)
				if err != nil {
					log.Error().Err(err).Msg("Error with client certificate")
					return nil, err
				}
				certs = append(certs, cert)
			}
			client.SetCertificates(certs...)
		}
	}

	r := client.R().SetHeaders(requestHeaders)

	enableTrace := config.GetBool("httpClient.enableTraceInfo")
	if enableTrace {
		r.EnableTrace()
	}

	if request.IsForm() {
		bodyAsMap, ok := request.BodyAsMap()
		if !ok {
			return nil, errors.New("cannot convert body to map")
		}
		r.SetFormData(bodyAsMap)
	} else if request.IsMultipartForm() {
		bodyAsMap, ok := request.BodyAsMap()
		if !ok {
			return nil, errors.New("cannot convert body to map")
		}
		formData := make(map[string]string)
		files := make(map[string]string)
		for k, v := range bodyAsMap {
			if strings.HasPrefix(v, "@") {
				files[k] = strings.TrimPrefix(v, "@")
			} else {
				formData[k] = v
			}
		}
		if len(formData) > 0 {
			r.SetFormData(formData)
		}
		if len(files) > 0 {
			r.SetFiles(files)
		}
	} else {
		r.SetBody(request.Body)
	}

	if len(request.Output) > 0 {
		r.SetOutput(request.Output)
	}

	resp, err := r.Execute(request.Method, request.Url)
	if err != nil {
		return nil, err
	}

	ti := resp.Request.TraceInfo()
	traceInfo := model.TraceInfo{}
	if enableTrace {
		traceInfo = model.TraceInfo{
			DNSLookup:      ti.DNSLookup,
			ConnTime:       ti.ConnTime,
			TCPConnTime:    ti.TCPConnTime,
			TLSHandshake:   ti.TLSHandshake,
			ServerTime:     ti.ServerTime,
			ResponseTime:   ti.ResponseTime,
			TotalTime:      ti.TotalTime,
			IsConnReused:   ti.IsConnReused,
			IsConnWasIdle:  ti.IsConnWasIdle,
			ConnIdleTime:   ti.ConnIdleTime,
			RequestAttempt: ti.RequestAttempt,
			RemoteAddr:     ti.RemoteAddr.String(),
		}
	}

	var body []byte
	if resp.IsSuccess() && len(request.Output) > 0 {
		body = []byte(fmt.Sprintf("Saved to file %s", request.Output))
	} else {
		body = resp.Body()
	}

	req := resp.Request
	requestBody := req.Body
	if requestBody == nil && req.FormData != nil {
		urlValues := req.FormData
		bodyAsMap := make(map[string][]string)
		for k, v := range urlValues {
			bodyAsMap[k] = v
		}
		requestBody = bodyAsMap
	}

	respReq := model.Request{
		Url:     req.URL,
		Method:  req.Method,
		Body:    requestBody,
		Headers: new(model.Headers).FromMap(req.Header),
	}

	response := model.Response{
		Headers:    new(model.Headers).FromMap(resp.Header()),
		Body:       body,
		Status:     resp.Status(),
		StatusCode: resp.StatusCode(),
		Proto:      resp.Proto(),
		Size:       resp.Size(),
		ReceivedAt: resp.ReceivedAt(),
		Time:       resp.Time(),
		TraceInfo:  traceInfo,
		Options:    request.Options,
		Request:    respReq,
	}

	log.Debug().Msgf("TraceInfo: %v", traceInfo)

	return &response, nil
}
