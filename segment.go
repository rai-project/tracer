package tracer

// Segment is analogous to an OpenTracing span
type Segment interface {
	// sync.RWMutex
	// Name       string `json:"name,omitempty"`
	// Type       string `json:"type,omitempty"`
	// ID         string `json:"id,omitempty"`
	// TraceID    string `json:"trace_id,omitempty"`
	// ParentID   string `json:"parent_id,omitempty"`
	// StartTime  int64  `json:"start_time,omitempty"`
	// EndTime    int64  `json:"end_time,omitempty"`
	// InProgress bool   `json:"in_progress,omitempty"`
	Finish()
	SetTag(key string, value interface{}) Segment
	SetHTTPPath(value string)
	SetHTTPHost(value string)
	Context() SegmentContext
}

// SegmentContext holds any internal state needed for managing Segments
type SegmentContext interface {
}
