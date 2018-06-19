// +build !darwin,!arm64
// +build nogpu

package observer

var (
	GPUMemInfo = NoOp
)
