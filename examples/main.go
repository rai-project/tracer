package main

import (
	"context"
	"fmt"
	"time"

	tracing "github.com/rai-project/tracer/zipkin"
)

func main() {
	// Create a new named tracer for this operation
	// Configure p3sr-trace to use it for spans that come from here
	tracer, closer := tracing.New("mlaas-client-inference")
	defer func() {
		closer.Close()
		fmt.Println("closed")
	}()
	tracing.InitGlobalTracer(tracer)

	// Create a root span for this command
	ctx := context.Background()
	sp, ctx := tracing.StartSpanFromContext(ctx, "inference_root")
	time.Sleep(time.Millisecond * 500)
	defer sp.Finish()

	readSp, _ := tracing.StartSpanFromContext(ctx, "read_files")
	time.Sleep(time.Second)
	defer readSp.Finish()

	// class, err := client.Inference(ctx, data, userID, modelID)
	// if err != nil {
	// 	log.WithError(err).Fatal()
	// }
	// log.Info("Class : ", class)
	return
}
