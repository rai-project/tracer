package main

import (
	"os"

	"github.com/rai-project/tracer/clibrary/env"
)

func GetTraceIdEnv() string {
	return os.Getenv(env.TraceIdEnv)
}
