package main

import (
	"fmt"
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

	fmt.Printf("Set env [%s] = %s....\n", traceEnvName, traceIDVal)

	os.Setenv(traceEnvName, traceIDVal)
}
