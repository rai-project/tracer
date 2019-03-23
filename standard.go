package tracer

import (
	"context"
	"runtime"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer/defaults"
	"github.com/rai-project/tracer/observer"
	"golang.org/x/sync/syncmap"
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
	mut.Lock()
	defer mut.Unlock()
	openTracers.Range(func(_ interface{}, value interface{}) bool {
		tr, ok := value.(Tracer)
		if !ok {
			return true
		}
		err := tr.Close()
		if err != nil {
			log.WithError(err).WithField("tracer", tr.Name()).Error("Failed to close tracer")
		}
		return true
	})
	openTracers = syncmap.Map{}
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
		std, err := New(config.App.Name)
		if err != nil {
			SetStd(noop)
			return
		}
		SetStd(std)

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
