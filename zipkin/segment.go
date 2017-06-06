package zipkin

import (
	opentracing "github.com/opentracing/opentracing-go"
)

type Segment struct {
	span opentracing.Span
}

func (s *Segment) Finish() {
	s.span.Finish()
}
