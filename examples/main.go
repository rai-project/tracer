package main

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	tracer "github.com/rai-project/tracer"
	zipkin "github.com/rai-project/tracer/zipkin"
)

func uServiceCall1(ctx context.Context) context.Context {
	time.Sleep(time.Millisecond * 125) // slow startup
	span, ctx := tracer.StartSpanFromContext(ctx, "service call 1")
	defer span.Finish()

	time.Sleep(time.Millisecond * 333) // some local work

	ctx = uServiceCall2(ctx)

	return ctx
}

func uServiceCall2(ctx context.Context) context.Context {
	time.Sleep(time.Millisecond * 125) // slow startup
	span, ctx := tracer.StartSpanFromContext(ctx, "service call 2")
	defer span.Finish()

	time.Sleep(time.Millisecond * 333) // some local work

	return ctx
}

func main() {
	config.Init(
		config.VerboseMode(true),
		config.DebugMode(true),
		config.ColorMode(true),
	)

	// make sure the tracer finishes its tracing when we're done
	defer tracer.Close()

	// Create some demo segments
	ctx := context.Background()
	rootSg, ctx := tracer.StartSpanFromContext(ctx, "root_segment")
	defer rootSg.Finish()
	time.Sleep(time.Millisecond * 500)

	childSg, _ := tracer.StartSpanFromContext(ctx, "child_segment")
	time.Sleep(time.Second)
	childSg.Finish()

	uServiceCall1(ctx)
}

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "tracer/examples")

		// choose which tracing backend to use
		tr, err := zipkin.NewTracer("test-tracer")
		if err != nil {
			log.Fatal(err)
		}

		// Use that tracer
		tracer.SetStd(tr)
	})
}
