package tracer

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
)

type Tracer interface {
	opentracing.Tracer
	io.Closer
	Init(serviceName string) error
	Name() string
	Endpoints() []string
}
