package flame

import (
	"encoding/json"

	"github.com/rai-project/tracer/convert"
	model "github.com/uber/jaeger/model/json"
)

type convertState struct {
	tree         *convert.IntervalTree
	jaegerTrace  model.Trace
	profile      *Profile
	root         convert.Interval
	nodes        []*Node
	childNodes   map[string][]*Node
	visitedNodes map[string]bool
}

func Marshal(trace model.Trace) ([]byte, error) {
	tr, err := Convert(trace)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(tr, "", "  ")
}

func Convert(tr model.Trace) (*Profile, error) {
	st, err := newConvertState(tr)
	if err != nil {
		return nil, err
	}

	nd := st.convertSpans(st.root)
	st.fixChildren()

	st.profile.RootNode = *nd

	return st.profile, nil
}

func newConvertState(tr model.Trace) (*convertState, error) {
	tr, err := convert.FixTrace(tr)
	if err != nil {
		return nil, err
	}
	tree, err := convert.NewIntervalTree(tr)
	if err != nil {
		return nil, err
	}
	err = tree.FilterOnlyChildrenOf("PredictImage")
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

	profile := &Profile{}

	return &convertState{
		tree:         tree,
		root:         tree.MaxInterval(),
		jaegerTrace:  jaegerTrace,
		profile:      profile,
		nodes:        []*Node{},
		childNodes:   map[string][]*Node{},
		visitedNodes: map[string]bool{},
	}, nil
}

func (st *convertState) getValue(sp convert.Interval) int {
	root := st.root
	rootDuration := 1000 * (float64(sp.Duration) / float64(root.Duration))
	return int(rootDuration)
}

func (st *convertState) convertSpans(root convert.Interval) *Node {
	rootID := string(root.SpanID)
	nd := &Node{
		ID:    rootID,
		Name:  root.OperationName,
		Value: st.getValue(root),
	}
	st.visitedNodes[rootID] = true

	for _, interval := range st.tree.GetIntervals() {
		intervalID := string(interval.SpanID)
		if _, ok := st.visitedNodes[intervalID]; ok {
			continue
		}
		parentID := string(interval.ParentSpanID)
		if _, ok := st.childNodes[parentID]; !ok {
			st.childNodes[parentID] = []*Node{}
		}
		st.childNodes[parentID] = append(st.childNodes[parentID], st.convertSpans(interval))
	}
	// for _, child := range st.tree.ChildrenOf(root) {
	// 	childID := string(child.SpanID)
	// 	if _, ok := st.visitedNodes[childID]; ok {
	// 		continue
	// 	}
	// 	parentID := string(child.ParentSpanID)
	// 	if _, ok := st.childNodes[parentID]; !ok {
	// 		st.childNodes[parentID] = []*Node{}
	// 	}
	// 	st.childNodes[parentID] = append(st.childNodes[parentID], st.convertSpans(child))
	// }
	st.nodes = append(st.nodes, nd)
	return nd
}

func (st *convertState) fixChildren() {
	for _, nd := range st.nodes {
		nd.Children = map[string]*Node{}
		for _, child := range st.childNodes[nd.ID] {
			nd.Children[nd.Name] = child
		}
	}
}
