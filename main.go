package main

import (
	"api-go/controllers"
	"api-go/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	cs "api-go/crudsql"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func main() {
	db, err := cs.OpenDatabase("./chinazes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	postService := services.DatabaseInit(db)
	userService := services.DatabaseInitU(db)

	postService.CreatePostsTable()
	userService.CreateUsersTable()

	// Создание контроллера
	postController := controllers.NewPostController(postService)
	userController := controllers.NewUserController(userService)

	// Настройка маршрутизациИ
	router := gin.Default()

	router.POST("/posts", postController.CreatePost)
	router.GET("/posts", postController.GetAllPosts)
	router.GET("/posts/:id", postController.GetPostByID)
	router.PUT("/posts/:id", postController.UpdatePost)
	router.DELETE("/posts/:id", postController.DeletePost)

	// Маршруты для регистрации и логина
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)

	// Маршруты для пользователей с аутентификацией
	auth := router.Group("/")
	auth.Use(controllers.AuthMiddleware())
	auth.DELETE("/users/:id", userController.DeleteUser)

	router.Run(":8080")

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
