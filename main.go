package main

import (
	. "colabnote/models"
	"database/sql"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	chars      = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	configPath string
)

func init() {
	flag.StringVar(&configPath, "-configpath", "../configs/colabnote.json", "path config file")
}
func main() {
	flag.Parse()
	conf := Config{}
	jsonFile, err := os.Open(configPath)
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	confjson, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(confjson, &conf)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("mysql", conf.DataSourceName)
	if err != nil {
		log.Printf("while connecting to db driver error happened %v", err)
		return
	}
	db.SetMaxOpenConns(5)
	err = db.Ping()
	if err != nil {
		log.Printf("while connecting to db error happened %v\n", err)
		return
	}
	handlers := &handlers{db: db}
	adminMux := mux.NewRouter()
	adminMux.HandleFunc("/api/getNote", handlers.getNote)
	adminMux.HandleFunc("/api/createNote", handlers.createNote)
	adminMux.HandleFunc("/api/deleteNoteById", handlers.deleteNoteById)
	adminMux.HandleFunc("/api/logIn", handlers.logIn)
	adminMux.HandleFunc("/api/register", handlers.register)
	adminMux.HandleFunc("/api/usersOnline", handlers.usersOnline)
	srv := &http.Server{
		Addr:         ":" + conf.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      adminMux,
	}
	srv.ListenAndServe()
}
