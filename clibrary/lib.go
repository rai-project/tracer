package main

// #include <stdlib.h>
// #include <stdlib.h>
// #cgo CFLAGS: -fPIC -O3
import (
	"C"
)

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"

	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	// _ "github.com/rai-project/tracer/noop"
	// _ "github.com/rai-project/tracer/zipkin"
)

type spanMap struct {
	spans map[uintptr]opentracing.Span
	sync.Mutex
}

type contextMap struct {
	contexts map[uintptr]context.Context
	sync.Mutex
}

var (
	spans = spanMap{
		spans: make(map[uintptr]opentracing.Span),
	}
	contexts = contextMap{
		contexts: make(map[uintptr]context.Context),
	}
	globalSpan opentracing.Span
	globalCtx  context.Context

	spanCounter uintptr = 1
	ctxCounter  uintptr = 1
)

func libInit() {
	globalSpan, globalCtx = tracer.StartSpanFromContext(
		context.Background(),
		tracer.APPLICATION_TRACE,
		"c_tracing",
	)
	// pp.Println("init lib")
}

func libDeinit() {
	time.Sleep(time.Second)
	if globalSpan != nil {
		globalSpan.Finish()
		pp.Println("closing global span")

		traceID := globalSpan.Context().(jaeger.SpanContext).TraceID()
		traceIDVal := traceID.String()

		pp.Println(fmt.Sprintf("http://%s:16686/trace/%v", "192.17.102.10", traceIDVal))

	}
}

//go:nosplit
func (s spanMap) Add(sp opentracing.Span) uintptr {
	s.Lock()
	defer s.Unlock()
	for {
		// id := uintptr(fastrand.Uint64n(math.MaxUint64))
		id := spanCounter
		spanCounter++
		if _, ok := s.spans[id]; ok {
			continue
		}
		s.spans[id] = sp
		return id
	}
	return 0
}

//go:nosplit
func (s spanMap) Get(id uintptr) opentracing.Span {
	s.Lock()
	res := s.spans[id]
	s.Unlock()
	return res
}

//go:nosplit
func (s spanMap) Delete(id uintptr) {
	s.Lock()
	delete(s.spans, id)
	s.Unlock()
}

//go:nosplit
func (s contextMap) Add(ctx context.Context) uintptr {
	s.Lock()
	defer s.Unlock()
	for {
		// id := uintptr(fastrand.Uint64n(math.MaxUint64))
		id := ctxCounter
		ctxCounter++
		if _, ok := s.contexts[id]; ok {
			continue
		}
		s.contexts[id] = ctx
		return id
	}
	return 0
}

//go:nosplit
func (s contextMap) Get(id uintptr) context.Context {
	if id == 0 {
		return globalCtx
	}
	s.Lock()
	res := s.contexts[id]
	s.Unlock()
	return res
}

//go:nosplit
func (s contextMap) Delete(id uintptr) {
	if id == 0 {
		return
	}
	s.Lock()
	delete(s.contexts, id)
	s.Unlock()
}

//export SpanStart
func SpanStart(lvl C.int32_t, cOperationName *C.char) uintptr {
	operationName := C.GoString(cOperationName)
	sp := tracer.StartSpan(tracer.Level(lvl), operationName)
	return spans.Add(sp)
}

//export SpanStartFromContext
func SpanStartFromContext(inCtx uintptr, lvl int32, cOperationName *C.char) (uintptr, uintptr) {
	operationName := C.GoString(cOperationName)
	sp, ctx := tracer.StartSpanFromContext(contexts.Get(inCtx), tracer.Level(lvl), operationName)
	spPtr, ctxPtr := spans.Add(sp), contexts.Add(ctx)

	return spPtr, ctxPtr
}

//export SpanAddTag
func SpanAddTag(spPtr uintptr, key *C.char, val *C.char) {
	sp := spans.Get(spPtr)
	if sp == nil {
		return
	}
	sp.SetTag(C.GoString(key), C.GoString(val))
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
		sp.SetTag(goKey, goVal)
	}
}

type Argument struct {
	Name  string `json:"n,omitempty"`
	Value string `json:"v,omitempty"`
}

//export SpanAddArgumentsTag
func SpanAddArgumentsTag(spPtr uintptr, length int, ckeys **C.char, cvals **C.char) {
	sp := spans.Get(spPtr)
	if sp == nil {
		pp.Println("span is nil")
		return
	}
	if length == 0 {
		pp.Println("go no arguments")
		return
	}
	keys := (*[1 << 28]*C.char)(unsafe.Pointer(ckeys))[:length:length]
	vals := (*[1 << 28]*C.char)(unsafe.Pointer(cvals))[:length:length]
	args := make([]Argument, length)
	for ii := 0; ii < length; ii++ {
		goKey := C.GoString(keys[ii])
		goVal := C.GoString(vals[ii])
		// if false && goKey == "function_name" {
		// 	pp.Println(goVal)
		// }
		args[ii] = Argument{
			Name:  goKey,
			Value: goVal,
		}
	}
	bts, err := json.Marshal(args)
	if err != nil {
		return
	}
	sp.SetTag("arguments", string(bts))
}

//export SpanFinish
func SpanFinish(spPtr uintptr) {
	sp := spans.Get(spPtr)
	if sp != nil {
		sp.Finish()
	}
	time.Sleep(100 * time.Microsecond)
}

//export SpanDelete
func SpanDelete(spPtr uintptr) {
	sp := spans.Get(spPtr)
	if sp != nil {
		spans.Delete(spPtr)
	}
}

//export SpanGetTraceID
func SpanGetTraceID(spPtr uintptr) *C.char {
	sp := spans.Get(spPtr)
	traceID := sp.Context().(jaeger.SpanContext).TraceID()
	return C.CString(strconv.FormatUint(traceID.Low, 16))
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
