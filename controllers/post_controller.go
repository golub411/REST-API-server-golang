package controllers

import (
	"fmt"
	"net/http"
	"reflect"
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
		Body string `json:"body"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID // Предполагаем, что в claims есть UserID

	numberID := strconv.FormatUint(uint64(userID), 10)

	post, err := c.service.CreatePost(req.Body, numberID)
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

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID // Предполагаем, что в claims есть UserID
	userRole := userClaims.Role // Предполагаем, что в claims есть Role

	// Получения id автора поста
	posT, err := c.service.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	postAuthorID := posT[0]["author_id"]

	// Проверка прав доступа
	if userRole != "admin" && postAuthorID != int64(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
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
	idStr := ctx.Param("id") // Это уже строка
	id, err := strconv.Atoi(idStr)
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

	fmt.Println(reflect.TypeOf(postAuthorID))
	fmt.Println(reflect.TypeOf(userID))

	// Проверка прав доступа
	if userRole != "admin" && postAuthorID != int64(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := c.service.DeletePost(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
