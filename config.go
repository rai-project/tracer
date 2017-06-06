package tracer

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type tracerConfig struct {
	Enabled bool          `json:"enabled" config:"tracer.enabled"`
	done    chan struct{} `json:"-" config:"-"`
}

var (
	Config = &tracerConfig{
		done: make(chan struct{}),
	}
)

func (tracerConfig) ConfigName() string {
	return "Tracer"
}

func (a *tracerConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *tracerConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
}

func (c tracerConfig) Wait() {
	<-c.done
}

func (c tracerConfig) String() string {
	return pp.Sprintln(c)
}

func (c tracerConfig) Debug() {
	log.Debug("Tracer Config = ", c)
}

func init() {
	config.Register(Config)
}
