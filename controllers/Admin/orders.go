package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func GetOrderDetails(ctx *gin.Context) {
	var orders []models.OrderItems

	type showOrders struct {
		OrderID        uint
		OrderCode      string
		Product_name   string
		Category_name  string
		Product_Price  float64
		TotalQuantity  int
		TotalPrice     float64
		Payment_Method string
		User_name      string
		User_Address   string
		Order_Date     time.Time
		Order_Status   string
	}
	if err := initializers.DB.Preload("Order").Preload("Product").Preload("Product.Category").Preload("Order.Address").Preload("Order.User").Find(&orders).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Fetch Items",
		})
		return
	}

	var List []showOrders

	for _, v := range orders {
		show := showOrders{
			OrderID:        v.ID,
			OrderCode:      v.Order.OrderCode,
			Product_name:   v.Product.Name,
			Product_Price:  v.Product.Price,
			TotalQuantity:  v.Order.TotalQuantity,
			TotalPrice:     v.Order.TotalAmount,
			Payment_Method: v.Order.PaymentMethod,
			Category_name:  v.Product.Category.Name,
			User_name:      v.Order.User.FirstName,
			User_Address:   v.Order.Address.Address,
			Order_Date:     v.Order.OrderDate,
			Order_Status:   v.Order.OrderStatus,
		}
		List = append(List, show)
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"Orders": List,
	})
}

func ChangeOrderStatus(ctx *gin.Context) {
	var order models.Order
	OrderID := ctx.Param("ID")
	Orderstatus := ctx.Request.FormValue("status")
	err := initializers.DB.First(&order, OrderID)
	if err.Error != nil {
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "Can't Find Order",
			"code":   404,
		})
		return
	}
	if order.OrderStatus == "Cancelled" {
		ctx.JSON(400, gin.H{
			"error": "This order is already cancelled",
		})
		return
	} else {
		switch Orderstatus {
		case "Delivered":
			if err := initializers.DB.Model(&order).Update("OrderStatus", "Delivered").Error; err != nil {
				ctx.JSON(500, gin.H{
					"status": "Fail",
					"Error":  "Failed to update order status",
				})
				return
			}
			ctx.JSON(200, gin.H{
				"message": "OrderStatus Changed to Delivered",
			})
		case "Pending":
			if err := initializers.DB.Model(&order).Update("OrderStatus", "Pending").Error; err != nil {
				ctx.JSON(500, gin.H{
					"status": "Fail",
					"Error":  "Failed to update order status",
				})
				return
			}
			ctx.JSON(200, gin.H{
				"message": "OrderStatus Changed to Pending",
			})
		default:
			ctx.JSON(400, gin.H{
				"message": "Change the status into 'Delivered','Pending'",
			})
		}
	}

}

func CancelOrder(ctx *gin.Context) {
	//userID:=ctx.Param("ID")
	OrderID := ctx.Request.FormValue("orderID")

	if OrderID == "" {
		ctx.JSON(400, gin.H{
			"error": "Order code is required",
		})
		return
	}

	var Orders models.Order

	if err := initializers.DB.Where("order_code=?", OrderID).First(&Orders).Error; err != nil {
		ctx.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}
	fmt.Println("===============================status: ", Orders.OrderStatus)

	if Orders.OrderStatus == "Cancelled" {
		ctx.JSON(400, gin.H{
			"error": "this order is already cancelled",
		})
		return
	} else {
		var orderItems []models.OrderItems

		fmt.Println("orderid===============================", Orders.ID)

		if err := initializers.DB.Where("order_id=?", Orders.ID).Find(&orderItems).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "failed to fetch order items",
			})
			return
		}

		for _, item := range orderItems {
			product := models.Product{}
			if err := initializers.DB.Where("id=?", item.ProductID).First(&product).Error; err != nil {
				ctx.JSON(500, gin.H{
					"error": "failed to fetch products",
				})
				return
			}

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

		Orders.OrderStatus = "Cancelled"

		if err := initializers.DB.Save(&Orders).Error; err != nil {
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
