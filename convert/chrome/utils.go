package chrome

import (
	"hash/fnv"
	"time"

	m "github.com/uber/jaeger/model"
)

func hash64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	res := h.Sum64()
	return res
}

func hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	res := h.Sum32()
	return res
}

func toTime(t uint64) time.Time {
	return m.EpochMicrosecondsAsTime(t)
}

func toDuration(d uint64) time.Duration {
	return m.MicrosecondsAsDuration(d)
}
