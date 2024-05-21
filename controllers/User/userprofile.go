package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

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
			Address:  value.Address,
			Phone_No: value.User.Phone,
			Town:     value.Town,
			District: value.District,
			Pincode:  value.Pincode,
			State:    value.State,
		}
		details = append(details, list)
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"Details": details,
	})
}
