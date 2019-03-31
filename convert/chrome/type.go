package chrome

import (
	"math"
	"time"

	"github.com/apex/log"

	"hash/fnv"

	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/pkg/errors"
)

var (
	initTime               time.Time
	DefaultDisplayTimeUnit = "ms"
)

func timeUnit(unit string) (time.Duration, error) {
	switch unit {
	case "ns":
		return time.Nanosecond, nil
	case "us":
		return time.Microsecond, nil
	case "ms":
		return time.Millisecond, nil
	case "":
		return time.Microsecond, nil
	default:
		return time.Duration(0), errors.Errorf("the display time unit %v is not valid", unit)
	}
}

func mustTimeUnit(u string) time.Duration {
	unit, err := timeUnit(u)
	if err != nil {
		panic(err)
	}
	return unit
}

// Trace is an entry of trace format.
// https://github.com/catapult-project/catapult/tree/master/tracing
type TraceEvent struct {
	Name      string      `json:"name,omitempty"`
	EventType string      `json:"ph"`
	Scope     string      `json:"s,omitempty"`
	Timestamp float64     `json:"ts"`
	Duration  float64     `json:"dur,omitempty"`
	ProcessID uint64      `json:"pid"`
	ThreadID  uint64      `json:"tid"`
	SpanID    uint64      `json:"id,omitempty"`
	Stack     int         `json:"sf,omitempty"`
	EndStack  int         `json:"esf,omitempty"`
	Arg       interface{} `json:"args,omitempty"`
	ColorName string      `json:"cname,omitempty"`
	Category  string      `json:"cat,omitempty"`

	Start     int64         `json:"start,omitempty"`
	End       int64         `json:"end,omitempty"`
	InitTime  time.Time     `json:"init_time_t,omitempty"`
	StartTime time.Time     `json:"start_time_t,omitempty"`
	EndTime   time.Time     `json:"end_time_t,omitempty"`
	Time      time.Time     `json:"time_t,omitempty"`
	TimeUnit  time.Duration `json:"timeUnit,omitempty"`
}

type EventFrame struct {
	Name   string `json:"name"`
	Parent int    `json:"parent,omitempty"`
}

type TraceEvents []TraceEvent

type Trace struct {
	ID              string                `json:"id,omitempty"`
	StartTime       time.Time             `json:"start_time,omitempty"`
	EndTime         time.Time             `json:"end_time,omitempty"`
	TraceEvents     TraceEvents           `json:"traceEvents,omitempty"`
	DisplayTimeUnit string                `json:"displayTimeUnit,omitempty"`
	Frames          map[string]EventFrame `json:"stackFrames"`
	TimeUnit        string                `json:"timeUnit,omitempty"`
}

// ID should be a unique ID representing this interval.  This
// is used to identify which interval to delete from the tree if
// there are duplicates.
func (x TraceEvent) ID() uint64 {
	h := fnv.New64a()
	h.Write([]byte(x.Name))
	return h.Sum64()
}

