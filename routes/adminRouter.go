package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
	//"github.com/mubashir/e-commerce/middleware"
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
	r.POST("/admin/product",controllers.AddProduct)
	r.GET("/admin/product",controllers.ListProducts)

	
	// category management
	r.POST("/admin/category", controllers.CreateCategory)
	r.GET("/admin/category", controllers.GetCategory)
	r.PATCH("/admin/category/:ID", controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID", controllers.DeleteCategory)
}
