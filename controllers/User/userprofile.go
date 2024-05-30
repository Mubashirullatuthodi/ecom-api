package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authotp "github.com/mubashir/e-commerce/AuthOTP"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"golang.org/x/crypto/bcrypt"
)

var ChangeConfirmation = false

func ListAddress(ctx *gin.Context) {
	var address []models.Address

	if err := initializers.DB.Find(&address).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  "Cant find products",
			"code":   500,
		})
		return
	}
	user_id := 0
	for _, value := range address {
		user_id = int(value.User_ID)
	}

	fmt.Println("user_id==============", user_id)

	if err := initializers.DB.Preload("User").Where("user_id=?", user_id).Find(&address).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"error":  "failed to list Address",
			"code":   500,
		})
		return
	}
	type UserDetails struct {
		Address_Id uint   `json:"address_id"`
		FirstName  string `json:"firstname"`
		LastName   string `json:"lastname"`
		//Email     string `json:"email"`
		//Gender    string `json:"gender"`
		Phone_No string `json:"phone_no"`
		Address  string `json:"address"`
		Town     string `json:"town"`
		District string `json:"district"`
		Pincode  string `json:"pincode"`
		State    string `json:"state"`
	}

	var details []UserDetails

	for _, value := range address {
		list := UserDetails{
			Address_Id: value.ID,
			FirstName:  value.User.FirstName,
			LastName:   value.User.LastName,
			Address:    value.Address,
			Phone_No:   value.User.Phone,
			Town:       value.Town,
			District:   value.District,
			Pincode:    value.Pincode,
			State:      value.State,
		}
		details = append(details, list)
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"Details": details,
	})
}

func ProfileChangePassword(ctx *gin.Context) {
	userid := ctx.GetUint("userid")
	var password struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := ctx.BindJSON(&password); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Failed to bind",
		})
		return
	}
	var user models.User
	result := initializers.DB.First(&user, userid)
	if result.Error != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.CurrentPassword))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "wrong old password",
		})
		return
	}

	if password.NewPassword != password.ConfirmPassword {
		ctx.JSON(400, gin.H{
			"error": "failssssss",
		})
	} else {
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(500, gin.H{
				"status": "failed to hash password",
			})
			return
		}

		user.Password = string(hashedpassword)
		r := initializers.DB.Save(&user)
		if r.Error != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to Change password",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "Password Changed Successfully",
		})
	}
}

func ProfileForgotPassword(ctx *gin.Context) {
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
		"message": "OTP for reset password is sent to your email,validate OTP",
	})
}


func EditProfile(ctx *gin.Context) {
	var useraddress models.Address
	var users models.User

	var editprofile struct {
		FirstName string `json:"firstname"`
		Gender    string `json:"gender"`
		Email     string `json:"email"`
		Phone_no  string `json:"phone_no"`
		Address   string `json:"address"`
		Pincode   string `json:"pincode"`
	}

	if err := ctx.ShouldBindJSON(&editprofile); err != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"error":  "failed to bind",
			"code":   500,
		})
		return
	}

	userid := ctx.GetUint("userid")
	fmt.Println("=================", userid)

	if err := initializers.DB.Where("user_id=?", userid).First(&useraddress).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "address not found",
		})
		return
	}

	if err := initializers.DB.Where("id=?", userid).First(&users).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	users.FirstName = editprofile.FirstName
	users.Gender = editprofile.Gender
	users.Email = editprofile.Email
	users.Phone = editprofile.Phone_no
	useraddress.Address = editprofile.Address
	useraddress.Pincode = editprofile.Pincode

	if err := initializers.DB.Save(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user",
		})
		return
	}

	if err := initializers.DB.Save(&useraddress).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "address": useraddress})

}



