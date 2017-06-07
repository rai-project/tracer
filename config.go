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
	// Config holds the data read by rai-project/config
	Config = &tracerConfig{
		done: make(chan struct{}),
	}
)

func (tracerConfig) ConfigName() string {
	return "Tracer"
}

func (c *tracerConfig) SetDefaults() {
	vipertags.SetDefaults(c)
}

func (c *tracerConfig) Read() {
	defer close(c.done)
	vipertags.Fill(c)
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
