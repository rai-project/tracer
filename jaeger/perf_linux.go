package jaeger

import (
	perfevents "github.com/opentracing-contrib/perfevents/go"
	jaeger "github.com/uber/jaeger-client-go"
)

var (
	perfeventsObserver jaeger.Observer = perfevents.NewObserver()
)
