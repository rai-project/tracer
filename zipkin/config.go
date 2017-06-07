package zipkin

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type zipkinConfig struct {
	Endpoints []string      `json:"endpoints" config:"zipkin.endpoints" default:"http://ec2-34-225-193-157.compute-1.amazonaws.com:9411/api/v1/spans"`
	done      chan struct{} `json:"-" config:"-"`
}

var (
	Config = &zipkinConfig{
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
