package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"startpoint/core/configuration"
	"startpoint/core/model"
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
	log.Debug().Msgf("Request %v -- %v -- %v", request.Url, request.Body, request.Method)

	// NOTE: creating new client for each request to simplify configuring it
	// (no need to reset to default values after request)
	var client *resty.Client = resty.New().SetLogger(newLogger(&log.Logger))
	// settings coming from the configuration to the client
	debug := configuration.GetBool("httpClient.debug")
	client.SetDebug(debug)

	timeoutSeconds, set := configuration.GetInt("httpClient.timeoutSeconds")
	if set && timeoutSeconds >= 0 {
		client.SetTimeout(time.Duration(timeoutSeconds * int(time.Second)))
	}

	proxy, set := configuration.GetString("httpClient.proxyUrl")
	if set && len(proxy) > 0 {
		client = client.SetProxy(proxy)
	}

	insecure := configuration.GetBool("httpClient.insecure")
	if insecure {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		rootCertificates, _ := configuration.GetStringSlice("httpClient.rootCertificates")
		if len(rootCertificates) > 0 {
			for _, cert := range rootCertificates {
				client.SetRootCertificate(cert)
			}
		}
		clientCertificates, set := configuration.GetSliceMapString("httpClient.clientCertificates")
		if set && len(clientCertificates) > 0 {
			var certs []tls.Certificate
			for _, pair := range clientCertificates {
				certFile := pair["certfile"]
				keyFile := pair["keyfile"]
				cert, err := tls.LoadX509KeyPair(certFile, keyFile)
				if err != nil {
					// TODO: error handling
					log.Fatal().Err(err).Msg("Error with client certificate")
				}
				certs = append(certs, cert)
			}
			client.SetCertificates(certs...)
		}
	}

	r := client.R().SetHeaders(requestHeaders)

	enableTrace := configuration.GetBool("httpClient.enableTraceInfo")
	if enableTrace {
		r.EnableTrace()
	}

	if request.IsForm() {
		bodyAsMap, ok := request.BodyAsMap()
		if !ok {
			return nil, errors.New("cannot convert body to map")
		}
		r.SetFormData(bodyAsMap)
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
	}

	log.Debug().Msgf("TraceInfo: %v", traceInfo)

	return &response, nil
}
