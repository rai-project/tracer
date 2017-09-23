package chrome

import (
	"fmt"
	"sort"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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
}

func (t TraceEvent) ID() string {
	return fmt.Sprintf("%s::%s/%v", t.Category, t.Name, t.ThreadID)
}

type TraceEvents []TraceEvent

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

type Trace struct {
	InitTime        time.Time              `json:"-"`
	StartTime       time.Time              `json:"-"`
	EndTime         time.Time              `json:"-"`
	TraceEvents     TraceEvents            `json:"traceEvents,omitempty"`
	DisplayTimeUnit string                 `json:"displayTimeUnit,omitempty"`
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

func (t Trace) Publish(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context, error) {

	var timeUnit time.Duration
	switch t.DisplayTimeUnit {
	case "ns":
		timeUnit = time.Nanosecond
	case "us":
		timeUnit = time.Microsecond
	case "ms":
		timeUnit = time.Millisecond
	case "":
		timeUnit = time.Microsecond
	default:
		return nil, nil, errors.Errorf("the display time unit %v is not valid", t.DisplayTimeUnit)
	}

	start := t.StartTime

	topOpts := append(
		[]opentracing.StartSpanOption{
			opentracing.StartTime(start),
			opentracing.Tags{
				"time_unit": t.DisplayTimeUnit,
			},
		},
		opts...,
	)
	span, newCtx := opentracing.StartSpanFromContext(ctx, operationName, topOpts...)
	if span == nil {
		return nil, ctx, errors.New("span not found in context")
	}
	defer span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: t.EndTime,
	})

	spans := map[string]*publishInfo{}

	sort.Sort(t)

	minTime := int64(0)
	events := []TraceEvent{}
	for _, event := range t.TraceEvents {
		if event.EventType != "B" && event.EventType != "E" {
			continue
		}
		t := t.InitTime.Add(time.Duration(event.Timestamp) * timeUnit)
		if start.After(t) {
			continue
		}
		events = append(events, event)
		if event.EventType != "B" {
			continue
		}
		if minTime != 0 && minTime < event.Timestamp {
			continue
		}
		minTime = event.Timestamp
	}

	for _, event := range events {
		id := event.ID()
		if event.EventType == "B" {
			t := time.Duration(event.Timestamp-minTime) * timeUnit
			startTime := start.Add(t)
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
			s, _ := opentracing.StartSpanFromContext(
				ctx,
				event.Name,
				opentracing.ChildOf(span.Context()),
				opentracing.StartTime(startTime),
				tags,
			)
			if s == nil {
				continue
			}
			spans[id] = &publishInfo{
				startEvent: event,
				startTime:  startTime,
				span:       s,
			}
			continue
		}
		startEntry, ok := spans[id]
		if !ok {
			continue
		}
		s := startEntry.span
		if event.Duration == 0 {
			event.Duration = time.Duration(event.Timestamp-startEntry.startEvent.Timestamp) * timeUnit
		}
		endTime := startEntry.startTime.Add(event.Duration)
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

	return span, newCtx, nil
}
