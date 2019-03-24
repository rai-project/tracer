package convert

type ViewerEvent struct {
	Name     string      `json:"name,omitempty"`
	Phase    string      `json:"ph"`
	Scope    string      `json:"s,omitempty"`
	Time     float64     `json:"ts"`
	Dur      float64     `json:"dur,omitempty"`
	Pid      uint64      `json:"pid"`
	Tid      uint64      `json:"tid"`
	ID       uint64      `json:"id,omitempty"`
	Stack    int         `json:"sf,omitempty"`
	EndStack int         `json:"esf,omitempty"`
	Arg      interface{} `json:"args,omitempty"`
	Cname    string      `json:"cname,omitempty"`
	Category string      `json:"cat,omitempty"`
}

type ViewerFrame struct {
	Name   string `json:"name"`
	Parent int    `json:"parent,omitempty"`
}

type NameArg struct {
	Name string `json:"name"`
}

type TaskArg struct {
	ID     uint64 `json:"id"`
	StartG uint64 `json:"start_g,omitempty"`
	EndG   uint64 `json:"end_g,omitempty"`
}

type RegionArg struct {
	TaskID uint64 `json:"taskid,omitempty"`
}

type SortIndexArg struct {
	Index int `json:"sort_index"`
}

// func (ctx *traceContext) makeSlice(ev *trace.Event, name string) *ViewerEvent {
// 	// If ViewerEvent.Dur is not a positive value,
// 	// trace viewer handles it as a non-terminating time interval.
// 	// Avoid it by setting the field with a small value.
// 	durationUsec := ctx.time(ev.Link) - ctx.time(ev)
// 	if ev.Link.Ts-ev.Ts <= 0 {
// 		durationUsec = 0.0001 // 0.1 nanoseconds
// 	}
// 	sl := &ViewerEvent{
// 		Name:     name,
// 		Phase:    "X",
// 		Time:     ctx.time(ev),
// 		Dur:      durationUsec,
// 		Tid:      ctx.proc(ev),
// 		Stack:    ctx.stack(ev.Stk),
// 		EndStack: ctx.stack(ev.Link.Stk),
// 	}

// 	// grey out non-overlapping events if the event is not a global event (ev.G == 0)
// 	if ctx.mode&modeTaskOriented != 0 && ev.G != 0 {
// 		// include P information.
// 		if t := ev.Type; t == trace.EvGoStart || t == trace.EvGoStartLabel {
// 			type Arg struct {
// 				P int
// 			}
// 			sl.Arg = &Arg{P: ev.P}
// 		}
// 		// grey out non-overlapping events.
// 		overlapping := false
// 		for _, task := range ctx.tasks {
// 			if _, overlapped := task.overlappingDuration(ev); overlapped {
// 				overlapping = true
// 				break
// 			}
// 		}
// 		if !overlapping {
// 			sl.Cname = colorLightGrey
// 		}
// 	}
// 	return sl
// }

// func (ctx *traceContext) emitTask(task *taskDesc, sortIndex int) {
// 	taskRow := uint64(task.id)
// 	taskName := task.name
// 	durationUsec := float64(task.lastTimestamp()-task.firstTimestamp()) / 1e3

// 	ctx.emitFooter(&ViewerEvent{Name: "thread_name", Phase: "M", Pid: tasksSection, Tid: taskRow, Arg: &NameArg{fmt.Sprintf("T%d %s", task.id, taskName)}})
// 	ctx.emit(&ViewerEvent{Name: "thread_sort_index", Phase: "M", Pid: tasksSection, Tid: taskRow, Arg: &SortIndexArg{sortIndex}})
// 	ts := float64(task.firstTimestamp()) / 1e3
// 	sl := &ViewerEvent{
// 		Name:  taskName,
// 		Phase: "X",
// 		Time:  ts,
// 		Dur:   durationUsec,
// 		Pid:   tasksSection,
// 		Tid:   taskRow,
// 		Cname: pickTaskColor(task.id),
// 	}
// 	targ := TaskArg{ID: task.id}
// 	if task.create != nil {
// 		sl.Stack = ctx.stack(task.create.Stk)
// 		targ.StartG = task.create.G
// 	}
// 	if task.end != nil {
// 		sl.EndStack = ctx.stack(task.end.Stk)
// 		targ.EndG = task.end.G
// 	}
// 	sl.Arg = targ
// 	ctx.emit(sl)

