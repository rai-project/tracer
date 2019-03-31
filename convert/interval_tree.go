package convert

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/Unknwon/com"
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/pkg/errors"
	"github.com/rai-project/evaluation"
	model "github.com/uber/jaeger/model/json"
	"github.com/ulule/deepcopier"
)

type IntervalTree struct {
	augmentedtree.Tree
	trace model.Trace
}

func NewIntervalTreeFromTraceFile(path string) (*IntervalTree, error) {
	if !com.IsFile(path) {
		return nil, errors.Errorf("trace %v does not exist", path)
	}
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read trace from %v", path)
	}
	return NewIntervalTreeFromTraceString(string(bts))
}

func NewIntervalTreeFromTraceString(data string) (*IntervalTree, error) {
	trace := evaluation.TraceInformation{}
	bts := []byte(data)
	err := json.Unmarshal(bts, &trace)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse trace")
	}
	return NewIntervalTree(trace.Traces[0])
}

func NewIntervalTree(trace model.Trace) (*IntervalTree, error) {
	tree := augmentedtree.New(1)
	spans := trace.Spans
	for _, span := range spans {
		tree.Add(spanToInterval(span))
	}
	return &IntervalTree{
		Tree:  tree,
		trace: trace,
	}, nil
}

func (t IntervalTree) ID() string {
	return string(t.trace.TraceID)
}

func (t IntervalTree) ChildrenOf(sp Interval) []Interval {
	elems := t.Query(sp)
	res := []Interval{}
	for _, e0 := range elems {
		e, ok := e0.(Interval)
		if !ok {
			continue
		}
		if e.ID() == sp.ID() {
			continue
		}
		if sp.Contains(e) {
			res = append(res, e)
		}
	}
	return res
}

func (t IntervalTree) DepthOf(sp Interval) int {
	elems := t.Query(sp)
	return len(elems)
}

func (t IntervalTree) ParentsOf(sp Interval) []Interval {
	elems := t.Query(sp)
	res := []Interval{}
	for _, e0 := range elems {
		e, ok := e0.(Interval)
		if !ok {
			continue
		}
		if e.ID() == sp.ID() {
			continue
		}
		res = append(res, e)
	}
	return res
}

func (t IntervalTree) ParentOf(sp Interval) Interval {
	elems := t.ParentsOf(sp)
	if len(elems) == 0 {
		return Interval{}
	}
	if len(elems) == 1 {
		return elems[0]
	}
	sort.Slice(elems, func(ii, jj int) bool {
		return elems[ii].LowAtDimension(0) > elems[jj].LowAtDimension(0)
	})
	return elems[0]
}

func (t IntervalTree) GetIntervals() []Interval {
	length := t.Len()
	intervals := make([]Interval, length)
	ii := 0
	t.Traverse(func(interval augmentedtree.Interval) {
		intervals[ii] = interval.(Interval)
		ii++
	})
	return intervals
}

func (t IntervalTree) MaxInterval() Interval {
	var maxInterval Interval

	t.Traverse(func(interval0 augmentedtree.Interval) {
		interval := interval0.(Interval)
		if maxInterval.IsNil() {
			maxInterval = interval
			return
		}
		if interval.Contains(maxInterval) {
			maxInterval = interval
			return
		}
	})
	return maxInterval
}

// func (t IntervalTree) MinInterval() Interval {
// 	var minInterval Interval

// 	t.Traverse(func(interval0 augmentedtree.Interval) {
// 		interval := interval0.(Interval)
// 		if minInterval.IsNil() {
// 			minInterval = interval
// 			return
// 		}
// 		if minInterval.Contans(interval) {
// 			minInterval = interval
// 			return
// 		}
// 	})
// 	return maxInterval
// }

// func (t IntervalTree) MaxInterval() Interval {
// 	var maxInterval Interval

// 	t.Traverse(func(interval0 augmentedtree.Interval) {
// 		interval := interval0.(Interval)
// 		if maxInterval.IsNil() {
// 			maxInterval = interval
// 			return
// 		}
// 		if maxInterval.End() > Interval.End() {
// 			maxInterval = interval
// 			return
// 		}
// 		if maxInterval.End() == Interval.End() && maxInterval.Start() < Interval.Start() {
// 			maxInterval = interval
// 			return
// 		}
// 	})
// 	return maxInterval
// }

// func (t IntervalTree) MinInterval() Interval {
// 	var maxInterval Interval

// 	t.Traverse(func(interval0 augmentedtree.Interval) {
// 		interval := interval0.(Interval)
// 		if maxInterval.IsNil() {
// 			maxInterval = interval
// 			return
// 		}
// 		if maxInterval.Start() < Interval.Start() {
// 			maxInterval = interval
// 			return
// 		}
// 		if maxInterval.Start() == Interval.Start() && maxInterval.End() > Interval.End() {
// 			maxInterval = interval
// 			return
// 		}
// 	})
// 	return maxInterval
// }

// func (t IntervalTree) GetIntervals() []Interval {
// 	spans := t.trace.Spans
// 	length := len(spans)
// 	intervals := make([]Interval, length)
// 	for ii, spans := range t.trace.Spans {
// 		intervals[ii] = spanToInterval(spans)
// 	}
// 	return intervals
// }

func (t IntervalTree) GetIntervalByIdx(idx int) Interval {
	return spanToInterval(t.trace.Spans[idx])
}

func (t IntervalTree) GetIntervalsByOperationName(name string) []Interval {
	res := []Interval{}
	for _, sp := range t.trace.Spans {
		if sp.OperationName == name {
			res = append(res, spanToInterval(sp))
		}
	}
	return res
}

func (t IntervalTree) FixParentRelationship() (model.Trace, error) {
	length := t.Len()
	newTrace := model.Trace{}
	err := deepcopier.Copy(t.trace).To(&newTrace)
	if err != nil {
		return model.Trace{}, nil
	}
	newSpans := make([]model.Span, length)
	ii := 0
	t.Traverse(func(iv0 augmentedtree.Interval) {
		interval, ok := iv0.(Interval)
		if !ok {
			return
		}
		parent := t.ParentOf(interval)
		if parent.Span == nil { // no parent
			return
		}
		newSpan := model.Span{}
		err := deepcopier.Copy(*interval.Span).To(&newSpan)
		if err != nil {
			return
		}
		newSpan.ParentSpanID = parent.SpanID
		newSpans[ii] = newSpan
		ii++
	})
	newTrace.Spans = newSpans
	return newTrace, nil
}

func (t IntervalTree) FilterByDepth(depth int) (model.Trace, error) {
	newTrace := model.Trace{}
	err := deepcopier.Copy(t.trace).To(&newTrace)
	if err != nil {
		return model.Trace{}, nil
	}
	newSpans := []model.Span{}
	t.Traverse(func(iv0 augmentedtree.Interval) {
		interval, ok := iv0.(Interval)
		if !ok {
			return
		}
		parents := t.ParentsOf(interval)
		if len(parents) > depth {
			return
		}
		newSpan := model.Span{}
		err := deepcopier.Copy(*interval.Span).To(&newSpan)
		if err != nil {
			return
		}
		newSpans = append(newSpans, newSpan)
	})
	newTrace.Spans = newSpans
	return newTrace, nil
}
