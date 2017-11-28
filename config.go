package tracer

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type tracerConfig struct {
	Enabled     bool          `json:"enabled" config:"tracer.enabled" default:"true"`
	Provider    string        `json:"provider" config:"tracer.provider" default:"zipkin"`
	LevelString string        `json:"level" config:"tracer.level"`
	Level       Level         `json:"-" config:"-"`
	done        chan struct{} `json:"-" config:"-"`
}

var (
	// Config holds the data read by rai-project/config
	Config = &tracerConfig{
		done:  make(chan struct{}),
		Level: NO_TRACE,
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
	c.Level = LevelFromName(c.LevelString)
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
