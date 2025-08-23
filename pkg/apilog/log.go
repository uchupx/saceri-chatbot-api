package apilog

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/uchupx/saceri-chatbot-api/pkg/helper"
)

type ApiLog struct {
	log *logrus.Entry
}

type Params struct {
	Level       uint64
	Version     string
	ServiceName string
}

func (l *ApiLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	fields = l.getContext(ctx, fields)
	entry := l.log.WithFields(fields)
	entry.Info(msg)
}

func (l *ApiLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	fields = l.getContext(ctx, fields)
	entry := l.log.WithFields(fields)
	entry.Debug(msg)
}

func (l *ApiLog) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	fields = l.getContext(ctx, fields)
	fields["error"] = err

	entry := l.log.WithFields(fields)
	entry.Error(msg)
}

func (l *ApiLog) Warn(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	fields = l.getContext(ctx, fields)
	if err != nil {
		fields["error"] = err
	}

	entry := l.log.WithFields(fields)
	entry.Warn(msg)

}

func (l *ApiLog) getContext(ctx context.Context, fields map[string]interface{}) map[string]interface{} {
	traceId := ctx.Value("trace_id")
	if traceId == nil {
		traceId = "none"
	}

	requestBody := ctx.Value("request_body")
	if requestBody == nil {
		requestBody = "none"
	} else {
		bytes, err := json.MarshalIndent(requestBody, "", "  ")
		if err != nil {
			requestBody = "none"
		}

		requestBody = string(bytes)
	}

	if fields == nil {
		fields = make(map[string]interface{})
	}

	fields["trace_id"] = traceId
	fields["request_body"] = requestBody

	return fields
}

func (l *ApiLog) CreateTrace(ctx context.Context) context.Context {
	traceId := l.generateTraceID()

	return context.WithValue(ctx, "trace_id", traceId)
}

func (l *ApiLog) generateTraceID() string {
	return "trace-" + helper.GenerateRandomString(16)
}

func (l *ApiLog) AttachBody(ctx context.Context, body interface{}) context.Context {
	return context.WithValue(ctx, "request_body", body)
}

func NewApiLog(params Params) *ApiLog {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
		PadLevelText:    true,
	})

	log.SetLevel(logrus.Level(params.Level))

	logEntry := log.WithFields(logrus.Fields{
		"service": params.ServiceName,
		"version": params.Version,
	})

	return &ApiLog{
		log: logEntry,
	}
}
