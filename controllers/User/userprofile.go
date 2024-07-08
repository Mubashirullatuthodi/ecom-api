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

	var details []UserDetails

	for _, value := range address {
		list := UserDetails{
			AddressId: value.ID,
			FirstName: value.User.FirstName,
			LastName:  value.User.LastName,
			Address:   value.Address,
			PhoneNo:   value.User.Phone,
			Town:      value.Town,
			District:  value.District,
			Pincode:   value.Pincode,
			State:     value.State,
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
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := ctx.BindJSON(&password); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "failed to bind")
		return
	}
	var user models.User
	result := initializers.DB.First(&user, userid)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to find user")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.CurrentPassword))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "wrong old password")
		return
	}

	if password.NewPassword != password.ConfirmPassword {
		utils.HandleError(ctx, http.StatusBadRequest, "Fails")
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
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to change password")
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}
	result := initializers.DB.Where("email = ?", Input.Email).First(&user)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to check email")
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to send via OTP")
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
		FirstName string `json:"firstName"`
		Gender    string `json:"gender"`
		Email     string `json:"email"`
		Phone_no  string `json:"phoneNo"`
		Address   string `json:"address"`
		Pincode   string `json:"pincode"`
	}

	if err := ctx.ShouldBindJSON(&editprofile); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}

	userid := ctx.GetUint("userid")
	fmt.Println("=================", userid)

	if err := initializers.DB.Where("user_id=?", userid).First(&useraddress).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "address not found")
		return
	}

	if err := initializers.DB.Where("id=?", userid).First(&users).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "User not found")
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to save")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "address": useraddress})

}
