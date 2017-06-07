package tracer

import "context"
import "net/http"

var stdTracer Tracer

type Tracer interface {
	SegmentFromContext(ctx context.Context) Segment
	// NewChildSegment(parent Segment) Segment
	// ContextWithSegment stores a Segment in a context
	ContextWithSegment(orig context.Context, s Segment) context.Context
	StartSegment(operationName string, opts ...StartSegmentOption) Segment
	StartSegmentFromContext(ctx context.Context, operationName string) (Segment, context.Context)
	Close()
	// Inject a Segment into a request
	Inject(c SegmentContext, r *http.Request) error
	// Extract a Segment from a request
	Extract()
}

type StartSegmentOption struct {
}

func SetGlobal(t Tracer) {
	stdTracer = t
}

func SegmentFromContext(ctx context.Context) Segment {
	return stdTracer.SegmentFromContext(ctx)
}

// func NewChildSegment(parent Segment) Segment {
// 	return stdTracer.NewChildSegment(parent)
// }

// ContextWithSegment stores a segment in a context
func ContextWithSegment(orig context.Context, s Segment) context.Context {
	return stdTracer.ContextWithSegment(orig, s)
}

func StartSegment(operationName string, opts ...StartSegmentOption) Segment {
	return stdTracer.StartSegment(operationName, opts...)
}

func StartSegmentFromContext(ctx context.Context, operationName string) (Segment, context.Context) {
	return stdTracer.StartSegmentFromContext(ctx, operationName)
}

func Close() {
	stdTracer.Close()
}
