package jaeger

import (
	"errors"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-lib/metrics"

	jaeger "github.com/uber/jaeger-client-go"

	zpk "github.com/uber/jaeger-client-go/zipkin"
)

var opentracingGlobalTracerIsSet bool

type Tracer struct {
	opentracing.Tracer
	closer      io.Closer
	endpoints   []string
	serviceName string
	initialized bool
}

func New(serviceName string) (*Tracer, error) {
	tracer := &Tracer{}
	err := tracer.Init(serviceName)
	if err != nil {
		return nil, nil
	}
	return tracer, nil
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
		zipkin.HTTPBatchSize(100),
		// zipkin.HTTPLogger(log),
	)
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
		return err
	}

	metricsFactory := metrics.NewLocalFactory(0)

	// Adds support for injecting and extracting Zipkin B3 Propagation HTTP headers, for use with other Zipkin collectors.
	zipkinPropagator := zpk.NewZipkinB3HTTPHeaderPropagator()

	tr, cl := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true /*sample all*/),
		jaeger.NewRemoteReporter(trans),
		jaeger.TracerOptions.Tag("app", config.App.Name),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Metrics(jaeger.NewMetrics(metricsFactory, map[string]string{"lib": "jaeger"})),
		jaeger.TracerOptions.Logger(log),
		jaeger.TracerOptions.Observer(perfeventsObserver),
		// jaeger.TracerOptions.ContribObserver(contribObserver),
		jaeger.TracerOptions.Gen128Bit(true),
		// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)

	t.closer = cl
	t.endpoints = endpoints
	t.Tracer = tr
	t.serviceName = serviceName

	return nil
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
	tracer.Register("jaeger", &Tracer{})
}
