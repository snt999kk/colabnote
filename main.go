package main

import (
	. "colabnote/models"
	"database/sql"
	_ "encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	flag.StringVar(&configPath, "-configpath", "conf.json", "path config file")
}
func main() {
	flag.Parse()
	conf := config{}
	jsonFile, _ := os.Open(configPath)
	fmt.Print(jsonFile)
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("mysql", "ramazan:Slamdunk19984pda!@tcp(localhost:3306)/mysql?&charset=utf8&interpolateParams=true")
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
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      adminMux,
	}
	go srv.ListenAndServe()
}
