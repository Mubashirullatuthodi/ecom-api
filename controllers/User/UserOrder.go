package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"

	// "github.com/mubashir/e-commerce/controllers/User"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func ViewOrder(ctx *gin.Context) {
	var order []models.Order
	var listOrder []gin.H
	UserID := ctx.GetUint("userid")

	if err := initializers.DB.Preload("User").Preload("Address").Where("user_id=?", UserID).Find(&order); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	for _, v := range order {
		var payment models.Payment
		initializers.DB.Where("receipt=?", v.OrderCode).First(&payment)
		fmt.Println("=================", payment.PaymentStatus)
		formattime := v.CreatedAt.Format("2006-01-02 15:04:05")

		offer := 0.0
		GrandTotal := 0
		total := 0

		var orders []models.OrderItems
		initializers.DB.Where("order_id=?", v.ID).Find(&orders)
		for _, d := range orders {

			offer += controllers.OfferCalc(d.ProductID) * float64(d.Quantity)
			total += int(d.SubTotal)
			GrandTotal = total - int(offer)
		}
		listOrder = append(listOrder, gin.H{
			"orderID":         v.ID,
			"userID":          v.UserId,
			"paymentMethod":   v.PaymentMethod,
			"orderDate":       formattime,
			"paymentStatus":   payment.PaymentStatus,
			"paidAmount":      payment.PaymentAmount,
			"offer_discount":  offer,
			"Grand_total":     GrandTotal-v.CouponDiscount,
			"Coupon_discount": v.CouponDiscount,
		})
	}
	ctx.JSON(200, gin.H{
		"data":   listOrder,
		"status": 200,
	})
}

func OrderDetails(ctx *gin.Context) {
	var orders []models.OrderItems

	type showOrders struct {
		OrderID       uint
		ProductID     uint
		OrderCode     string
		Product_name  string
		Product_Price float64
		OrderQuantity int
		TotalPrice    float64
		//CouponDiscount     int //want to add
		//TotalAfterDiscount int //want to add
		Order_Date    string
		Order_Status  string
		OfferDiscount int
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
		var coupon models.Coupons
		initializers.DB.Where("coupon_code=?", v.Order.CouponCode).First(&coupon)

		offer := controllers.OfferCalc(v.ProductID) * float64(v.Quantity)

		//Format Date
		formatdate := v.Order.CreatedAt.Format("2006-01-02 15:04:05")

		show := showOrders{
			OrderID:       v.OrderID,
			ProductID:     v.ProductID,
			OrderCode:     v.Order.OrderCode,
			Product_name:  v.Product.Name,
			Product_Price: v.Product.Price,
			OrderQuantity: v.Quantity,
			TotalPrice:    v.SubTotal,
			// CouponDiscount:     int(coupon.Discount),
			// TotalAfterDiscount: int(v.Order.OrderAmount),
			Order_Date:    formatdate,
			Order_Status:  v.OrderStatus,
			OfferDiscount: int(offer),
		}
		List = append(List, show)
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"Orders": List,
	})
}

func CancelOrder(ctx *gin.Context) {

	var orderitem models.OrderItems
	//var wallet models.Wallet

	orderID := ctx.Param("ID")
	convorderid, _ := strconv.ParseUint(orderID, 10, 64)
	productid := ctx.Request.FormValue("productid")
	convproductid, _ := strconv.ParseUint(productid, 10, 64)
	fmt.Println("converted product id: ", convproductid)
	//finding product quantity
	var product models.Product
	if err := initializers.DB.First(&product, uint(convproductid)).Error; err != nil {
		ctx.JSON(401, gin.H{
			"error": "failed to fetch the product to return quantity",
		})
		return
	}

	beforeCancellationQuantity, _ := strconv.Atoi(product.Quantity)
	fmt.Println("before quantity------------------------------>", beforeCancellationQuantity)

	if err := initializers.DB.Where("order_id=? AND product_id=?", uint(convorderid), uint(convproductid)).First(&orderitem); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error":  "Order not Exist",
			"status": 401,
		})
	} else {
		if orderitem.OrderStatus == "Cancelled" {
			ctx.JSON(200, gin.H{
				"message": "Order aready Cancelled",
				"status":  200,
			})
			return
		}
		var order models.Order
		if err := initializers.DB.Where("id=?", uint(convorderid)).First(&order).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error": "failed to find order code!!",
			})
		}

		var paymentid models.Payment
		initializers.DB.Where("receipt=?", order.OrderCode).First(&paymentid)

		cancelAmount := paymentid.PaymentAmount
		fmt.Println("-------------------------->", cancelAmount)
		fmt.Println("payedpaisa-------------------------->", paymentid.PaymentAmount)

		if err := initializers.DB.Model(&orderitem).Updates(&models.OrderItems{
			OrderStatus: "Cancelled",
		}); err.Error != nil {
			ctx.JSON(401, gin.H{
				"error":  "order not cancelled",
				"status": 401,
			})
		} else {
			ctx.JSON(200, gin.H{
				"message": "Order Cancelled Succesfully",
			})
			beforeCancellationQuantity += orderitem.Quantity
			fmt.Println("after quantity--------------------->", beforeCancellationQuantity)
			convQuantity := strconv.Itoa(beforeCancellationQuantity)

			product.Quantity = convQuantity
			if err := initializers.DB.Save(&product).Error; err != nil {
				log.Fatalf("Failed to save product: %v", err)
			}

			userid := ctx.GetUint("userid")
			initializers.DB.Create(&models.Wallet{
				Balance: float64(cancelAmount),
				UserID:  userid,
			})

		}
	}
}
