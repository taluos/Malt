package discovery

import "time"

type builderOptions struct {
	timeout  time.Duration
	insecure bool
}

type BuilderOptions func(o *builderOptions)

func WithTimeout(timeout time.Duration) BuilderOptions {
	return func(o *builderOptions) {
		o.timeout = timeout
	}
}

func WithInsecure(insecure bool) BuilderOptions {
	return func(o *builderOptions) {
		o.insecure = insecure
	}
}
