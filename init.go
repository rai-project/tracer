package tracer

import (
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log *logrus.Entry = logger.New().WithField("pkg", "tracer")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "tracer")
	})
}
