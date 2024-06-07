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
		OrderID              uint
		ProductID            uint
		OrderCode            string
		Product_name         string
		Category_name        string
		Product_Price        float64
		TotalQuantity        int
		TotalPrice           float64
		Payment_Method       string
		Order_Date           string
		Order_Status         string
		Order_Product_status string
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
			OrderID:              v.OrderID,
			ProductID:            v.ProductID,
			OrderCode:            v.Order.OrderCode,
			Product_name:         v.Product.Name,
			Product_Price:        v.Product.Price,
			TotalQuantity:        v.Quantity,
			TotalPrice:           v.SubTotal,
			Payment_Method:       v.Order.PaymentMethod,
			Category_name:        v.Product.Category.Name,
			Order_Date:           formatdate,
			Order_Status:         v.Order.OrderStatus,
			Order_Product_status: v.ProductOrderStatus,
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

		fmt.Println("orderid===============================", order.ID)

		if err := initializers.DB.Where("order_id=?", order.ID).Find(&orderItems).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "failed to fetch order items",
			})
			return
		}

		//returning amount
		var grandTotal float64

		for _, item := range orderItems {
			product := models.Product{}
			if err := initializers.DB.Where("id=?", item.ProductID).First(&product).Error; err != nil {
				ctx.JSON(500, gin.H{
					"error": "failed to fetch products",
				})
				return
			}

			grandTotal += item.SubTotal

			productQty, _ := strconv.ParseUint(product.Quantity, 10, 64)
			productQty += uint64(item.Quantity)
			product.Quantity = strconv.FormatUint(productQty, 10)
			fmt.Println("result======================", productQty)

			if err := initializers.DB.Save(&product).Error; err != nil {
				ctx.JSON(500, gin.H{
					"error": "failed to update product",
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

		userid := ctx.GetUint("userid")

		wallet := models.Wallet{
			Balance: grandTotal,
			UserID:  userid,
		}

		if err := initializers.DB.Create(&wallet).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "failed to return cash to wallet",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "Order cancelled successfully",
		})
	}

}

type orderid struct {
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Status    string `json:"status"`
}

func CancelSingleProduct(ctx *gin.Context) {
	var orderid orderid
	var orderitems []models.OrderItems
	userid := ctx.GetUint("userid")

	if err := ctx.ShouldBindJSON(&orderid); err != nil {
		ctx.JSON(404, gin.H{
			"error": "Failed to Bind",
		})
		return
	}
	//handlig status
	if orderid.Status != "Cancel" {
		ctx.JSON(500, gin.H{
			"error": "enter the status 'Cancel'",
		})
		return
	}
	if err := initializers.DB.Joins("Order").Where("order_items.order_id=?", orderid.OrderID).Find(&orderitems).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid order id",
		})
	}

	for _, v := range orderitems {
		fmt.Println("-----------", v.Order.UserId)
		fmt.Println("-----------", v.ProductID)

		if userid != v.Order.UserId && orderid.ProductID != v.ProductID {
			ctx.JSON(500, gin.H{
				"error": "invalid product",
			})
			return
		} else {
			if err := initializers.DB.Model(&orderitems).Where("product_id=?", orderid.ProductID).Update("product_order_status", orderid.Status).Error; err != nil {
				ctx.JSON(500, gin.H{"error": "failed to update status"})
				return
			}
		}

	}

	ctx.JSON(200, gin.H{
		"status": "success",
	})

}
