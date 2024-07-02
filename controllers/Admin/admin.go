package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/middleware"
	"github.com/mubashir/e-commerce/models"
)

var user models.User

var RoleAdmin = "Admin"

func AdminSignUp(ctx *gin.Context) {
	var adminSignUp models.Admin
	err := ctx.ShouldBindJSON(&adminSignUp)
	if err != nil {
		ctx.JSON(406, gin.H{
			"status": "Fail",
			"error":  "Json binding error",
			"code":   406,
		})
		return
	}
	er := initializers.DB.Create(&adminSignUp)
	if er.Error != nil {
		ctx.JSON(500, gin.H{
			"status":  "Fail",
			"message": "Failed to signUp",
			"code":    500,
		})
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
		ctx.JSON(501, gin.H{
			"status": "Fail",
			"error":  "Fail to Bind json",
			"code":   501,
		})
		return
	}

	var existingAdmin models.Admin
	result := initializers.DB.Where("email = ?", admin.Email).First(&existingAdmin)

	if result.Error != nil {
		ctx.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Invalid email or Password",
			"code":   401,
		})
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
		ctx.JSON(401, gin.H{
			"status": "fail",
			"error":  "invalid email or password",
			"code":   401,
		})
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
		Phone_no  string `json:"phoneNo"`
		Status    string `json:"status"`
	}

	var List []list

	if err := initializers.DB.Find(&listUser).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"error":  err.Error(),
			"code":   500,
		})
		return
	}

	for _, value := range listUser {
		listUsers := list{
			Id:        int(value.ID),
			FirstName: value.FirstName,
			LastName:  value.LastName,
			Email:     value.Email,
			Gender:    value.Gender,
			Phone_no:  value.Phone,
			Status:    value.Status,
		}
		List = append(List, listUsers)
	}

	fmt.Println("list", List)

	ctx.JSON(http.StatusOK, List)
}

func DeleteUser(ctx *gin.Context) {

	id := ctx.Param("ID")
	fmt.Println("=============", id)
	if err := initializers.DB.Where("ID = ?", id).First(&user).Error; err != nil {
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "User not found",
			"code":   404,
		})
	}

	//soft delete
	initializers.DB.Delete(&user)

	ctx.JSON(204, gin.H{
		"status":  "success",
		"message": "user delete succesfully",
	})

}

func UpdateUser(ctx *gin.Context) {

	id := ctx.Param("ID")

	if err := initializers.DB.First(&user, id).Error; err != nil {
		fmt.Println("id", id)
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "user not found",
			"code":   404,
		})
		return
	}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to bind json",
			"code":   400,
		})
		return
	}

	if err := initializers.DB.Model(&user).Updates(user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":   "success",
		"messsage": "Succesfully updated",
	})
}

func Status(ctx *gin.Context) {
	var check models.User
	user := ctx.Param("ID")
	err := initializers.DB.First(&check, user)
	if err.Error != nil {
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "Can't Find User",
			"code":   404,
		})
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
