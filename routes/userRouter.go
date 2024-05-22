package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/User"
	"github.com/mubashir/e-commerce/middleware"
)

var roleUser = "User"

func UserGroup(r *gin.RouterGroup) {
	//user authentication
	r.POST("/user/signup", controllers.Signup)
	r.POST("/user/signup/verify-otp", controllers.PostOtp)
	r.POST("/user/signup/resend-otp", controllers.ResendOtp)
	r.POST("/user/login", controllers.UserLogin)
	r.GET("/user/logout", controllers.Logout)
	r.POST("/user/forgotpassword", controllers.ForgotPassword)
	r.POST("/user/forgototpcheck", controllers.OtpCheck)
	r.POST("/user/resetpassword", controllers.ResetPassword)

	//product page
	r.GET("/user/product", controllers.ProductPage)
	r.GET("/user/product/:ID", middleware.AuthMiddleware(roleUser), controllers.ProductDetail)

	//User Address
	r.POST("/user/address", middleware.AuthMiddleware(roleUser), controllers.AddAddress)
	r.PATCH("/user/address/:ID", middleware.AuthMiddleware(roleUser), controllers.EditAddress)
	r.DELETE("/user/address/:ID", middleware.AuthMiddleware(roleUser), controllers.DeleteAddress)

	//User Profile
	r.GET("/user/profile/address", middleware.AuthMiddleware(roleUser), controllers.ListAddress)

	//cart management
	r.POST("/user/cart", middleware.AuthMiddleware(roleUser), controllers.AddtoCart)
	r.GET("/user/cart", middleware.AuthMiddleware(roleUser), controllers.ListCart)
	//r.POST("/user/cart/reducing/:ID", middleware.AuthMiddleware(roleUser), controllers.ReducingQuantity)
	r.DELETE("/user/cart/:ID", middleware.AuthMiddleware(roleUser), controllers.RemoveCart)

	//search filter
	r.GET("/user/search", middleware.AuthMiddleware(roleUser), controllers.SearchProduct)

	//checkout page
	r.GET("/user/cartcheckout", middleware.AuthMiddleware(roleUser), controllers.CheckoutCart)

}
