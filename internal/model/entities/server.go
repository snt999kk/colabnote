package entities

import (
	"colabnote/internal/config"
	"colabnote/internal/database"
	"colabnote/internal/logger"
	"context"
	"log"
	"net/http"
)

type Server struct {
	server http.Server
}

func NewServer(conf config.Config, handler http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    "localhost:" + conf.Port,
			Handler: handler,
		},
	}
}
func (s *Server) Run() error {
	err := database.Database.Ping()
	if err != nil {
		log.Println(err)
	}
	defer database.Database.Close()
	err = s.server.ListenAndServe()
	if err != nil {
		logger.Log(err)
		return err
	}
	logger.Info("server started at port " + config.Conf.Port)
	return nil
}
func (s *Server) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
