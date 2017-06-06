package tracer

import "context"

var stdTracer Tracer

type Tracer interface {
	SegmentFromContext(ctx context.Context) Segment
	NewChildSegment(parent Segment) Segment
	ContextWithSegment(orig context.Context, s Segment) context.Context
	StartSegment() Segment
	StartSegmentFromContext(ctx context.Context, operationName string) (Segment, context.Context)
	Close()
}

func SetGlobal(t Tracer) {
	stdTracer = t
}

func SegmentFromContext(ctx context.Context) Segment {
	return stdTracer.SegmentFromContext(ctx)
}

func NewChildSegment(parent Segment) Segment {
	return stdTracer.NewChildSegment(parent)
}

func ContextWithSegment(orig context.Context, s Segment) context.Context {
	return stdTracer.ContextWithSegment(orig, s)
}
func StartSegment() Segment {
	return stdTracer.StartSegment()
}

func StartSegmentFromContext(ctx context.Context, operationName string) (Segment, context.Context) {
	return stdTracer.StartSegmentFromContext(ctx, operationName)
}

func Close() {
	stdTracer.Close()
}
