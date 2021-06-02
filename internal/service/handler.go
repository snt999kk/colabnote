package service

import (
	"colabnote/internal/database"
	"colabnote/internal/model/entities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func GetNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]

	list, err := entities.NotesList(token)

	resp, err := json.Marshal(list)
	if err != nil {
		log.Printf("While marshaling respose in db b URL:%v error happened %v", r.URL, err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(resp)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	item := entities.Note{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = entities.CreateNote(token, item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func DeleteNoteById(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	id := entities.Note{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = entities.DeleteNoteById(token, id.Id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	user := entities.User{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	login := entities.Login{}
	login.Exists, login.Token, err = entities.LogIn(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp, err := json.Marshal(login)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(resp)
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Print(1)
	ruser := entities.User{}
	token := ""
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ruser)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rows, _ := database.Database.Query("SELECT `token` FROM appusers WHERE `login` = ? ", ruser.Login)
	exists := true
	if rows.Next() == false {
		exists = false
	}
	login := entities.Login{Exists: exists, Token: ""}
	fmt.Println(exists)
	if exists == false {
		for true {
			var b strings.Builder
			n := len(chars)
			for i := 0; i < 20; i++ {
				b.WriteRune(chars[rand.Intn(n)])
			}
			token = b.String()
			rows, _ := database.Database.Query("SELECT * FROM appusers WHERE `login` = ? AND `token` = ?", ruser.Login, token)
			if rows.Next() == false {
				break
			}
		}
		_, _ = database.Database.Exec("INSERT INTO `appusers` (`login`, `password`, `token`) VALUES (?, ?, ?)", ruser.Login, ruser.Password, token)
	}
	login.Token = token
	resp, _ := json.Marshal(&login)
	count := ""
	rows, _ = database.Database.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
	w.Write(resp)
}

/*func (s *Server) usersOnline(writer http.ResponseWriter, request *http.Request) {
	conn, _ := upgrader.Upgrade(writer, request, nil)
	s.websocket = conn
	count := ""
	rows, _ := s.db.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
	s.websocket.WriteMessage(1, []byte("users online: "+count))
	//	h.websocket.Close()
	fmt.Println('a')
}*/
