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
		return err
	}
	defer database.Database.Close()
	logger.Info("server started at port " + config.Conf.Port)
	err = s.server.ListenAndServe()
	if err != nil {
		logger.Log(err)
		return err
	}
	return nil
}
func (s *Server) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