// 	if task.create != nil && task.create.Type == trace.EvUserTaskCreate && task.create.Args[1] != 0 {
// 		ctx.arrowSeq++
// 		ctx.emit(&ViewerEvent{Name: "newTask", Phase: "s", Tid: task.create.Args[1], ID: ctx.arrowSeq, Time: ts, Pid: tasksSection})
// 		ctx.emit(&ViewerEvent{Name: "newTask", Phase: "t", Tid: taskRow, ID: ctx.arrowSeq, Time: ts, Pid: tasksSection})
// 	}
// }

// func (ctx *traceContext) emitRegion(s regionDesc) {
// 	if s.Name == "" {
// 		return
// 	}

// 	if !tsWithinRange(s.firstTimestamp(), ctx.startTime, ctx.endTime) &&
// 		!tsWithinRange(s.lastTimestamp(), ctx.startTime, ctx.endTime) {
// 		return
// 	}

// 	ctx.regionID++
// 	regionID := ctx.regionID

// 	id := s.TaskID
// 	scopeID := fmt.Sprintf("%x", id)
// 	name := s.Name

// 	sl0 := &ViewerEvent{
// 		Category: "Region",
// 		Name:     name,
// 		Phase:    "b",
// 		Time:     float64(s.firstTimestamp()) / 1e3,
// 		Tid:      s.G, // only in goroutine-oriented view
// 		ID:       uint64(regionID),
// 		Scope:    scopeID,
// 		Cname:    pickTaskColor(s.TaskID),
// 	}
// 	if s.Start != nil {
// 		sl0.Stack = ctx.stack(s.Start.Stk)
// 	}
// 	ctx.emit(sl0)

// 	sl1 := &ViewerEvent{
// 		Category: "Region",
// 		Name:     name,
// 		Phase:    "e",
// 		Time:     float64(s.lastTimestamp()) / 1e3,
// 		Tid:      s.G,
// 		ID:       uint64(regionID),
// 		Scope:    scopeID,
// 		Cname:    pickTaskColor(s.TaskID),
// 		Arg:      RegionArg{TaskID: s.TaskID},
// 	}
// 	if s.End != nil {
// 		sl1.Stack = ctx.stack(s.End.Stk)
// 	}
// 	ctx.emit(sl1)
// }

// type heapCountersArg struct {
// 	Allocated uint64
// 	NextGC    uint64
// }

// func (ctx *traceContext) emitHeapCounters(ev *trace.Event) {
// 	if ctx.prevHeapStats == ctx.heapStats {
// 		return
// 	}
// 	diff := uint64(0)
// 	if ctx.heapStats.nextGC > ctx.heapStats.heapAlloc {
// 		diff = ctx.heapStats.nextGC - ctx.heapStats.heapAlloc
// 	}
// 	if tsWithinRange(ev.Ts, ctx.startTime, ctx.endTime) {
// 		ctx.emit(&ViewerEvent{Name: "Heap", Phase: "C", Time: ctx.time(ev), Pid: 1, Arg: &heapCountersArg{ctx.heapStats.heapAlloc, diff}})
// 	}
// 	ctx.prevHeapStats = ctx.heapStats
// }

// func (ctx *traceContext) emitThreadCounters(ev *trace.Event) {
// 	if ctx.prevThreadStats == ctx.threadStats {
// 		return
// 	}
// 	if tsWithinRange(ev.Ts, ctx.startTime, ctx.endTime) {
// 		ctx.emit(&ViewerEvent{Name: "Threads", Phase: "C", Time: ctx.time(ev), Pid: 1, Arg: &threadCountersArg{
// 			Running:   ctx.threadStats.prunning,
// 			InSyscall: ctx.threadStats.insyscall}})
// 	}
// 	ctx.prevThreadStats = ctx.threadStats
// }

// func (ctx *traceContext) emitInstant(ev *trace.Event, name, category string) {
// 	if !tsWithinRange(ev.Ts, ctx.startTime, ctx.endTime) {
// 		return
// 	}

// 	cname := ""
// 	if ctx.mode&modeTaskOriented != 0 {
// 		taskID, isUserAnnotation := isUserAnnotationEvent(ev)

// 		show := false
// 		for _, task := range ctx.tasks {
// 			if isUserAnnotation && task.id == taskID || task.overlappingInstant(ev) {
// 				show = true
// 				break
// 			}
// 		}
// 		// grey out or skip if non-overlapping instant.
// 		if !show {
// 			if isUserAnnotation {
// 				return // don't display unrelated user annotation events.
// 			}
// 			cname = colorLightGrey
// 		}
// 	}
// 	var arg interface{}
// 	if ev.Type == trace.EvProcStart {
// 		type Arg struct {
// 			ThreadID uint64
// 		}
// 		arg = &Arg{ev.Args[0]}
// 	}
// 	ctx.emit(&ViewerEvent{
// 		Name:     name,
// 		Category: category,
// 		Phase:    "I",
// 		Scope:    "t",
// 		Time:     ctx.time(ev),
// 		Tid:      ctx.proc(ev),
// 		Stack:    ctx.stack(ev.Stk),
// 		Cname:    cname,
// 		Arg:      arg})
// }

