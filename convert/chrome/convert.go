package chrome

import (
	"encoding/json"
	"time"

	"github.com/imdario/mergo"
	"github.com/rai-project/tracer/convert"
	model "github.com/uber/jaeger/model/json"
)

type convertState struct {
	tree        *convert.IntervalTree
	jaegerTrace model.Trace
	trace       *Trace
}

func Marshal(trace model.Trace) ([]byte, error) {
	tr, err := Convert(trace)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(tr, "", "  ")
}

func Convert(tr model.Trace) (*Trace, error) {
	st, err := newConvertState(tr)
	if err != nil {
		return nil, err
	}

	err = st.convertSpans()
	if err != nil {
		return nil, err
	}

	// st.trace.ZeroOut()

	return st.trace, nil
}

func newConvertState(tr model.Trace) (*convertState, error) {
	// tr, err := convert.FixTrace(tr)
	// if err != nil {
	// 	return nil, err
	// }
	tree, err := convert.NewIntervalTree(tr)
	if err != nil {
		return nil, err
	}
	_, err = tree.FilterOnlyChildrenOf("PredictImage")
	if err != nil {
		return nil, err
	}

	jaegerTrace, err := tree.FixParentRelationship()
	if err != nil {
		return nil, err
	}
	tree, err = convert.NewIntervalTree(jaegerTrace)
	if err != nil {
		return nil, err
	}

	trace := &Trace{
		ID:              string(jaegerTrace.TraceID),
		DisplayTimeUnit: "ms",
	}

	return &convertState{
		tree:        tree,
		jaegerTrace: jaegerTrace,
		trace:       trace,
	}, nil
}

func (st *convertState) convertSpans() error {
	spans := st.jaegerTrace.Spans
	events := []TraceEvent{}
	for _, span := range spans {
		spanEvents, err := st.convertSpan(span)
		if err != nil {
			return err
		}
		events = append(events, spanEvents...)
	}

	st.trace.TraceEvents = events

	return nil
}

func (st *convertState) convertSpan(sp model.Span) ([]TraceEvent, error) {
	cat := convert.Classify(sp)
	color := colorName(cat)
	depth := st.tree.DepthOf(convert.ToInterval(sp))
	// pp.Println(sp.StartTime)

	_ = depth
	args := map[string]interface{}{
		"depth":      depth,
		"start_time": sp.StartTime, //toTime(sp.StartTime)
		"duration":   sp.Duration,  //toTime(sp.StartTime + sp.Duration)
	}
	for _, tag := range sp.Tags {
		args[tag.Key] = tag.Value
	}
	common := TraceEvent{
		Name:      sp.OperationName,
		SpanID:    hash64(string(sp.SpanID)),
		Category:  cat,
		ColorName: color,
		Start:     int64(sp.StartTime),
		End:       int64(sp.StartTime + sp.Duration),
		StartTime: toTime(sp.StartTime),
		EndTime:   toTime(sp.StartTime + sp.Duration),
		Arg:       args, // map[string]interface{}{},
	}
	// region := TraceEvent{
	// 	EventType: "X",
	// }
	begin := TraceEvent{
		EventType: "B",
		Timestamp: formatTime(sp.StartTime),
	}
	end := TraceEvent{
		EventType: "E",
		Timestamp: formatTime(sp.StartTime + sp.Duration),
	}

	// if err := mergo.Merge(&region, common); err != nil {
	// 	return nil, err
	// }

	if err := mergo.Merge(&begin, common); err != nil {
		return nil, err
	}

	if err := mergo.Merge(&end, common); err != nil {
		return nil, err
	}

	return []TraceEvent{ /*region,*/ begin, end}, nil
}

func formatTime(t0 uint64) float64 {
	d := toDuration(t0)
	return float64(d) / float64(time.Millisecond)
	// return float64(t0)
}
