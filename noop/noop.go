package noop

import (
	opentracing "github.com/opentracing/opentracing-go"
)

type Tracer struct {
	opentracing.NoopTracer
}

func New(serviceName string) (*Tracer, error) {
	return &Tracer{opentracing.NoopTracer{}}, nil
}

func (*Tracer) Name() string {
	return "Noop"
}

func (*Tracer) Endpoint() string {
	return ""
}

func (*Tracer) Close() error {
	return nil
}
