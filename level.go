package tracer

import (
	"github.com/gogo/protobuf/proto"
)

type Level int32

const (
	NO_TRACE          Level = 0
	APPLICATION_TRACE Level = 1
	MODEL_TRACE       Level = 2
	FRAMEWORK_TRACE   Level = 3
	LIBRARY_TRACE     Level = 4
	HARDWARE_TRACE    Level = 5
	FULL_TRACE        Level = 6
)

var Level_name = map[int32]string{
	0: "NO_TRACE",
	1: "APPLICATION_TRACE",
	2: "MODEL_TRACE",
	3: "FRAMEWORK_TRACE",
	4: "LIBRARY_TRACE",
	5: "HARDWARE_TRACE",
	6: "FULL_TRACE",
}
var Level_value = map[string]int32{
	"NO_TRACE":          0,
	"APPLICATION_TRACE": 1,
	"MODEL_TRACE":       2,
	"FRAMEWORK_TRACE":   3,
	"LIBRARY_TRACE":     4,
	"HARDWARE_TRACE":    5,
	"FULL_TRACE":        6,
}
var Level_get = map[string]Level{
	"NO_TRACE":          NO_TRACE,
	"APPLICATION_TRACE": APPLICATION_TRACE,
	"MODEL_TRACE":       MODEL_TRACE,
	"FRAMEWORK_TRACE":   FRAMEWORK_TRACE,
	"LIBRARY_TRACE":     LIBRARY_TRACE,
	"HARDWARE_TRACE":    HARDWARE_TRACE,
	"FULL_TRACE":        FULL_TRACE,
}

func LevelFromName(s string) Level {
	if s == "" {
		return Config.Level
	}
	lvl, ok := Level_get[s]
	if !ok {
		log.Errorf("invalid level spec %v", s)
		return NO_TRACE
	}
	return lvl
}

func (x Level) String() string {
	return proto.EnumName(Level_name, int32(x))
}
