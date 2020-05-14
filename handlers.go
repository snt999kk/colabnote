package main

import (
	. "ToDoList/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type handlers struct {
	db        *sql.DB
	websocket *websocket.Conn
}

func (h *handlers) getNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	//fmt.Println(table, "*")

	rows, err := h.db.Query("SELECT `id`, `color`, `name`, `done`, `text`, `date` FROM `table1` WHERE `token` = ?", token)
	if err != nil {
		log.Printf("while making query to db in %v error happened %v", r.URL, err)
		return
	}
	list := []Note{}
	for rows.Next() {
		item := Note{}
		err = rows.Scan(&item.Id, &item.Color, &item.Name, &item.Done, &item.Text, &item.Date)
		if err != nil {
			log.Printf("while getting user in db by URL:%v error happened %v", r.URL, err)
			return
		}
		list = append(list, item)
		//	fmt.Println("%v", item)
	}

	resp, err := json.Marshal(list)
	if err != nil {
		log.Printf("While marshaling respose in db b URL:%v error happened %v", r.URL, err)
		return
	}
	w.Write(resp)
	rows.Close()
}

func (h *handlers) createNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	item := Note{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, _ = h.db.Exec("INSERT INTO `table1` (`name`, `text`, `date`, `done`, `color`, `token`) VALUES (?, ?, ?, 0, ?, ?)", item.Name, item.Text, item.Date, item.Color, token)
}

func (h *handlers) deleteNoteById(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	id := Note{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, _ = h.db.Exec("DELETE  FROM `table1` WHERE `id` = ? AND `token` = ?", id.Id, token)
}

func (h *handlers) logIn(w http.ResponseWriter, r *http.Request) {
	user := User{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rows, _ := h.db.Query("SELECT `token` FROM appusers WHERE `login` = ? AND `password` =?", user.Login, user.Password)
	defer rows.Close()
	exists := true
	login := Login{}
	i := 0
	for rows.Next() {
		i++
		_ = rows.Scan(&user.Token)
		fmt.Println(user.Token)
	}
	if i == 0 {
		exists = false
		fmt.Println("No users found with given login and password")
	}
	login.Exists = exists
	login.Token = user.Token
	resp, _ := json.Marshal(login)
	w.Write(resp)
}

func (h *handlers) register(w http.ResponseWriter, r *http.Request) {
	ruser := User{}
	token := ""
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ruser)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rows, _ := h.db.Query("SELECT `token` FROM appusers WHERE `login` = ? ", ruser.Login)
	exists := true
	if rows.Next() == false {
		exists = false
	}
	login := Login{Exists: exists, Token: ""}
	fmt.Println(exists)
	if exists == false {
		for true {
			var b strings.Builder
			n := len(chars)
			for i := 0; i < 20; i++ {
				b.WriteRune(chars[rand.Intn(n)])
			}
			token = b.String()
			rows, _ := h.db.Query("SELECT * FROM appusers WHERE `login` = ? AND `token` = ?", ruser.Login, token)
			if rows.Next() == false {
				break
			}
		}
		_, _ = h.db.Exec("INSERT INTO `appusers` (`login`, `password`, `token`) VALUES (?, ?, ?)", ruser.Login, ruser.Password, token)
	}
	login.Token = token
	resp, _ := json.Marshal(&login)
	w.Write(resp)
	count := ""
	rows, _ = h.db.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
	//	fmt.Println(h.websocket)
	h.websocket.WriteMessage(1, []byte("users online "+count))
}

func (h *handlers) usersOnline(writer http.ResponseWriter, request *http.Request) {
	conn, _ := upgrader.Upgrade(writer, request, nil)
	h.websocket = conn
	count := ""
	rows, _ := h.db.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
	h.websocket.WriteMessage(1, []byte("users online: "+count))
	//	h.websocket.Close()
	fmt.Println('a')
}
