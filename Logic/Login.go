package Logic

type Login struct {
	Exists bool   `json:"exists, bool"`
	Token  string `json:"token, string"`
}