package entities

import (
	"colabnote/internal/database"
	"fmt"
	"math/rand"
	"strings"
)

var (
	Chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string
}

func RegisterUser(user User) (User, error) {
	rows, err := database.Database.Query("SELECT token FROM public.appusers WHERE login = $1", user.Login)
	if err != nil {
		return User{}, err
	}
	if rows.Next() != false {
		return User{}, fmt.Errorf("User with given login already exists")
	}
	exists := true
	token := ""
	for exists {
		var b strings.Builder
		n := len(Chars)
		for i := 0; i < 20; i++ {
			b.WriteRune(Chars[rand.Intn(n)])
		}
		token = b.String()
		rows, err := database.Database.Query("SELECT * FROM public.appusers WHERE login = $1 AND token = $2", user.Login, token)
		if err != nil {
			return User{}, err
		}
		if rows.Next() == false {
			exists = false
		}
	}
	_, err = database.Database.Exec("INSERT INTO public.appusers (login, password, token) VALUES ($1, $2, $3)", user.Login, user.Password, token)
	if err != nil {
		return User{}, err
	}
	user.Token = token
	return user, err
}
