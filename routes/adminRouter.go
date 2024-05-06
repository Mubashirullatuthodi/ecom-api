package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
)

var roleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {
	//======================= admin authentication ==========================
	r.POST("/admin/login", controllers.AdminLogin)

	//user management
	r.GET("/admin/usermanagement", controllers.ListUsers)
	r.PATCH("/admin/block", controllers.Status)
	r.PATCH("/admin/:ID", controllers.UpdateUser)
	r.DELETE("/admin/delete/:ID", controllers.DeleteUser)

	//product management

	//category management
	r.POST("/admin/category", controllers.CreateCategory)
	r.GET("/admin/category", controllers.GetCategory)
	r.PATCH("/admin/category/:ID", controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID", controllers.DeleteCategory)
}
