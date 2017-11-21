package tracer

import "github.com/gogo/protobuf/proto"

type Level int32

const (
	NO_TRACE        Level = 0
	STEP_TRACE      Level = 1
	FRAMEWORK_TRACE Level = 2
	CPU_ONLY_TRACE  Level = 3
	HARDWARE_TRACE  Level = 4
	FULL_TRACE      Level = 5
)

var Level_name = map[int32]string{
	0: "NO_TRACE",
	1: "STEP_TRACE",
	2: "FRAMEWORK_TRACE",
	3: "CPU_ONLY_TRACE",
	4: "HARDWARE_TRACE",
	5: "FULL_TRACE",
}
var Level_value = map[string]int32{
	"NO_TRACE":        0,
	"STEP_TRACE":      1,
	"FRAMEWORK_TRACE": 2,
	"CPU_ONLY_TRACE":  3,
	"HARDWARE_TRACE":  4,
	"FULL_TRACE":      5,
}
var Level_get = map[string]Level{
	"NO_TRACE":        NO_TRACE,
	"STEP_TRACE":      STEP_TRACE,
	"FRAMEWORK_TRACE": FRAMEWORK_TRACE,
	"CPU_ONLY_TRACE":  CPU_ONLY_TRACE,
	"HARDWARE_TRACE":  HARDWARE_TRACE,
	"FULL_TRACE":      FULL_TRACE,
}

func LevelFromName(s string) Level {
	if s == "" {
		return NO_TRACE
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
