//+build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go-sourcepath"
	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer/convert/flame"
)

func main() {
	pp.WithLineInfo = true
	tr, err := flame.ConvertTraceFile(os.Args[1])
	if err != nil {
		pp.Println(err)
		os.Exit(1)
	}
	outputJSONFile := filepath.Join(sourcepath.MustAbsoluteDir(), "_fixtures", "example_trace.json")
	err = ioutil.WriteFile(outputJSONFile, tr, 0600)
	if err != nil {
		pp.Println(err)
		os.Exit(1)
	}

	outputHTMLFile := filepath.Join(sourcepath.MustAbsoluteDir(), "_fixtures", "example_trace.html")
	wr := &bytes.Buffer{}
	flame.GenerateHtml(wr, "flame", string(tr))
	err = ioutil.WriteFile(outputHTMLFile, wr.Bytes(), 0600)
	if err != nil {
		pp.Println(err)
		os.Exit(1)
	}
}
