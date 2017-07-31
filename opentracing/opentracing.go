package opentracing

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
)

type Tracer struct {
	initialized bool
	opentracing.Tracer
	closer      io.Closer
	endpoints   []string
	serviceName string
}
