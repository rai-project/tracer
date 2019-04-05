package convert

import model "github.com/uber/jaeger/model/json"

type TraceInformation struct {
	Traces []model.Trace     `bson:"data,omitempty" json:"data,omitempty"`
	Total  int               `bson:"total,omitempty" json:"total,omitempty"`
	Limit  int               `bson:"limit,omitempty" json:"limit,omitempty"`
	Offset int               `bson:"offset,omitempty" json:"offset,omitempty"`
	Errors []structuredError `bson:"errors,omitempty" json:"errors,omitempty"`
}

type structuredError struct {
	Code    int           `json:"code,omitempty" bson:"code,omitempty"`
	Msg     string        `json:"msg,omitempty" bson:"msg,omitempty"`
	TraceID model.TraceID `json:"traceID,omitempty" bson:"traceID,omitempty"`
}
