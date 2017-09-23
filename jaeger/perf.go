package jaeger

import (
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

var (
	perfeventsObserver jaeger.Observer = noopObserver{}
)

type noopObserver struct{}

func (o noopObserver) OnStartSpan(operationName string, options opentracing.StartSpanOptions) jaeger.SpanObserver {
	return nil
}
