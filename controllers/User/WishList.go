package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func AddToWishlist(ctx *gin.Context) {
	userid := ctx.GetUint("userid")
	id := ctx.Param("ID")
	ProductID, _ := strconv.ParseUint(id, 10, 32)

	var WishList models.WishList

	if err := initializers.DB.Where("user_id = ? AND product_id = ?", userid, id).First(&WishList).Error; err == nil {
		ctx.JSON(200, gin.H{
			"message": "Product Already Exist in the Wishist",
		})
		return
	}

	wishlist := models.WishList{
		UserID:    userid,
		ProductID: uint(ProductID),
	}

	if err := initializers.DB.Create(&wishlist).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "could not add Product to wishlist",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Product Added to Wishlist",
	})
}

func RemoveWishlist(ctx *gin.Context) {
	userid := ctx.GetUint("userid")
	productid := ctx.Param("ID")
	convID, _ := strconv.ParseUint(productid, 10, 32)

	var WishList models.WishList

	if err := initializers.DB.Where("product_id=? AND user_id=?", uint(convID), userid).Delete(&WishList).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to remove item",
		})
		return
	}
	ctx.JSON(204, gin.H{
		"status":  "success",
		"message": "Item remove successfuly",
	})

}

func ListWishList(ctx *gin.Context) {
	var wishlist []models.WishList

	userid := ctx.GetUint("userid")

	if err := initializers.DB.Where("user_id=?", userid).Preload("Product").Find(&wishlist).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error fetching wishlist",
		})
		return
	}
	type show struct {
		Productname  string
		Productprice float64
	}

	var list []show

	for _, v := range wishlist {
		new := show{
			Productname:  v.Product.Name,
			Productprice: v.Product.Price,
		}
		list = append(list, new)
	}

	ctx.JSON(200, gin.H{
		"wishlist": list,
	})
}
