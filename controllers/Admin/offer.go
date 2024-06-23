package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

type addoffer struct {
	ProductID uint    `json:"productID"`
	OfferName string  `json:"offername"`
	Discount  float64 `json:"discount"`
	Created   string  `json:"created"`
	Expire    string  `json:"expire"`
}

func CreateOffer(ctx *gin.Context) {
	var addoffer addoffer
	var product models.Product
	ctx.ShouldBindJSON(&addoffer)
	if err := initializers.DB.Where("id=?", addoffer.ProductID).First(&product); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error":  "Product not available",
			"status": 401,
		})
		return
	}
	startDate, _ := time.Parse("2006-01-02", addoffer.Created)
	EndDate, _ := time.Parse("2006-01-02", addoffer.Expire)

	if err := initializers.DB.Create(&models.Offer{
		ProductID: addoffer.ProductID,
		OfferName: addoffer.OfferName,
		Discount:  addoffer.Discount,
		Created:   startDate,
		Expire:    EndDate,
	}); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error":  "offer already exist",
			"status": 401,
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "offer added for the product",
			"status":  200,
		})
	}
}

func ListOffer(ctx *gin.Context) {
	var offers []models.Offer
	var offerList []gin.H
	if err := initializers.DB.Find(&offers); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error":  "Offer not Found",
			"status": 401,
		})
		return
	}
	for _, v := range offers {
		offerList = append(offerList, gin.H{
			"offerName":   v.OfferName,
			"offerAmount": v.Discount,
			"ProductID":   v.ProductID,
			"Created":     v.Created,
			"Expires":     v.Expire,
		})
	}

	ctx.JSON(200, gin.H{
		"data":   offerList,
		"status": 200,
	})
}

func OfferCalc(productID uint) float64 {
	var offerCheck models.Product
	var Discount float64 = 0
	if err := initializers.DB.Joins("Offer").First(&offerCheck, "products.id = ?", productID); err.Error != nil {
		return Discount
	} else {
		percentage := offerCheck.Offer.Discount
		fmt.Println("%:  ", percentage)
		ProductAmount := offerCheck.Price
		fmt.Println("product amount: ", ProductAmount)
		Discount = (percentage * float64(ProductAmount)) / 100
		fmt.Println("discount: ", Discount)
	}
	return Discount
}

func OfferApply(ctx *gin.Context) {
	var offer models.Offer
	offerid := ctx.Param("ID")
	convID, _ := strconv.ParseUint(offerid, 32, 10)

	if err := initializers.DB.Unscoped().First(&offer, "id=?", uint(convID)).Error; err != nil {
		ctx.JSON(400, gin.H{
			"Error": "Invalid OfferID",
		})
		return
	}

	action := ctx.PostForm("action")

	switch action {
	case "list":
		if err := initializers.DB.Unscoped().Model(&offer).Where("id=?", uint(convID)).Update("deleted_at", nil).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to restore offer",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"Message": "Offer Listed Successfully",
		})

	case "unlist":
		if err := initializers.DB.Where("id=?", uint(convID)).Delete(&offer).Error; err != nil {
			ctx.JSON(500, gin.H{
				"Error": "Failed to delete offer",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Message": "Offer Unisted succesfully",
		})
	default:
		ctx.JSON(400, gin.H{
			"Error": "Invalid Action",
		})
	}
}
