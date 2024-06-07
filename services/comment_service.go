package services

import (
	"api-go/crudsql"
	"api-go/models"
	"fmt"
	"strconv"
)

type CommentService struct {
	db *crudsql.Database
}

func DatabaseInitC(db *crudsql.Database) *CommentService {
	return &CommentService{db: db}
}

func (p *CommentService) CreateCommentsTable() error {
	err := p.db.CreateTable("comments", []string{"id INTEGER PRIMARY KEY AUTOINCREMENT", "body TEXT", "author_id INTEGER", "post_id INTEGER"})
	return err
}

func (p *CommentService) CreateComment(body string, author_id string, post_id string) (*models.Comment, error) {
	author_id_num, err := strconv.Atoi(author_id)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать author_id в число: %w", err)
	}

	post_id_num, err := strconv.Atoi(post_id)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать post_id в число: %w", err)
	}

	comment := &models.Comment{Body: body, Author_id: author_id_num, Post_id: post_id_num}

	err = p.db.InsertValue("comments", []string{"body", "author_id", "post_id"}, []interface{}{comment.Body, comment.Author_id, comment.Post_id})
	if err != nil {
		return nil, fmt.Errorf("не удалось вставить комментарий в базу данных: %w", err)
	}

	return comment, nil
}

func (p *CommentService) GetCommentsByPostId(postID string) ([]map[string]interface{}, error) {
	// Преобразование идентификатора поста в число.
	postIDNum, err := strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("неверный идентификатор поста: %w", err)
	}

	// Выборка комментариев из базы данных.
	comments, err := p.db.SelectValueWhere("comments", []string{"body", "author_id"}, fmt.Sprintf("post_id = %d", postIDNum))
	if err != nil {
		return nil, fmt.Errorf("не удалось достать комментарии из базы: %w", err)
	}
	return comments, nil
}

func (p *CommentService) DeleteCommentById(id string) error {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("неверный идентификатор комментария: %w", err)
	}

	err = p.db.DeleteValue("comments", map[string]interface{}{"id": idNum})
	return err
}
