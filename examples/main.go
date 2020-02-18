package main

import (
	"context"
	"time"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	_ "github.com/rai-project/tracer/noop"
	"github.com/sirupsen/logrus"
)

var (
	log   *logrus.Entry
	conf2 = `tracer:
  enabled: true
  backend: zipkin
  endpoints:
    - http://localhost:9411/api/v1/spans
`
	conf = `tracer:
  enabled: true
  provider: jaeger
  endpoints:
    - localhost
  level: FULL_TRACE
`
)

func uServiceCall1(ctx context.Context) context.Context {
	time.Sleep(time.Millisecond * 125) // slow startup
	span, ctx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, "service call 1")
	defer span.Finish()

	time.Sleep(time.Millisecond * 333) // some local work

	ctx = uServiceCall2(ctx)

	return ctx
}

func uServiceCall2(ctx context.Context) context.Context {
	time.Sleep(time.Millisecond * 125) // slow startup
	span, ctx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, "service call 2")
	defer span.Finish()

	time.Sleep(time.Millisecond * 333) // some local work

	return ctx
}

func main() {

	log = logger.New().WithField("pkg", "tracer/examples")

	// choose which tracing backend to use
	tr, err := tracer.New("test-tracer")
	if err != nil {
		log.Fatal(err)
	}

	// Use that tracer
	tracer.SetStd(tr)

	// make sure the tracer finishes its tracing when we're done
	defer tracer.Close()

	// Create some demo segments
	ctx := context.Background()
	rootSg, ctx := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, "root_segment")
	defer rootSg.Finish()
	time.Sleep(time.Millisecond * 500)

	childSg, _ := tracer.StartSpanFromContext(ctx, tracer.APPLICATION_TRACE, "child_segment")
	time.Sleep(time.Second)
	childSg.Finish()

	uServiceCall1(ctx)
}

func init() {
	config.Init(
		config.VerboseMode(true),
		config.DebugMode(true),
		config.ColorMode(true),
		config.AppName("carml"),
		config.ConfigString(conf),
	)
}
