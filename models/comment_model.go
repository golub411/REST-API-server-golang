package models

type Comment struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_id int    `json:"author_id"`
	Post_id   int    `json:"post_id"`
}
