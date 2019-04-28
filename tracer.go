package tracer

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
)

type Tracer interface {
	opentracing.Tracer
	ID() string
	StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context)
	io.Closer
	Init(serviceName string, opts ...Option) error
	Name() string
	Level() Level
	SetLevel(Level)
	Endpoints() []string
}
