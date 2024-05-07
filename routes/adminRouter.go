package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
	"github.com/mubashir/e-commerce/middleware"
)

var RoleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {
	// admin authentication
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/admin/signup", controllers.AdminSignUp)

	//user management
	r.GET("/admin/usermanagement", controllers.ListUsers)
	r.PATCH("/admin/block/:ID", controllers.Status)
	r.PATCH("/admin/:ID", controllers.UpdateUser)
	r.DELETE("/admin/delete/:ID", controllers.DeleteUser)

	// product management

	// category management
	r.POST("/admin/category", middleware.AuthMiddleware(RoleAdmin), controllers.CreateCategory)
	r.GET("/admin/category", middleware.AuthMiddleware(RoleAdmin), controllers.GetCategory)
	r.PATCH("/admin/category/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.DeleteCategory)
}
