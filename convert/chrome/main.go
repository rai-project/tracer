//+build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/evaluation"
	"github.com/rai-project/tracer/convert/chrome"
)

func convert(path string) (string, error) {

	if !com.IsFile(path) {
		return "", errors.Errorf("trace %v does not exist", path)
	}
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read trace from %v", path)
	}

	trace := evaluation.TraceInformation{}
	err = json.Unmarshal(bts, &trace)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse trace")
	}
	bts, err = chrome.Marshal(trace.Traces[0])
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func main() {
	tr, err := convert(os.Args[1])
	if err != nil {
		pp.Println(err)
		os.Exit(1)
	}
	fmt.Println(tr)
}
