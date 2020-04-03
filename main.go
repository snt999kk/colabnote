package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"math/rand"
	"net/http"
	_ "net/smtp"
	"strings"
	_ "strings"
	"time"
)

type Note struct {
	Id    int
	Color string `json:"color, string"`
	Name  string `json:"title, string"`
	Done  int    `json:"done, int"`
	Text  string `json:"text, string"`
	Date  string `json:"date, string"`
}
type User struct {
	Login    string `json:"login, string"`
	Password string `json:"password, string"`
	token    string
}

/**
Структура для логина
*/
type Login struct {
	Exists bool   `json:"exists, bool"`
	Token  string `json:"token, string"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	db, err := sql.Open("mysql", "reg-user:qwerty123123@tcp(localhost:3306)/mysql?&charset=utf8&interpolateParams=true")
	fmt.Println("%v", err)
	//panic(err)
	db.SetMaxOpenConns(2)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println(db.Stats().OpenConnections)
	//fmt.Println("%v", list[0])
	http.HandleFunc("/api/getNote", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Token"][0]
		//fmt.Println(table, "*")
		rows, err := db.Query("SELECT `id`, `color`, `name`, `done`, `text`, `date` FROM `table1` WHERE `token` = ?", token)
		if err != nil {
			panic(err)
		}
		list := []Note{}
		for rows.Next() {
			item := Note{}
			_ = rows.Scan(&item.Id, &item.Color, &item.Name, &item.Done, &item.Text, &item.Date)
			list = append(list, item)
			fmt.Println("%v", item)
		}
		resp, err := json.Marshal(list)
		//	fmt.Println(resp, "&")
		if err != nil {
			panic(err)
		}
		w.Write(resp)
		rows.Close()
	})
	http.HandleFunc("/api/createNote", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Token"][0]
		item := Note{}
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		err = json.Unmarshal(body, &item)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, _ = db.Exec("INSERT INTO `table1` (`name`, `text`, `date`, `done`, `color`, `token`) VALUES (?, ?, ?, 0, ?, ?)", item.Name, item.Text, item.Date, item.Color, token)
	})
	http.HandleFunc("/api/deleteNoteById", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Token"][0]
		id := Note{}
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		err = json.Unmarshal(body, &id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, _ = db.Exec("DELETE  FROM `table1` WHERE `id` = ? AND `token` = ?", id.Id, token)
	})
	http.HandleFunc("/api/logIn", func(w http.ResponseWriter, r *http.Request) {
		user := User{}
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		err = json.Unmarshal(body, &user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		rows, _ := db.Query("SELECT `token` FROM appusers WHERE `login` = ? AND `password` =?", user.Login, user.Password)
		defer rows.Close()
		exists := true
		login := Login{}
		i := 0
		for rows.Next() {
			i++
			_ = rows.Scan(&user.token)
			fmt.Println(user.token)
		}
		if i == 0 {
			exists = false
			fmt.Println("No users found with given login and password")
		}
		login.Exists = exists
		login.Token = user.token
		resp, _ := json.Marshal(login)
		w.Write(resp)
	})
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		ruser := User{}
		token := ""
		body, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &ruser)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		rows, _ := db.Query("SELECT `token` FROM appusers WHERE `login` = ? ", ruser.Login)
		exists := true
		if rows.Next() == false {
			exists = false
		}
		login := Login{Exists: exists, Token: ""}
		if exists == false {
			for true {
				var b strings.Builder
				n := len(chars)
				for i := 0; i < 20; i++ {
					b.WriteRune(chars[rand.Intn(n)])
				}
				token = b.String()
				rows, _ := db.Query("SELECT * FROM appusers WHERE `login` = ? AND `token` = ?", ruser.Login, token)
				if rows.Next() == false {
					break
				}
			}
			_, _ = db.Exec("INSERT INTO `appusers` (`login`, `password`, `token`) VALUES (?, ?, ?)", ruser.Login, ruser.Password, token)
		}
		login.Token = token
		resp, _ := json.Marshal(&login)
		w.Write(resp)
	})
	http.HandleFunc("/api/usersOnline", func(writer http.ResponseWriter, request *http.Request) {
		conn, _ := upgrader.Upgrade(writer, request, nil)
		for true {
			conn.WriteMessage(2, []byte("25"))
			time.Sleep(5 * time.Second)
		}
		fmt.Println('a')
	})
	http.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("dwd"))
	})
	http.ListenAndServe(":8080", nil)
}
