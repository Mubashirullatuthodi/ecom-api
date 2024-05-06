package controllers

import (
	//"crypto/rand"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"golang.org/x/crypto/bcrypt"
)

var user models.User

func Signup(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existingUser models.User
	result := initializers.DB.Where("email=?", user.Email).First(&existingUser)
	if result.Error == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "this user already exists",
		})
		return
	}
	// if err := initializers.DB.Model(&user).Update("deleted_at",nil).Error;err!=nil{
	// 	ctx.JSON(http.StatusInternalServerError,gin.H{
	// 		"error":"failed to recover",
	// 	})
	// 	return
	// }
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash",
		})
		return
	}
	user.Password = string(hashedpassword)

	otp := generateOTP()

	otpRecord := models.OTP{
		Otp:    otp,
		Exp:    time.Now().Add(5 * time.Minute),
		UserID: user.UserID,
	}
	initializers.DB.Create(&otpRecord)

	errr := sendEmail(user.Email, otp)

	if errr != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to send OTP via email",
		})
		return
	}

	initializers.DB.Create(&user)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Please check your email and enter the OTP",
	})
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", otp)

	//otp:=int(b[0])<<24|int(b[1])<<16|int(b[2])<<8|int(b[3])

	// otp= otp%1000000
	// if otp<0{
	// 	otp=-otp
	// }
	// return otp
}

func sendEmail(email string, otp string) error {

	from := "mubashirullatuthodi@gmail.com"
	password := "jhzv jkkb ewgr sbrh"
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("Subject: Your OTP for Sign Up\n\n OTP is: " + otp)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
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

func UserLogin(ctx *gin.Context) {
	var postinguser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	hashedpassword,_:=bcrypt.GenerateFromPassword([]byte(postinguser.Password),bcrypt.DefaultCost)

	postinguser.Password = string(hashedpassword)

	if err := ctx.ShouldBindJSON(&postinguser); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}


	//var existingUser models.User
	

	result := initializers.DB.Where("email=?", postinguser.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid name or password",
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login Successfully",
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

	result := initializers.DB.Where("user_id = ? AND exp < ?", user.UserID, time.Now()).First(&existOTP)
	if result.RowsAffected > 0 {
		otp := generateOTP()

		existOTP.Otp = otp
		existOTP.Exp = time.Now().Add(5 * time.Minute)
		initializers.DB.Save(&existOTP)

		err := sendEmail(user.Email, otp)
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
