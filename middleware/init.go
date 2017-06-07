// +build !develop

package middleware

import (
	"github.com/Sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log *logrus.Entry
)

func init() {
	config.BeforeInit(func() {
		log = logger.New().WithField("pkg", "tracer/zipkin")
	})
}
