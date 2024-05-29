package services

import (
	"api-go/crudsql"
	"api-go/models"
	"fmt"
)

type PostService struct {
	db *crudsql.Database
}

func DatabaseInit(db *crudsql.Database) *PostService {
	return &PostService{db: db}
}

func (p *PostService) CreatePostsTable() error {
	err := p.db.CreateTable("posts", []string{"post_id INTEGER PRIMARY KEY AUTOINCREMENT", "body TEXT"})
	return err
}

func (p *PostService) CreatePost(body string) (*models.Post, error) {
	post := &models.Post{Body: body}
	// Убедитесь, что post.Body имеет правильный тип, который должен быть помещен в срез интерфейсов
	err := p.db.InsertValue("posts", []string{"body"}, []interface{}{post.Body})
	return post, err
}

func (p *PostService) GetAllPosts() ([]map[string]interface{}, error) {
	// Assuming SelectValue returns ([]models.Post, error)
	posts, err := p.db.SelectValue("posts", []string{"post_id", "body"})
	if err != nil {
		// Handle the error appropriately, e.g., return it or log it
		return nil, err
	}
	return posts, nil
}

func (p *PostService) GetById(id int) ([]map[string]interface{}, error) {
	return p.db.SelectValueWhere("posts", []string{"post_id", "body"}, fmt.Sprintf("post_id = %d", id))
}

func (p *PostService) UpdatePost(id int, body string) (map[string]interface{}, error) {
	set := map[string]interface{}{"body": body}
	where := map[string]interface{}{"id": id}
	err := p.db.UpdateValue("your_table", set, where)
	if err != nil {
		// Handle the error
		return nil, err
	}
	// If there's no error, return a success message
	return map[string]interface{}{"message": "Post updated successfully"}, nil
}

func (p *PostService) DeletePost(id int) error {
	err := p.db.DeleteValue("posts", map[string]interface{}{"id": id})
	return err
}
