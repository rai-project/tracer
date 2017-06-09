package sourcepath

import (
	"strings"
	"testing"
)

func TestAbsoluteDir(t *testing.T) {
	expectedPathSuffix := `github.com/GeertJohan/go-sourcepath`

	path, err := AbsoluteDir()
	if err != nil {
		t.Fatalf("error getting absolute dir: %v", err)
	}

	if !strings.HasSuffix(path, expectedPathSuffix) {
		t.Fatalf("expected path ending with `%s` but got `%s`", expectedPathSuffix, path)
	}
}
