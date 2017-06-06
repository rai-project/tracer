package zipkin

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func NewTracer(serviceName string) *Tracer {
	trans, err := zipkin.NewHTTPTransport(
		endpoint,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
	}
	tr, cl := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true), // sample all traces
		jaeger.NewRemoteReporter(trans),
	)

	opentracing.SetGlobalTracer(tr)

	return &Tracer{tracer: tr, closer: cl}
}

func (t *Tracer) SegmentFromContext(ctx context.Context) tracer.Segment {
	panic("Unimplemented")
}
func (t *Tracer) NewChildSegment(parent tracer.Segment) tracer.Segment {
	panic("Unimplemented")
}

func (t *Tracer) ContextWithSegment(orig context.Context, s tracer.Segment) context.Context {
	panic("Unimplemented")
}

func (t *Tracer) StartSegment() tracer.Segment {
	panic("Unimplemented")
}

// StartSegmentFromContext starts and returns a Segment with `operationName`,
// using any Segment in the ctx as its parent. If none can be found,
// StartSegmentFromContext creates a root (parentless) Segment
//
// The second return value is a context.Context object built around the
// returned Segment.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, ...) {
//        sp, ctx := opentracing.StartSpanFromContext(ctx, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func (t *Tracer) StartSegmentFromContext(ctx context.Context, operationName string) (tracer.Segment, context.Context) {

	// use opentracing to extract or create a span from the context
	sp, ctx := opentracing.StartSpanFromContext(ctx, operationName)

	sg := &Segment{span: sp}
	return sg, ctx
}

func (t *Tracer) Close() {
	t.closer.Close()
}
