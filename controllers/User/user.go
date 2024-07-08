package controllers

import (
	//"crypto/rand"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authotp "github.com/mubashir/e-commerce/AuthOTP"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/middleware"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
	"golang.org/x/crypto/bcrypt"
)

var user models.User

var RoleUser = "User"
var Confirmation = false

var OTPverification = false

type NewUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Phone     string `json:"phoneNo"`
	Password  string `json:"password"`
	Status    string `json:"status"`
}

var newUserInstance NewUser

func Signup(ctx *gin.Context) {

	OTPverification = false
	if err := ctx.ShouldBindJSON(&newUserInstance); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Please ensure that all required fields are correctly filled out and try again")
		return
	}

	errors := utils.ValidateUserInstance(utils.NewUser(newUserInstance))
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"validation errors": errors})
		return
	}
	newUserInstance.FirstName = utils.CapitalizeFirstLetter(newUserInstance.FirstName)

	fmt.Println("user: ", newUserInstance)
	var existingUser models.User
	result := initializers.DB.Where("email=?", newUserInstance.Email).First(&existingUser)
	if result.Error == nil {
		utils.HandleError(ctx, http.StatusConflict, "This user already Exist")
		return
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(newUserInstance.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to Hash password")
		return
	}
	newUserInstance.Password = string(hashedpassword)

	otp := authotp.GenerateOTP()

	otpRecord := models.OTP{
		Otp:    otp,
		Email:  newUserInstance.Email,
		Exp:    time.Now().Add(5 * time.Minute),
		UserID: user.ID,
	}
	initializers.DB.Create(&otpRecord)

	if err := authotp.SendEmail(newUserInstance.Email, otp); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to send OTP via email")
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Please check your email and enter the OTP",
	})
}

func PostOtp(ctx *gin.Context) {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Failed to bind")
		return
	}

	var otp models.OTP
	if err := initializers.DB.Where("otp = ?", input.OTP).First(&otp).Error; err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid OTP")
		return
	}

	if time.Now().After(otp.Exp) {
		utils.HandleError(ctx, http.StatusBadRequest, "OTP has expired. Please request a new otp.")
		return
	}
	//var existingUser models.User

	if input.OTP == otp.Otp {
		OTPverification = true //if the otp success it will become true and create user
	}

	if OTPverification {

		usernew := models.User{
			FirstName: newUserInstance.FirstName,
			LastName:  newUserInstance.LastName,
			Email:     newUserInstance.Email,
			Gender:    newUserInstance.Gender,
			Phone:     newUserInstance.Phone,
			Password:  newUserInstance.Password,
			Status:    newUserInstance.Status,
		}

		initializers.DB.Create(&usernew)

		initializers.DB.Delete(&otp)
		ctx.JSON(201, gin.H{
			"message": "OTP verified Succesfully. User account created",
		})
		if existUser, err := ReactivateUser(); err == nil {
			saveUser(existUser)
			log.Println("Account reactivated !")
		}
		// res := initializers.DB.Unscoped().Where("email=?", usernew.Email).First(&existingUser)
		// if res.Error == nil && existingUser.DeletedAt.Valid {
		// 	existingUser.FirstName = newUserInstance.FirstName
		// 	existingUser.LastName = newUserInstance.LastName
		// 	existingUser.Email = newUserInstance.Email
		// 	existingUser.Gender = newUserInstance.Gender
		// 	existingUser.Phone = newUserInstance.Phone
		// 	existingUser.Password = newUserInstance.Password
		// 	existingUser.DeletedAt.Time = time.Time{}
		// 	existingUser.DeletedAt.Valid = false
		// 	if err := initializers.DB.Save(&existingUser).Error; err != nil {
		// 		utils.HandleError(ctx, http.StatusInternalServerError, "Account Reactivated")
		// 		return
		// 	}
		// 	fmt.Println("recovered!!!!")
		//}
	} else {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to signup")
	}
	newUserInstance = NewUser{}
}

func ResendOtp(ctx *gin.Context) {
	var existOTP models.OTP

	result := initializers.DB.Where("email=?", newUserInstance.Email).First(&existOTP)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to resend")
		return
	}

	newOTP := authotp.GenerateOTP()
	fmt.Println("=================otp:", newOTP)

	fmt.Println("=========================existotp: ", existOTP.Otp)

	fmt.Println("===========================email: ", newUserInstance.Email)
	fmt.Println("===========================otpemail: ", existOTP.Email)
	if existOTP.Email == newUserInstance.Email {
		existOTP.Otp = newOTP
		existOTP.Email = newUserInstance.Email
		existOTP.Exp = time.Now().Add(5 * time.Minute)
		if err := initializers.DB.Save(&existOTP).Error; err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to save update OTP")
			return
		}
	}
	err := authotp.SendEmail(newUserInstance.Email, newOTP)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to send via Email")
		return
	}

	ctx.JSON(200, gin.H{
		"message": "new OTP sent successfully,please check your email",
	})

}

