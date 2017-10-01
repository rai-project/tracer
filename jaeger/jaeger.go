package jaeger

import (
	"errors"
	"io"
	"runtime"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/defaults"
	"github.com/rai-project/tracer/observer"
	"github.com/rai-project/uuid"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-lib/metrics"
	context "golang.org/x/net/context"

	jaeger "github.com/uber/jaeger-client-go"

	zpk "github.com/uber/jaeger-client-go/zipkin"
)

type Tracer struct {
	opentracing.Tracer
	id          string
	closer      io.Closer
	endpoints   []string
	serviceName string
	initialized bool
	usingPerf   bool
}

func New(serviceName string) (tracer.Tracer, error) {
	tracer := &Tracer{}
	err := tracer.Init(serviceName)
	if err != nil {
		return nil, nil
	}
	return tracer, nil
}

func (t *Tracer) ID() string {
	return t.id
}

func (t *Tracer) Init(serviceName string) error {
	if t.initialized {
		return nil
	}
	defer func() {
		t.initialized = true
	}()
	Config.Wait()
	endpoints := Config.Endpoints
	if len(endpoints) == 0 {
		return errors.New("no endpoints defined for jaeger tracer")
	}
	trans, err := zipkin.NewHTTPTransport(
		endpoints[0],
		zipkin.HTTPBatchSize(10),
		// zipkin.HTTPLogger(log),
	)
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
		return err
	}

	metricsFactory := metrics.NewLocalFactory(0)

	// Adds support for injecting and extracting Zipkin B3 Propagation HTTP headers, for use with other Zipkin collectors.
	zipkinPropagator := zpk.NewZipkinB3HTTPHeaderPropagator()

	tracerOpts := []jaeger.TracerOption{
		jaeger.TracerOptions.Tag("app", config.App.Name),
		jaeger.TracerOptions.Tag("perfevents", defaults.PerfEvents),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Metrics(jaeger.NewMetrics(metricsFactory, map[string]string{"lib": "jaeger"})),
		jaeger.TracerOptions.Logger(log),
		// jaeger.TracerOptions.ContribObserver(contribObserver),
		jaeger.TracerOptions.Gen128Bit(false),
		// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	}

	for _, observer := range observer.Config.Observers {
		tracerOpts = append(tracerOpts, jaeger.TracerOptions.ContribObserver(&wrapObserver{observer}))
	}

	tr, cl := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true /*sample all*/),
		jaeger.NewRemoteReporter(trans),
		tracerOpts...,
	)

	t.id = uuid.NewV4()
	t.closer = cl
	t.endpoints = endpoints
	t.Tracer = tr
	t.serviceName = serviceName

	t.usingPerf = false
	if runtime.GOOS == "linux" {
		for _, o := range observer.Config.ObserverNames {
			if o == "perf" || o == "perf_events" {
				t.usingPerf = true
				break
			}
		}
	}

	return nil
}

// startSpanFromContextWithTracer is factored out for testing purposes.
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if t.usingPerf {
		opts = append([]opentracing.StartSpanOption{opentracing.Tag{"perfevents", defaults.PerfEvents}}, opts...)
	}
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = t.StartSpan(operationName, opts...)
	} else {
		span = t.StartSpan(operationName, opts...)
	}
	return span, opentracing.ContextWithSpan(ctx, span)
}

func (t *Tracer) Close() error {
	return t.closer.Close()
}

func (t *Tracer) Endpoints() []string {
	return t.endpoints
}

func (t *Tracer) Name() string {
	return "jaeger::" + t.serviceName
}

func init() {
	tracer.Register("jaeger", &Tracer{}, New)
}
