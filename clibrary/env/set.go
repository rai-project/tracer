package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/uber/jaeger-client-go"
)

func Set(sc jaeger.SpanContext) {
	if os.Getenv(TraceIdEnv) != "" {
		return
	}

	traceIDVal := sc.TraceID().String()

	fmt.Printf("Set env [%s] = %s....\n", TraceIdEnv, traceIDVal)

	os.Setenv(TraceIdEnv, traceIDVal)

	if sc.ParentID() != 0 {
		os.Setenv(ParentSpanIdEnv, strconv.FormatUint(uint64(sc.ParentID()), 16))
	}

	os.Setenv(SpanIdEnv, strconv.FormatUint(uint64(sc.SpanID()), 16))

}
