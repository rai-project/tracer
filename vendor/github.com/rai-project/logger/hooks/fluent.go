package hooks

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/evalphobia/logrus_fluent"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/spf13/viper"
)

func init() {
	config.OnInit(func() {
		logger.Config.Wait()

		if !logger.UsingHook("fluent") {
			return
		}

		logger.Config.Wait()

		host := viper.GetString("fluent.host")
		port := viper.GetInt("fluent.port")
		tags := viper.GetStringSlice("fluent.tags")

		h, err := logrus_fluent.New(host, port)
		if err != nil {
			fmt.Println("failed to load fluent logger hook ", err)
      return
		}

		for _, tag := range tags {
			h.SetTag(tag)
		}

		h.SetLevels([]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		})

		logger.RegisterHook("fluent", h)
	})
}
