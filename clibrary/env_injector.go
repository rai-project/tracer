package main

import (
	"os"

	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

var (
	traceEnvName = "api_tracer_trace_id"
)

func GetTraceEnv() string {
	return os.Getenv(traceEnvName)
}

func SetTraceEnv(span opentracing.Span) {

	traceID := span.Context().(jaeger.SpanContext).TraceID()
	traceIDVal := traceID.String()

	os.Setenv(traceEnvName, traceIDVal)
}
