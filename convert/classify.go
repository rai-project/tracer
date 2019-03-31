package convert

import model "github.com/uber/jaeger/model/json"

func Classify(sp model.Span) string {
	operationName := sp.OperationName
	_ = operationName
	return "U"
}
