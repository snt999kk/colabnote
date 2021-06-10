package entities

import (
	"colabnote/internal/database"
	"math/rand"
	"strings"
)

var (
	Chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
)

type User struct {
	Login    string `json:"login,string"`
	Password string `json:"password,string"`
	Token    string
}

func RegisterUser(user User) (User, error) {
	rows, err := database.Database.Query("SELECT token FROM appusers WHERE login = ? ", user.Login)
	if err != nil {
		return User{}, err
	}
	exists := true
	if rows.Next() == false {
		exists = false
	}
	token := ""
	for exists {
		var b strings.Builder
		n := len(Chars)
		for i := 0; i < 20; i++ {
			b.WriteRune(Chars[rand.Intn(n)])
		}
		token = b.String()
		rows, _ := database.Database.Query("SELECT * FROM appusers WHERE login = ? AND token = ?", user.Login, token)
		if rows.Next() == false {
			exists = false
		}
	}
	_, _ = database.Database.Exec("INSERT INTO appusers (login, password, token) VALUES (?, ?, ?)", user.Login, user.Password, token)
	user.Token = token
	return user, err
}
