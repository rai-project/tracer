package chrome

import (
	"sort"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Trace is an entry of trace format.
// https://github.com/catapult-project/catapult/tree/master/tracing
type TraceEvent struct {
	Name      string                 `json:"name"`
	Category  string                 `json:"cat"`
	EventType string                 `json:"ph"`
	Timestamp int                    `json:"ts"`  // displayTimeUnit
	Duration  int                    `json:"dur"` // displayTimeUnit
	ProcessID int                    `json:"pid"`
	ThreadID  int                    `json:"tid"`
	Args      map[string]interface{} `json:"args"`
}

func (t TraceEvent) ID() string {
	return t.Name + "/" + strconv.Itoa(t.ThreadID)
}

type TraceEvents []TraceEvent

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

type Trace struct {
	StartTime       time.Time   `json:"-"`
	EndTime         time.Time   `json:"-"`
	TraceEvents     TraceEvents `json:"traceEvents"`
	DisplayTimeUnit string      `json:"displayTimeUnit"`
	OtherData       struct {
		Version string `json:"version"`
	} `json:"otherData"`
}

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

type publishInfo struct {
	startEvent TraceEvent
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

	span, newCtx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	if span == nil {
		return nil, ctx, errors.New("span not found in context")
	}
	defer span.Finish()

	start := t.StartTime

	spans := map[string]*publishInfo{}

	sort.Sort(t)

	for _, event := range t.TraceEvents {
		if event.EventType != "B" && event.EventType != "E" {
			continue
		}
		id := event.ID()
		if event.EventType == "B" {
			s := opentracing.StartSpan(
				event.Name,
				opentracing.ChildOf(span.Context()),
				opentracing.StartTime(start.Add(time.Duration(event.Timestamp)*timeUnit)),
				opentracing.Tags{
					"category":   event.Category,
					"process_id": event.ProcessID,
					"thread_id":  event.ThreadID,
				},
			)
			spans[id] = &publishInfo{
				startEvent: event,
				span:       s,
			}
			continue
		}
		startEntry, ok := spans[id]
		if !ok {
			continue
		}
		if event.Duration != 0 {
			event.Duration = event.Timestamp - startEntry.startEvent.Timestamp
		}
		finishTime := start.Add(time.Duration(startEntry.startEvent.Timestamp+event.Duration) * timeUnit)
		startEntry.span.FinishWithOptions(opentracing.FinishOptions{
			FinishTime: finishTime,
		})
		delete(spans, id)
	}

	return span, newCtx, nil
}
