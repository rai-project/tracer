package main

import (
	"os"
	"runtime"

	"github.com/k0kubun/pp"
	"github.com/rai-project/go-cupti"
	"github.com/spf13/cast"
)

var (
	cuptiInstance *cupti.CUPTI
)

func initCupti() {
	if false {
		pp.Println("CUPTI_TRACE = ", os.Getenv("CUPTI_TRACE"))
		pp.Println("CUDA_TRACE = ", os.Getenv("CUDA_TRACE"))
	}
	if runtime.GOOS != "linux" {
		return
	}
	enableCupti := cast.ToBool(os.Getenv("CUPTI_TRACE")) || cast.ToBool(os.Getenv("CUDA_TRACE"))
	if !enableCupti {
		return
	}

	// pp.Println("initializing cupti")

	cu, err := cupti.New(cupti.Context(globalCtx), cupti.SamplingPeriod(0))
	if err != nil {
		panic(err)
	}
	cuptiInstance = cu
}

func deinitCupti() {
	if cuptiInstance == nil {
		return
	}
	cuptiInstance.Close()
}
