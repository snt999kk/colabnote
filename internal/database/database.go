package database

import (
	"colabnote/internal/config"
	"database/sql"
	"fmt"
)

var Database *sql.DB

func InitDB(conf config.Config) error {
	DBsettings := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", conf.Host, conf.DBport, conf.User, conf.Password, conf.DBname)
	data, err := sql.Open("postgres", DBsettings)
	Database = data
	return err
}
