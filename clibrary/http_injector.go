package main

import (
	"strconv"
	"strings"

	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

// Option is a function that sets an option on Propagator
type PropagatorOption func(propagator *Propagator)

// BaggagePrefix is a function that sets baggage prefix on Propagator
func BaggagePrefix(prefix string) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.baggagePrefix = prefix
	}
}

// Propagator is an Injector and Extractor
type Propagator struct {
	baggagePrefix string
}

// Baggage is by default enabled and uses prefix 'baggage-'.
func NewEnvPropagator(opts ...PropagatorOption) *Propagator {
	p := &Propagator{baggagePrefix: "baggage-"}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Propagator) Inject(
	sc jaeger.SpanContext, abstractCarrier interface{}) error {
	textMapWriter, ok := abstractCarrier.(opentracing.TextMapWriter)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}

	textMapWriter.Set("x-a3-traceid", sc.TraceID().String())
	if sc.ParentID() != 0 {
		textMapWriter.Set("x-a3-parentspanid", strconv.FormatUint(uint64(sc.ParentID()), 16))
	}
	textMapWriter.Set("x-a3-spanid", strconv.FormatUint(uint64(sc.SpanID()), 16))
	if sc.IsSampled() {
		textMapWriter.Set("x-a3-sampled", "1")
	} else {
		textMapWriter.Set("x-a3-sampled", "0")
	}
	sc.ForeachBaggageItem(func(k, v string) bool {
		textMapWriter.Set(p.baggagePrefix+k, v)
		return true
	})
	return nil
}

// Extract conforms to the Extractor interface for encoding Zipkin HTTP a3 headers
func (p Propagator) Extract(abstractCarrier interface{}) (jaeger.SpanContext, error) {
	textMapReader, ok := abstractCarrier.(opentracing.TextMapReader)
	if !ok {
		return jaeger.SpanContext{}, opentracing.ErrInvalidCarrier
	}
	var traceID jaeger.TraceID
	var spanID uint64
	var parentID uint64
	sampled := false
	var baggage map[string]string
	err := textMapReader.ForeachKey(func(rawKey, value string) error {
		key := strings.ToLower(rawKey) // TODO not necessary for plain TextMap
		var err error
		if key == "x-a3-traceid" {
			traceID, err = jaeger.TraceIDFromString(value)
		} else if key == "x-a3-parentspanid" {
			parentID, err = strconv.ParseUint(value, 16, 64)
		} else if key == "x-a3-spanid" {
			spanID, err = strconv.ParseUint(value, 16, 64)
		} else if key == "x-a3-sampled" && (value == "1" || value == "true") {
			sampled = true
		} else if strings.HasPrefix(key, p.baggagePrefix) {
			if baggage == nil {
				baggage = make(map[string]string)
			}
			baggage[key[len(p.baggagePrefix):]] = value
		}
		return err
	})

	if err != nil {
		return jaeger.SpanContext{}, err
	}
	if !traceID.IsValid() {
		return jaeger.SpanContext{}, opentracing.ErrSpanContextNotFound
	}
	return jaeger.NewSpanContext(
		traceID,
		jaeger.SpanID(spanID),
		jaeger.SpanID(parentID),
		sampled, baggage), nil
}
