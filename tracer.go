package tracer

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer/noop"
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
	opentracing.SetGlobalTracer(t)
}

func Std() Tracer {
	return stdTracer
}

func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return stdTracer.StartSpan(operationName, opts...)
}

// StartSpanFromContext starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a ChildOfRef. If no such parent could be
// found, StartSpanFromContext creates a root (parentless) Span.
//
// The second return value is a context.Context object built around the
// returned Span.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, ...) {
//        sp, ctx := opentracing.StartSpanFromContext(ctx, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if _, ok := opentracing.GlobalTracer().(opentracing.NoopTracer); ok {
		log.Error("tracer is using a No-op tracer")
	}
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
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

func init() {
	t, _ := noop.New("")
	SetStd(t)
}
