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
	config.Init()

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

func main2() {
	// Create a new named tracer for this operation
	tr, closer := zipkin.New("mlaas-client-inference")
	defer func() {
		closer.Close()
		fmt.Println("closed")
	}()
	zipkin.InitGlobalTracer(tr)

	// Create a root span for this command
	ctx := context.Background()
	sp, ctx := zipkin.StartSpanFromContext(ctx, "inference_root")
	time.Sleep(time.Millisecond * 500)
	defer sp.Finish()

	readSp, _ := zipkin.StartSpanFromContext(ctx, "read_files")
	time.Sleep(time.Second)
	defer readSp.Finish()

	// class, err := client.Inference(ctx, data, userID, modelID)
	// if err != nil {
	// 	log.WithError(err).Fatal()
	// }
	// log.Info("Class : ", class)
	return
}

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "tracer/examples")
	})
}
