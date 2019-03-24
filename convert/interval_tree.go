package convert

import (
	"sort"

	"github.com/Workiva/go-datastructures/augmentedtree"
	model "github.com/uber/jaeger/model/json"
)

type SpanTree struct {
	augmentedtree.Tree
}

func NewSpanTree(spans []model.Span) SpanTree {
	tree := augmentedtree.New(1)
	for _, s := range spans {
		tree.Add(spanToInterval(s))
	}
	return SpanTree{tree}
}

func (t SpanTree) ParentsOf(sp Interval) []Interval {
	elems := t.Query(sp)
	res := []Interval{}
	for _, e := range elems {
		if e.ID() != sp.ID() {
			continue
		}
		res = append(res, e.(Interval))
	}
	return res
}

func (t SpanTree) ParentOf(sp Interval) Interval {
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
