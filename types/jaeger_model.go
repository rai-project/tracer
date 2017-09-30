package types

import (
	"github.com/uber/jaeger/model"
)

type JaegerSpan struct {
	model.Span
}

type JaegerTrace struct {
	model.Trace
}
