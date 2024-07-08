package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func AdminViewOrder(ctx *gin.Context) {
	var order []models.Order
	var orderData []gin.H
	count := 0
	if err := initializers.DB.Preload("Address.User").Find(&order); err.Error != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "failed to fetch order")
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
		OrderID       uint    `json:"orderID"`
		OrderCode     string  `json:"orderCode"`
		Productname   string  `json:"ProductName"`
		Categoryname  string  `json:"categoryName"`
		ProductPrice  float64 `json:"productPrice"`
		TotalQuantity int     `json:"totalQuantity"`
		TotalPrice    float64 `json:"totalPrice"`
		Username      string  `json:"userName"`
		UserAddress   string  `json:"userAddress"`
		UserAddressID uint    `json:"userAddressID"`
		OrderDate     string  `json:"orderDate"`
		OrderStatus   string  `json:"orderStatus"`
	}
	if err := initializers.DB.Preload("Order").Preload("Product").Preload("Product.Category").Preload("Order.Address").Preload("Order.User").Find(&orders).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to fetch Items")
		return
	}

	var List []showOrders

	for _, v := range orders {
		formatdate := v.CreatedAt.Format("2006-01-02 15:04:05")
		show := showOrders{
			OrderID:       v.ID,
			OrderCode:     v.Order.OrderCode,
			Productname:   v.Product.Name,
			ProductPrice:  v.Product.Price,
			TotalQuantity: v.Order.TotalQuantity,
			TotalPrice:    v.Order.OrderAmount,
			Categoryname:  v.Product.Category.Name,
			Username:      v.Order.User.FirstName,
			UserAddress:   v.Order.Address.Address,
			OrderStatus:   v.OrderStatus,
			UserAddressID: v.Order.AddressID,
			OrderDate:     formatdate,
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
	OrderStatus := ctx.Request.FormValue("status")
	//productID := ctx.Request.FormValue("productID")
	//convID, _ := strconv.ParseUint(productID, 10, 64)
	err := initializers.DB.Where("id=?", uint(convOrderID)).First(&order)
	if err.Error != nil {
		utils.HandleError(ctx, http.StatusNotFound, "cant find the order")
		return
	}
	if order.OrderStatus == "Cancelled" {
		utils.HandleError(ctx, http.StatusBadRequest, "this order is already cancelled")
		return
	} else {
		switch OrderStatus {
		case "Delivered":
			if err := initializers.DB.Model(&order).Update("OrderStatus", "Delivered").Error; err != nil {
				utils.HandleError(ctx, http.StatusInternalServerError, "failed to update order status")
				return
			}
			ctx.JSON(200, gin.H{
				"message": "OrderStatus Changed to Delivered",
			})
		case "Pending":
			if err := initializers.DB.Model(&order).Update("OrderStatus", "Pending").Error; err != nil {
				utils.HandleError(ctx, http.StatusInternalServerError, "failed to update order status")
				return
			}
			ctx.JSON(200, gin.H{
				"message": "OrderStatus Changed to Pending",
			})
		default:
			utils.HandleError(ctx, http.StatusBadRequest, "Change the status into 'Delivered','Pending'")
			return
		}
	}

}
