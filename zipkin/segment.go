package zipkin

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
)

type Segment struct {
	span opentracing.Span
}

func (s *Segment) Finish() {
	if tracer.Config.Enabled {
		s.span.Finish()
	}
}
