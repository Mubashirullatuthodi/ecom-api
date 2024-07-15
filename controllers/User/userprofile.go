package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authotp "github.com/mubashir/e-commerce/AuthOTP"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
	"golang.org/x/crypto/bcrypt"
)

var ChangeConfirmation = false

type UserDetails struct {
	AddressId uint   `json:"addressID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhoneNo   string `json:"phoneNo"`
	Address   string `json:"address"`
	Town      string `json:"town"`
	District  string `json:"district"`
	Pincode   string `json:"pincode"`
	State     string `json:"state"`
}

func ListAddress(ctx *gin.Context) {
	var address []models.Address

	if err := initializers.DB.Find(&address).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Cant find products")
		return
	}
	user_id := 0
	for _, value := range address {
		user_id = int(value.User_ID)
	}

	fmt.Println("user_id==============", user_id)

	if err := initializers.DB.Preload("User").Where("user_id=?", user_id).Find(&address).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to list address")
		return
	}
	var details []UserDetails

	for _, value := range address {
		details = append(details, UserDetails{
			AddressId: value.ID,
			FirstName: value.User.FirstName,
			LastName:  value.User.LastName,
			Address:   value.Address,
			PhoneNo:   value.User.Phone,
			Town:      value.Town,
			District:  value.District,
			Pincode:   value.Pincode,
			State:     value.State,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"Details": details,
	})
}

func ProfileChangePassword(ctx *gin.Context) {
	userID := ctx.GetUint("userid")
	var password struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := ctx.BindJSON(&password); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "failed to bind")
		return
	}
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to find User")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.CurrentPassword)); err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Wrong old password")
		return
	}

	if password.NewPassword != password.ConfirmPassword {
		utils.HandleError(ctx, http.StatusBadRequest, "New password and confirm password do not match")
		return
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to hash Password")
		return
	}

	user.Password = string(hashedpassword)
	if err := initializers.DB.Save(&user).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to change password")
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Password changed successfully",
	})

}

func ProfileForgotPassword(ctx *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}
	var user models.User
	result := initializers.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to check email")
		return
	}
	otp := authotp.GenerateOTP()

	otpRecord := models.OTP{
		Otp:    otp,
		Email:  input.Email,
		Exp:    time.Now().Add(5 * time.Minute),
		UserID: user.ID,
	}
	if err := initializers.DB.Create(&otpRecord).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to create OTP record")
		return
	}

	if err := authotp.SendEmail(input.Email, otp); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to send OTP via email")
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "OTP for reset password is sent to your email. validate OTP.",
	})
}

func EditProfile(ctx *gin.Context) {
	var useraddress models.Address
	var user models.User

	var input struct {
		FirstName string `json:"firstName"`
		Gender    string `json:"gender"`
		Email     string `json:"email"`
		Phone_no  string `json:"phoneNo"`
		Address   string `json:"address"`
		Pincode   string `json:"pincode"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}

	userID := ctx.GetUint("userid")
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "User not found")
		return
	}

	if err := initializers.DB.Where("user_id=?", userID).First(&useraddress).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "address not found")
		return
	}

	user.FirstName = input.FirstName
	user.Gender = input.Gender
	user.Email = input.Email
	user.Phone = input.Phone_no
	useraddress.Address = input.Address
	useraddress.Pincode = input.Pincode

	if err := initializers.DB.Save(&user).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to save user")
		return
	}

	if err := initializers.DB.Save(&useraddress).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to save")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Address updated successfully",
		"address": useraddress,
	})

}
