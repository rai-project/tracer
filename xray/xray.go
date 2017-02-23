package xray

import (
	"encoding/json"
	"net"

	aaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/rai-project/aws"
	"github.com/rai-project/tracer"
)

type xrayClient struct {
	options Options
}

var (
	TraceHeader = "X-Amzn-Trace-Id"
)

func New(opts ...Option) (*xrayClient, error) {
	sess, err := aws.NewSession()
	if err != nil {
		return nil, err
	}
	options := Options{
		client: xray.New(sess),
		daemon: "localhost:2000",
	}

	for _, o := range opts {
		o(&options)
	}

	return &xrayClient{
		options: options,
	}, nil
}

func (x *xrayClient) Record(s tracer.Segment) error {
	s.RLock()
	b, err := json.Marshal(s)
	if err != nil {
		s.RUnlock()
		return err
	}
	s.RUnlock()

	// Use XRay Client if available
	if x.options.client != nil {
		_, err := x.options.client.PutTraceSegments(&xray.PutTraceSegmentsInput{
			TraceSegmentDocuments: []*string{
				aaws.String("TraceSegmentDocument"),
				aaws.String(string(b)),
			},
		})
		return err
	}
	// Use Daemon
	c, err := net.Dial("udp", x.options.daemon)
	if err != nil {
		return err
	}

	header := append([]byte(`{"format": "json", "version": 1}`), byte('\n'))
	_, err = c.Write(append(header, b...))
	return err
}

func (x *xrayClient) Options() Options {
	return x.options
}
