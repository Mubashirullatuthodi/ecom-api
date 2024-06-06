package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	user := r.Group("/")
	routes.UserGroup(user)

	admin := r.Group("/")
	routes.AdminGroup(admin)

	r.Run()
}
