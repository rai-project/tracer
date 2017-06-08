package zipkin

import (
	"errors"
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

var opentracingGlobalTracerIsSet bool

type Tracer struct {
	tracer      opentracing.Tracer
	closer      io.Closer
	endpoint    string
	serviceName string
}

func NewTracer(serviceName string) (tracer.Tracer, error) {
	endpoint := Config.Endpoints[0]
	trans, err := zipkin.NewHTTPTransport(
		endpoint,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
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
	return &Tracer{tracer: tr, closer: cl, endpoint: endpoint, serviceName: serviceName}, nil
}

func (t *Tracer) Close() error {
	return t.closer.Close()
}

func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(operationName, opts...)
}

func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	if req, ok := carrier.(*http.Request); ok {
		return t.tracer.Inject(
			sm,
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}
	return t.tracer.Inject(sm, format, carrier)
}

func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	if req, ok := carrier.(*http.Request); ok {
		wireContext, err := t.tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
		return wireContext, err
	}
	return t.tracer.Extract(format, carrier)
}

func (t *Tracer) Endpoint() string {
	return t.endpoint
}

func (t *Tracer) Name() string {
	return "zipkin::" + t.serviceName
}
