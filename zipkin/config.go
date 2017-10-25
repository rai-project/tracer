package zipkin

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/utils"
	"github.com/rai-project/vipertags"
)

type zipkinConfig struct {
	Endpoints   []string      `json:"endpoints" config:"tracer.endpoints" env:"TRACER_ENDPOINTS"`
	LevelString string        `json:"level" config:"tracer.level"`
	Level       tracer.Level  `json:"-" config:"-"`
	done        chan struct{} `json:"-" config:"-"`
}

var (
	fixEndpoints = utils.FixEndpoints("http://", "9411", "/api/v1/spans")
	Config       = &zipkinConfig{
		done: make(chan struct{}),
	}
)

func (zipkinConfig) ConfigName() string {
	return "Zipkin"
}

func (a *zipkinConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *zipkinConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	a.Endpoints = fixEndpoints(a.Endpoints)
	a.Level = tracer.LevelFromName(a.LevelString)
}

func (c zipkinConfig) Wait() {
	<-c.done
}

func (c zipkinConfig) String() string {
	return pp.Sprintln(c)
}

func (c zipkinConfig) Debug() {
	log.Debug("Zipkin Config = ", c)
}

func init() {
	config.Register(Config)
}
