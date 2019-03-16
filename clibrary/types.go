package main

import "github.com/opentracing/opentracing-go"

type cTags map[string]string

// Apply satisfies the StartSpanOption interface.
func (t cTags) Apply(o *opentracing.StartSpanOptions) {
	if o.Tags == nil {
		o.Tags = make(map[string]interface{})
	}
	for k, v := range t {
		o.Tags[k] = v
	}
}
