package main

import (
	"C"
)

import (
	"context"
	"math"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	"gitlab.com/NebulousLabs/fastrand"
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
)

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
func SpanStart(lvl int32, operationName string) uintptr {
	sp := tracer.StartSpan(tracer.Level(lvl), operationName)
	return spans.Add(sp)
}

//export SpanStartFromContext
func SpanStartFromContext(inCtx uintptr, lvl int32, operationName string) (uintptr, uintptr) {
	sp, ctx := tracer.StartSpanFromContext(contexts.Get(inCtx), tracer.Level(lvl), operationName)
	return spans.Add(sp), contexts.Add(ctx)
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
	if sp != nil {
		sp.Finish()
	}
	spans.Delete(spPtr)
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
