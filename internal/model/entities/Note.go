package entities

import (
	"colabnote/internal/database"
	"colabnote/internal/logger"
)

type Note struct {
	Id    int
	Color string `json:"color, string"`
	Name  string `json:"title, string"`
	Done  int    `json:"done, int"`
	Text  string `json:"text, string"`
	Date  string `json:"date, string"`
}

func NotesList(token string) ([]Note, error) {
	rows, err := database.Database.Query("SELECT id, color, name, done, text, date FROM table1 WHERE token = ?", token)
	if err != nil {
		logger.Info("while getting user in db error happened %v" + err.Error())
		return nil, err
	}
	list := []Note{}
	for rows.Next() {
		item := Note{}
		err = rows.Scan(&item.Id, &item.Color, &item.Name, &item.Done, &item.Text, &item.Date)
		if err != nil {
			logger.Info("while getting user in db error happened %v" + err.Error())
			return nil, err
		}
		list = append(list, item)
		//	fmt.Println("%v", item)
	}
	rows.Close()
	return list, err
}
func CreateNote(token string, item Note) error {
	_, err := database.Database.Exec("INSERT INTO table1 (name, text, date, done, color, token) VALUES (?, ?, ?, 0, ?, ?)", item.Name, item.Text, item.Date, item.Color, token)
	if err != nil {
		logger.Info("while getting user in db error happened %v" + err.Error())
		return err
	}
	return nil
}
func DeleteNoteById(token string, id int) error {
	_, err := database.Database.Exec("DELETE  FROM `table1` WHERE `id` = ? AND `token` = ?", id, token)
	if err != nil {
		logger.Info("while getting user in db error happened %v" + err.Error())
		return err
	}
	return nil
}
