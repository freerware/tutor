package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/tabwriter"

	"github.com/freerware/tutor/config"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MuxConfiguration struct {
	PathPrefix string
	Handlers   []HandlerConfiguration
}

type HandlerConfiguration struct {
	Path        string
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Methods     []string
}

type ServerParameters struct {
	fx.In

	Configuration    config.Configuration
	MuxConfiguration MuxConfiguration
	Logger           *zap.Logger
}

type Server struct {
	host       string
	port       int
	httpServer http.Server
	logger     *zap.Logger
}

func New(parameters ServerParameters) Server {

	serverConfig := parameters.Configuration.Server

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler: newMux(parameters.Logger, parameters.MuxConfiguration),
	}

	s := Server{
		port:       serverConfig.Port,
		host:       serverConfig.Host,
		httpServer: httpServer,
		logger:     parameters.Logger,
	}

	return s
}

func newMux(logger *zap.Logger, m MuxConfiguration) *mux.Router {

	r := mux.NewRouter()
	sr := r.PathPrefix(m.PathPrefix).Subrouter()
	for _, h := range m.Handlers {
		sr.HandleFunc(h.Path, h.HandlerFunc).Methods(h.Methods...)
	}
	printMux(logger, m)

	return r
}

func printMux(logger *zap.Logger, m MuxConfiguration) {
	logger.Info("~~~~~~~~~~ PATHS ~~~~~~~~~~")
	for _, h := range m.Handlers {
		for _, method := range h.Methods {
			b := bytes.NewBufferString("")
			w := tabwriter.NewWriter(b, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "%s\t%s%s", method, m.PathPrefix, h.Path)
			w.Flush()
			logger.Info(b.String())
		}
	}
	logger.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

func (s *Server) Start() error {
	s.logger.Info(
		"Starting HTTP server",
		zap.String("host", s.host),
		zap.Int("port", s.port),
	)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Host() string {
	return s.host
}
