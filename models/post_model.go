package models

type Post struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_id int    `json:"author_id"`
}
