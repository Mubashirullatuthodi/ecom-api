package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

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
	for _, value := range Total {
		result += value
	}
	fmt.Println("sum===================", result)

	ctx.JSON(200, gin.H{
		"status":       "success",
		"total amount": fmt.Sprintf("%.2f rs", result),
	})
}

func PlaceOrder(ctx *gin.Context) {
	var checkout struct {
		Address_id   uint   `json:"address_id"`
		Payment_type string `json:"payment_type"`
	}
	if err := ctx.ShouldBind(&checkout); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Failed to bind",
		})
		return
	}

	userid := ctx.GetUint("userid")
	var cart []models.Cart

	initializers.DB.Preload("Product").Where("user_id", userid).Find(&cart)

	//appending the amount of each product with the multiple of quantity
	var Total []float64
	var Quantity int

	for _, v := range cart {
		Quantity += int(v.Quantity)
		fmt.Println("================", v.Product.Price)
		fmt.Println("================", v.Quantity)
		Total = append(Total, v.Product.Price*float64(v.Quantity))
	}

	//total of carts amount
	sum := 0.0

	for _, v := range Total {
		sum += v
	}

	fmt.Println("total=====================", Total)

	orderCode := GenerateOrderID(10)

	order := models.Order{
		OrderCode:     orderCode,
		UserId:        userid,
		PaymentMethod: checkout.Payment_type,
		AddressID:     checkout.Address_id,
		TotalQuantity: Quantity,
		TotalAmount:   sum,
		OrderDate:     time.Now(),
	}

	initializers.DB.Create(&order)

	for _, v := range cart {
		orderitems := models.OrderItems{
			OrderID:   order.ID,
			ProductID: v.Product_ID,
			Quantity:  int(v.Quantity),
			SubTotal:  v.Product.Price * float64(v.Quantity),
		}
		initializers.DB.Create(&orderitems)
	}

	/////////////////////////////////////
	//var newProductQuantity string
	for _, p := range cart {
		fmt.Println("product Quantity====================", p.Product.Quantity)
		var products models.Product
		if err := initializers.DB.First(&products, p.Product_ID).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error": "Product not found",
			})
			return
		}

		fmt.Println("Product Quantity Before:", products.Quantity)
		qty, _ := strconv.ParseUint(products.Quantity, 10, 32)

		newQuantity := uint64(qty) - uint64(p.Quantity)
		finalQuantity := strconv.FormatFloat(float64(newQuantity), 'f', -1, 64)
		fmt.Println("final===============================", finalQuantity)
		products.Quantity = finalQuantity

		if err := initializers.DB.Save(&products).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Could not update product quantity"})
			return
		}
		fmt.Println("Product Quantity After:", products.Quantity)
	}

	initializers.DB.Delete(&cart)

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "ordered successfullly",
	})

}

const charset = "123456789ASDQWEZXC"

func GenerateOrderID(length int) string {
	rand.Seed(time.Now().UnixNano())
	orderID := "ORD_ID"

	for i := 0; i < length; i++ {
		orderID += string(charset[rand.Intn(len(charset))])
	}
	return orderID
}
