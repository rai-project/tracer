package observer

import (
	"github.com/opentracing-contrib/go-observer"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	NoOp otobserver.Observer = noopObserver{}
)

type noopObserver struct{}

func (o noopObserver) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) (otobserver.SpanObserver, bool) {
	return nil, false
}
