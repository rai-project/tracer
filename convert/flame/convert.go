package flame

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"time"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/tracer/convert"
	cnv "github.com/rai-project/tracer/convert"
	model "github.com/uber/jaeger/model/json"
)

type convertState struct {
	tree        *convert.IntervalTree
	jaegerTrace model.Trace
	root        convert.Interval
	nodes       map[string]*Node
}

func Marshal(trace model.Trace) ([]byte, error) {
	tr, err := Convert(trace)
	if err != nil {
		return nil, err
	}
	return json.Marshal(tr)
}

func Convert(tr model.Trace) (*Node, error) {
	st, err := newConvertState(tr)
	if err != nil {
		return nil, err
	}

	nd := st.convertSpans(nil, st.root, 0)

	return nd, nil
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

	rootInterval, err := tree.FilterOnlyChildrenOf("PredictStep")
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

	return &convertState{
		tree:        tree,
		root:        *rootInterval,
		jaegerTrace: jaegerTrace,
		nodes:       map[string]*Node{},
	}, nil
}

func (st *convertState) getValue(sp convert.Interval) int {
	root := st.root
	rootDuration := float64(time.Millisecond) * (float64(sp.Duration) / float64(root.Duration))
	return int(rootDuration)
}

func (st *convertState) convertSpans(rootNode *Node, root convert.Interval, depth int) *Node {
	rootID := string(root.SpanID)

	if root.OperationName == "cupti_new" {
		return nil
	}

	if _, ok := st.nodes[rootID]; ok {
		return nil
	}

	nd := &Node{
		ID:       rootID,
		Name:     cleanName(root.OperationName),
		Value:    0,
		Interval: &root,
		Children: []*Node{},
	}
	st.nodes[rootID] = nd

	children := convert.Intervals(st.tree.ImmediateChildrenOf(root))
	sort.Sort(children)

	for ii := len(children) - 1; ii >= 0; ii-- {
		child := children[ii]
		// for _, child := range children {
		e := st.convertSpans(nd, child, depth+1)
		if e == nil {
			continue
		}
		nd.Children = append(nd.Children, e)
		nd.Value += e.Value
	}

	if len(children) == 0 {
		nd.Value = st.getValue(root)
	}

	return nd
}

func ConvertTraceFile(path string) ([]byte, error) {

	if !com.IsFile(path) {
		return nil, errors.Errorf("trace %v does not exist", path)
	}
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read trace from %v", path)
	}

	trace := cnv.TraceInformation{}
	err = json.Unmarshal(bts, &trace)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse trace")
	}
	bts, err = Marshal(trace.Traces[0])
	if err != nil {
		return nil, err
	}
	return bts, nil
}
