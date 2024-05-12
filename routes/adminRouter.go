package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
	"github.com/mubashir/e-commerce/middleware"
	//"github.com/mubashir/e-commerce/middleware"
)

var RoleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {
	// admin authentication
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/admin/signup", controllers.AdminSignUp)
	r.GET("/admin/logout", middleware.AuthMiddleware(RoleAdmin), controllers.AdminLogout)

	//user management
	r.GET("/admin/usermanagement", middleware.AuthMiddleware(RoleAdmin), controllers.ListUsers)
	r.PATCH("/admin/block/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.Status)
	r.PATCH("/admin/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.UpdateUser)
	r.DELETE("/admin/delete/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.DeleteUser)

	// product management
	r.POST("/admin/product", controllers.AddProduct)
	r.GET("/admin/product", controllers.ListProducts)
	r.PATCH("/admin/product/:ID", controllers.EditProduct)
	r.PATCH("/admin/product/image/:ID", controllers.ImageUpdate)
	r.DELETE("/admin/product/:ID", controllers.DeleteProduct)

	// category management
	r.POST("/admin/category", controllers.CreateCategory)
	r.GET("/admin/category", controllers.GetCategory)
	r.PATCH("/admin/category/:ID", controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID", controllers.DeleteCategory)
}
