package ctimer

import (
	"encoding/json"
	"fmt"
	"time"

	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

type TraceEvent struct {
	Name               string    `json:"name,omitempty"`
	Metadata           string    `json:"metadata,omitempty"`
	LayerType          string    `json:"layer_type,omitempty"`
	Start              int64     `json:"start,omitempty"`
	End                int64     `json:"end,omitempty"`
	ProcessID          uint64    `json:"process_id,omitempty"`
	ThreadID           uint64    `json:"thread_id,omitempty"`
	LayerSequenceIndex int64     `json:"layer_sequence_index,omitempty"`
	Shapes             [][]int64 `json:"shapes,omitempty"`
	StartTime          time.Time `json:"-"`
	EndTime            time.Time `json:"-"`
}

func (t TraceEvent) ID() string {
	return fmt.Sprintf("%s/%v", t.Name, t.ThreadID)
}

type TraceEvents []TraceEvent

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Start < t[j].Start }

type Trace struct {
	Name        string      `json:"name,omitempty"`
	Metadata    string      `json:"metadata,omitempty"`
	Start       int64       `json:"start,omitempty"`
	End         int64       `json:"end,omitempty"`
	StartTime   time.Time   `json:"-"`
	EndTime     time.Time   `json:"-"`
	TraceEvents TraceEvents `json:"elements,omitempty"`
}

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

func New(data string) (*Trace, error) {
	trace := new(Trace)
	err := json.Unmarshal([]byte(data), trace)
	if err != nil {
		return nil, err
	}
	trace.StartTime = time.Unix(0, trace.Start)
	trace.EndTime = time.Unix(0, trace.End)
	for ii, event := range trace.TraceEvents {
		trace.TraceEvents[ii].StartTime = time.Unix(0, event.Start)
		trace.TraceEvents[ii].EndTime = time.Unix(0, event.End)
		if event.LayerType == "" {
			trace.TraceEvents[ii].LayerType = trace.Metadata
		}
	}
	return trace, nil
}

func (event *TraceEvent) Publish(ctx context.Context, opts ...opentracing.StartSpanOption) error {
	tags := opentracing.Tags{
		"metadata":   event.Metadata,
		"process_id": event.ProcessID,
		"thread_id":  event.ThreadID,
	}
	if event.LayerType != "" {
		tags["layer_type"] = event.LayerType
	}
	if event.LayerSequenceIndex != 0 {
		tags["layer_sequence_index"] = event.LayerSequenceIndex
	}
	if len(event.Shapes) != 0 {
		bts, err := json.Marshal(event.Shapes)
		if err == nil {
			tags["shapes"] = string(bts)
		}
	}
	s, _ := opentracing.StartSpanFromContext(
		ctx,
		event.Name,
		opentracing.StartTime(event.StartTime),
		tags,
	)
	if s == nil {
		log.WithField("event_name", event.Name).
			WithField("tags", tags).
			Error("failed to create span from context")
		return nil
	}
	s.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: event.EndTime,
	})
	return nil
}

func (t *Trace) Publish(ctx context.Context, opts ...opentracing.StartSpanOption) error {
	for _, event := range t.TraceEvents {
		if err := event.Publish(ctx, opts...); err != nil {
			return err
		}
	}
	return nil
}
