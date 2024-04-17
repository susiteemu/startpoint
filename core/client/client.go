package client

import (
	"errors"
	"startpoint/core/model"

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

var client *resty.Client = resty.New().SetDebug(true).SetLogger(newLogger(&log.Logger))

func DoRequest(request model.Request) (*model.Response, error) {
	requestHeaders := request.Headers.ToMap()
	log.Debug().Msgf("Request %v -- %v -- %v", request.Url, request.Body, request.Method)

	// TODO have enable trace come from config
	r := client.R().SetHeaders(requestHeaders).EnableTrace()

	if request.IsForm() {
		bodyAsMap, ok := request.BodyAsMap()
		if !ok {
			return nil, errors.New("cannot convert body to map")
		}
		r = r.SetFormData(bodyAsMap)
	} else {
		r = r.SetBody(request.Body)
	}

	resp, err := r.Execute(request.Method, request.Url)
	if err != nil {
		return nil, err
	}

	ti := resp.Request.TraceInfo()
	traceInfo := model.TraceInfo{
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

	response := model.Response{
		Headers:    new(model.Headers).FromMap(resp.Header()),
		Body:       resp.Body(),
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
