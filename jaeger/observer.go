package jaeger

import (
	"github.com/opentracing-contrib/go-observer"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

type wrapObserver struct {
	otobserver.Observer
}

func (w *wrapObserver) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) (jaeger.ContribSpanObserver, bool) {
	res, ok := w.Observer.OnStartSpan(sp, operationName, options)
	return res, ok
}
