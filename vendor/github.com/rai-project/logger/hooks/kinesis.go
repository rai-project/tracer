package hooks

import (
	"fmt"

	"github.com/sirupsen/logrus"
	aaws "github.com/aws/aws-sdk-go/aws"
	"github.com/evalphobia/logrus_kinesis"
	"github.com/rai-project/aws"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

func init() {
	config.OnInit(func() {
		logger.Config.Wait()

		if !logger.UsingHook("kinesis") {
			return
		}
		println("setting up kinesis logger hook @", config.App.Name)

		aws.Config.Wait()

		c, err := aws.NewConfig()
		if err != nil {
			println("failed to load kinesis logger hook ", err)
			return
		}
		cred, err := c.Credentials.Get()
		if err != nil {
			println("failed to load kinesis logger hook ", err)
			return
		}
		h, err := logrus_kinesis.New(config.App.Name, logrus_kinesis.Config{
			AccessKey: cred.AccessKeyID,
			SecretKey: cred.SecretAccessKey,
			Region:    aaws.StringValue(c.Region),
			Endpoint:  aaws.StringValue(c.Endpoint),
		})
		if err != nil {
			fmt.Println("failed to load kinesis logger hook ", err)
			return
		}
		h.SetLevels([]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		})
		logger.RegisterHook("kinesis", h)
	})
}
