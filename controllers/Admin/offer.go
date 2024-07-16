package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

type addoffer struct {
	ProductID uint    `json:"productID"`
	OfferName string  `json:"offerName"`
	Discount  float64 `json:"discount"`
	Created   string  `json:"created"`
	Expire    string  `json:"expire"`
}

func CreateOffer(ctx *gin.Context) {
	var addOffer addoffer
	var product models.Product
	ctx.ShouldBindJSON(&addOffer)
	if err := initializers.DB.Where("id=?", addOffer.ProductID).First(&product); err.Error != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "failed to bind")
		return
	}
	startDate, _ := time.Parse("2006-01-02", addOffer.Created)
	EndDate, _ := time.Parse("2006-01-02", addOffer.Expire)

	if err := initializers.DB.Create(&models.Offer{
		ProductID: addOffer.ProductID,
		OfferName: addOffer.OfferName,
		Discount:  addOffer.Discount,
		Created:   startDate,
		Expire:    EndDate,
	}); err.Error != nil {
		utils.HandleError(ctx, http.StatusConflict, "offer already exist")
		return
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
		utils.HandleError(ctx, http.StatusNotFound, "Offer not Found")
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
	var offer models.Offer
	if err := initializers.DB.Where("product_id=?", productID).First(&offer).Error; err != nil {
		return 0
	}
	if time.Now().After(offer.Expire) {
		return 0
	}
	var Discount float64 = 0
	if err := initializers.DB.Joins("Offer").First(&offerCheck, "products.id = ?", productID); err.Error != nil {
		return Discount
	}
	discountPercentage := offerCheck.Offer.Discount
	ProductAmount := offerCheck.Price
	Discount = (discountPercentage * float64(ProductAmount)) / 100

	fmt.Println("Discount percentage:", discountPercentage)
	fmt.Println("Product amount:", ProductAmount)
	fmt.Println("Discount:", Discount)

	return Discount
}

func OfferApply(ctx *gin.Context) {
	var offer models.Offer
	offerID := ctx.Param("ID")
	convID, _ := strconv.ParseUint(offerID, 32, 10)

	if err := initializers.DB.Unscoped().First(&offer, "id=?", uint(convID)).Error; err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid OfferID")
		return
	}

	action := ctx.PostForm("action")

	switch action {
	case "list":
		if err := initializers.DB.Unscoped().Model(&offer).Where("id=?", uint(convID)).Update("deleted_at", nil).Error; err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to restore offer")
			return
		}
		ctx.JSON(200, gin.H{
			"Message": "Offer Listed Successfully",
		})

	case "unlist":
		if err := initializers.DB.Where("id=?", uint(convID)).Delete(&offer).Error; err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to delete offer")
			return
		}

		ctx.JSON(200, gin.H{
			"Message": "Offer Unisted succesfully",
		})
	default:
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid Action")
		return
	}
}
