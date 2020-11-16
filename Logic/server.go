package Logic

import (
	"context"
	"database/sql"
	"log"
	"net/http"
)

type server struct {
	server http.Server
	db     *sql.DB
}

func NewServer(conf Config) (*server, error) {
	data, err := sql.Open("mysql", conf.DataSourceName)
	return &server{
		server: http.Server{
			Addr:    ":" + conf.Port,
			Handler: nil,
		},
		db: data,
	}, err
}
func (s *server) Run() error {
	err := s.db.Ping()
	if err != nil {
		log.Println(err, "taima")
	}
	defer s.db.Close()
	s.initHandler()
	err = s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
func (s *server) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
