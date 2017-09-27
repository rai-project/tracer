package tracer

import (
	"context"
	"runtime"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer/defaults"
)

var (
	stdTracer Tracer
)

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

func MustNew(serviceName string) Tracer {
	backendName := Config.Provider
	if backendName == "" || !Config.Enabled {
		backendName = "noop"
	}
	tr, err := NewFromName(serviceName, backendName)
	if err != nil {
		// just use the noop tracer
		tr, err = NewFromName(serviceName, "noop")
		if err != nil {
			panic(err)
		}
	}
	return tr
}

func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	if stdTracer == nil {
		return nil
	}
	if runtime.GOOS == "linux" {
		opts = append(opts, opentracing.Tag{"perfevents", defaults.PerfEvents})
	}
	return stdTracer.StartSpan(operationName, opts...)
}

func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if stdTracer == nil {
		return nil, nil
	}
	if runtime.GOOS == "linux" {
		opts = append(opts, opentracing.Tag{"perfevents", defaults.PerfEvents})
	}
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

func Enabled() bool {
	if stdTracer == nil {
		return false
	}
	return Config.Enabled
}

func Close() error {
	openTracers.Range(func(_ interface{}, value interface{}) bool {
		tr, ok := value.(Tracer)
		if !ok {
			return true
		}
		tr.Close()
		return true
	})
	return nil
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
