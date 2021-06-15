package main

import (
	"colabnote/internal/config"
	"colabnote/internal/database"
	"colabnote/internal/logger"
	"colabnote/internal/model/entities"
	"colabnote/internal/service"
	"context"
	"flag"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "-configpath", "configs/conf.json", "path config file")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}
func main() {
	flag.Parse()
	err := config.ParseConf(configPath)
	if err != nil {
		logger.Log(err)
	}
	err = database.InitDB(config.Conf)
	if err != nil {
		logger.Log(err)
	}
	adminMux := mux.NewRouter()
	adminMux.HandleFunc("/api/getNote", service.GetNote).Methods("GET")
	adminMux.HandleFunc("/api/createNote", service.CreateNote).Methods("POST")
	adminMux.HandleFunc("/api/deleteNoteById", service.DeleteNoteById).Methods("DELETE")
	adminMux.HandleFunc("/api/logIn", service.LogIn).Methods("POST")
	adminMux.HandleFunc("/api/register", service.Register).Methods("POST")
	myServer := entities.NewServer(config.Conf, adminMux)
	go func() {
		err = myServer.Run()
		if err != nil {
			logger.Log(err)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	if err != nil {
		cancel()
	}
	go func(cancel context.CancelFunc) {
		sign := make(chan os.Signal)
		signal.Notify(sign, os.Interrupt)
		for sn := range sign {
			if sn == os.Interrupt {
				cancel()
			}
		}
	}(cancel)
	<-ctx.Done()
	ctxShutDown, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	myServer.Shutdown(ctxShutDown)
	cancel1()
}
