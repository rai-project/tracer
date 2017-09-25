package chrome

import (
	"fmt"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/tracer"
	"golang.org/x/net/context"
)

// Trace is an entry of trace format.
// https://github.com/catapult-project/catapult/tree/master/tracing
type TraceEvent struct {
	Name      string                 `json:"name,omitempty"`
	Category  string                 `json:"cat,omitempty"`
	EventType string                 `json:"ph,omitempty"`
	Timestamp int64                  `json:"ts,omitempty"`  // displayTimeUnit
	Duration  time.Duration          `json:"dur,omitempty"` // displayTimeUnit
	ProcessID int64                  `json:"pid,omitempty"`
	ThreadID  int64                  `json:"tid,omitempty"`
	Args      map[string]interface{} `json:"args,omitempty"`
	Time      time.Time              `json:"-"`
}

func (t TraceEvent) ID() string {
	return fmt.Sprintf("%s::%s/%v", t.Category, t.Name, t.ThreadID)
}

type TraceEvents []TraceEvent

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

type Trace struct {
	StartTime       time.Time              `json:"-"`
	EndTime         time.Time              `json:"-"`
	TraceEvents     TraceEvents            `json:"traceEvents,omitempty"`
	DisplayTimeUnit string                 `json:"displayTimeUnit,omitempty"`
	TimeUnit        string                 `json:"timeUnit,omitempty"`
	OtherData       map[string]interface{} `json:"otherData,omitempty"`
}

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

type publishInfo struct {
	startEvent TraceEvent
	startTime  time.Time
	span       opentracing.Span
}

func (t Trace) Publish(ctx context.Context, tracer tracer.Tracer, opts ...opentracing.StartSpanOption) error {

	var timeUnit time.Duration
	switch t.TimeUnit {
	case "ns":
		timeUnit = time.Nanosecond
	case "us":
		timeUnit = time.Microsecond
	case "ms":
		timeUnit = time.Millisecond
	case "":
		timeUnit = time.Microsecond
	default:
		return errors.Errorf("the display time unit %v is not valid", t.DisplayTimeUnit)
	}

	spans := map[string]*publishInfo{}

	for _, event := range t.TraceEvents {
		id := event.ID()
		if event.EventType == "B" {
			tags := opentracing.Tags{
				"category":   event.Category,
				"process_id": event.ProcessID,
				"thread_id":  event.ThreadID,
				// "start_timestamp": timeUnit * time.Duration(event.Timestamp),
				// "start_time":      startTime,
			}
			for k, v := range event.Args {
				tags[k] = v
			}
			s, _ := tracer.StartSpanFromContext(
				ctx,
				event.Name,
				opentracing.StartTime(event.Time),
				tags,
			)
			if s == nil {
				continue
			}
			spans[id] = &publishInfo{
				startEvent: event,
				startTime:  event.Time,
				span:       s,
			}
			continue
		}
		startEntry, ok := spans[id]
		if !ok {
			continue
		}
		s := startEntry.span
		endTime := event.Time
		if event.Duration != 0 {
			endTime = startEntry.startTime.Add(event.Duration * timeUnit)
		}
		// duration := endTime.Sub(startEntry.startTime).Nanoseconds()
		s.
			// SetTag("end_timestamp", timeUnit*time.Duration(event.Timestamp)).
			// SetTag("endtime", endTime).
			// SetTag("duration(ns)", duration).
			FinishWithOptions(opentracing.FinishOptions{
				FinishTime: endTime,
			})
		delete(spans, id)
	}

	return nil
}
