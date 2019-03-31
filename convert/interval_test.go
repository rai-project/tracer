package convert

import (
	"path/filepath"
	"testing"

	"github.com/GeertJohan/go-sourcepath"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

var (
	sourceDir    = sourcepath.MustAbsoluteDir()
	fixturesDir  = filepath.Join(sourceDir, "_fixtures")
	exampleTrace = filepath.Join(fixturesDir, "example_trace.json")
)

func TestNewIntervalTree(t *testing.T) {
	tree, err := NewIntervalTreeFromTraceFile(exampleTrace)
	assert.NoError(t, err)
	assert.NotEmpty(t, tree)

	if !assert.Equal(t, int(tree.Len()), len(tree.trace.Spans)) {
		t.FailNow()
	}

	sp := tree.GetIntervalsByOperationName("Dropout")
	assert.NotEmpty(t, sp)

	firstDropoutSpan := sp[0]

	children := tree.ChildrenOf(firstDropoutSpan)
	assert.NotEmpty(t, children)
	for _, c := range children {
		assert.NotEmpty(t, c)
	}

	intervals := tree.GetIntervals()
	for _, interval := range intervals {
		if firstDropoutSpan.Contains(interval) {
			assert.NotEmpty(t, interval)
		}
	}

	spans := tree.trace.Spans
	for _, span := range spans {
		if firstDropoutSpan.StartTime <= span.StartTime &&
			(span.StartTime+span.Duration) <= (firstDropoutSpan.StartTime+firstDropoutSpan.Duration) {
			assert.NotEmpty(t, span)
		}
	}

	firstDropoutChild := tree.ChildrenOf(firstDropoutSpan)[0]
	parents := tree.ParentsOf(firstDropoutChild)
	assert.NotEmpty(t, parents)
	for _, p := range parents {
		assert.NotEmpty(t, p)
	}

	firstDropoutChildParent := tree.ParentOf(firstDropoutChild)
	assert.NotEmpty(t, firstDropoutChildParent)
	pp.Println(firstDropoutChildParent.OperationName)
	// pp.Println(firstDropoutChildParent.OperationName)

	immediateDropoutChilden := tree.ImmediateChildrenOf(firstDropoutSpan)
	assert.NotEmpty(t, immediateDropoutChilden)
	for _, c := range immediateDropoutChilden {
		assert.NotEmpty(t, c)
	}

	// pp.Println(int64(firstDropoutSpan.Start()))
	// pp.Println(int64(firstDropoutSpan.Duration))
}
