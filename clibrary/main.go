package main

import "C"

import (
	"path/filepath"

	"github.com/k0kubun/pp"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/sirupsen/logrus"
)

var (
	IsDebug   bool
	IsVerbose bool
	AppSecret string
	CfgFile   string
	log       *logrus.Entry = logrus.New().WithField("pkg", "tracer/clibrary")
)

//export TracerInit
func TracerInit() {

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
	if IsDebug || IsVerbose {
		pp.WithLineInfo = true
	}
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
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}