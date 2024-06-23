package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func AdminViewOrder(ctx *gin.Context) {
	var order []models.Order
	var orderData []gin.H
	count := 0
	if err := initializers.DB.Preload("Address.User").Find(&order); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error": "failed to fetch order",
		})
		return
	}

	for _, v := range order {
		formatdate := v.CreatedAt.Format("2006-01-02 15:04:05")
		orderData = append(orderData, gin.H{
			"id":            v.ID,
			"user":          v.Address.User.FirstName,
			"address":       v.Address.Address,
			"appliedCoupon": v.CouponCode,
			"orderPrice":    v.OrderAmount,
			"PaymentMethod": v.PaymentMethod,
			"orderDate":     formatdate,
		})
		count++
	}
	ctx.JSON(200, gin.H{
		"data":        orderData,
		"totalOrders": count,
		"status":      200,
	})
}

func GetOrderDetails(ctx *gin.Context) {
	var orders []models.OrderItems

	type showOrders struct {
		OrderID        uint    //
		OrderCode      string  //
		Product_name   string  //
		Category_name  string  //
		Product_Price  float64 //
		TotalQuantity  int     //
		TotalPrice     float64 //
		User_name      string  //
		User_Address   string
		User_AddressID uint
		Order_Date     string
		Order_Status   string //
	}
	if err := initializers.DB.Preload("Order").Preload("Product").Preload("Product.Category").Preload("Order.Address").Preload("Order.User").Find(&orders).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Fetch Items",
		})
		return
	}

	var List []showOrders

	for _, v := range orders {
		formatdate := v.CreatedAt.Format("2006-01-02 15:04:05")
		show := showOrders{
			OrderID:        v.ID,
			OrderCode:      v.Order.OrderCode,
			Product_name:   v.Product.Name,
			Product_Price:  v.Product.Price,
			TotalQuantity:  v.Order.TotalQuantity,
			TotalPrice:     v.Order.OrderAmount,
			Category_name:  v.Product.Category.Name,
			User_name:      v.Order.User.FirstName,
			User_Address:   v.Order.Address.Address,
			Order_Status:   v.OrderStatus,
			User_AddressID: v.Order.AddressID,
			Order_Date:     formatdate,
		}
		List = append(List, show)
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"Orders": List,
	})
}

func ChangeOrderStatus(ctx *gin.Context) {
	var order models.OrderItems
	OrderID := ctx.Param("ID")
	convOrderID, _ := strconv.ParseUint(OrderID, 10, 64)
	Orderstatus := ctx.Request.FormValue("status")
	productID := ctx.Request.FormValue("productID")
	convID, _ := strconv.ParseUint(productID, 10, 64)
	err := initializers.DB.Where("order_id=? AND product_id=?", uint(convOrderID), uint(convID)).First(&order)
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

// func CancelOrder(ctx *gin.Context) {
// 	//userID:=ctx.Param("ID")
// 	OrderID := ctx.Request.FormValue("orderID")

// 	if OrderID == "" {
// 		ctx.JSON(400, gin.H{
// 			"error": "Order code is required",
// 		})
// 		return
// 	}

// 	var Orders models.Order

// 	if err := initializers.DB.Where("order_code=?", OrderID).First(&Orders).Error; err != nil {
// 		ctx.JSON(404, gin.H{
// 			"error": "Order not found",
// 		})
// 		return
// 	}
// 	fmt.Println("===============================status: ", Orders.OrderStatus)

// 	if Orders.OrderStatus == "Cancelled" {
// 		ctx.JSON(400, gin.H{
// 			"error": "this order is already cancelled",
// 		})
// 		return
// 	} else {
// 		var orderItems []models.OrderItems

// 		fmt.Println("orderid===============================", Orders.ID)

// 		if err := initializers.DB.Where("order_id=?", Orders.ID).Find(&orderItems).Error; err != nil {
// 			ctx.JSON(500, gin.H{
// 				"error": "failed to fetch order items",
// 			})
// 			return
// 		}

// 		for _, item := range orderItems {
// 			product := models.Product{}
// 			if err := initializers.DB.Where("id=?", item.ProductID).First(&product).Error; err != nil {
// 				ctx.JSON(500, gin.H{
// 					"error": "failed to fetch products",
// 				})
// 				return
// 			}

// 			productQty, _ := strconv.ParseUint(product.Quantity, 10, 64)
// 			productQty += uint64(item.Quantity)
// 			product.Quantity = strconv.FormatUint(productQty, 10)
// 			fmt.Println("result======================", productQty)

// 			if err := initializers.DB.Save(&product).Error; err != nil {
// 				ctx.JSON(500, gin.H{
// 					"error": "failed to update product",
// 				})
// 				return
// 			}
// 		}

// 		Orders.OrderStatus = "Cancelled"

// 		if err := initializers.DB.Save(&Orders).Error; err != nil {
// 			ctx.JSON(500, gin.H{
// 				"error": "Failed to cancel the order",
// 			})
// 			return
// 		}

// 		ctx.JSON(200, gin.H{
// 			"status":  "success",
// 			"message": "Order cancelled successfully",
// 		})

// 	}
// }
