package tracer

import (
	"errors"
	"strings"

	"golang.org/x/sync/syncmap"
)

var (
	tracers     syncmap.Map
	openTracers syncmap.Map
)

type tracerRegistryItem struct {
	tracer Tracer
	new    func(serviceName string) (Tracer, error)
}

func FromName(s string) (Tracer, error) {
	s = strings.ToLower(s)
	val, ok := tracers.Load(s)
	if !ok {
		log.WithField("tracer", s).
			Warn("cannot find tracer")
		return nil, errors.New("cannot find tracer")
	}
	tracer, ok := val.(tracerRegistryItem)
	if !ok {
		log.WithField("tracer", s).
			Warn("invalid tracer")
		return nil, errors.New("invalid tracer")
	}
	return tracer.tracer, nil
}

func NewFromName(serviceName, backendName string) (Tracer, error) {
	s := strings.ToLower(backendName)
	val, ok := tracers.Load(s)
	if !ok {
		log.WithField("tracer", s).
			Warn("cannot find tracer")
		return nil, errors.New("cannot find tracer")
	}
	tracer, ok := val.(tracerRegistryItem)
	if !ok {
		log.WithField("tracer", s).
			Warn("invalid tracer")
		return nil, errors.New("invalid tracer")
	}
	tr, err := tracer.new(serviceName)
	if err != nil {
		return nil, err
	}
	openTracers.Store(tr.ID(), tr)
	return tr, nil
}

func Register(name string, s Tracer, newFunc func(serviceName string) (Tracer, error)) {
	tracers.Store(strings.ToLower(name), tracerRegistryItem{tracer: s, new: newFunc})
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
