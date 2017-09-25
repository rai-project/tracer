package noop

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	context "golang.org/x/net/context"
)

type Tracer struct {
	opentracing.NoopTracer
}

func New(serviceName string) (tracer.Tracer, error) {
	return &Tracer{opentracing.NoopTracer{}}, nil
}

func (*Tracer) Init(_ string) error {
	return nil
}

func (*Tracer) Name() string {
	return "Noop"
}

// startSpanFromContextWithTracer is factored out for testing purposes.
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return t.NoopTracer.StartSpan(operationName, opts...), ctx
}

func (*Tracer) Endpoints() []string {
	return []string{}
}

func (*Tracer) Close() error {
	return nil
}

func init() {
	tracer.Register("disabled", &Tracer{}, New)
	tracer.Register("noop", &Tracer{}, New)
}
