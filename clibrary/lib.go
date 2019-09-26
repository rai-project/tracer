package main

// #include <stdint.h>
// #cgo CFLAGS: -fPIC -O3
import (
	"C"
)

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/opentracing/opentracing-go"
	"gitlab.com/NebulousLabs/fastrand"

	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	// _ "github.com/rai-project/tracer/noop"
)

type spanInfo struct {
	ctx           context.Context
	level         tracer.Level
	operationName string
	startTime     time.Time
	endTime       time.Time
	tags          opentracing.Tags
}

type spanMap struct {
	spans map[uintptr]*spanInfo
	sync.RWMutex
}

type contextMap struct {
	contexts map[uintptr]context.Context
	sync.RWMutex
}

var (
	spans = &spanMap{
		spans: map[uintptr]*spanInfo{},
	}
	contexts = &contextMap{
		contexts: map[uintptr]context.Context{},
	}
)

//go:nosplit
func (s *spanMap) Add(sp *spanInfo) uintptr {
	s.Lock()
	defer s.Unlock()
	if s.spans == nil {
		s.spans = map[uintptr]*spanInfo{}
	}
	for {
		id := uintptr(fastrand.Uint64n(math.MaxUint64))
		if _, ok := s.spans[id]; ok {
			continue
		}
		s.spans[id] = sp
		return id
	}
	return 0
}

//go:nosplit
func (s *spanMap) Get(id uintptr) *spanInfo {
	s.RLock()
	res := s.spans[id]
	s.RUnlock()
	return res
}

//go:nosplit
func (s *spanMap) Delete(id uintptr) {
	s.Lock()
	delete(s.spans, id)
	s.Unlock()
}

//go:nosplit
func (s *contextMap) Add(ctx context.Context) uintptr {
	s.Lock()
	defer s.Unlock()
	for {
		id := uintptr(fastrand.Uint64n(math.MaxUint64))
		if _, ok := s.contexts[id]; ok {
			continue
		}
		s.contexts[id] = ctx
		return id
	}
	return 0
}

//go:nosplit
func (s *contextMap) Get(id uintptr) context.Context {
	if id == 0 {
		return globalCtx
	}
	s.RLock()
	res := s.contexts[id]
	s.RUnlock()
	return res
}

//go:nosplit
func (s *contextMap) Delete(id uintptr) {
	if id == 0 {
		return
	}
	s.Lock()
	delete(s.contexts, id)
	s.Unlock()
}

//export SpanStart
func SpanStart(lvl C.int32_t, cOperationName *C.char) uintptr {
	initLib()
	now := time.Now()
	operationName := C.GoString(cOperationName)
	sp := &spanInfo{
		ctx:           contexts.Get(0),
		level:         tracer.Level(lvl),
		operationName: operationName,
		startTime:     now,
		tags:          opentracing.Tags{},
	}
	spPtr := spans.Add(sp)

	return spPtr
}

//export SpanStartFromContext
func SpanStartFromContext(inCtx uintptr, lvl int32, cOperationName *C.char) (uintptr, uintptr) {
	initLib()
	now := time.Now()
	operationName := C.GoString(cOperationName)
	// pp.Println("starting span  ", operationName)
	sp := &spanInfo{
		ctx:           contexts.Get(inCtx),
		level:         tracer.Level(lvl),
		operationName: operationName,
		startTime:     now,
		tags:          opentracing.Tags{},
	}
	spPtr := spans.Add(sp)

	return spPtr, 0
}

//export SpanAddTag
func SpanAddTag(spPtr uintptr, key *C.char, val *C.char) {
	sp := spans.Get(spPtr)
	if sp == nil {
		return
	}
	sp.tags[C.GoString(key)] = C.GoString(val)
}

//export SpanAddTags
func SpanAddTags(spPtr uintptr, length int, ckeys **C.char, cvals **C.char) {
	sp := spans.Get(spPtr)
	if sp == nil {
		pp.Println("span is nil")
		return
	}
	if length == 0 {
		pp.Println("got no tags")
		return
	}
	keys := (*[1 << 28]*C.char)(unsafe.Pointer(ckeys))[:length:length]
	vals := (*[1 << 28]*C.char)(unsafe.Pointer(cvals))[:length:length]
	for ii := 0; ii < length; ii++ {
		goKey := C.GoString(keys[ii])
		goVal := C.GoString(vals[ii])
		// if goKey == "function_name" {
		// 	pp.Println(goVal)
		// }
		sp.tags[goKey] = goVal
	}
}

//export SpanAddArgumentsTag
func SpanAddArgumentsTag(spPtr uintptr, length int, ckeys **C.char, cvals **C.char) {
	sp := spans.Get(spPtr)
	if sp == nil {
		pp.Println("span is nil")
		return
	}
	if length == 0 {
		return
	}
	keys := (*[1 << 28]*C.char)(unsafe.Pointer(ckeys))[:length:length]
	vals := (*[1 << 28]*C.char)(unsafe.Pointer(cvals))[:length:length]
	// args := make([]Argument, length)
	for ii := 0; ii < length; ii++ {
		goKey := C.GoString(keys[ii])
		goVal := C.GoString(vals[ii])
		sp.tags["argument/"+strconv.Itoa(ii)+"/"+goKey] = goVal
	}
}

type Argument struct {
	Name  string `json:"n,omitempty"`
	Value string `json:"v,omitempty"`
}

//export SpanAddArgumentsTagJSON
func SpanAddArgumentsTagJSON(spPtr uintptr, length int, ckeys **C.char, cvals **C.char) {
	sp := spans.Get(spPtr)
	if sp == nil {
		pp.Println("span is nil")
		return
	}
	if length == 0 {
		return
	}
	keys := (*[1 << 28]*C.char)(unsafe.Pointer(ckeys))[:length:length]
	vals := (*[1 << 28]*C.char)(unsafe.Pointer(cvals))[:length:length]
	args := make([]Argument, length)
	for ii := 0; ii < length; ii++ {
		goKey := C.GoString(keys[ii])
		goVal := C.GoString(vals[ii])
		args[ii] = Argument{
			Name:  goKey,
			Value: goVal,
		}
	}
	bts, err := json.Marshal(args)
	if err != nil {
		return
	}
	sp.tags["arguments"] = string(bts)
}

//export SpanFinish
func SpanFinish(spPtr uintptr) {
	sp := spans.Get(spPtr)
	if sp != nil {
		sp.endTime = time.Now()
	}
}

//export SpanDelete
func SpanDelete(spPtr uintptr) {
	sp := spans.Get(spPtr)
	if sp != nil {
		span, _ := tracer.StartSpanFromContext(
			sp.ctx,
			sp.level,
			sp.operationName,
			opentracing.StartTime(sp.startTime),
			sp.tags,
		)
		if span == nil {
			pp.Println("nil span")
			return
		}
		span.FinishWithOptions(opentracing.FinishOptions{
			FinishTime: sp.endTime,
		})
		spans.Delete(spPtr)
	}
}

//export SpanGetTraceID
func SpanGetTraceID(spPtr uintptr) *C.char {
	// sp := spans.Get(spPtr)
	// traceID := sp.Context().(jaeger.SpanContext).TraceID()
	// return C.CString(strconv.FormatUint(traceID.Low, 16))
	return nil
}

//export ContextNewBackground
func ContextNewBackground() uintptr {
	ctx := context.Background()
	return contexts.Add(ctx)
}

//export ContextDelete
func ContextDelete(ctx uintptr) {
	contexts.Delete(ctx)
}
