package setup

import (
	"github.com/osfx/snail/middleware"
	"github.com/osfx/snail/middleware/proxy"
)

// Proxy configures a new Proxy middleware instance.
func Proxy(c *Controller) (middleware.Middleware, error) {
	if upstreams, err := proxy.NewStaticUpstreams(c.Dispenser); err == nil {
		return func(next middleware.Handler) middleware.Handler {
			return proxy.Proxy{Next: next, Upstreams: upstreams}
		}, nil
	} else {
		return nil, err
	}
}
