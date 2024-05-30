package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func AddAddress(ctx *gin.Context) {
	var user models.User

	userID := ctx.GetUint("userid")

	var inputAddress struct {
		User_ID  uint   `json:"user_id"`
		Address  string `json:"address"`
		Town     string `json:"town"`
		District string `json:"district"`
		Pincode  string `json:"pincode"`
		State    string `json:"state"`
	}

	if err := ctx.ShouldBindJSON(&inputAddress); err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "Failed to Bind",
			"code":   400,
		})
		return
	}

	if err := initializers.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "user ID not found to add address",
			"code":   400,
		})
		return
	}

	AddressUser := models.Address{

		User_ID:  userID,
		Address:  inputAddress.Address,
		District: inputAddress.District,
		Town:     inputAddress.Town,
		State:    inputAddress.State,
		Pincode:  inputAddress.Pincode,
	}

	if err := initializers.DB.Create(&AddressUser).Error; err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "Failed to add Address",
			"code":   400,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "Success",
		"Message": "Address added Successlly",
	})
}

func EditAddress(ctx *gin.Context) {
	userID:=ctx.GetUint("userid")

	var editAddress struct {
		User_id  uint   `json:"user_id"`
		Address  string `json:"address"`
		Town     string `json:"town"`
		District string `json:"district"`
		Pincode  string `json:"pincode"`
		State    string `json:"state"`
	}

	if err := ctx.BindJSON(&editAddress); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to bind",
		})
		return
	}

	id := ctx.Param("ID")

	var address models.Address
	if err := initializers.DB.First(&address, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "address not found",
		})
		return
	}
	if address.User_ID == userID {
		address.Address = editAddress.Address
		address.District = editAddress.District
		address.Pincode = editAddress.Pincode
		address.State = editAddress.Town
		address.Town = editAddress.Town

		if err := initializers.DB.Save(&address).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "address": address})
	} else {
		ctx.JSON(400, gin.H{
			"error": "user_id not found",
		})
		return
	}

}

func DeleteAddress(ctx *gin.Context) {
	user_id, exist := ctx.Get("userid")
	if !exist {
		ctx.JSON(500, gin.H{
			"error": "user_id not found",
		})
		return
	}

	fmt.Println("user:======================", user_id)
	var address models.Address

	id := ctx.Param("ID")
	fmt.Println("=============", id)
	if err := initializers.DB.Where("ID = ?", id).First(&address).Error; err != nil {
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "User not found",
			"code":   404,
		})
	}

	//soft delete
	if user_id == address.User_ID {
		initializers.DB.Delete(&address)

		ctx.JSON(204, gin.H{
			"status":  "success",
			"message": "address delete succesfully",
		})
	} else {
		ctx.JSON(400, gin.H{
			"status":  "fail",
			"message": "user not found",
		})
	}

}
