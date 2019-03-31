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

	sp := tree.GetIntervalsByOperationName("Dropout")
	assert.NotEmpty(t, sp)

	firstDropoutSpan := sp[0]

	children := tree.ChildrenOf(firstDropoutSpan)
	for _, c := range children {
		pp.Printf("children = %v\n", c.SpanID)
	}

	intervals := tree.GetIntervals()
	for _, interval := range intervals {
		if firstDropoutSpan.Contains(interval) {
			pp.Printf("e children = %v\n", interval.OperationName)
		}
	}

	spans := tree.trace.Spans
	for _, span := range spans {
		if firstDropoutSpan.StartTime <= span.StartTime &&
			(span.StartTime+span.Duration) <= (firstDropoutSpan.StartTime+firstDropoutSpan.Duration) {
			pp.Printf("manual children = %v\n", span.OperationName)
		}
	}

	pp.Println(firstDropoutSpan.OperationName)

	// pp.Println(int64(firstDropoutSpan.Start()))
	// pp.Println(int64(firstDropoutSpan.Duration))
}
