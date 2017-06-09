package hooks

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/sebest/logrusly"
	"github.com/spf13/viper"
)

func init() {
	config.OnInit(func() {
		logger.Config.Wait()

		if !logger.UsingHook("loggly") {
			return
		}

		logger.Config.Wait()

		token := viper.GetString("loggly.token")
		host := viper.GetString("loggly.host")
		tags := viper.GetStringSlice("loggly.tags")

		level, err := logrus.ParseLevel(logger.Config.Level)
		if err != nil {
			fmt.Println("cannot parse level from logger config ", err)
			return
		}

		h := logrusly.NewLogglyHook(token, host, level, tags...)

		logger.RegisterHook("loggly", h)
	})
}
