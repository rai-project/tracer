package jaeger

import (
	"context"
	"errors"
	"io"
	"runtime"
	"strings"

	//machineinfo "github.com/rai-project/machine/info"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	zpk "github.com/uber/jaeger-client-go/zipkin"

	"github.com/rai-project/config"
	osinfo "github.com/rai-project/machine/os"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/defaults"
	"github.com/rai-project/tracer/observer"
	raiutils "github.com/rai-project/utils"
	"github.com/rai-project/uuid"
)

type Tracer struct {
	opentracing.Tracer
	id          string
	closer      io.Closer
	endpoints   []string
	serviceName string
	initialized bool
	usingPerf   bool
	level       tracer.Level
}

func New(serviceName string, opts ...tracer.Option) (tracer.Tracer, error) {
	tracer := &Tracer{
		level: Config.Level,
	}
	err := tracer.Init(serviceName)
	if err != nil {
		return nil, nil
	}
	return tracer, nil
}

func (t *Tracer) ID() string {
	return t.id
}

func (t *Tracer) Level() tracer.Level {
	return t.level
}

func (t *Tracer) SetLevel(lvl tracer.Level) {
	t.level = lvl
}

func (t *Tracer) Init(serviceName string, opts ...tracer.Option) error {
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
	var trans jaeger.Transport
	var err error
	if strings.HasPrefix(endpoints[0], "udp://") {
		trans, err = jaeger.NewUDPTransport(
			strings.TrimPrefix(endpoints[0], "udp://"),
			0,
		)
	} else {
		trans, err = zipkin.NewHTTPTransport(
			endpoints[0],
			zipkin.HTTPBatchSize(200),
			zipkin.HTTPLogger(log),
		)
	}
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
		return err
	}

	// Adds support for injecting and extracting Zipkin B3 Propagation HTTP headers, for use with other Zipkin collectors.
	zipkinPropagator := zpk.NewZipkinB3HTTPHeaderPropagator()
	_ = zipkinPropagator

	tracerOpts := []jaeger.TracerOption{
		jaeger.TracerOptions.Tag("app", config.App.Name),
		jaeger.TracerOptions.Tag("arch", runtime.GOARCH),
		jaeger.TracerOptions.Tag("os", runtime.GOOS),
		jaeger.TracerOptions.Tag("os_info", osinfo.Info()),
		jaeger.TracerOptions.Tag("host", raiutils.GetHostIP()),
		jaeger.TracerOptions.Tag("commit_id", config.App.Version.GitCommit),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Logger(log),
		// jaeger.TracerOptions.ContribObserver(contribObserver),
		jaeger.TracerOptions.Gen128Bit(true),
		// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
		jaeger.TracerOptions.MaxTagValueLength(1024),
		// jaeger.TracerOptions.PoolSpans(true),
	}

	tracerOpts = append(tracerOpts, opts...)

	//if machineinfo.Info != nil {
	//	buf := new(bytes.Buffer)
	//	if err := json.Marshal(machineinfo.Info, buf); err == nil {
	//		tracerOpts = append(tracerOpts, jaeger.TracerOptions.Tag("machine", string(buf)))
	//	}
	//}

	t.usingPerf = false
	if runtime.GOOS == "linux" {
		for _, o := range observer.Config.ObserverNames {
			if o == "perf" || o == "perf_events" || o == "perfevents" {
				t.usingPerf = true
				break
			}
		}
	}
	if t.usingPerf {
		tracerOpts = append(tracerOpts, jaeger.TracerOptions.Tag("perfevents", defaults.PerfEvents))
	}

	for _, observer := range observer.Config.Observers {
		tracerOpts = append(tracerOpts, jaeger.TracerOptions.ContribObserver(&wrapObserver{observer}))
	}

	tr, cl := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true /*sample all*/),
		jaeger.NewRemoteReporter(trans,
			jaeger.ReporterOptions.QueueSize(10000),
			// jaeger.ReporterOptions.Logger(log),
			// jaeger.ReporterOptions.BufferFlushInterval(time.Second),
		),
		tracerOpts...,
	)

	t.id = uuid.NewV4()
	t.closer = cl
	t.endpoints = endpoints
	t.Tracer = tr
	t.serviceName = serviceName

	return nil
}

// startSpanFromContextWithTracer is factored out for testing purposes.
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := t.StartSpan(operationName, opts...)
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
