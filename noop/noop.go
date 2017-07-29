package noop

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
)

type Tracer struct {
	opentracing.NoopTracer
}

func New(serviceName string) (*Tracer, error) {
	return &Tracer{opentracing.NoopTracer{}}, nil
}

func (*Tracer) Init(_ string) error {
	return nil
}

func (*Tracer) Name() string {
	return "Noop"
}

func (*Tracer) Endpoints() []string {
	return []string{}
}

func (*Tracer) Close() error {
	return nil
}

func init() {
	tracer.AddTracer("disabled", &Tracer{})
	tracer.AddTracer("noop", &Tracer{})
}
