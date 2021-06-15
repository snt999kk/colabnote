package entities

import (
	"colabnote/internal/database"
	"colabnote/internal/logger"
	"time"
)

type Note struct {
	Id    int
	Color string `json:"color, string"`
	Name  string `json:"title, string"`
	Done  bool   `json:"done, bool"`
	Text  string `json:"text, string"`
	Date  string `json:"date, string"`
}

func NotesList(token string) ([]Note, error) {
	rows, err := database.Database.Query("SELECT id, color, name, done, text, date FROM public.notes WHERE token = $1", token)
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
	date, err := time.Parse("02.01.2006", item.Date)
	if err != nil {
		return err
	}
	_, err = database.Database.Exec("INSERT INTO public.notes (name, text, date, done, color, token) VALUES ($1, $2, $3, false, $4, $5)", item.Name, item.Text, date, item.Color, token)
	if err != nil {
		logger.Info("while getting user in db error happened %v" + err.Error())
		return err
	}
	return nil
}
func DeleteNoteById(token string, id int) error {
	_, err := database.Database.Exec("DELETE FROM public.notes WHERE id = $1 AND token = $2", id, token)
	if err != nil {
		logger.Info("while getting user in db error happened %v" + err.Error())
		return err
	}
	return nil
}