func UserLogin(ctx *gin.Context) {
	var postinguser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	hashedpassword, _ := bcrypt.GenerateFromPassword([]byte(postinguser.Password), bcrypt.DefaultCost)

	postinguser.Password = string(hashedpassword)

	if err := ctx.ShouldBindJSON(&postinguser); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Failed to bind Json")
		return
	}

	result := initializers.DB.Where("email=?", postinguser.Email).First(&user)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "invalid email or password")
		return
	}
	password := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(postinguser.Password))
	if password != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid Password")
		return
	}

	if user.Status == "Active" {
		tokenstring, err := middleware.JwtToken(ctx, user.ID, postinguser.Email, RoleUser)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to create token")
			return
		}
		ctx.SetCookie("Authorization"+RoleUser, tokenstring, int((time.Hour * 1).Seconds()), "", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login Successfully",
		})
	} else {
		utils.HandleError(ctx, http.StatusForbidden, "you are blocked by admin")
	}
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("Authorization"+RoleUser, "", -1, "", "", false, true)
	ctx.JSON(200, gin.H{
		"message": "Successfully Logout",
	})
}

func ForgotPassword(ctx *gin.Context) {
	type input struct {
		Email string `json:"email"`
	}
	var Input input
	if err := ctx.ShouldBindJSON(&Input); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to bind input")
		return
	}
	result := initializers.DB.Where("email = ?", Input.Email).First(&user)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to check email")
		return
	}
	otp := authotp.GenerateOTP()
	var OTP models.OTP
	initializers.DB.Where("email=?", Input.Email).First(&OTP)

	if Input.Email == OTP.Email {
		OTP.Otp = otp
		OTP.Exp = time.Now().Add(5 * time.Minute)

		initializers.DB.Save(&OTP)
	}

	otpRecord := models.OTP{
		Otp:    otp,
		Email:  Input.Email,
		Exp:    time.Now().Add(5 * time.Minute),
		UserID: user.ID,
	}
	initializers.DB.Create(&otpRecord)

	if err := authotp.SendEmail(Input.Email, otp); err != nil {
		utils.HandleError(ctx, http.StatusForbidden, "Failed to send OTP via email")
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "OTP for reset password is sent to your email,validate OTP",
	})
}

func OtpCheck(ctx *gin.Context) {
	type OTP struct {
		Otp string `json:"otp"`
	}
	var newOTP OTP
	if err := ctx.ShouldBindJSON(&newOTP); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Json Binding Error")
		return
	}

	var existingOTP models.OTP

	if err := initializers.DB.First(&existingOTP, "otp=?", newOTP.Otp); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid OTP")
		return
	}

	if time.Now().After(existingOTP.Exp) {
		utils.HandleError(ctx, http.StatusForbidden, "OTP has expired. Please request a new otp.")
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
		Newpassword string `json:"newPassword"`
	}

	var input Input

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Failed to Bind input")
		return
	}

	if !Confirmation {
		utils.HandleError(ctx, http.StatusInternalServerError, "Validate the OTP first")
		return
	} else {
		errr := initializers.DB.Where("email=?", input.Email).First(&user)
		if errr.Error != nil {
			utils.HandleError(ctx, http.StatusNotFound, "Email account not exist")
			return
		}

		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(input.Newpassword), bcrypt.DefaultCost)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		if err := initializers.DB.Model(&user).Update("password", string(hashedpassword)).Error; err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to update password")
			return
		}
		Confirmation = false
		ctx.JSON(200, gin.H{
			"status":  "success",
			"Message": "Password reset successfull",
		})
	}
}

func RefreshToken(ctx *gin.Context) {
	returnObject := gin.H{
		"status":  "OK",
		"message": "Refresh Token route",
	}

	email, exists := ctx.Get("email")
	if !exists {
		log.Println("Email key not found")

		returnObject["message"] = "Email not found."
		ctx.JSON(401, returnObject)
		return
	}
	var user models.User

	initializers.DB.First(&user, "email=?", email)

	if user.ID == 0 {
		returnObject["message"] = "User not found"
		ctx.JSON(400, returnObject)
		return
	}

	token, err := middleware.JwtToken(ctx, user.ID, user.Email, RoleUser)

	if err != nil {
		returnObject["message"] = "Token Creation Error."
		ctx.JSON(401, returnObject)
		return
	}

	ctx.SetCookie("Authorization"+RoleUser, token, int((time.Hour * 1).Seconds()), "", "", false, true)

	returnObject["token"] = token
	returnObject["user"] = user

	ctx.JSON(200, returnObject)
}

func ReactivateUser() (*models.User, error) {
	var existingUser models.User
	res := initializers.DB.Unscoped().Where("email=?", newUserInstance.Email).First(&existingUser)
	if res.Error == nil && existingUser.DeletedAt.Valid {
		existingUser.FirstName = newUserInstance.FirstName
		existingUser.LastName = newUserInstance.LastName
		existingUser.Email = newUserInstance.Email
		existingUser.Gender = newUserInstance.Gender
		existingUser.Phone = newUserInstance.Phone
		existingUser.Password = newUserInstance.Password
		existingUser.DeletedAt.Time = time.Time{}
		existingUser.DeletedAt.Valid = false
		return &existingUser, nil
	}
	return nil, res.Error
}

func saveUser(user *models.User) {
	initializers.DB.Save(user)
}
