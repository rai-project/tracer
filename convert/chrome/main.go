//+build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go-sourcepath"
	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/evaluation"
	"github.com/rai-project/tracer/convert/chrome"
)

func convert(path string) ([]byte, error) {

	if !com.IsFile(path) {
		return nil, errors.Errorf("trace %v does not exist", path)
	}
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read trace from %v", path)
	}

	trace := evaluation.TraceInformation{}
	err = json.Unmarshal(bts, &trace)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse trace")
	}
	bts, err = chrome.Marshal(trace.Traces[0])
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func main() {
	tr, err := convert(os.Args[1])
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
