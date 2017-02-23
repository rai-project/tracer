package tracer

import "sync"

type Segment struct {
	sync.RWMutex
	Name       string  `json:"name,omitempty"`
	Type       string  `json:"type,omitempty"`
	ID         string  `json:"id,omitempty"`
	TraceID    string  `json:"trace_id,omitempty"`
	ParentID   string  `json:"parent_id,omitempty"`
	StartTime  float64 `json:"start_time,omitempty"`
	EndTime    float64 `json:"end_time,omitempty"`
	InProgress bool    `json:"in_progress,omitempty"`
}
