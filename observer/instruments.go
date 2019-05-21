// +build !darwin !cgo !instruments

package observer

var (
	Instruments = NoOp
)
