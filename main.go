package main

import (
	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/controllers"
	"github.com/mubashir/e-commerce/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()

	// store := cookie.NewStore([]byte("secret"))
	// r.Use(sessions.Sessions("mysession",store))

	//r.GET("/", controllers.GetHome)

	//Admin
	r.POST("/admin/login", controllers.AdminLogin)
	r.GET("/admin/usermanagement", controllers.ListUsers)
	r.PATCH("/admin/block/:ID", controllers.Status)
	r.PATCH("/admin/:ID",controllers.UpdateUser)
	r.DELETE("/admin/delete/:ID",controllers.DeleteUser)
	r.POST("/admin/category",controllers.CreateCategory)
	r.GET("/admin/category",controllers.GetCategory)
	r.PATCH("/admin/category/:ID",controllers.UpdateCategory)
	r.DELETE("/admin/category/:ID",controllers.DeleteCategory)
	// r.POST("/admin/products",controllers.CreateProduct)
	
	

	//user
	r.POST("/user/signup", controllers.Signup)
	r.POST("/user/signup/verify-otp", controllers.PostOtp)
	r.POST("/user/signup/resend-otp", controllers.ResendOtp)
	r.POST("/user/login", controllers.UserLogin)
	//r.GET("/users/productview",)
	//r.GET("/users/searchproduct",)

	r.Run()

}
