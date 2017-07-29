package jaeger

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type jaegerConfig struct {
	Endpoints []string      `json:"endpoints" config:"tracer.endpoints"`
	done      chan struct{} `json:"-" config:"-"`
}

var (
	Config = &jaegerConfig{
		done: make(chan struct{}),
	}
)

func (jaegerConfig) ConfigName() string {
	return "Jaeger"
}

func (a *jaegerConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *jaegerConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
}

func (c jaegerConfig) Wait() {
	<-c.done
}

func (c jaegerConfig) String() string {
	return pp.Sprintln(c)
}

func (c jaegerConfig) Debug() {
	log.Debug("Jaeger Config = ", c)
}

func init() {
	config.Register(Config)
}
