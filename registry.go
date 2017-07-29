package tracer

import (
	"errors"
	"strings"

	"golang.org/x/sync/syncmap"
)

var tracers syncmap.Map

func FromName(s string) (Tracer, error) {
	s = strings.ToLower(s)
	val, ok := tracers.Load(s)
	if !ok {
		log.WithField("tracer", s).
			Warn("cannot find tracer")
		return nil, errors.New("cannot find tracer")
	}
	tracer, ok := val.(Tracer)
	if !ok {
		log.WithField("tracer", s).
			Warn("invalid tracer")
		return nil, errors.New("invalid tracer")
	}
	return tracer, nil
}

func AddTracer(name string, s Tracer) {
	tracers.Store(strings.ToLower(name), s)
}

func Tracers() []string {
	names := []string{}
	tracers.Range(func(key, _ interface{}) bool {
		if name, ok := key.(string); ok {
			names = append(names, name)
		}
		return true
	})
	return names
}
