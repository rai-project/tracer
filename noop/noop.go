package noop

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
)

type Tracer struct {
	opentracing.NoopTracer
	id    string
	level tracer.Level
}

func New(serviceName string, opts ...tracer.Option) (tracer.Tracer, error) {
	return &Tracer{
		NoopTracer: opentracing.NoopTracer{},
		id:         uuid.NewV4(),
	}, nil
}

func (t *Tracer) ID() string {
	return t.id
}

func (t *Tracer) Level() tracer.Level {
	return tracer.NO_TRACE
}

func (t *Tracer) SetLevel(lvl tracer.Level) {
	t.level = lvl
}

func (*Tracer) Init(_ string, opts ...tracer.Option) error {
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
