// Package server implements a configurable, general-purpose web server.
package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/bradfitz/http2"
)

// Server represents an instance of a server, which serves
// static content at a particular address (host and port).
type Server struct {
	HTTP2   bool
	address string
	tls     bool
	vhosts  map[string]virtualHost
}

func New(addr string, configs []Config) (*Server, error) {
	var tls bool
	if len(configs) > 0 {
		tls = configs[0].TLS.Enabled
	}

	s := &Server{
		address: addr,
		tls:     tls,
		vhosts:  make(map[string]virtualHost),
	}

	for _, conf := range configs {
		if _, exists := s.vhosts[conf.Host]; exists {
			return nil, fmt.Errorf("cannot serve %s - host already defined for address %s", conf.Address(), s.address)
		}

		vh := virtualHost{config: conf}

		// Build middleware stack
		err := vh.buildStack()
		if err != nil {
			return nil, err
		}

		s.vhosts[conf.Host] = vh
	}

	return s, nil
}

func (s *Server) Serve() error {
	server := &http.Server{
		Addr:    s.address,
		Handler: s,
	}

	if s.HTTP2 {

		http2.ConfigureServer(server, nil)
	}

	for _, vh := range s.vhosts {

		for _, start := range vh.config.Startup {
			err := start()
			if err != nil {
				return err
			}
		}

		// Execute shutdown commands on exit
		if len(vh.config.Shutdown) > 0 {
			go func(vh virtualHost) {
				// Wait for signal
				interrupt := make(chan os.Signal, 1)
				signal.Notify(interrupt, os.Interrupt, os.Kill)
				<-interrupt

				// Run callbacks
				exitCode := 0
				for _, shutdownFunc := range vh.config.Shutdown {
					err := shutdownFunc()
					if err != nil {
						exitCode = 1
						log.Println(err)
					}
				}
				os.Exit(exitCode)
			}(vh)
		}
	}

	if s.tls {
		var tlsConfigs []TLSConfig
		for _, vh := range s.vhosts {
			tlsConfigs = append(tlsConfigs, vh.config.TLS)
		}
		return ListenAndServeTLSWithSNI(server, tlsConfigs)
	}
	return server.ListenAndServe()
}

func ListenAndServeTLSWithSNI(srv *http.Server, tlsConfigs []TLSConfig) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	config := new(tls.Config)
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, len(tlsConfigs))
	for i, tlsConfig := range tlsConfigs {
		config.Certificates[i], err = tls.LoadX509KeyPair(tlsConfig.Certificate, tlsConfig.Key)
		if err != nil {
			return err
		}
	}
	config.BuildNameToCertificate()

	config.MinVersion = tlsConfigs[0].ProtocolMinVersion
	config.MaxVersion = tlsConfigs[0].ProtocolMaxVersion
	config.CipherSuites = tlsConfigs[0].Ciphers
	config.PreferServerCipherSuites = tlsConfigs[0].PreferServerCipherSuites

	err = setupClientAuth(tlsConfigs, config)
	if err != nil {
		return err
	}

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(conn, config)

	return srv.Serve(tlsListener)
}

func setupClientAuth(tlsConfigs []TLSConfig, config *tls.Config) error {
	var clientAuth bool
	for _, cfg := range tlsConfigs {
		if len(cfg.ClientCerts) > 0 {
			clientAuth = true
			break
		}
	}

	if clientAuth {
		pool := x509.NewCertPool()
		for _, cfg := range tlsConfigs {
			for _, caFile := range cfg.ClientCerts {
				caCrt, err := ioutil.ReadFile(caFile)
				if err != nil {
					return err
				}
				if !pool.AppendCertsFromPEM(caCrt) {
					return fmt.Errorf("error loading client certificate '%s': no certificates were successfully parsed", caFile)
				}
			}
		}
		config.ClientCAs = pool
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {

		if rec := recover(); rec != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}()

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	if _, ok := s.vhosts[host]; !ok {
		if _, ok2 := s.vhosts["0.0.0.0"]; ok2 {
			host = "0.0.0.0"
		} else if _, ok2 := s.vhosts[""]; ok2 {
			host = ""
		}
	}

	if vh, ok := s.vhosts[host]; ok {
		w.Header().Set("Server", "Snail")

		status, _ := vh.stack.ServeHTTP(w, r)

		if status >= 400 {
			DefaultErrorFunc(w, r, status)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No such host at %s", s.address)
	}
}

func DefaultErrorFunc(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%d %s", status, http.StatusText(status))
}
