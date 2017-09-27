package noop

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	context "golang.org/x/net/context"
)

type Tracer struct {
	opentracing.NoopTracer
	id string
}

func New(serviceName string) (tracer.Tracer, error) {
	return &Tracer{
		NoopTracer: opentracing.NoopTracer{},
		id:         uuid.NewV4(),
	}, nil
}

func (t *Tracer) ID() string {
	return t.id
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
