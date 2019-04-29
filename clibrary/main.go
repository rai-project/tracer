package main

import "C"

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/clibrary/env"
	"github.com/rai-project/utils"
	jaeger "github.com/uber/jaeger-client-go"

	"github.com/k0kubun/pp"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	IsDebug        bool
	IsVerbose      bool
	AppSecret      string
	CfgFile        string
	tracerInitOnce sync.Once
	globalSpan     opentracing.Span
	globalCtx      context.Context = context.Background()
	log            *logrus.Entry   = logrus.New().WithField("pkg", "tracer/clibrary")
)

//export TracerSetLevel
func TracerSetLevel(lvl int32) {
	tracer.SetLevel(tracer.Level(lvl))
}

//export TracerClose
func TracerClose() {
	libDeinit()
	tracer.Close()
}

//export TracerInit
func TracerInit() {
	tracerInitOnce.Do(doTracerInit)
}

func doTracerInit() {
	pp.Println("initializing library")
	log.Level = logrus.DebugLevel
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "tracer/clibrary")
	})

	color.NoColor = false
	opts := []config.Option{
		config.AppName("carml"),
		config.ColorMode(true),
		config.DebugMode(IsDebug),
		config.VerboseMode(IsVerbose),
	}

	pp.WithLineInfo = true

	if c, err := homedir.Expand(CfgFile); err == nil {
		CfgFile = c
	}
	if c, err := filepath.Abs(CfgFile); err == nil {
		CfgFile = c
	}
	opts = append(opts, config.ConfigFileAbsolutePath(CfgFile))

	if AppSecret != "" {
		opts = append(opts, config.AppSecret(AppSecret))
	}

	config.Init(opts...)

	// tracer.ResetStd(
	// 	jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, NewEnvPropagator(BaggagePrefix("rai:)"))),
	// 	jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, NewEnvPropagator(BaggagePrefix("rai:)"))),
	// )

	tracer.SetLevel(tracer.FULL_TRACE)
	libInit()
}

func initLib() {
	TracerInit()
}

func libInit() {
	extraOpts := []opentracing.StartSpanOption{}

	parentCtx, err := env.GetSpanContext()
	if err == nil {
		extraOpts = append(extraOpts, opentracing.ChildOf(parentCtx))
	}

	globalSpan, globalCtx = tracer.StartSpanFromContext(
		globalCtx,
		tracer.APPLICATION_TRACE,
		"api_tracing",
		extraOpts...,
	)

	traceIDVal := globalSpan.Context().(jaeger.SpanContext).TraceID().String()

	fmt.Printf("Got traceid = %s....\n", traceIDVal)

	if err != nil {
		env.Set(globalSpan.Context().(jaeger.SpanContext))
		initCupti()
	}

}

func libDeinit() {
	if false {
		pp.Println("deinit")
	}
	deinitCupti()
	if globalSpan != nil {
		globalSpan.Finish()
		pp.Println("closing global span")

		traceID := globalSpan.Context().(jaeger.SpanContext).TraceID()
		traceIDVal := traceID.String()

		ip, _ := utils.GetExternalIp()
		pp.Println(fmt.Sprintf("http://%s:16686/trace/%v", ip, traceIDVal))
		globalSpan = nil
	}
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}
