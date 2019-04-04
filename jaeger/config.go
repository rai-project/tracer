package jaeger

import (
	"strings"

	"github.com/rai-project/tracer"

	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer/utils"
	"github.com/rai-project/vipertags"
)

type jaegerConfig struct {
	Endpoints   []string      `json:"endpoints" config:"tracer.endpoints"`
	LevelString string        `json:"level" config:"tracer.level"`
	Level       tracer.Level  `json:"-" config:"-"`
	done        chan struct{} `json:"-" config:"-"`
}

var (
	fixEndpoints = utils.FixEndpoints("http://", "9411", "/api/v1/spans")
	Config       = &jaegerConfig{
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
	if len(a.Endpoints) == 0 {
		return
	}
	if strings.HasPrefix(a.Endpoints[0], "udp://") {
		return
	}
	a.Endpoints = fixEndpoints(a.Endpoints)
	a.Level = tracer.LevelFromName(a.LevelString)
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
