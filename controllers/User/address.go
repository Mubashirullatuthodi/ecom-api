package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}

	if err := initializers.DB.First(&user, userID).Error; err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "User_ID not found to add address")
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
		utils.HandleError(ctx, http.StatusBadRequest, "failed to add address")
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "Success",
		"Message": "Address added Successlly",
	})
}

func EditAddress(ctx *gin.Context) {
	userID := ctx.GetUint("userid")

	var editAddress struct {
		UserID   uint   `json:"userID"`
		Address  string `json:"address"`
		Town     string `json:"town"`
		District string `json:"district"`
		Pincode  string `json:"pincode"`
		State    string `json:"state"`
	}

	if err := ctx.BindJSON(&editAddress); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}

	id := ctx.Param("ID")

	var address models.Address
	if err := initializers.DB.First(&address, id).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "address not found")
		return
	}
	if address.User_ID == userID {
		address.Address = editAddress.Address
		address.District = editAddress.District
		address.Pincode = editAddress.Pincode
		address.State = editAddress.State
		address.Town = editAddress.Town

		if err := initializers.DB.Save(&address).Error; err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "failed to save")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "address": address})
	} else {
		utils.HandleError(ctx, http.StatusBadRequest, "User_ID not found")
		return
	}

}

func DeleteAddress(ctx *gin.Context) {
	userID, exist := ctx.Get("userid")
	if !exist {
		utils.HandleError(ctx, http.StatusInternalServerError, "User_ID not found")
		return
	}

	fmt.Println("user:======================", userID)
	var address models.Address

	id := ctx.Param("ID")
	fmt.Println("=============", id)
	if err := initializers.DB.Where("ID = ?", id).First(&address).Error; err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "User not found")
		return
	}

	//soft delete
	if userID == address.User_ID {
		if err := initializers.DB.Delete(&address); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "failed to delete address")
			return
		}

		ctx.JSON(204, gin.H{
			"status":  "success",
			"message": "address delete succesfully",
		})
	} else {
		utils.HandleError(ctx, http.StatusBadRequest, "user not found")
		return
	}

}
