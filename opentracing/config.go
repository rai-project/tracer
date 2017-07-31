package opentracing

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer/utils"
	"github.com/rai-project/vipertags"
)

type opentracingConfig struct {
	Endpoints []string      `json:"endpoints" config:"tracer.endpoints"`
	done      chan struct{} `json:"-" config:"-"`
}

var (
	fixEndpoints = utils.FixEndpoints("http://", "9411", "/api/v1/spans")
	Config       = &opentracingConfig{
		done: make(chan struct{}),
	}
)

func (opentracingConfig) ConfigName() string {
	return "OpenTracing"
}

func (a *opentracingConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *opentracingConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	a.Endpoints = fixEndpoints(a.Endpoints)
}

func (c opentracingConfig) Wait() {
	<-c.done
}

func (c opentracingConfig) String() string {
	return pp.Sprintln(c)
}

func (c opentracingConfig) Debug() {
	log.Debug("OpenTracing Config = ", c)
}

func init() {
	config.Register(Config)
}
