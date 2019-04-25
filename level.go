package tracer

type Level int32

const (
	NO_TRACE Level = iota
	APPLICATION_TRACE
	MODEL_TRACE
	FRAMEWORK_TRACE
	ML_LIBRARY_TRACE
	SYSTEM_LIBRARY_TRACE
	HARDWARE_TRACE
	FULL_TRACE

	LIBRARY_TRACE = ML_LIBRARY_TRACE
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

func LevelToName(l Level) string {
	switch l {
	case NO_TRACE:
		return "NO_TRACE"
	case APPLICATION_TRACE:
		return "APPLICATION_TRACE"
	case MODEL_TRACE:
		return "MODEL_TRACE"
	case FRAMEWORK_TRACE:
		return "FRAMEWORK_TRACE"
	case ML_LIBRARY_TRACE:
		return "ML_LIBRARY_TRACE"
	case SYSTEM_LIBRARY_TRACE:
		return "SYSTEM_LIBRARY_TRACE"
	case HARDWARE_TRACE:
		return "HARDWARE_TRACE"
	case FULL_TRACE:
		return "FULL_TRACE"
	default:
		panic("unknow trace level")
	}
}
