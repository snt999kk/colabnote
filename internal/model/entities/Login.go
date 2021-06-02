package entities

import (
	"colabnote/internal/database"
	"colabnote/internal/logger"
	"fmt"
)

type Login struct {
	Exists bool   `json:"exists, bool"`
	Token  string `json:"token, string"`
}

func LogIn(user User) (bool, string, error) {
	rows, err := database.Database.Query("SELECT token FROM appusers WHERE login = ? AND password =?", user.Login, user.Password)
	if err != nil {
		return false, "", err
	}
	rows.Close()
	exists := true
	i := 0
	for rows.Next() {
		i++
		err = rows.Scan(&user.Token)
		if err != nil {
			return false, "", err
		}
		fmt.Println(user.Token)
	}
	if i == 0 {
		exists = false
		logger.Info("No users found with given login and password")
	}
	return exists, user.Token, nil
}
