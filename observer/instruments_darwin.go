// +build ignore
// +build cgo,darwin

package observer

import (
	"github.com/nicolai86/instruments"
	otobserver "github.com/opentracing-contrib/go-observer"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	Instruments otobserver.Observer = instrumentsObserver{}
)

type instrumentsObserver struct {
}

// SpanObserver collects perfevent metrics
type instrumentsSpanObserver struct {
	sp     opentracing.Span
	region instruments.Region
}

func (t instrumentsObserver) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) (otobserver.SpanObserver, bool) {
	return newInstrumentsSpanObserver(sp, options)
}

func newInstrumentsSpanObserver(sp opentracing.Span, options opentracing.StartSpanOptions) (otobserver.SpanObserver, bool) {
	region := instruments.StartWithArguments(43, 0, 0, 0, instruments.ColorPurple)
	so := &instrumentsSpanObserver{
		sp:     sp,
		region: region,
	}
	return so, true
}

// Callback called from opentracing.Span.SetTag()
func (so *instrumentsSpanObserver) OnSetTag(key string, value interface{}) {

}

// Callback called from opentracing.Span.SetOperationName()
func (so *instrumentsSpanObserver) OnSetOperationName(operationName string) {

}

// Callback called from opentracing.Span.Finish()
func (so *instrumentsSpanObserver) OnFinish(options opentracing.FinishOptions) {
	so.region.End()
}
