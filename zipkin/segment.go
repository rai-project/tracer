package zipkin

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
)

type Segment struct {
	span opentracing.Span
}

type SegmentContext struct {
	sc opentracing.SpanContext
}

func (s *Segment) Finish() {
	if tracer.Config.Enabled {
		s.span.Finish()
	}
}

func (s *Segment) Context() tracer.SegmentContext {
	return &SegmentContext{sc: s.span.Context()}
}

func (s *Segment) SetTag(key string, value interface{}) tracer.Segment {
	s.span = s.span.SetTag(key, value)
	return s
}
