package tracer

type Level int32

const (
	NO_TRACE Level = iota
	APPLICATION_TRACE
	MODEL_TRACE
	FRAMEWORK_TRACE
	LIBRARY_TRACE
	HARDWARE_TRACE
	FULL_TRACE
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
