package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
)

func GetSpanContext() (jaeger.SpanContext, error) {

	var traceID jaeger.TraceID
	var spanID uint64
	var parentID uint64

	if value := os.Getenv(TraceIdEnv); value != "" {
		m, err := jaeger.TraceIDFromString(value)
		if err != nil {
			return jaeger.SpanContext{}, errors.Wrapf(err, "failed to get env[%v]", TraceIdEnv)
		}
		traceID = m
	}

	if value := os.Getenv(ParentSpanIdEnv); value != "" {
		m, err := strconv.ParseUint(value, 16, 64)
		if err != nil {
			return jaeger.SpanContext{}, errors.Wrapf(err, "failed to get env[%v]", ParentSpanIdEnv)
		}
		parentID = m
	}

	if value := os.Getenv(SpanIdEnv); value != "" {
		m, err := strconv.ParseUint(value, 16, 64)
		if err != nil {
			return jaeger.SpanContext{}, errors.Wrapf(err, "failed to get env[%v]", SpanIdEnv)
		}
		spanID = m
	}

	if !traceID.IsValid() {
		return jaeger.SpanContext{}, opentracing.ErrSpanContextNotFound
	}

	fmt.Printf("traceID = %v\n", traceID)

	return jaeger.NewSpanContext(
		traceID,
		jaeger.SpanID(spanID),
		jaeger.SpanID(parentID),
		true,
		map[string]string{}), nil
}
