package server

import (
	"net"

	"github.com/osfx/snail/middleware"
)

// Configuration for a single server.
type Config struct {

	Host string
	BindHost string
	Port string
	Root string

	// HTTPS configuration
	TLS TLSConfig

	// Middleware stack; map of path scope to middleware
	Middleware map[string][]middleware.Middleware

	// Functions (or methods) to execute at server start
	Startup []func() error

	// Functions (or methods) to execute when the server quits;
	Shutdown []func() error
	
	// The path to the configuration file from which this was loaded
	ConfigFile string
	AppName string
	AppVersion string
}

// Address returns the host:port of c as a string.
func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

// TLSConfig describes how TLS should be configured and used,
type TLSConfig struct {
	Enabled                  bool
	Certificate              string
	Key                      string
	Ciphers                  []uint16
	ProtocolMinVersion       uint16
	ProtocolMaxVersion       uint16
	PreferServerCipherSuites bool
	ClientCerts              []string
}
