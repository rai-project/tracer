// +build !develop

package zipkin

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"

	config "github.com/rai-project/config"
)

var stdCloser io.Closer

const (
	endpoint = "http://ec2-34-225-193-157.compute-1.amazonaws.com:9411/api/v1/spans"
	enable   = true
)

// New returns a new named tracer and closer
func New(serviceName string) (opentracing.Tracer, io.Closer) {
	trans, err := zipkin.NewHTTPTransport(
		endpoint,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
	}
	return jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true), // sample all traces
		jaeger.NewRemoteReporter(trans),
	)
}

// New2 returns a new named tracer and closer
func New2(serviceName string) (opentracing.Tracer, io.Closer) {
	return jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true), // sample all traces
		jaeger.NewLoggingReporter(jaeger.StdLogger),
	)
}

// func StartSpan(operationName string) opentracing.Span {
// 	return Tracer.StartSpan(operationName)
// }

// Globals returns the global tracer and closer
func Globals() (opentracing.Tracer, io.Closer) {
	return opentracing.GlobalTracer(), stdCloser
}

func StartSpan(operationName string) opentracing.Span {
	return opentracing.StartSpan(operationName)
}

func ContextWithSpan(c context.Context, sp opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(c, sp)
}

func InitGlobalTracer(tracer opentracing.Tracer) {
	opentracing.InitGlobalTracer(tracer)
}

func StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName)
}

func init() {
	config.AfterInit(func() {
		if !enable {
			log.Debug("Tracing disabled.")
			return
		}

		log.Debug("Zipkin endpoint: ", endpoint)
		trans, err := zipkin.NewHTTPTransport(
			endpoint,
			zipkin.HTTPBatchSize(1),
			zipkin.HTTPLogger(jaeger.StdLogger),
		)
		if err != nil {
			log.WithError(err).Error("Cannot initialize HTTP transport")
		}

		var tracer opentracing.Tracer
		tracer, stdCloser = jaeger.NewTracer(
			"global-tracer",
			jaeger.NewConstSampler(true), // sample all traces
			jaeger.NewRemoteReporter(trans),
		)
		opentracing.InitGlobalTracer(tracer)

	})

}
