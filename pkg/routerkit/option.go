package routerkit

import "github.com/opentracing/opentracing-go"

type routerConfig struct {
	serviceName string
	spanOpts    []opentracing.StartSpanOption // additional span options to be applied
}

// RouterOption represents an option that can be passed to NewRouter.
type RouterOption func(*routerConfig)

func defaults(cfg *routerConfig) {
	cfg.serviceName = "mux.router"
}

// WithServiceName sets the given service name for the router.
func WithServiceName(name string) RouterOption {
	return func(cfg *routerConfig) {
		cfg.serviceName = name
	}
}

// WithSpanOptions applies the given set of options to the spans started
// by the router.
func WithSpanOptions(opts ...opentracing.StartSpanOption) RouterOption {
	return func(cfg *routerConfig) {
		cfg.spanOpts = opts
	}
}
