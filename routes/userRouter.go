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
	r.PATCH("/user/profile", middleware.AuthMiddleware(roleUser), controllers.EditProfile) //---
	r.POST("/user/profile/changepassword", middleware.AuthMiddleware(roleUser), controllers.ProfileChangePassword)

	//orders
	r.POST("/user/profile/orderscancel/:ID", middleware.AuthMiddleware(roleUser), controllers.CancelOrder)
	//r.POST("/user/profile/ordercancelsingle", middleware.AuthMiddleware(roleUser), controllers.CancelSingleProduct)

	//forgotPassword
	r.POST("/user/profile/forgotpassword", middleware.AuthMiddleware(roleUser), controllers.ProfileForgotPassword)
	r.POST("/user/profile/forgototpcheck", middleware.AuthMiddleware(roleUser), controllers.OtpCheck)
	r.POST("/user/profile/forgotresetpassword", middleware.AuthMiddleware(roleUser), controllers.ResetPassword)

	//cart management
	r.POST("/user/cart", middleware.AuthMiddleware(roleUser), controllers.AddtoCart)
	r.GET("/user/cart", middleware.AuthMiddleware(roleUser), controllers.ListCart)
	r.DELETE("/user/cart/:ID", middleware.AuthMiddleware(roleUser), controllers.RemoveCart)

	//search filter
	r.GET("/user/search", middleware.AuthMiddleware(roleUser), controllers.SearchProduct)

	//checkout page
	r.POST("/user/cartcheckout", middleware.AuthMiddleware(roleUser), controllers.PlaceOrder)

	//vieworder
	r.GET("/user/vieworder", middleware.AuthMiddleware(roleUser), controllers.ViewOrder)
	r.GET("/user/orderdetails", middleware.AuthMiddleware(roleUser), controllers.OrderDetails)

	//Wishlist
	r.POST("/user/addwishlist/:ID", middleware.AuthMiddleware(roleUser), controllers.AddToWishlist)
	r.DELETE("/user/removewishlist/:ID", middleware.AuthMiddleware(roleUser), controllers.RemoveWishlist)
	r.GET("/user/wishlist", middleware.AuthMiddleware(roleUser), controllers.ListWishList)

	//payment
	r.GET("/payment", func(ctx *gin.Context) {
		token := ctx.GetString("token")
		ctx.HTML(200, "Razorpay.html", gin.H{
			"Token": token,
		})
	})
	r.POST("/payment/submit", controllers.CreatePayment)

	//wallet
	r.GET("/user/wallet", middleware.AuthMiddleware(roleUser), controllers.GetWallet)
	r.GET("/user/wallethistory", middleware.AuthMiddleware(roleUser), controllers.WalletHistory)

	//invoice
	r.POST("/user/invoice/:ID", middleware.AuthMiddleware(roleUser), controllers.GenerateInvoice)

	r.POST("/refresh-token", middleware.AuthMiddleware(roleUser), controllers.RefreshToken)

}
