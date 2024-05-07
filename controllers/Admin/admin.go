package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mubashir/e-commerce/initializers"
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
			"Error":  "Json binding error",
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

	if admin.Password != existingAdmin.Password {
		ctx.JSON(401, gin.H{
			"status": "fail",
			"error":  "invalid email or password",
			"code":   401,
		})
		return
	}

	//token:=middleware.JwtTokenStart(ctx, existingAdmin.ID, existingAdmin.Email, RoleAdmin)
	//ctx.SetCookie("jwtToken"+RoleAdmin, token, int((time.Hour * 1).Seconds()), "/", "Audvision.online", false, false)
	ctx.JSON(202, gin.H{
		"status":  "success",
		"message": "successfully Logged to adminpanel",
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
		fmt.Println("id",id)
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
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
