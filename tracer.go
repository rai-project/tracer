package tracer

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
)

type Tracer interface {
	opentracing.Tracer
	io.Closer
	Name() string
	Endpoint() string
}

var stdTracer Tracer

func SetStd(t Tracer) {
	stdTracer = t
}

func Std() Tracer {
	return stdTracer
}

func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return stdTracer.StartSpan(operationName, opts...)
}

func Close() {
	stdTracer.Close()
}

func Enabled() bool {
	return Config.Enabled
}

func Backend() string {
	return Config.Backend
}
