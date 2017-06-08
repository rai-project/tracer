package zipkin

import (
	"errors"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

var opentracingGlobalTracerIsSet bool

type Tracer struct {
	opentracing.Tracer
	closer      io.Closer
	endpoint    string
	serviceName string
}

func NewTracer(serviceName string) (tracer.Tracer, error) {
	endpoint := Config.Endpoints[0]
	trans, err := zipkin.NewHTTPTransport(
		endpoint,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(log),
	)
	if err != nil {
		log.WithError(err).Error("Cannot initialize HTTP transport")
		return nil, err
	}
	tr, cl := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true /*sample all*/),
		jaeger.NewRemoteReporter(trans),
	)

	if _, ok := opentracing.GlobalTracer().(opentracing.NoopTracer); !ok {
		log.Error("Expecting global tracer to be uninitialized")
		return nil, errors.New("expecting global tracer to be uninitialized")
	}
	opentracing.SetGlobalTracer(tr)
	return &Tracer{Tracer: tr, closer: cl, endpoint: endpoint, serviceName: serviceName}, nil
}

func (t *Tracer) Close() error {
	return t.closer.Close()
}

func (t *Tracer) Endpoint() string {
	return t.endpoint
}

func (t *Tracer) Name() string {
	return "zipkin::" + t.serviceName
}