// func (ctx *traceContext) emitArrow(ev *trace.Event, name string) {
// 	if ev.Link == nil {
// 		// The other end of the arrow is not captured in the trace.
// 		// For example, a goroutine was unblocked but was not scheduled before trace stop.
// 		return
// 	}
// 	if ctx.mode&modeGoroutineOriented != 0 && (!ctx.gs[ev.Link.G] || ev.Link.Ts < ctx.startTime || ev.Link.Ts > ctx.endTime) {
// 		return
// 	}

// 	if ev.P == trace.NetpollP || ev.P == trace.TimerP || ev.P == trace.SyscallP {
// 		// Trace-viewer discards arrows if they don't start/end inside of a slice or instant.
// 		// So emit a fake instant at the start of the arrow.
// 		ctx.emitInstant(&trace.Event{P: ev.P, Ts: ev.Ts}, "unblock", "")
// 	}

// 	color := ""
// 	if ctx.mode&modeTaskOriented != 0 {
// 		overlapping := false
// 		// skip non-overlapping arrows.
// 		for _, task := range ctx.tasks {
// 			if _, overlapped := task.overlappingDuration(ev); overlapped {
// 				overlapping = true
// 				break
// 			}
// 		}
// 		if !overlapping {
// 			return
// 		}
// 	}

// 	ctx.arrowSeq++
// 	ctx.emit(&ViewerEvent{Name: name, Phase: "s", Tid: ctx.proc(ev), ID: ctx.arrowSeq, Time: ctx.time(ev), Stack: ctx.stack(ev.Stk), Cname: color})
// 	ctx.emit(&ViewerEvent{Name: name, Phase: "t", Tid: ctx.proc(ev.Link), ID: ctx.arrowSeq, Time: ctx.time(ev.Link), Cname: color})
// }

// func (ctx *traceContext) stack(stk []*trace.Frame) int {
// 	return ctx.buildBranch(ctx.frameTree, stk)
// }

// // buildBranch builds one branch in the prefix tree rooted at ctx.frameTree.
// func (ctx *traceContext) buildBranch(parent frameNode, stk []*trace.Frame) int {
// 	if len(stk) == 0 {
// 		return parent.id
// 	}
// 	last := len(stk) - 1
// 	frame := stk[last]
// 	stk = stk[:last]

// 	node, ok := parent.children[frame.PC]
// 	if !ok {
// 		ctx.frameSeq++
// 		node.id = ctx.frameSeq
// 		node.children = make(map[uint64]frameNode)
// 		parent.children[frame.PC] = node
// 		ctx.consumer.consumeViewerFrame(strconv.Itoa(node.id), ViewerFrame{fmt.Sprintf("%v:%v", frame.Fn, frame.Line), parent.id})
// 	}
// 	return ctx.buildBranch(node, stk)
// }

// type jsonWriter struct {
// 	w   io.Writer
// 	enc *json.Encoder
// }

// func viewerDataTraceConsumer(w io.Writer, start, end int64) traceConsumer {
// 	frames := make(map[string]ViewerFrame)
// 	enc := json.NewEncoder(w)
// 	written := 0
// 	index := int64(-1)

// 	io.WriteString(w, "{")
// 	return traceConsumer{
// 		consumeTimeUnit: func(unit string) {
// 			io.WriteString(w, `"displayTimeUnit":`)
// 			enc.Encode(unit)
// 			io.WriteString(w, ",")
// 		},
// 		consumeViewerEvent: func(v *ViewerEvent, required bool) {
// 			index++
// 			if !required && (index < start || index > end) {
// 				// not in the range. Skip!
// 				return
// 			}
// 			if written == 0 {
// 				io.WriteString(w, `"traceEvents": [`)
// 			}
// 			if written > 0 {
// 				io.WriteString(w, ",")
// 			}
// 			enc.Encode(v)
// 			// TODO: get rid of the extra \n inserted by enc.Encode.
// 			// Same should be applied to splittingTraceConsumer.
// 			written++
// 		},
// 		consumeViewerFrame: func(k string, v ViewerFrame) {
// 			frames[k] = v
// 		},
// 		flush: func() {
// 			io.WriteString(w, `], "stackFrames":`)
// 			enc.Encode(frames)
// 			io.WriteString(w, `}`)
// 		},
// 	}
// }
