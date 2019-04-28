package tracer

import (
	"context"
	"runtime"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer/defaults"
	"github.com/rai-project/tracer/observer"
)

var (
	stdTracer Tracer
	mut       sync.Mutex
	noop      Tracer
	usingPerf bool
)

func SetStd(t Tracer) {
	stdTracer = t
	opentracing.SetGlobalTracer(t)
}

func Std() Tracer {
	return stdTracer
}

func ResetStd(options ...Option) Tracer {
	if stdTracer != nil {
		stdTracer.Close()
	}
	std, err := New(config.App.Name, options...)
	if err != nil {
		SetStd(noop)
		return nil
	}
	SetStd(std)
	return std
}

func New(serviceName string, options ...Option) (Tracer, error) {
	backendName := Config.Provider
	if backendName == "" || !Config.Enabled {
		backendName = "noop"
	}
	return NewFromName(serviceName, backendName, options...)
}

func MustNew(serviceName string, options ...Option) Tracer {
	backendName := Config.Provider
	if backendName == "" || !Config.Enabled {
		backendName = "noop"
	}
	tr, err := NewFromName(serviceName, backendName, options...)
	if err != nil {
		// just use the noop tracer
		tr, err = NewFromName(serviceName, "noop")
		if err != nil {
			panic(err)
		}
	}
	return tr
}

func StartSpan(lvl Level, operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	if stdTracer == nil {
		return nil
	}
	if lvl > stdTracer.Level() {
		return noop.StartSpan(operationName, opts...)
	}
	opts = append(opts, opentracing.Tag{"trace_level", lvl.String()})
	if usingPerf {
		opts = append(opts, opentracing.Tag{"perfevents", defaults.PerfEvents})
	}
	return stdTracer.StartSpan(operationName, opts...)
}

func StartSpanFromContext(ctx context.Context, lvl Level, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if stdTracer == nil {
		return nil, ctx
	}
	if lvl > stdTracer.Level() {
		return noop.StartSpanFromContext(ctx, operationName, opts...)
	}
	opts = append(opts, opentracing.Tag{"trace_level", lvl.String()})
	if usingPerf {
		opts = append(opts, opentracing.Tag{"perfevents", defaults.PerfEvents})
	}
	return stdTracer.StartSpanFromContext(ctx, operationName, opts...)
}

func Enabled() bool {
	if stdTracer == nil {
		return false
	}
	return Config.Enabled
}

func Close() error {
	if stdTracer != nil {
		err := stdTracer.Close()
		stdTracer = nil
		return err
	}
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

func GetLevel() Level {
	if stdTracer == nil {
		return NO_TRACE
	}
	return stdTracer.Level()
}

func SetLevel(lvl Level) {
	if stdTracer == nil {
		return
	}
	stdTracer.SetLevel(lvl)
}

func init() {
	loadNoop := func(name string) {
		if name == "" {
			name = "tracer"
		}
		no, err := NewFromName("tracer", "noop")
		if err != nil {
			return
		}
		noop = no
	}
	config.AfterInit(func() {
		loadNoop(config.App.Name)
		ResetStd()

		if runtime.GOOS == "linux" {
			for _, o := range observer.Config.ObserverNames {
				if o == "perf" || o == "perf_events" || o == "perfevents" {
					usingPerf = true
					break
				}
			}
		}
	})
	loadNoop("")
}
