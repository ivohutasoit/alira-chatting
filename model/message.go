package model

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Text     string `json:"text"`
}
