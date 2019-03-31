package convert

import (
	"errors"
	"strings"
)

func (t *IntervalTree) FilterOnlyChildrenOf(operationName string) error {
	var rootParent Interval
	allIntervals := t.GetIntervals()
	operationName = strings.ToLower(operationName)
	for _, interval := range allIntervals {
		if strings.ToLower(interval.OperationName) == operationName {
			rootParent = interval
			break
		}
	}
	if rootParent.IsNil() {
		return errors.New("unable to find root node")
	}
	children := t.ChildrenOf(rootParent)
	childrenIds := map[string]bool{}
	for _, child := range children {
		childrenIds[string(child.SpanID)] = true
	}

	for _, interval := range allIntervals {
		if interval.SpanID == rootParent.SpanID {
			continue
		}
		if ok := childrenIds[string(interval.SpanID)]; ok {
			continue
		}
		t.Delete(interval)
	}

	return nil
}
