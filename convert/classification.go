package convert

import "math"

type Classification int

const (
	ClassificationGenericAPI Classification = iota + 1
	ClassificationDeepScope
	ClassificationAPIRequest
	ClassificationAPITracing
	ClassificationMXNetCAPI
	ClassificationOpen
	ClassificationURLs
	ClassificationClose
	ClassificationFrameworkLayer
	ClassificationUnknown = Classification(math.MaxInt32)
)
