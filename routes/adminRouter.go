package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
	"github.com/mubashir/e-commerce/middleware"
)

var roleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {
	// admin authentication
	r.POST("/admin/login", controllers.AdminLogin)

	//user management
	r.GET("/admin/usermanagement", controllers.ListUsers)
	r.PATCH("/admin/block/:ID", middleware.AuthMiddleware(roleAdmin), controllers.Status)
	r.PATCH("/admin/:ID", middleware.AuthMiddleware(roleAdmin), controllers.UpdateUser)
	r.DELETE("/admin/delete/:ID", middleware.AuthMiddleware(roleAdmin), controllers.DeleteUser)

	// product management

	// category management
	r.POST("/admin/category", middleware.AuthMiddleware(roleAdmin), controllers.CreateCategory)
	r.GET("/admin/category", middleware.AuthMiddleware(roleAdmin), controllers.GetCategory)
	r.PATCH("/admin/category/:ID", middleware.AuthMiddleware(roleAdmin), controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID", middleware.AuthMiddleware(roleAdmin), controllers.DeleteCategory)
}
