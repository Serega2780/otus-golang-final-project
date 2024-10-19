package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/config"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/logger"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/service"
	"github.com/gorilla/mux"
)

type Server struct {
	ctx     context.Context
	ip      string
	port    string
	log     *logger.Logger
	srv     *http.Server
	service service.ImageService
}

func NewServer(ctx context.Context, logger *logger.Logger, conf *config.HTTPServerConfig,
	service service.ImageService,
) *Server {
	return &Server{ctx: ctx, log: logger, ip: conf.IP, port: conf.Port, service: service}
}

func (s *Server) Start(ctx context.Context) {
	h := NewProxyHandler(ctx, s.log, s.service)

	r := mux.NewRouter().UseEncodedPath()
	r.PathPrefix("/fill/").Handler(http.StripPrefix("/fill/", h))
	server := &http.Server{
		Addr:              strings.Join([]string{s.ip, s.port}, ":"),
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}
	s.srv = server
	s.ctx = ctx

	s.log.Infof("http server start on port %s", s.port)
	go func() {
		_ = server.ListenAndServe()
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("could not shutdown http server %v", err)
		return
	}
	s.log.Info("http server stopped")
}
