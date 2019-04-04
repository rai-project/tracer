//+build ignore

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go-sourcepath"
	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/tracer/convert/flame"
	model "github.com/uber/jaeger/model/json"
)

func convert(path string) ([]byte, error) {

	if !com.IsFile(path) {
		return nil, errors.Errorf("trace %v does not exist", path)
	}
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read trace from %v", path)
	}

	trace := model.Trace{}
	err = json.Unmarshal(bts, &trace)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse trace")
	}
	bts, err = flame.Marshal(trace)
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
