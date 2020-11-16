package Logic

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var (
	chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
)

func (s *server) getNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]

	rows, err := s.db.Query("SELECT `id`, `color`, `name`, `done`, `text`, `date` FROM `table1` WHERE `token` = ?", token)
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

func (s *server) createNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	item := Note{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, _ = s.db.Exec("INSERT INTO `table1` (`name`, `text`, `date`, `done`, `color`, `token`) VALUES (?, ?, ?, 0, ?, ?)", item.Name, item.Text, item.Date, item.Color, token)
}

func (s *server) deleteNoteById(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]
	id := Note{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, _ = s.db.Exec("DELETE  FROM `table1` WHERE `id` = ? AND `token` = ?", id.Id, token)
}

func (s *server) logIn(w http.ResponseWriter, r *http.Request) {
	user := User{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rows, err := s.db.Query("SELECT token FROM appusers WHERE login = ? AND password =?", user.Login, user.Password)
	if err != nil {
		log.Println(err)
	}
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

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	fmt.Print(1)
	ruser := User{}
	token := ""
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ruser)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Print(5)
	rows, _ := s.db.Query("SELECT `token` FROM appusers WHERE `login` = ? ", ruser.Login)
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
			rows, _ := s.db.Query("SELECT * FROM appusers WHERE `login` = ? AND `token` = ?", ruser.Login, token)
			if rows.Next() == false {
				break
			}
		}
		_, _ = s.db.Exec("INSERT INTO `appusers` (`login`, `password`, `token`) VALUES (?, ?, ?)", ruser.Login, ruser.Password, token)
	}
	login.Token = token
	resp, _ := json.Marshal(&login)
	w.Write(resp)
	count := ""
	rows, _ = s.db.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
	//	fmt.Println(h.websocket)
}
func (s *server) initHandler() {
	adminMux := mux.NewRouter()
	adminMux.HandleFunc("/api/getNote", s.getNote)
	adminMux.HandleFunc("/api/createNote", s.createNote)
	adminMux.HandleFunc("/api/deleteNoteById", s.deleteNoteById)
	adminMux.HandleFunc("/api/logIn", s.logIn)
	adminMux.HandleFunc("/api/register", s.register)
	s.server.Handler = adminMux
}

/*func (s *server) usersOnline(writer http.ResponseWriter, request *http.Request) {
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
