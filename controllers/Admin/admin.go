package controllers

import (
	"fmt"
	"net/http"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	//"github.com/icza/session"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

type Admin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var user models.User

const RoleAdmin = "admin"

func AdminLogin(ctx *gin.Context) {
	var admin Admin

	// session := sessions.Default(ctx)
	// check := session.Get(RoleAdmin)
	// if check == nil {

	// }

	if err := ctx.BindJSON(&admin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	var existingAdmin Admin
	result := initializers.DB.Where("email = ?", admin.Email).First(&existingAdmin)

	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or Password",
		})
		return
	}

	if admin.Password != existingAdmin.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully Login to Admin panel",
	})
}

func ListUsers(ctx *gin.Context) {
	var listuser []models.User

	type list struct {
		Id        int
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Gender    string `json:"gender"`
		Phone_no  string `json:"phone_no"`
		Status    string `json:"status"`
	}

	var List []list

	if err := initializers.DB.Find(&listuser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, value := range listuser {
		listusers := list{
			Id:        int(value.ID),
			FirstName: value.FirstName,
			LastName:  value.LastName,
			Email:     value.Email,
			Gender:    value.Gender,
			Phone_no:  value.Phone,
			Status:    value.Status,
		}
		List = append(List, listusers)
	}

	fmt.Println("list", List)

	ctx.JSON(http.StatusOK, List)
}

func DeleteUser(ctx *gin.Context) {

	id := ctx.Param("ID")
	fmt.Println("=============", id)
	initializers.DB.Where("ID = ?", id).First(&user)

	//soft delete
	initializers.DB.Delete(&user)

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "user delete succesfully",
	})

}

func UpdateUser(ctx *gin.Context) {

	id := ctx.Param("ID")

	if err := initializers.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "usere not found",
		})
		return
	}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := initializers.DB.Model(&user).Updates(user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"messsage": "Succesfully updated",
	})
}

func Status(ctx *gin.Context) {
	var check models.User
	user := ctx.Param("ID")
	initializers.DB.First(&check, user)
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
