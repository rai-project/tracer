// +build ignore

package trace

import (
	"net/url"
	"time"

	"github.com/apex/log"
	chrm "github.com/sensepost/gowitness/chrome"
	"github.com/sensepost/gowitness/storage"
	"github.com/sensepost/gowitness/utils"
)

var (
	chrome       chrm.Chrome
	waitTimeout  int
	screenshotDB storage.Storage
)

func (t Trace) ScreenshotURL() string {
	panic("todo")
	return ""
}

func (t Trace) Screenshot() string {
	screenshotURL := t.ScreenshotURL()
	startTime := time.Now()

	u, err := url.ParseRequestURI(screenshotURL)
	if err != nil {
		log.WithField("url", screenshotURL).Fatal("Invalid URL specified")
	}

	// Process this URL
	utils.ProcessURL(u, &chrome, &screenshotDB, waitTimeout)

	log.WithField("run-time", time.Since(startTime)).Info("Complete")

	return ""
}

func InitChrome() {

	// Init Google Chrome
	chrome = chrm.Chrome{
		Resolution:    "1440,900",
		ChromeTimeout: 90,
		Path:          "",
	}
	chrome.Setup()

	screenshotDB = storage.Storage{}
	screenshotDB.Open("micro.db")

	waitTimeout = 3
}
