package models

type Note struct {
	Id    int
	Color string `json:"color,string"`
	Name  string `json:"title,string"`
	Done  int    `json:"done,int"`
	Text  string `json:"text,string"`
	Date  string `json:"date,string"`
}
