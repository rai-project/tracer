package flame

import (
	"encoding/json"

	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer/convert"
	model "github.com/uber/jaeger/model/json"
)

type convertState struct {
	tree         *convert.IntervalTree
	jaegerTrace  model.Trace
	profile      *Profile
	root         convert.Interval
	nodes        map[string]*Node
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

	st.convertSpans(st.root)

	for _, nd := range st.nodes {
		println(nd.ID, " ", st.nodes[string(nd.Interval.ParentSpanID)].ID)
	}

	st.profile.RootNode = *st.nodes[string(st.root.SpanID)]
	pp.Println(len(st.nodes))

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

	jaegerTrace, err := tree.FixParentRelationship()
	if err != nil {
		return nil, err
	}

	tree, err = convert.NewIntervalTree(jaegerTrace)
	if err != nil {
		return nil, err
	}

	err = tree.FilterOnlyChildrenOf("PredictImage")
	if err != nil {
		return nil, err
	}

	jaegerTrace, err = tree.FixParentRelationship()
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
		nodes:        map[string]*Node{},
		visitedNodes: map[string]bool{},
	}, nil
}

func (st *convertState) getValue(sp convert.Interval) int {
	root := st.root
	rootDuration := 1000 * (float64(sp.Duration) / float64(root.Duration))
	return int(rootDuration)
}

func (st *convertState) convertSpans(root convert.Interval) *Node {
	nd := st.convertSpan(root)
	for _, interval := range st.tree.GetIntervals() {
		// pp.Println(interval.OperationName, st.visitedNodes[string(interval.SpanID)])
		st.convertSpan(interval)
	}

	return nd
}

func (st *convertState) convertSpan(root convert.Interval) *Node {
	rootID := string(root.SpanID)

	if val, ok := st.visitedNodes[rootID]; ok && val {
		return st.nodes[rootID]
	}

	nd := &Node{
		ID:       rootID,
		Name:     root.OperationName,
		Value:    st.getValue(root),
		Interval: &root,
		Children: map[string]*Node{},
	}

	parentID := string(root.ParentSpanID)
	pp.Println(parentID)
	if _, ok := st.nodes[parentID]; !ok {
		st.nodes[parentID] = &Node{
			Children: map[string]*Node{},
		}
	}

	st.nodes[parentID].Children[rootID] = nd

	// mergo.Merge(&nd, st.nodes[rootID])

	if _, ok := st.nodes[rootID]; ok {
		nd.Children = st.nodes[rootID].Children
	}

	st.nodes[rootID] = nd
	st.visitedNodes[rootID] = true
	return nd
}
