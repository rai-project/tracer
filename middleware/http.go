package middleware

import (
	"github.com/labstack/echo"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/rai-project/tracer"
)

func ToHTTPRequest(tr tracer.Tracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			sg := opentracing.SpanFromContext(ctx)
			if sg == nil {
				return next(c)
			}
			sg.SetTag(zipkincore.HTTP_HOST, req.Host)
			sg.SetTag(zipkincore.HTTP_PATH, req.URL.String())
			sg.SetTag(zipkincore.HTTP_METHOD, req.Method)
			carrier := opentracing.HTTPHeadersCarrier(req.Header)
			if err := tr.Inject(sg.Context(), opentracing.HTTPHeaders, carrier); err != nil {
				log.Errorf("error encountered while trying to inject span: %+v", err)
			}
			return next(c)
		}
	}
}

func ToHTTPResponse(tracer tracer.Tracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
			req := c.Request()
			ctx := req.Context()
			sg := opentracing.SpanFromContext(ctx)
			if sg == nil {
				return next(c)
			}
			tracer.Inject(
				sg.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Response().Header()))
			return next(c)
		}
	}
}

func FromHTTPRequest(tracer tracer.Tracer, operationName string) echo.MiddlewareFunc {
	// Try to join to a trace propagated in `req`.
	log.WithField("tracer_name", tracer.Name()).
		WithField("tracer_operation", operationName).
		Infof("added from http request tracer")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
			req := c.Request()
			ctx := req.Context()

			startSpanOpts := []opentracing.StartSpanOption{}

			carrier := opentracing.HTTPHeadersCarrier(req.Header)
			wireContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

			if err == nil && wireContext != nil {
				startSpanOpts = append(startSpanOpts, opentracing.ChildOf(wireContext))
			}

			// create segment
			sg := opentracing.StartSpan(operationName, startSpanOpts...)
			if sg == nil {
				log.WithError(err).Error("Unable to start segment.")
				return next(c)
			}
			if requestID := c.Response().Header().Get(echo.HeaderXRequestID); requestID != "" {
				sg.SetTag("request_id", requestID)
			}
			// record HTTP method
			ext.HTTPMethod.Set(sg, req.Method)
			// record HTTP url
			ext.HTTPUrl.Set(sg, req.URL.String())
			// record HTTP status code
			defer ext.HTTPStatusCode.Set(sg, uint16(c.Response().Status))

			defer sg.Finish()

			// store span in context
			ctx = opentracing.ContextWithSpan(req.Context(), sg)

			// update request context to include our new span
			req = req.WithContext(ctx)
			c.SetRequest(req)

			// next middleware or actual request handler
			return next(c)
		}
	}
}
