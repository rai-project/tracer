// +build linux

package observer

import (
	"github.com/opentracing-contrib/go-observer"
	perfevents "github.com/opentracing-contrib/perfevents/go"
)

var (
	PerfEvents otobserver.Observer = perfevents.NewObserver()
)
