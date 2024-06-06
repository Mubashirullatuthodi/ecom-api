package controllers

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func OrderDetails(ctx *gin.Context) {
	var orders []models.OrderItems

	type showOrders struct {
		OrderCode      string
		Product_name   string
		Category_name  string
		Product_Price  float64
		TotalQuantity  int
		TotalPrice     float64
		Payment_Method string
		Order_Date     string
		Order_Status   string
	}

	userid := ctx.GetUint("userid")

	if err := initializers.DB.Preload("Order").Preload("Product").Preload("Product.Category").Preload("Order.Address").Preload("Order.User").Joins("JOIN orders ON orders.id = order_items.order_id").Where("orders.user_id=?", userid).Find(&orders).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Fetch Items",
		})
		return
	}

	var List []showOrders

	for _, v := range orders {
		//Format Date
		formatdate := v.Order.OrderDate.Format("2006-01-02 15:04:05")

		show := showOrders{
			OrderCode:      v.Order.OrderCode,
			Product_name:   v.Product.Name,
			Product_Price:  v.Product.Price,
			TotalQuantity:  v.Quantity,
			TotalPrice:     v.SubTotal,
			Payment_Method: v.Order.PaymentMethod,
			Category_name:  v.Product.Category.Name,
			Order_Date:     formatdate,
			Order_Status:   v.Order.OrderStatus,
		}
		List = append(List, show)
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"Orders": List,
	})
}

func CancelOrder(ctx *gin.Context) {
	ordercode := ctx.Request.FormValue("orderID")

	if ordercode == "" {
		ctx.JSON(400, gin.H{
			"error": "Order code is required",
		})
		return
	}

	var order models.Order

	if err := initializers.DB.Where("order_code=?", ordercode).First(&order).Error; err != nil {
		ctx.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.OrderStatus == "Cancelled" {
		ctx.JSON(400, gin.H{
			"error": "this order is already cancelled",
		})
		return
	} else {

		var orderItems []models.OrderItems

		fmt.Println("orderid===============================",order.ID)

		if err := initializers.DB.Where("order_id=?", order.ID).Find(&orderItems).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "failed to fetch order items",
			})
			return
		}

		for _,item:=range orderItems{
			product:=models.Product{}
			if err:=initializers.DB.Where("id=?",item.ProductID).First(&product).Error;err!=nil{
				ctx.JSON(500,gin.H{
					"error":"failed to fetch products",
				})
				return
			}

			productQty,_:=strconv.ParseUint(product.Quantity,10,64)
			productQty+=uint64(item.Quantity)
			product.Quantity=strconv.FormatUint(productQty,10)
			fmt.Println("result======================",productQty)

			if err:=initializers.DB.Save(&product).Error;err!=nil{
				ctx.JSON(500,gin.H{
					"error":"failed to update product",
				})
				return
			}
		}

		order.OrderStatus = "Cancelled"

		if err := initializers.DB.Save(&order).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to cancel the order",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "Order cancelled successfully",
		})
	}

}