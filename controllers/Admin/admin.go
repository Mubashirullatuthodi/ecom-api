package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/middleware"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

var user models.User

var RoleAdmin = "Admin"

func AdminSignUp(ctx *gin.Context) {
	var adminSignUp models.Admin
	err := ctx.ShouldBindJSON(&adminSignUp)
	if err != nil {
		utils.HandleError(ctx, http.StatusNotAcceptable, "failed to bind")
		return
	}
	if err := initializers.DB.Create(&adminSignUp); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to signup")
		return
	}

	ctx.JSON(201, gin.H{
		"status":  "Success",
		"message": "Admin Added Succesfully",
	})

}

type Admin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminLogin(ctx *gin.Context) {
	var admin Admin

	if err := ctx.BindJSON(&admin); err != nil {
		utils.HandleError(ctx, http.StatusNotImplemented, "failed to bind")
		return
	}

	var existingAdmin models.Admin
	result := initializers.DB.Where("email = ?", admin.Email).First(&existingAdmin)

	if result.Error != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if admin.Email == existingAdmin.Email || admin.Password == existingAdmin.Password {
		tokenstring, _ := middleware.JwtToken(ctx, existingAdmin.ID, existingAdmin.Email, RoleAdmin)
		ctx.SetCookie("Authorization"+RoleAdmin, tokenstring, int((time.Hour * 1).Seconds()), "", "", false, true)
		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "successfully Logged to adminpanel",
		})
	} else {
		utils.HandleError(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}
}

func ListUsers(ctx *gin.Context) {
	var listUser []models.User

	type list struct {
		Id        int
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Gender    string `json:"gender"`
		PhoneNo   string `json:"phoneNo"`
		Status    string `json:"status"`
	}

	var List []list

	if err := initializers.DB.Find(&listUser).Error; err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	for _, value := range listUser {
		listUsers := list{
			Id:        int(value.ID),
			FirstName: value.FirstName,
			LastName:  value.LastName,
			Email:     value.Email,
			Gender:    value.Gender,
			PhoneNo:   value.Phone,
			Status:    value.Status,
		}
		List = append(List, listUsers)
	}

	fmt.Println("list", List)

	ctx.JSON(http.StatusOK, List)
}

func DeleteUser(ctx *gin.Context) {

	id := ctx.Param("ID")
	convID, _ := strconv.ParseUint(id, 10, 32)
	fmt.Println("=============", id)
	if err := initializers.DB.Where("id = ?", uint(convID)).Delete(&user).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "Failed to Delete")
		return
	}

	ctx.JSON(204, gin.H{
		"status":  "success",
		"message": "user delete succesfully",
	})

}

func UpdateUser(ctx *gin.Context) {

	id := ctx.Param("ID")
	convID, _ := strconv.ParseUint(id, 10, 32)
	if err := initializers.DB.First(&user, "id =?", uint(convID)).Error; err != nil {
		fmt.Println("id", id)
		utils.HandleError(ctx, http.StatusNotFound, "user not found")
		return
	}

	if err := ctx.BindJSON(&user); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Failed to bind json")
		return
	}

	if err := initializers.DB.Model(&user).Updates(user).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"status":   "success",
		"messsage": "Succesfully updated",
	})
}

func Status(ctx *gin.Context) {
	var check models.User
	userid := ctx.Param("ID")
	convID, _ := strconv.ParseUint(userid, 10, 32)
	err := initializers.DB.First(&check, "id=?", uint(convID))
	if err.Error != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "can't find user")
		return
	}
	if check.Status == "Active" {
		initializers.DB.Model(&check).Update("status", "Blocked")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "user Blocked",
		})
	} else {
		initializers.DB.Model(&check).Update("status", "Active")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User Unblocked",
		})
	}

}

func AdminLogout(ctx *gin.Context) {
	ctx.SetCookie("Authorization"+RoleAdmin, "", -1, "", "", false, true)
	ctx.JSON(200, gin.H{
		"Message": "Admin LOGOUT Successfully",
	})
}
