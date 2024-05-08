package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/User"
)

//var roleUser = "User"

func UserGroup(r *gin.RouterGroup) {
	//===================== user authentication ====================
	r.POST("/user/signup", controllers.Signup)
	r.POST("/user/signup/verify-otp", controllers.PostOtp)
	r.POST("/user/signup/resend-otp", controllers.ResendOtp)
	r.POST("/user/login", controllers.UserLogin)

	//product page
	r.GET("/user/product", controllers.ProductPage)
	r.GET("/user/product/:ID", controllers.ProductDetail)

}
