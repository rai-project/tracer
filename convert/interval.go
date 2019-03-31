package convert

import (
	"hash/fnv"
	"math"

	"github.com/Workiva/go-datastructures/augmentedtree"
	model "github.com/uber/jaeger/model/json"
)

type Interval struct {
	*model.Span
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	res := h.Sum64()
	return res
}

func (i Interval) Start() uint64 {
	if i.Span == nil {
		return 0
	}
	return uint64(i.StartTime)
}

func (i Interval) End() uint64 {
	if i.Span == nil {
		return math.MaxUint64
	}
	return uint64(i.StartTime + i.Duration)
}

func (i Interval) Contains(iv Interval) bool {
	if i.Span == nil {
		return false
	}
	if iv.Span == nil {
		return false
	}
	return i.Start() <= iv.Start() && i.End() >= iv.End()
}

func (i Interval) LowAtDimension(uint64) int64 {
	return int64(i.Start())
}

func (i Interval) HighAtDimension(uint64) int64 {
	return int64(i.End())
}

func (i Interval) ID() uint64 {
	return hash(string(i.SpanID))
}

// func (i Interval) OverlapsAtDimension(iv augmentedtree.Interval, dimension uint64) bool {
// 	return i.HighAtDimension(dimension) > iv.HighAtDimension(dimension) &&
// 		i.LowAtDimension(dimension) < iv.LowAtDimension(dimension)
// }

// func (i Interval) OverlapsAtDimension(iv augmentedtree.Interval, dimension uint64) bool {
// 	if i.Span == nil {
// 		return false
// 	}
// 	if iv.(Interval).Span == nil {
// 		return false
// 	}
// 	return i.HighAtDimension(dimension) > iv.LowAtDimension(dimension) &&
// 		i.LowAtDimension(dimension) < iv.HighAtDimension(dimension)
// }

func (ci Interval) OverlapsAtDimension(iv0 augmentedtree.Interval, dim uint64) bool {
	if ci.Span == nil {
		return false
	}
	iv := iv0.(Interval)
	if iv.Span == nil {
		return false
	}

	if (uint64(iv.LowAtDimension(dim)) <= ci.Start()) && (ci.End() <= uint64(iv.HighAtDimension(dim))) {
		// self       ================
		// other   =====================
		return true
	} else if (ci.Start() <= uint64(iv.LowAtDimension(dim))) && (uint64(iv.LowAtDimension(dim)) <= ci.End()) {
		// self      ================
		// other         ===============
		return true
	} else if (ci.Start() <= uint64(iv.HighAtDimension(dim))) && (uint64(iv.HighAtDimension(dim)) <= ci.End()) {
		// self      ===============
		// other  =================
		return true
	}
	return false
}

func spanToInterval(s model.Span) Interval {
	return Interval{&s}
}
