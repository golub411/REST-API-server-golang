# ⭐REST-API-server-golang

🔥I made this api like a foundation for my and yours projects. This api use my own library CRUD-sql, gin and packege http/net for servers. Api corresponds meets all architectural requirements. 

```
  func main() {
  	db, err := cs.OpenDatabase("./chinazes.db")
  	if err != nil {
  		log.Fatal(err)
  	}
  	defer db.Close()
  
  	postService := services.DatabaseInit(db)
  
  	postService.CreatePostsTable()
  
  	// Создание контроллера
  	postController := controllers.NewPostController(postService)
  
  	// Настройка маршрутизациИ
  	router := gin.Default()
  
  	router.POST("/posts", postController.CreatePost)
  	router.GET("/posts", postController.GetAllPosts)
  	router.GET("/posts/:id", postController.GetPostByID)
  	router.PUT("/posts/:id", postController.UpdatePost)
  	router.DELETE("/posts/:id", postController.DeletePost)
  
  	router.Run(":8080")
  
  	fmt.Println("Server is running at http://localhost:8080")
  	log.Fatal(http.ListenAndServe(":8080", nil))
}

```
