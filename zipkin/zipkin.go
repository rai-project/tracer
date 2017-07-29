package zipkin

import (
	"errors"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer"
	"github.com/rai-project/utils"
)

type Tracer struct {
	opentracing.Tracer
	closer      io.Closer
	endpoints   []string
	serviceName string
}

func neverSample(_ uint64) bool { return false }

func alwaysSample(_ uint64) bool { return true }

func New(serviceName string) (*Tracer, error) {
	tracer := &Tracer{
	}
	err := tracer.Init(serviceName)
	if err != nil {
		return nil, nil
	}
	return tracer, nil
}

func (t *Tracer) Init(serviceName string) error {
	inDebugModeQ := config.App.IsDebug
	endpoints := Config.Endpoints
	if len(endpoints) == 0 {
		return errors.New("no endpoints defined for zipkin tracer")
	}
	collector, err := zipkin.NewHTTPCollector(
		endpoints[0],
		zipkin.HTTPBatchSize(1),
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
		zipkin.WithLogger(log),
		zipkin.TraceID128Bit(true),
		zipkin.DebugMode(inDebugModeQ),
	)

	t.Tracer = tr
	t.closer = collector
	t.endpoints = endpoints
	t.serviceName := serviceName

	return nil
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
	tracer.AddTracer("zipkin", &Tracer{})
}
