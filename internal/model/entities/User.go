package entities

type User struct {
	Login    string `json:"login,string"`
	Password string `json:"password,string"`
	Token    string
}
