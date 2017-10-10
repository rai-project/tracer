// +build !darwin

package observer

import (
	"github.com/opentracing-contrib/go-observer"
	opentracing "github.com/opentracing/opentracing-go"
	nvml "github.com/rai-project/nvml-go"
	"github.com/spf13/cast"

	olog "github.com/opentracing/opentracing-go/log"
)

var (
	GPUMemInfo otobserver.Observer = newGPUMemInfo()
)

type gpuMemInfo struct {
	count   int
	handles []nvml.DeviceHandle
}

func newGPUMemInfo() gpuMemInfo {
	count, err := nvml.DeviceCount()
	if err != nil {
		panic(err)
	}
	handles := make([]nvml.DeviceHandle, count)
	for ii := range handles {
		handle, err := nvml.DeviceGetHandleByIndex(ii)
		if err != nil {
			panic(err)
		}
		handles[ii] = handle
	}
	return &gpuMemInfo{
		count:   count,
		handles: handles,
	}
}

// OnStartSpan creates a new gpuMemInfo for the span
func (o *gpuMemInfo) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) (otobserver.SpanObserver, bool) {
	return newMemInfoSpan(o, sp, options)
}

// SpanDummy collects perfevent metrics
type gpuMemInfoSpan struct {
	gpuMemInfo
	sp opentracing.Span
}

// NewSpanDummy creates a new SpanDummy that can emit perfevent
// metrics
func newMemInfoSpan(info *gpuMemInfo, s opentracing.Span, opts opentracing.StartSpanOptions) (*gpuMemInfoSpan, bool) {
	so := &gpuMemInfoSpan{
		gpuMemInfo: info,
		sp:         s,
	}
	for _, handle := range o.handles {
		meminfo, err := nvml.DeviceMemoryInformation(handle)
		if err != nil {
			return nil, false
		}
		sp.LogFields(
			olog.String("start_gpu_mem_used", cast.ToString(meminfo.Used)),
			olog.String("start_gpu_mem_free", cast.ToString(meminfo.Free)),
			olog.String("start_gpu_mem_total", cast.ToString(meminfo.Total)),
		)
	}

	return so, true
}

func (so *gpuMemInfoSpan) OnSetOperationName(operationName string) {
}

func (so *gpuMemInfoSpan) OnSetTag(key string, value interface{}) {
}

func (so *gpuMemInfoSpan) OnFinish(options opentracing.FinishOptions) {
	for _, handle := range so.handles {
		meminfo, err := nvml.DeviceMemoryInformation(handle)
		if err != nil {
			return nil, false
		}
		so.sp.LogFields(
			olog.String("finish_gpu_mem_used", cast.ToString(meminfo.Used)),
			olog.String("finish_gpu_mem_free", cast.ToString(meminfo.Free)),
			olog.String("finish_gpu_mem_total", cast.ToString(meminfo.Total)),
		)
	}
}
