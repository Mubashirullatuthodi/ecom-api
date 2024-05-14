package controllers

import (
	//"crypto/rand"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authotp "github.com/mubashir/e-commerce/AuthOTP"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/middleware"
	"github.com/mubashir/e-commerce/models"
	"golang.org/x/crypto/bcrypt"
)

var user models.User

var RoleUser = "User"
var Confirmation = false

func Signup(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(406, gin.H{
			"status": "Fail",
			"error":  "json Binding Error",
			"code":   406,
		})
		return
	}

	var existingUser models.User
	result := initializers.DB.Where("email=?", user.Email).First(&existingUser)
	if result.Error == nil {
		ctx.JSON(409, gin.H{
			"status": "Fail",
			"error":  "this user already exists",
			"code":   409,
		})
		return
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash",
		})
		return
	}
	user.Password = string(hashedpassword)

	otp := authotp.GenerateOTP()

	otpRecord := models.OTP{
		Otp:    otp,
		Email:  user.Email,
		Exp:    time.Now().Add(5 * time.Minute),
		UserID: user.ID,
	}
	initializers.DB.Create(&otpRecord)

	errr := authotp.SendEmail(user.Email, otp)

	if errr != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to send OTP via email",
			"code":   400,
		})
		return
	}

	initializers.DB.Create(&user)

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Please check your email and enter the OTP",
	})

	res := initializers.DB.Unscoped().Where("email=?", user.Email).First(&existingUser)
	if res.Error == nil && existingUser.DeletedAt.Valid {
		existingUser.DeletedAt.Time = time.Time{}
		existingUser.DeletedAt.Valid = false
		if err := initializers.DB.Save(&existingUser).Error; err != nil {
			ctx.JSON(500, gin.H{
				"status": "Fail",
				"Error":  "Failed To reactive account",
				"code":   500,
			})
			return
		}
		fmt.Println("helloooiii")
	}
}

func PostOtp(ctx *gin.Context) {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var otp models.OTP
	if err := initializers.DB.Where("otp = ?", input.OTP).First(&otp).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid OTP",
		})
		return
	}

	if time.Now().After(otp.Exp) {
		ctx.JSON(400, gin.H{
			"error": "OTP has expired. Please request a new otp.",
		})
		return
	}
	initializers.DB.Delete(&otp)

	ctx.JSON(200, gin.H{
		"message": "OTP verified Succesfully. User account activated",
	})
}

func ResendOtp(ctx *gin.Context) {
	var userOTP models.OTP
	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existOTP models.OTP

	result := initializers.DB.Where("user_id = ? AND exp < ?", user.ID, time.Now()).First(&existOTP)
	if result.RowsAffected > 0 {
		otp := authotp.GenerateOTP()

		existOTP.Otp = otp
		existOTP.Exp = time.Now().Add(5 * time.Minute)
		initializers.DB.Save(&existOTP)

		err := authotp.SendEmail(user.Email, otp)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to send OTP via Email",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "new OTP sent successfully,please check your email",
		})
		return
	}

}

func UserLogin(ctx *gin.Context) {
	var postinguser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	hashedpassword, _ := bcrypt.GenerateFromPassword([]byte(postinguser.Password), bcrypt.DefaultCost)

	postinguser.Password = string(hashedpassword)

	if err := ctx.ShouldBindJSON(&postinguser); err != nil {
		ctx.JSON(400, gin.H{
			"status": "fail",
			"error":  err.Error(),
			"code":   400,
		})
		return
	}

	result := initializers.DB.Where("email=?", postinguser.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"error":  "Invalid name or password",
			"code":   500,
		})
		return
	}
	password := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(postinguser.Password))
	if password != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid password",
		})
		return
	}

	tokenstring,_:=middleware.JwtToken(ctx, user.ID, postinguser.Email, RoleUser)
	ctx.SetCookie("Authorization"+RoleUser, tokenstring, 3600*24*30, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login Successfully",
	})
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("Authorization"+RoleUser,"",-1,"","",false,true)
	ctx.JSON(200, gin.H{
		"message":   "Successfully Logout",
	})
}

func ForgotPassword(ctx *gin.Context) {
	type input struct {
		Email string `json:"email"`
	}
	var Input input
	if err := ctx.ShouldBindJSON(&Input); err != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"error":  "failed to bind",
			"code":   500,
		})
		return
	}
	result := initializers.DB.Where("email = ?", Input.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"Error":  "failed to check email",
			"code":   500,
		})
		return
	}
	otp := authotp.GenerateOTP()

	otpRecord := models.OTP{
		Otp:    otp,
		Email:  Input.Email,
		Exp:    time.Now().Add(5 * time.Minute),
		UserID: user.ID,
	}
	initializers.DB.Create(&otpRecord)

	errr := authotp.SendEmail(Input.Email, otp)

	if errr != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to send OTP via email",
			"code":   400,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "otp for reset password is sent to your email,reset your password",
	})
}

func OtpCheck(ctx *gin.Context) {
	type OTP struct {
		Otp string `json:"otp"`
	}
	var newOTP OTP
	if err := ctx.ShouldBindJSON(&newOTP); err != nil {
		ctx.JSON(406, gin.H{
			"status": "Fail",
			"error":  "json Binding Error",
			"code":   406,
		})
		return
	}

	var existingOTP models.OTP

	result := initializers.DB.Where("otp = ?", newOTP.Otp).First(&existingOTP)
	if result.Error != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"Error":  "Invalid OTP",
			"code":   500,
		})
		return
	}

	if time.Now().After(existingOTP.Exp) {
		ctx.JSON(400, gin.H{
			"error": "OTP has expired. Please request a new otp.",
		})
		return
	}

	Confirmation = true

	initializers.DB.Delete(&existingOTP)

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Enter new password.",
	})
}

func ResetPassword(ctx *gin.Context) {
	type Input struct {
		Email       string `json:"email"`
		Newpassword string `json:"newpassword"`
	}

	var input Input

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"Error":  "failed to bind",
			"code":   500,
		})
		return
	}

	if !Confirmation {
		ctx.JSON(500, gin.H{
			"status":  "fail",
			"message": "validate the otp",
			"code":    500,
		})
	} else {
		errr := initializers.DB.Where("email=?", input.Email).First(&user)
		if errr.Error != nil {
			ctx.JSON(404, gin.H{
				"status": "fail",
				"Error":  "email not found",
				"code":   404,
			})
			return
		}

		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(input.Newpassword), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(500, gin.H{
				"status": "failed to hash password",
			})
			return
		}

		//user.Password = string(hashedpassword)

		if err := initializers.DB.Model(&user).Update("password", string(hashedpassword)).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "failed to update password"})
			return
		}
		Confirmation = false
		ctx.JSON(200, gin.H{
			"status":  "success",
			"Message": "Password reset successfull",
		})
	}
}