// LowAtDimension returns an integer representing the lower bound
// at the requested dimension.
func (x TraceEvent) LowAtDimension(d uint64) int64 {
	if d != 1 {
		return 0
	}
	return x.Start
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (x TraceEvent) HighAtDimension(d uint64) int64 {
	if d != 1 {
		return 0
	}
	return x.End
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (mi TraceEvent) OverlapsAtDimension(iv augmentedtree.Interval, dimension uint64) bool {
	return mi.HighAtDimension(dimension) > iv.LowAtDimension(dimension) &&
		mi.LowAtDimension(dimension) < iv.HighAtDimension(dimension)
}

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

func (x Trace) Adjust() (Trace, error) {
	tr, err := x.DeleteIgnoredEvents()
	if err != nil {
		log.WithError(err).Error("failed to delete ignored events")
		tr = x
	}
	return tr.ZeroOut(), nil
}

func durationOf(t TraceEvent) time.Duration {
	return t.EndTime.Sub(t.StartTime)
}

func (x Trace) DeleteIgnoredEvents() (Trace, error) {
	var minTimeStamp float64
	var adjustedEvent TraceEvent
	events := TraceEvents{}
	for _, event := range x.TraceEvents { // assumes that there is only one thing to ignore
		if event.Category == "ignore" {
			if event.EventType == "E" {
				adjustedEvent = event
			}
		}
		if minTimeStamp < event.Timestamp {
			minTimeStamp = event.Timestamp
		}
	}
	if adjustedEvent.Name == "" {
		return x, nil
	}
	// pp.Println(adjustedEvent)
	for _, event := range x.TraceEvents {
		timeUnit := event.TimeUnit
		// initTime, _ := time.Parse(time.RFC3339Nano, event.Init)
		// pp.Println(timeAdjustmentI, "   ", event.Timestamp, "   ", event.Timestamp-timeAdjustmentI)
		if event.Category == "ignore" {
			continue
		}
		if event.EndTime.After(adjustedEvent.Time) && event.StartTime.Before(adjustedEvent.Time) {
			event.Duration = event.Duration - adjustedEvent.Duration
		}
		if event.EndTime.Before(adjustedEvent.Time) {
			events = append(events, event)
			continue
		}
		// if event.Name == "load_nd_array" {
		// 	continue
		// }

		if event.EventType == "B" || event.EventType == "E" {
			if event.Time.After(adjustedEvent.Time) {
				event.Time = event.Time.Add(-durationOf(adjustedEvent))
				event.Timestamp = event.Timestamp - float64(durationOf(adjustedEvent))/float64(timeUnit)
				// pp.Println(event.Timestamp, "   ", adjustedEvent.Timestamp, "  ", int64(adjustedEvent.Duration), "   ", event.Timestamp-adjustedEvent.Timestamp+minTimeStamp)
			}
			if event.StartTime.After(adjustedEvent.StartTime) {
				event.Start = event.Start - adjustedEvent.Start
				event.StartTime = event.StartTime.Add(-durationOf(adjustedEvent))
			}

			if event.EndTime.After(adjustedEvent.EndTime) {
				event.End = event.End - adjustedEvent.End
				event.EndTime = event.EndTime.Add(-durationOf(adjustedEvent))
			}
		}

		events = append(events, event)
	}
	x.TraceEvents = events
	return x, nil
}

func (x Trace) MaxEvent() TraceEvent {
	var maxEvent TraceEvent
	maxTimeStamp := math.SmallestNonzeroFloat64
	for _, event := range x.TraceEvents { // assumes that there is only one thing to ignore
		if event.Category == "ignore" {
			continue
		}
		if event.EventType != "E" {
			continue
		}
		if maxTimeStamp < event.Timestamp {
			maxTimeStamp = event.Timestamp
			maxEvent = event
		}
	}
	return maxEvent
}

func (x Trace) MinEvent() TraceEvent {
	var minEvent TraceEvent
	minTimeStamp := math.MaxFloat64
	for _, event := range x.TraceEvents { // assumes that there is only one thing to ignore
		if event.Category == "ignore" {
			continue
		}
		if event.EventType != "B" {
			continue
		}
		if minTimeStamp > event.Timestamp {
			minTimeStamp = event.Timestamp
			minEvent = event
		}
	}
	return minEvent
}

func (x Trace) ZeroOut() Trace {
	minEvent := x.MinEvent()
	minTimeStamp := minEvent.Timestamp

	td := minEvent.StartTime.Sub(x.StartTime)

	return x.AddTimestampOffset(int64(-minTimeStamp)).AddDurationOffset(td)
}

func (x Trace) AddTimestampOffset(ts int64) Trace {
	events := make([]TraceEvent, len(x.TraceEvents))
	for ii, event := range x.TraceEvents {
		if event.EventType == "B" || event.EventType == "E" {
			event.Timestamp = event.Timestamp + float64(ts)
			event.Start = event.Start + ts
			event.End = event.End + ts
		}
		events[ii] = event
	}
	x.TraceEvents = events
	return x
}

func (x Trace) AddDurationOffset(td time.Duration) Trace {
	events := make([]TraceEvent, len(x.TraceEvents))
	for ii, event := range x.TraceEvents {
		if event.EventType == "B" || event.EventType == "E" {
			event.Time = event.Time.Add(td)
			event.StartTime = event.StartTime.Add(td)
			event.EndTime = event.EndTime.Add(td)
		}
		events[ii] = event
	}
	x.TraceEvents = events
	return x
}

func (x Trace) HashID() int64 {
	h := fnv.New32a()
	h.Write([]byte(x.ID))
	return int64(h.Sum32())
}
