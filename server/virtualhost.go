package server

import (
	"net/http"

	"github.com/osfx/snail/middleware"
)

// virtualHost represents a virtual host/server.
type virtualHost struct {
	config     Config
	fileServer middleware.Handler
	stack      middleware.Handler
}

// buildStack builds the server's middleware stack based
// on its config. 
func (vh *virtualHost) buildStack() error {
	vh.fileServer = FileServer(http.Dir(vh.config.Root), []string{vh.config.ConfigFile})
	vh.compile(vh.config.Middleware["/"])

	return nil
}

func (vh *virtualHost) compile(layers []middleware.Middleware) {
	vh.stack = vh.fileServer // core app layer
	for i := len(layers) - 1; i >= 0; i-- {
		vh.stack = layers[i](vh.stack)
	}
}
