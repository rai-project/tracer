//+build ignore

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go-sourcepath"
	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer/convert/chrome"
)

func main() {
	pp.WithLineInfo = true
	tr, err := chrome.ConvertFile(os.Args[1])
	if err != nil {
		pp.Println(err)
		os.Exit(1)
	}
	outputFile := filepath.Join(sourcepath.MustAbsoluteDir(), "_fixtures", "example_trace.json")
	err = ioutil.WriteFile(outputFile, tr, 0600)
	if err != nil {
		pp.Println(err)
		os.Exit(1)
	}
}
