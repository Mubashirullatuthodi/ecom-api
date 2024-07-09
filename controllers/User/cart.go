package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func AddtoCart(ctx *gin.Context) {
	var addcart struct {
		Userid    uint `json:"userID"`
		Productid uint `json:"productID"`
		Quantity  uint `json:"quantity"`
	}

	if err := ctx.BindJSON(&addcart); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to binds")
		return
	}
	id := ctx.GetUint("userid")
	fmt.Println("id=====================", id)

	var product models.Product
	result := initializers.DB.First(&product, addcart.Productid)
	if result.Error != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "product not found")
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
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to update cart")
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
			utils.HandleError(ctx, http.StatusBadRequest, "Failed to add cart")
			return
		}

		ctx.JSON(200, gin.H{
			"status":  "Success",
			"Message": "cart added Successfully",
		})
	}
}

func ListCart(ctx *gin.Context) {
	var listCart []models.Cart

	if err := initializers.DB.Preload("User").Preload("Product").Find(&listCart).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to fetch items")
		return
	}

	type Showcart struct {
		CartId             uint   `json:"cartid"`
		Userid             uint   `json:"userid"`
		ProductName        string `json:"productName"`
		ProductImage       string `json:"productImage"`
		ProductDescription string `json:"productDescription"`
		Quantity           string `json:"quantity"`
		AvailableQuantity  string `json:"stockAvailable"`
		Price              string `json:"price"`
	}

	var List []Showcart

	var Grandtotal int

	for _, value := range listCart {

		qty := strconv.FormatUint(uint64(value.Quantity), 10)
		fmt.Println("============================", qty)
		total := value.Product.Price * float64(value.Quantity)
		fmt.Println("total=============================", total)
		totalPrice := strconv.FormatFloat(total, 'f', -1, 64)
		list := Showcart{
			CartId:             value.ID,
			Userid:             value.User_ID,
			ProductName:        value.Product.Name,
			ProductImage:       value.Product.ImagePath[0],
			ProductDescription: value.Product.Description,
			Price:              totalPrice,
			Quantity:           qty,
			AvailableQuantity:  value.Product.Quantity,
		}
		List = append(List, list)
		Grandtotal += int(total)
	}
	token, _ := ctx.Get("token")
	fmt.Println("jwt----------------------------", token)
	//fmt.Println("=======================", List)

	ctx.JSON(200, gin.H{
		"status":     "success",
		"products":   List,
		"GrandTotal": Grandtotal,
	})
}

func RemoveCart(ctx *gin.Context) {
	var carts models.Cart

	id := ctx.Param("ID")

	if err := initializers.DB.Where("ID = ?", id).First(&carts).Error; err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "cart not found")
		return
	}

	if err := initializers.DB.Delete(&carts); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to remove cart")
		return
	}

	ctx.JSON(204, gin.H{
		"status":  "success",
		"message": "cart removed successfully",
	})
}
