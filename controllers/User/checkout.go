package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func CheckoutCart(ctx *gin.Context) {
	var carts []models.Cart
	userid := ctx.GetUint("userid")

	initializers.DB.Where("user_id", userid).First(&carts)

	if err := initializers.DB.Preload("User").Preload("Product").Find(&carts).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Fetch Items",
		})
		return
	}

	var Total []float64

	for _, v := range carts {
		fmt.Println("================", v.Product.Price)
		fmt.Println("================", v.Quantity)
		Total = append(Total, v.Product.Price*float64(v.Quantity))
	}
	fmt.Println("total===============", Total)

	var result float64
	for i := 0; i < len(Total)-1; i++ {
		sum := Total[i] + Total[i+1]
		result = sum
	}
	fmt.Println("sum===================", result)

	ctx.JSON(200, gin.H{
		"status":       "success",
		"total amount": fmt.Sprintf("%.2f/-", result),
	})
}
