package tracer

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	context "golang.org/x/net/context"
)

type Tracer interface {
	opentracing.Tracer
	ID() string
	StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context)
	io.Closer
	Init(serviceName string) error
	Name() string
	Level() Level
	Endpoints() []string
}
