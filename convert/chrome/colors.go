package chrome

// Mapping from more reasonable color names to the reserved color names in
// https://github.com/catapult-project/catapult/blob/master/tracing/tracing/base/color_scheme.html#L50
// The chrome trace viewer allows only those as cname values.
const (
	colorLightMauve     = "thread_state_uninterruptible" // 182, 125, 143
	colorOrange         = "thread_state_iowait"          // 255, 140, 0
	colorSeafoamGreen   = "thread_state_running"         // 126, 200, 148
	colorVistaBlue      = "thread_state_runnable"        // 133, 160, 210
	colorTan            = "thread_state_unknown"         // 199, 155, 125
	colorIrisBlue       = "background_memory_dump"       // 0, 180, 180
	colorMidnightBlue   = "light_memory_dump"            // 0, 0, 180
	colorDeepMagenta    = "detailed_memory_dump"         // 180, 0, 180
	colorBlue           = "vsync_highlight_color"        // 0, 0, 255
	colorGrey           = "generic_work"                 // 125, 125, 125
	colorGreen          = "good"                         // 0, 125, 0
	colorDarkGoldenrod  = "bad"                          // 180, 125, 0
	colorPeach          = "terrible"                     // 180, 0, 0
	colorBlack          = "black"                        // 0, 0, 0
	colorLightGrey      = "grey"                         // 221, 221, 221
	colorWhite          = "white"                        // 255, 255, 255
	colorYellow         = "yellow"                       // 255, 255, 0
	colorOlive          = "olive"                        // 100, 100, 0
	colorCornflowerBlue = "rail_response"                // 67, 135, 253
	colorSunsetOrange   = "rail_animation"               // 244, 74, 63
	colorTangerine      = "rail_idle"                    // 238, 142, 0
	colorShamrockGreen  = "rail_load"                    // 13, 168, 97
	colorGreenishYellow = "startup"                      // 230, 230, 0
	colorDarkGrey       = "heap_dump_stack_frame"        // 128, 128, 128
	colorTawny          = "heap_dump_child_node_arrow"   // 204, 102, 0
	colorLemon          = "cq_build_running"             // 255, 255, 119
	colorLime           = "cq_build_passed"              // 153, 238, 102
	colorPink           = "cq_build_failed"              // 238, 136, 136
	colorSilver         = "cq_build_abandoned"           // 187, 187, 187
	colorManzGreen      = "cq_build_attempt_runnig"      // 222, 222, 75
	colorKellyGreen     = "cq_build_attempt_passed"      // 108, 218, 35
	colorAnotherGrey    = "cq_build_attempt_failed"      // 187, 187, 187
)

var colorForTask = []string{
	colorLightMauve,
	colorOrange,
	colorSeafoamGreen,
	colorVistaBlue,
	colorTan,
	colorMidnightBlue,
	colorIrisBlue,
	colorDeepMagenta,
	colorGreen,
	colorDarkGoldenrod,
	colorPeach,
	colorOlive,
	colorCornflowerBlue,
	colorSunsetOrange,
	colorTangerine,
	colorShamrockGreen,
	colorTawny,
	colorLemon,
	colorLime,
	colorPink,
	colorSilver,
	colorManzGreen,
	colorKellyGreen,
}

func pickTaskColor(id uint64) string {
	idx := id % uint64(len(colorForTask))
	return colorForTask[idx]
}

func colorName(cat string) string {
	id := hash64(cat)
	return pickTaskColor(id)
}
