package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func AddtoCart(ctx *gin.Context) {
	var addcart struct {
		Userid    uint `json:"user_id"`
		Productid uint `json:"product_id"`
		Quantity  uint `json:"quantity"`
	}

	if err := ctx.BindJSON(&addcart); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Failed to bind",
		})
		return
	}
	id := ctx.GetUint("userid")
	fmt.Println("id=====================", id)

	var product models.Product
	result := initializers.DB.First(&product, addcart.Productid)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"mesage": "Product not found",
			"code":   400,
		})
		return
	}

	qty, _ := strconv.ParseUint(product.Quantity, 10, 32)

	if addcart.Quantity > uint(qty) {
		ctx.JSON(400, gin.H{
			"status":          "Fail",
			"Error":           "Out of stock",
			"available stock": qty,
			"code":            400,
		})
		return
	}

	var existingCart models.Cart
	result = initializers.DB.Where("user_id=? AND product_id=?", id, addcart.Productid).First(&existingCart)
	if result.Error == nil {
		newQuantity := existingCart.Quantity + addcart.Quantity
		if newQuantity > uint(qty) {
			ctx.JSON(400, gin.H{
				"status":          "Fail",
				"Error":           "Out Of Stock",
				"available stock": qty,
				"code":            400,
			})
			return
		}
		existingCart.Quantity = newQuantity
		if err := initializers.DB.Save(&existingCart).Error; err != nil {
			ctx.JSON(400, gin.H{
				"status": "Fail",
				"Error":  "Failede to update Cart",
				"code":   400,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "Cart Updated Successfully",
		})
	} else {

		Addcart := models.Cart{
			User_ID:    id,
			Product_ID: addcart.Productid,
			Quantity:   addcart.Quantity,
		}

		if err := initializers.DB.Create(&Addcart).Error; err != nil {
			ctx.JSON(400, gin.H{
				"status": "Fail",
				"Error":  "Failed to add cart",
				"code":   400,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":  "Success",
			"Message": "cart added Successfully",
		})
	}
}

func ListCart(ctx *gin.Context) {
	var listcart []models.Cart

	if err := initializers.DB.Preload("User").Preload("Product").Find(&listcart).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Fetch Items",
		})
		return
	}

	type Showcart struct {
		CartId              uint   `json:"cartid"`
		Userid              uint   `json:"userid"`
		Product_name        string `json:"product_name"`
		Product_image       string `json:"product_image"`
		Product_description string `json:"product_description"`
		Quantity            string `json:"quantity"`
		Available_Quantity  string `json:"stock_available"`
		Price               string `json:"price"`
	}

	var List []Showcart

	for _, value := range listcart {
		qty := strconv.FormatUint(uint64(value.Quantity), 10)
		fmt.Println("============================", qty)
		total := value.Product.Price * float64(value.Quantity)
		fmt.Println("total=============================", total)
		totalPrice := strconv.FormatFloat(total, 'f', -1, 64)
		list := Showcart{
			CartId:              value.ID,
			Userid:              value.User_ID,
			Product_name:        value.Product.Name,
			Product_image:       value.Product.ImagePath[0],
			Product_description: value.Product.Description,
			Price:               totalPrice,
			Quantity:            qty,
			Available_Quantity:  value.Product.Quantity,
		}
		List = append(List, list)
	}
	token, _ := ctx.Get("token")
	fmt.Println("jwt----------------------------", token)
	//fmt.Println("=======================", List)

	ctx.JSON(200, gin.H{
		"status":   "success",
		"products": List,
	})
}

func RemoveCart(ctx *gin.Context) {
	var carts models.Cart

	id := ctx.Param("ID")

	if err := initializers.DB.Where("ID = ?", id).First(&carts).Error; err != nil {
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "cart not found",
			"code":   404,
		})
	}

	initializers.DB.Delete(&carts)

	ctx.JSON(204, gin.H{
		"status":  "success",
		"message": "cart removed successfully",
	})
}

// func ReducingQuantity(ctx *gin.Context){
// 	id := ctx.Param("ID")

// }
