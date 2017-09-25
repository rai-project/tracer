package tracer

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
)

var stdTracer Tracer

func SetStd(t Tracer) {
	stdTracer = t
	opentracing.SetGlobalTracer(t)
}

func Std() Tracer {
	return stdTracer
}

func New(serviceName string) (Tracer, error) {
	backendName := Config.Provider
	if backendName == "" || !Config.Enabled {
		backendName = "noop"
	}
	return NewFromName(serviceName, backendName)
}

func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	if stdTracer == nil {
		return nil
	}
	return stdTracer.StartSpan(operationName, opts...)
}

func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if stdTracer == nil {
		return nil, nil
	}
	if _, ok := opentracing.GlobalTracer().(opentracing.NoopTracer); ok {
		log.Error("tracer is using a No-op tracer")
		return nil, nil
	}
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

func Enabled() bool {
	if stdTracer == nil {
		return false
	}
	return Config.Enabled
}

func Endpoints() []string {
	if stdTracer == nil {
		return []string{}
	}
	return stdTracer.Endpoints()
}

func Provider() string {
	if stdTracer == nil {
		return Config.Provider
	}
	return stdTracer.Name()
}

func Close() {
	if stdTracer == nil {
		return
	}
	stdTracer.Close()
}

func init() {
	config.AfterInit(func() {
		std, err := New(config.App.Name)
		if err != nil {
			// just use the noop tracer
			std, err := NewFromName(config.App.Name, "noop")
			if err != nil {
				return
			}
			SetStd(std)
			return
		}
		SetStd(std)
	})
}
