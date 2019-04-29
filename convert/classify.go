package convert

import (
	"errors"
	"strconv"

	model "github.com/uber/jaeger/model/json"
)

var (
	operationClassifications = map[string]Classification{
		"deepscope":   ClassificationDeepScope,
		"api_request": ClassificationAPIRequest,
		"api_tracing": ClassificationAPITracing,
		"/mlmodelscope.org.dlframework.Predict/Open":  ClassificationOpen,
		"/mlmodelscope.org.dlframework.Predict/URLs":  ClassificationURLs,
		"/mlmodelscope.org.dlframework.Predict/Close": ClassificationClose,
	}
)

func FrameworkLayerIndex(sp model.Span) (int, error) {
	val := getSpanTagByKey(sp, "layer_sequence_index")
	if val == nil {
		return -1, errors.New("not a framework layer")
	}
	e, ok := val.(string)
	if !ok {
		return -1, errors.New("not a framework layer")
	}
	return strconv.Atoi(e)
}

func isAPISpan(sp model.Span) bool {
	val := getSpanTagByKey(sp, "library_name")
	if val == nil {
		return false
	}
	e, ok := val.(string)
	if !ok {
		return false
	}
	return e != ""
}

func isLayerSpan(sp model.Span) bool {
	val := getSpanTagByKey(sp, "layer_sequence_index")
	if val == nil {
		return false
	}
	e, ok := val.(string)
	if !ok {
		return false
	}
	return e != ""
}

func isMXNETCAPISpan(sp model.Span) bool {
	val := getSpanTagByKey(sp, "category")
	if val == nil {
		return false
	}
	e, ok := val.(string)
	if !ok {
		return false
	}
	return e == "MXNET_C_API"
}

func Classify(sp model.Span) string {
	operationName := sp.OperationName
	if val, ok := operationClassifications[operationName]; ok {
		return val.String()
	}
	if isMXNETCAPISpan(sp) {
		return ClassificationMXNetCAPI.String()
	}
	if isLayerSpan(sp) {
		return ClassificationFrameworkLayer.String()
	}
	if isAPISpan(sp) {
		if e, ok := getSpanTagByKey(sp, "library_name").(string); ok {
			return e
		}
	}
	return ClassificationUnknown.String()
}
