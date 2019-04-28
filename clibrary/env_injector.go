package main

// #cgo CFLAGS: -I${SRCDIR} -O3 -Wall -g
// #cgo LDFLAGS: -ldl
// void enviable_setenv(char *line);
import "C"

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
	C.enviable_setenv(fmt.Sprintf("%s=%s", traceEnvName, traceIDVal))
}
