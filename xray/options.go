package xray

import "github.com/aws/aws-sdk-go/service/xray"

type Options struct {
	// XRay Client when using API
	client *xray.XRay
	// Daemon address when using UDP
	daemon string
}

type Option func(o *Options)

func Client(x *xray.XRay) Option {
	return func(o *Options) {
		o.client = x
	}
}

func Daemon(x string) Option {
	return func(o *Options) {
		o.daemon = x
	}
}
