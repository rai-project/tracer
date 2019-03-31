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
		if e.ID() != sp.ID() {
			continue
		}
		if sp.Contains(e) {
			res = append(res, e)
		}
	}
	return res
}

func (t IntervalTree) ParentsOf(sp Interval) []Interval {
	elems := t.Query(sp)
	res := []Interval{}
	for _, e0 := range elems {
		e, ok := e0.(Interval)
		if !ok {
			continue
		}
		if e.ID() != sp.ID() {
			continue
		}
		res = append(res, e)
	}
	return res
}

func (t IntervalTree) ParentOf(sp Interval) Interval {
	elems := t.Query(sp)
	if len(elems) == 0 {
		return Interval{}
	}
	if len(elems) == 1 {
		return elems[0].(Interval)
	}
	sort.Slice(elems, func(ii, jj int) bool {
		return elems[ii].LowAtDimension(0) < elems[jj].LowAtDimension(0)
	})
	return elems[0].(Interval)
}

// func (t IntervalTree) GetIntervals() []Interval {
// 	length := t.Len()
// 	intervals := make([]Interval, length)
// 	ii := 0
// 	t.Traverse(func(interval augmentedtree.Interval) {
// 		intervals[ii] = interval.(Interval)
// 	})
// 	return intervals
// }

func (t IntervalTree) GetIntervals() []Interval {
	spans := t.trace.Spans
	length := len(spans)
	intervals := make([]Interval, length)
	for ii, spans := range t.trace.Spans {
		intervals[ii] = spanToInterval(spans)
	}
	return intervals
}

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
