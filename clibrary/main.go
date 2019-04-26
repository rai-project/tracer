package main

import "C"

import (
	"path/filepath"
	"sync"

	"github.com/rai-project/tracer"

	"github.com/k0kubun/pp"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	IsDebug        bool
	IsVerbose      bool
	AppSecret      string
	CfgFile        string
	tracerInitOnce sync.Once
	log            *logrus.Entry = logrus.New().WithField("pkg", "tracer/clibrary")
)

//export TracerSetLevel
func TracerSetLevel(lvl int32) {
	tracer.SetLevel(tracer.Level(lvl))
}

//export TracerClose
func TracerClose() {
	libDeinit()
	tracer.Close()
}

//export TracerInit
func TracerInit() {
	tracerInitOnce.Do(doTracerInit)
}

func doTracerInit() {
	log.Level = logrus.DebugLevel
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "tracer/clibrary")
	})

	color.NoColor = false
	opts := []config.Option{
		config.AppName("carml"),
		config.ColorMode(true),
		config.DebugMode(IsDebug),
		config.VerboseMode(IsVerbose),
	}

	pp.WithLineInfo = true

	if c, err := homedir.Expand(CfgFile); err == nil {
		CfgFile = c
	}
	if c, err := filepath.Abs(CfgFile); err == nil {
		CfgFile = c
	}
	opts = append(opts, config.ConfigFileAbsolutePath(CfgFile))

	if AppSecret != "" {
		opts = append(opts, config.AppSecret(AppSecret))
	}

	config.Init(opts...)

	tracer.SetLevel(tracer.FULL_TRACE)
	libInit()
}

func initLib() {
	TracerInit()
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}
