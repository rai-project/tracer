package main

import (
	"C"
)

import (
	"context"
	"sync"
	"unsafe"

	"github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	// _ "github.com/rai-project/tracer/noop"
	// _ "github.com/rai-project/tracer/zipkin"
)

type spanMap struct {
	spans map[uintptr]opentracing.Span
	sync.Mutex
}

var spans = spanMap{}

func (s spanMap) Add(sp opentracing.Span) uintptr {
	id := fromSpan(sp)
	s.Lock()
	s.spans[id] = sp
	s.Unlock()
	return id
}

func (s spanMap) Get(id uintptr) opentracing.Span {
	s.Lock()
	res := s.spans[id]
	s.Unlock()
	return res
}

func (s spanMap) Delete(id uintptr) {
	s.Lock()
	delete(s.spans, id)
	s.Unlock()
}

//go:nosplit
func fromSpan(sp opentracing.Span) uintptr {
	return (uintptr)(unsafe.Pointer(&sp))
}

//go:nosplit
func fromContext(ctx context.Context) uintptr {
	return (uintptr)(unsafe.Pointer(&ctx))
}

//go:nosplit
func toContext(ctx uintptr) context.Context {
	return *((*context.Context)(unsafe.Pointer(ctx)))
}

//export SpanStart
func SpanStart(lvl int32, operationName string) uintptr {
	sp := tracer.StartSpan(tracer.Level(lvl), operationName)
	return spans.Add(sp)
}

//export SpanAddTag
func SpanAddTag(spPtr uintptr, key, val string) {
	sp := spans.Get(spPtr)
	sp.SetTag(key, val)
}

//export SpanAddTags
func SpanAddTags(spPtr uintptr, len int, keys []string, vals []string) {
	sp := spans.Get(spPtr)
	for ii := 0; ii < len; ii++ {
		sp.SetTag(keys[ii], vals[ii])
	}
}

//export SpanAddArgumentsTag
func SpanAddArgumentsTag(spPtr uintptr, len int, keys []string, vals []string) {
	sp := spans.Get(spPtr)
	args := make(map[string]string, len)
	for ii := 0; ii < len; ii++ {
		args[keys[ii]] = vals[ii]
	}
	sp.SetTag("arguments", args)
}

//export SpanFinish
func SpanFinish(spPtr uintptr) {
	sp := spans.Get(spPtr)
	sp.Finish()
	spans.Delete(spPtr)
}

//export StartSpanFromContext
func StartSpanFromContext(inCtx uintptr, lvl int32, operationName string, tags map[string]string) (uintptr, uintptr) {
	sp, ctx := tracer.StartSpanFromContext(toContext(inCtx), tracer.Level(lvl), operationName, cTags(tags))
	return spans.Add(sp), fromContext(ctx)
}
