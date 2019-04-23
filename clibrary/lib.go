package main

// #include <stdlib.h>
// #include <stdlib.h>
// #cgo CFLAGS: -fPIC -O3
import (
	"C"
)

import (
	"context"
	"math"
	"strconv"
	"sync"
	"unsafe"

	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"gitlab.com/NebulousLabs/fastrand"

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
  globalCtx context.Context
)

func init() {
  ctx := context.Background()
  globalSpan, globalCtx = tracer.StartSpanFromContext(ctx, tracer.LIBRARY_TRACE, "c_tracing")
  contexts.contexts[0] = globalCtx
}

//go:nosplit
func (s spanMap) Add(sp opentracing.Span) uintptr {
	id := uintptr(fastrand.Uint64n(math.MaxUint64))
	s.Lock()
	s.spans[id] = sp
	s.Unlock()
	return id
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
	id := uintptr(fastrand.Uint64n(math.MaxUint64))
	s.Lock()
	s.contexts[id] = ctx
	s.Unlock()
	return id
}

//go:nosplit
func (s contextMap) Get(id uintptr) context.Context {
	if id == 0 {
		return context.Background()
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
  pp.Println(operationName)
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
    return
  }
	keys := (*[1 << 28]*C.char)(unsafe.Pointer(ckeys))[:length:length]
	vals := (*[1 << 28]*C.char)(unsafe.Pointer(cvals))[:length:length]
	for ii := 0; ii < length; ii++ {
		sp.SetTag(C.GoString(keys[ii]), C.GoString(vals[ii]))
	}
}

//export SpanAddArgumentsTag
func SpanAddArgumentsTag(spPtr uintptr, length int, ckeys **C.char, cvals **C.char) {
  sp := spans.Get(spPtr)
  if sp == nil {
    return
  }
	keys := (*[1 << 28]*C.char)(unsafe.Pointer(ckeys))[:length:length]
	vals := (*[1 << 28]*C.char)(unsafe.Pointer(cvals))[:length:length]
	args := make(map[string]string, length)
	for ii := 0; ii < length; ii++ {
		args[C.GoString(keys[ii])] = C.GoString(vals[ii])
	}
	sp.SetTag("arguments", args)
}

//export SpanFinish
func SpanFinish(spPtr uintptr) {
  sp := spans.Get(spPtr)
	if sp != nil {
		sp.Finish()
	}
	spans.Delete(spPtr)
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
