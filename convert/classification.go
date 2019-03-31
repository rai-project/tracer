package convert

import "math"

type Classification int

const (
	ClassificationGenericAPI Classification = iota + 1
	ClassificationAPIRequest
	ClassificationOpen
	ClassificationURLs
	ClassificationClose
	ClassificationFrameworkLayer
	ClassificationUnknown = Classification(math.MaxInt32)
)
