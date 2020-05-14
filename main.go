package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	_ "log"
	"math/rand"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("mysql", "reg-user:qwerty123123@tcp(localhost:3306)/mysql?&charset=utf8&interpolateParams=true")
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
