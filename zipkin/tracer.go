package zipkin

import (
	"context"
	"errors"
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rai-project/tracer"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

var opentracingGlobalTracerIsSet bool

type Tracer struct {
	tracer   opentracing.Tracer
	closer   io.Closer
	endpoint string
	name     string
}

func NewTracer(serviceName string) *Tracer {
	endpoint := Config.Endpoints[0]
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
		jaeger.NewConstSampler(true /*sample all*/),
		jaeger.NewRemoteReporter(trans),
	)

	if _, ok := opentracing.GlobalTracer().(opentracing.NoopTracer); !ok {
		log.Error("Expecting global tracer to be uninitialized")
	}
	opentracing.SetGlobalTracer(tr)
	return &Tracer{tracer: tr, closer: cl, endpoint: endpoint, name: serviceName}
}

func (t *Tracer) SegmentFromContext(ctx context.Context) tracer.Segment {
	panic("tracer/zipkin/Segment from Context Unimplemented")
}
func (t *Tracer) NewChildSegment(parent tracer.Segment) tracer.Segment {
	panic("NewChildSegment Unimplemented")
}

func (t *Tracer) ContextWithSegment(orig context.Context, s tracer.Segment) context.Context {

	sg, ok := s.(*Segment)
	if !ok {
		return orig
	}

	return opentracing.ContextWithSpan(orig, sg.span)
}

func (t *Tracer) StartSegment(operationName string, sc tracer.SegmentContext) (tracer.Segment, error) {
	wireContext, ok := sc.(*SegmentContext)
	if !ok {
		return nil, errors.New("Starting segment from something that is not a zipkin.SegmentContext")
	}
	span := t.tracer.StartSpan(operationName, ext.RPCServerOption(wireContext.sc))
	return &Segment{span: span}, nil
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

func (t *Tracer) Inject(c tracer.SegmentContext, req *http.Request) error {

	sm, ok := c.(*SegmentContext)
	if !ok {
		return errors.New("Injecting something that is not a zipkin.SegmentContext")
	}

	return t.tracer.Inject(
		sm.sc,
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

func (t *Tracer) Extract(req *http.Request) (tracer.SegmentContext, error) {
	wireContext, err := t.tracer.Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	return &SegmentContext{sc: wireContext}, err
}

func (t *Tracer) Endpoint() string {
	return t.endpoint
}

func (t *Tracer) Name() string {
	return t.name
}
