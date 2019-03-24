package convert

import (
	"hash/fnv"

	"github.com/Workiva/go-datastructures/augmentedtree"
	model "github.com/uber/jaeger/model/json"
)

type Interval struct {
	model.Span
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (i Interval) LowAtDimension(uint64) int64 {
	return int64(i.StartTime)
}

func (i Interval) HighAtDimension(uint64) int64 {
	return int64(i.StartTime + i.Duration)
}

func (i Interval) ID() uint64 {
	return hash(string(i.SpanID))
}

func (i Interval) OverlapsAtDimension(iv augmentedtree.Interval, dimension uint64) bool {
	return i.HighAtDimension(dimension) > iv.LowAtDimension(dimension) &&
		i.LowAtDimension(dimension) < iv.HighAtDimension(dimension)
}

func spanToInterval(s model.Span) Interval {
	return Interval{s}
}
