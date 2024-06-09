package controllers

import (
	"net/http"
	"strconv"

	"api-go/services"

	"api-go/utils"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
	return &PostController{service: service}
}

// CREATE POST

func (c *PostController) CreatePost(ctx *gin.Context) {
	var req struct {
		Body      string `json:"body"`
		Author_id string `json:"author_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := c.service.CreatePost(req.Body, req.Author_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}

// GET ALL POSTS

func (c *PostController) GetAllPosts(ctx *gin.Context) {
	posts, err := c.service.GetAllPosts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

// GET POST BY ID

func (c *PostController) GetPostByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := c.service.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// UPDATE POST

func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := c.service.UpdatePost(id, req.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

//

func (c *PostController) DeletePost(ctx *gin.Context) {
	id, err := strconv.Atoi("123") //ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID // Предполагаем, что в claims есть UserID
	userRole := userClaims.Role // Предполагаем, что в claims есть Role

	// Получения id автора поста
	post, err := c.service.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	postAuthorID := post[0]["author_id"]

	// Проверка прав доступа
	if userRole != "admin" && postAuthorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := c.service.DeletePost(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
