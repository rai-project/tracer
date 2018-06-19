package zipkin

import (
	"errors"
	"io"

	context "context"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/rai-project/tracer"
	"github.com/rai-project/utils"
	"github.com/rai-project/uuid"
)

type Tracer struct {
	initialized bool
	id          string
	opentracing.Tracer
	closer      io.Closer
	endpoints   []string
	serviceName string
	level       tracer.Level
}

func neverSample(_ uint64) bool { return false }

func alwaysSample(_ uint64) bool { return true }

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

func (t *Tracer) Level() tracer.Level {
	return Config.Level
}

func (t *Tracer) SetLevel(lvl tracer.Level) {
	t.level = lvl
}

func (t *Tracer) Init(serviceName string) error {
	if t.initialized {
		return nil
	}
	defer func() {
		t.initialized = true
	}()
	Config.Wait()

	inDebugModeQ := false
	endpoints := Config.Endpoints
	if len(endpoints) == 0 {
		return errors.New("no endpoints defined for zipkin tracer")
	}
	collector, err := zipkin.NewHTTPCollector(
		endpoints[0],
		zipkin.HTTPBatchSize(100),
		zipkin.HTTPLogger(log),
	)
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
		return err
	}

	recorder := zipkin.NewRecorder(collector, inDebugModeQ, utils.GetHostIP(), serviceName)

	// Create our tracer.
	tr, err := zipkin.NewTracer(
		recorder,
		zipkin.WithSampler(alwaysSample),
		// zipkin.WithLogger(log),
		zipkin.TraceID128Bit(true),
		zipkin.DebugMode(inDebugModeQ),
	)

	t.id = uuid.NewV4()
	t.Tracer = tr
	t.closer = collector
	t.endpoints = endpoints
	t.serviceName = serviceName

	return nil
}

// startSpanFromContextWithTracer is factored out for testing purposes.
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	var span opentracing.Span
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
	return "zipkin::" + t.serviceName
}

func init() {
	tracer.Register("zipkin", &Tracer{}, New)
}
