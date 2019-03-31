package chrome

import "github.com/Workiva/go-datastructures/augmentedtree"

func NewTree(tr Trace) augmentedtree.Tree {
	tree := augmentedtree.New(1)
	for _, event := range tr.TraceEvents {
		tree.Add(event)
	}
	return tree
}
