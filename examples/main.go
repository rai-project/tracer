package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	tracer "github.com/rai-project/tracer"
	zipkin "github.com/rai-project/tracer/zipkin"
)

func uServiceCall(ctx context.Context) {
	time.Sleep(time.Millisecond * 125) // slow startup
	sg, _ := tracer.StartSegmentFromContext(ctx, "service call")
	defer sg.Finish()

	time.Sleep(time.Millisecond * 333) // some work
}

func main() {
	config.Init(
		config.VerboseMode(true),
		config.DebugMode(true),
		config.ColorMode(true),
	)

	fmt.Println("here")
	// choose which tracing backend to use
	tr := zipkin.NewTracer("test-tracer")

	// Use that tracer
	tracer.SetGlobal(tr)

	// make sure the tracer finishes its tracing when we're done
	defer tr.Close()

	// Create some demo segments
	ctx := context.Background()
	rootSg, ctx := tracer.StartSegmentFromContext(ctx, "root_segment")
	defer rootSg.Finish()
	time.Sleep(time.Millisecond * 500)

	childSg, _ := tracer.StartSegmentFromContext(ctx, "child_segment")
	time.Sleep(time.Second)
	childSg.Finish()

	uServiceCall(ctx)
}

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "tracer/examples")
	})
}
