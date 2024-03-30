package main

import (
	"github.com/venkata-dra/crud-go/initializers"
	"github.com/venkata-dra/crud-go/models"
)

func init() {
	initializers.LoadEnvInitializers()
	initializers.ConnectToDb()
}
func main() {
	initializers.DB.AutoMigrate(&models.Post{})
	initializers.DB.AutoMigrate(&models.User{})
}
