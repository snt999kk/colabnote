package entities

import (
	"colabnote/internal/database"
	"encoding/json"
	"fmt"
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

func Register(user User, token string) error {
	rows, err := database.Database.Query("SELECT `token` FROM appusers WHERE `login` = ? ", Username)
	if err != nil {
		return err
	}
	exists := true
	if rows.Next() == false {
		exists = false
	}
	login := Login{Exists: exists, Token: ""}
	fmt.Println(exists)
	if exists == false {
		for true {
			var b strings.Builder
			n := len(Chars)
			for i := 0; i < 20; i++ {
				b.WriteRune(Chars[rand.Intn(n)])
			}
			token = b.String()
			rows, _ := database.Database.Query("SELECT * FROM appusers WHERE `login` = ? AND `token` = ?", user.Login, token)
			if rows.Next() == false {
				break
			}
		}
		_, _ = database.Database.Exec("INSERT INTO `appusers` (`login`, `password`, `token`) VALUES (?, ?, ?)", user.Password, user.Password, token)
	}
	login.Token = token
	resp, _ := json.Marshal(&login)
	count := ""
	rows, _ = database.Database.Query("SELECT COUNT(*) FROM mysql.appusers")
	rows.Next()
	rows.Scan(&count)
}
