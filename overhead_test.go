package tracer

import (
	"context"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
	_ "github.com/rai-project/tracer/jaeger"
)

func BenchmarkTracer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var span opentracing.Span
		ctx := context.Background()
		span, ctx = StartSpanFromContext(ctx, FULL_TRACE, "test_run")
		time.Sleep(time.Second)
		span.Finish()
		Close()
	}
}

func init() {
	config.Init(
		config.AppName("carml"),
		config.DebugMode(true),
		config.VerboseMode(true),
	)
}
