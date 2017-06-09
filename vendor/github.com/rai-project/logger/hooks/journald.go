// +build linux

package hooks

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/wercker/journalhook"
)

func init() {
	config.OnInit(func() {
		logger.Config.Wait()

		if !logger.UsingHook("journald") {
			return
		}

		h := &journalhook.JournalHook{}
		logger.RegisterHook("journald", h)
	})
}
