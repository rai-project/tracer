package observer

import (
	"github.com/k0kubun/pp"
	"github.com/opentracing-contrib/go-observer"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type observerConfig struct {
	ObserverNames []string              `json:"observers" config:"tracer.observers"`
	Observers     []otobserver.Observer `json:"-" config:"-"`
	done          chan struct{}         `json:"-" config:"-"`
}

var (
	Config = &observerConfig{
		done: make(chan struct{}),
	}
)

func (observerConfig) ConfigName() string {
	return "tracer/observer"
}

func (a *observerConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *observerConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	if len(a.ObserverNames) == 0 {
		a.ObserverNames = Default
	}
	for _, observer := range a.ObserverNames {
		switch observer {
		case "perf", "perf_events", "perfevents":
			a.Observers = append(a.Observers, PerfEvents)
			continue
		case "instruments":
			a.Observers = append(a.Observers, Instruments)
			continue
		case "meminfo":
			a.Observers = append(a.Observers, MemInfo)
		case "gpu_meminfo":
			a.Observers = append(a.Observers, GPUMemInfo)
		case "dummy":
			a.Observers = append(a.Observers, Dummy)
		}
	}
}

func (c observerConfig) Wait() {
	<-c.done
}

func (c observerConfig) String() string {
	return pp.Sprintln(c)
}

func (c observerConfig) Debug() {
	log.Debug("observer Config = ", c)
}

func init() {
	config.Register(Config)
}
