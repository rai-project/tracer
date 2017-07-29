package opentracing

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
)

type Tracer struct {
	opentracing.Tracer
	closer      io.Closer
	endpoints   []string
	serviceName string
}
