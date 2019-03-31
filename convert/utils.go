package convert

import (
	"time"

	"github.com/pkg/errors"
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
