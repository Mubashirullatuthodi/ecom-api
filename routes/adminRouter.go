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
	r.GET("/admin/listusers", middleware.AuthMiddleware(RoleAdmin), controllers.ListUsers)
	r.PATCH("/admin/block/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.Status)
	r.PATCH("/admin/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.UpdateUser)
	r.DELETE("/admin/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.DeleteUser)

	// product management
	r.POST("/admin/product", middleware.AuthMiddleware(RoleAdmin), controllers.AddProduct)
	r.GET("/admin/product", middleware.AuthMiddleware(RoleAdmin), controllers.ListProducts)
	r.PATCH("/admin/product/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.EditProduct)
	r.PATCH("/admin/product/image/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.ImageUpdate)
	r.DELETE("/admin/product/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.DeleteProduct)

	// category management
	r.POST("/admin/category", middleware.AuthMiddleware(RoleAdmin), controllers.CreateCategory)
	r.GET("/admin/category", middleware.AuthMiddleware(RoleAdmin), controllers.GetCategory)
	r.PATCH("/admin/category/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID", middleware.AuthMiddleware(RoleAdmin), controllers.DeleteCategory)

	//order management
	r.GET("/admin/orderdetails",middleware.AuthMiddleware(RoleAdmin),controllers.GetOrderDetails)
	r.POST("/admin/order/:ID",middleware.AuthMiddleware(RoleAdmin),controllers.ChangeOrderStatus)
	r.POST("/admin/cancelorder",middleware.AuthMiddleware(RoleAdmin),controllers.CancelOrder)

	//Coupon 
	r.GET("/admin/coupons",middleware.AuthMiddleware(RoleAdmin),controllers.ListCoupon)
	r.POST("/admin/coupons",middleware.AuthMiddleware(RoleAdmin),controllers.CreateCoupon)
	r.DELETE("/admin/coupons",middleware.AuthMiddleware(RoleAdmin),controllers.DeleteCoupon)

}
