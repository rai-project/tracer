package convert

import (
	"time"

	"github.com/pkg/errors"
	model "github.com/uber/jaeger/model/json"
)

func timeUnit(unit string) (time.Duration, error) {
	switch unit {
	case "ns":
		return time.Nanosecond, nil
	case "us":
		return time.Microsecond, nil
	case "ms":
		return time.Millisecond, nil
	case "":
		return time.Microsecond, nil
	default:
		return time.Duration(0), errors.Errorf("the display time unit %v is not valid", unit)
	}
}

func getSpanTags(sp model.Span) map[string]interface{} {
	tags := sp.Tags
	res := map[string]interface{}{}
	for _, tag := range tags {
		res[tag.Key] = tag.Value
	}
	return res
}

func getSpanTagByKey(sp model.Span, key string) interface{} {
	tags := sp.Tags
	for _, tag := range tags {
		if tag.Key == key {
			return tag.Value
		}
	}
	return nil
}
