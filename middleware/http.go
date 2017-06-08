package middleware

import (
	"fmt"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing/_thrift/gen-go/zipkincore"
	"github.com/rai-project/tracer"
)

type tracingSpanContextKey struct{}

// RequestFunc is a middleware function for outgoing HTTP requests.
type RequestFunc func(req *http.Request) *http.Request

// ToHTTPRequest returns a RequestFunc that injects an OpenTracing Span found in
// context into the HTTP Headers. If no such Span can be found, the RequestFunc
// is a noop.
func ToHTTPRequest(tr tracer.Tracer) RequestFunc {
	return func(req *http.Request) *http.Request {
		ctx := req.Context()
		// Retrieve the Span from context.
		if sg := opentracing.SpanFromContext(ctx); sg != nil {

			// We are going to use this span in a client request, so mark as such.
			// sg.SetKind(tracer.RPCClient) // TODO?
			// ext.SpanKindRPCClient.Set(span)

			// Add some standard OpenTracing tags, useful in an HTTP request.
			// ext.HTTPMethod.Set(span, req.Method)
			// sg.SetHTTPMethod(req.Method) // TODO?

			sg.SetTag(zipkincore.HTTP_HOST, req.Host)
			sg.SetTag(zipkincore.HTTP_PATH, req.URL.String())
			sg.SetTag(zipkincore.HTTP_METHOD, req.Method)

			// ext.HTTPUrl.Set(
			// 	span,
			// 	fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path),
			// )
			// sg.SetHTTPUrl(fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)) // TODO?

			// Add information on the peer service we're about to contact.
			// if host, portString, err := net.SplitHostPort(req.URL.Host); err == nil {
			// ext.PeerHostname.Set(span, host)
			// sg.SetPeerHostname(host) // TODO?
			// if port, err := strconv.Atoi(portString); err != nil {
			// ext.PeerPort.Set(span, uint16(port))
			// sg.SetPeerPort(uint16(port)) // TODO?
			// }
			// } else {
			// ext.PeerHostname.Set(span, req.URL.Host)
			// sg.SetPeerHostname(req.URL.host) // TODO?
			// }

			// Inject the Span context into the outgoing HTTP Request.
			// if err := tracer.Inject(
			// 	sg.Context(),
			// 	opentracing.TextMap,
			// 	opentracing.HTTPHeadersCarrier(req.Header),
			if err := tr.Inject(sg.Context(), opentracing.TextMap, req); err != nil {
				fmt.Printf("error encountered while trying to inject span: %+v", err)
			}
		}

		return req
	}
}

// HandlerFunc is a middleware function for incoming HTTP requests.
type HandlerFunc func(next http.Handler) http.Handler

// FromHTTPRequest returns a Middleware HandlerFunc that tries to join with an
// rai-project trace found in the HTTP request headers and starts a new Segment
// called `operationName`. If no trace could be found in the HTTP request
// headers, the Segment will be a trace root. The Segment is incorporated in the
// HTTP Context object and can be retrieved with
// tracer.SegmentFromContext(ctx).
func FromHTTPRequest(tracer tracer.Tracer, operationName string) HandlerFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Try to join to a trace propagated in `req`.
			fmt.Println(tracer.Name(), tracer.Endpoint())

			wireContext, err := tracer.Extract(opentracing.HTTPHeaders, req)
			if err != nil {
				log.WithError(err).Error("error while trying to extract span: %+v\n", err)
				return
			}

			// create segment
			sg := opentracing.StartSpan(operationName, opentracing.ChildOf(wireContext))
			if sg == nil {
				log.WithError(err).Error("Unable to start segment.")
				return
			}
			sg.SetTag("serverSide", "here")
			defer sg.Finish()

			// store span in context
			ctx := opentracing.ContextWithSpan(req.Context(), sg)

			// update request context to include our new span
			req = req.WithContext(ctx)

			// next middleware or actual request handler
			next.ServeHTTP(w, req)
		})
	}
}
