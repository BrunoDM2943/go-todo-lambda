package model

type Item struct {
	ID    int64
	Title string `json:"title"`
	Text  string `json:"text"`
}
