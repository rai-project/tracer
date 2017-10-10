//+build !darwin

package observer

import (
	"runtime"

	opentracing "github.com/opentracing/opentracing-go"
)

var (
	Instruments = NoOp
)
