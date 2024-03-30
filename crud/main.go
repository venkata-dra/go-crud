package main

import (
	"github.com/gin-gonic/gin"
	"github.com/venkata-dra/crud-go/controllers"
	"github.com/venkata-dra/crud-go/initializers"
	"github.com/venkata-dra/crud-go/middleware"
)

func init() {
	initializers.LoadEnvInitializers()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.GET("/validate", middleware.RequiredAuth, controllers.Validate)
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.SignUp)
	r.POST("/post", controllers.PostsCreate)
	r.GET("/getPosts", middleware.RequiredAuth, controllers.PostIndex)
	r.GET("/getPosts/:id", controllers.PostsShow)
	r.PUT("/posts/:id", controllers.PostsUpdate)
	r.DELETE("/delete/:id", controllers.PostDelete)
	r.Run() // listen and serve on 0.0.0.0:8080
}
