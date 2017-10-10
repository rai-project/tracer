package observer

import (
	"github.com/opentracing-contrib/go-observer"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	Dummy otobserver.Observer = dummy{}
)

type dummy struct{}

// OnStartSpan creates a new Dummy for the span
func (o dummy) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) (otobserver.SpanObserver, bool) {
	return NewSpanDummy(sp, options)
}

// SpanDummy collects perfevent metrics
type SpanDummy struct {
	sp opentracing.Span
}

// NewSpanDummy creates a new SpanDummy that can emit perfevent
// metrics
func NewSpanDummy(s opentracing.Span, opts opentracing.StartSpanOptions) (*SpanDummy, bool) {
	so := &SpanDummy{
		sp: s,
	}
	return so, true
}

func (so *SpanDummy) OnSetOperationName(operationName string) {
}

func (so *SpanDummy) OnSetTag(key string, value interface{}) {
}

func (so *SpanDummy) OnFinish(options opentracing.FinishOptions) {
}
