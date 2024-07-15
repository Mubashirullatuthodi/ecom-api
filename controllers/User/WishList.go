package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func AddToWishlist(ctx *gin.Context) {
	userID := ctx.GetUint("userid")
	productIDStr := ctx.Param("ID")
	ProductID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var WishList models.WishList

	if err := initializers.DB.Where("user_id = ? AND product_id = ?", userID, uint(ProductID)).First(&WishList).Error; err == nil {
		ctx.JSON(200, gin.H{
			"message": "Product Already Exist in the Wishist",
		})
		return
	}

	wishlist := models.WishList{
		UserID:    userID,
		ProductID: uint(ProductID),
	}

	if err := initializers.DB.Create(&wishlist).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "could not add product to wishlist")
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Product Added to Wishlist",
	})
}

func RemoveWishlist(ctx *gin.Context) {
	userID := ctx.GetUint("userid")
	productIDStr := ctx.Param("ID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var WishList models.WishList

	if err := initializers.DB.Where("product_id=? AND user_id=?", uint(productID), userID).Delete(&WishList).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to remove item")
		return
	}
	ctx.JSON(204, gin.H{
		"status":  "success",
		"message": "Item remove successfuly",
	})

}

func ListWishList(ctx *gin.Context) {
	var wishlist []models.WishList

	userID := ctx.GetUint("userid")

	if err := initializers.DB.Where("user_id=?", userID).Preload("Product").Find(&wishlist).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Error fetching wishlist")
		return
	}
	type show struct {
		Productname  string
		Productprice float64
	}

	var list []show

	for _, item := range wishlist {
		list = append(list, show{
			Productname:  item.Product.Name,
			Productprice: item.Product.Price,
		})
	}

	ctx.JSON(200, gin.H{
		"wishlist": list,
	})
}
