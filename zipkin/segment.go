package zipkin

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing/_thrift/gen-go/zipkincore"
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

func (s *Segment) SetHTTPHost(value string) {
	s.SetTag(zipkincore.HTTP_HOST, value)
}
func (s *Segment) SetHTTPPath(value string) {
	s.SetTag(zipkincore.HTTP_PATH, value)
}
