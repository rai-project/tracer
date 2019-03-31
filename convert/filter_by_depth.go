package convert

import (
	"github.com/pkg/errors"
	model "github.com/uber/jaeger/model/json"
)

func FilterTraceByDepth(trace model.Trace, depth int) (model.Trace, error) {
	tree, err := NewIntervalTree(trace)
	if err != nil {
		return model.Trace{}, errors.Wrap(err, "failed to create interval tree for trace")
	}
	return tree.FilterByDepth(depth)
}
