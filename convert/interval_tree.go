package convert

import (
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
