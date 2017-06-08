package zipkin

import (
	"github.com/Sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

type loggerWrapper struct {
	*logrus.Entry
}

var (
	log *loggerWrapper
)

func (l *loggerWrapper) Error(s string) {
	l.Entry.Error(s)
}

func init() {
	config.AfterInit(func() {
		log = &loggerWrapper{logger.New().WithField("pkg", "tracer/zipkin")}
	})
}
