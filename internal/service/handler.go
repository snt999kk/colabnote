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

var (
	chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
)

func GetNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Token"][0]

	rows, err := database.Database.Query("SELECT `id`, `color`, `name`, `done`, `text`, `date` FROM `table1` WHERE `token` = ?", token)
	if err != nil {
		log.Printf("while making query to db in %v error happened %v", r.URL, err)
		return
	}
	list := []entities.Note{}
	for rows.Next() {
		item := entities.Note{}
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
	_, _ = database.Database.Exec("INSERT INTO `table1` (`name`, `text`, `date`, `done`, `color`, `token`) VALUES (?, ?, ?, 0, ?, ?)", item.Name, item.Text, item.Date, item.Color, token)
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
	_, _ = database.Database.Exec("DELETE  FROM `table1` WHERE `id` = ? AND `token` = ?", id.Id, token)
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
	rows, err := database.Database.Query("SELECT token FROM appusers WHERE login = ? AND password =?", user.Login, user.Password)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	exists := true
	login := entities.Login{}
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
	w.Write(resp)
	count := ""
	rows, _ = database.Database.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
	//	fmt.Println(h.websocket)
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
