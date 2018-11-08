package tracer

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

func LevelFromName(s string) Level {
	if s == "" {
		return Config.Level
	}
	lvl, err := LevelString(s)
	if err != nil {
		log.WithError(err).Errorf("invalid level spec %v", s)
		return NO_TRACE
	}
	return lvl
}
