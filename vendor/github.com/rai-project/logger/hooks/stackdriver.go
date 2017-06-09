package hooks

import (
	"github.com/k0kubun/pp"
	"github.com/knq/sdhook"
	"github.com/rai-project/config"
	"github.com/rai-project/googlecloud"
	"github.com/rai-project/logger"
)

func init() {
	config.OnInit(func() {
		logger.Config.Wait()

		if !logger.UsingHook("stackdriver") {
			return
		}

		googlecloud.Config.Wait()

		opts := googlecloud.NewOptions()

		h, err := sdhook.New(
			sdhook.GoogleLoggingAgent(),
			sdhook.ProjectID(opts.ProjectID),
			sdhook.GoogleServiceAccountCredentialsJSON(opts.Bytes()),
			sdhook.ErrorReportingService(config.App.Name),
			sdhook.ErrorReportingLogName("error_log"),
		)
		if err != nil {
			pp.Println(err)
			return
		}
		logger.RegisterHook("stackdriver", h)
	})
}
