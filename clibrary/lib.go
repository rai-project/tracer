package main

import (
	"C"
)

import (
	"context"
	"unsafe"

	"github.com/opentracing/opentracing-go"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	// _ "github.com/rai-project/tracer/noop"
	// _ "github.com/rai-project/tracer/zipkin"
)

//go:nosplit
func fromSpan(sp opentracing.Span) uintptr {
	return (uintptr)(unsafe.Pointer(&sp))
}

//go:nosplit
func toSpan(sp uintptr) opentracing.Span {
	return *((*opentracing.Span)(unsafe.Pointer(sp)))
}

//go:nosplit
func fromContext(ctx context.Context) uintptr {
	return (uintptr)(unsafe.Pointer(&ctx))
}

//go:nosplit
func toContext(ctx uintptr) context.Context {
	return *((*context.Context)(unsafe.Pointer(ctx)))
}

//export StartSpan
func StartSpan(lvl int32, operationName string, tags map[string]string) uintptr {
	// sp := tracer.StartSpan(tracer.Level(lvl), operationName, opentracing.Tags(tags))
	// return (uintptr)(unsafe.Pointer(&sp))
	return fromSpan(tracer.StartSpan(tracer.Level(lvl), operationName, cTags(tags)))
}

//export SpanAddTag
func SpanAddTag(spPtr uintptr, key, val string) {
	sp := toSpan(spPtr)
	sp.SetTag(key, val)
}

//export SpanFinish
func SpanFinish(spPtr uintptr) {
	sp := toSpan(spPtr)
	sp.Finish()
}

//export StartSpanFromContext
func StartSpanFromContext(inCtx uintptr, lvl int32, operationName string, tags map[string]string) (uintptr, uintptr) {
	sp, ctx := tracer.StartSpanFromContext(toContext(inCtx), tracer.Level(lvl), operationName, cTags(tags))
	return fromSpan(sp), fromContext(ctx)
}
