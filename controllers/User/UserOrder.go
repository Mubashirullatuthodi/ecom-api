package controllers

import (
	"fmt"
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
			"Grand_total":     GrandTotal - v.CouponDiscount,
			"Coupon_discount": v.CouponDiscount,
			"ordPayID":        v.PayOrdID,
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
		OrderItemID   uint
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
	orderID := ctx.Param("ID")
	convID, _ := strconv.ParseUint(orderID, 10, 32)

	if err := initializers.DB.Preload("Order").Preload("Product").Preload("Product.Category").Preload("Order.Address").Preload("Order.User").Joins("JOIN orders ON orders.id = order_items.order_id").Where("orders.user_id=? AND order_id=?", userid, uint(convID)).Find(&orders).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Fetch Items",
		})
		return
	}

	var List []showOrders
	offer := 0.0
	couponOffer := 0
	grandTotal := 0
	cancelAmount := 0

	for _, v := range orders {
		var coupon models.Coupons
		initializers.DB.Where("coupon_code=?", v.Order.CouponCode).First(&coupon)

		if v.OrderStatus == "Cancelled" {
			cancelAmount = int(v.SubTotal)
		}
		offer = controllers.OfferCalc(v.ProductID) * float64(v.Quantity)
		couponOffer = v.Order.CouponDiscount
		grandTotal += int(v.SubTotal)
		//Format Date
		formatdate := v.Order.CreatedAt.Format("2006-01-02 15:04:05")

		show := showOrders{
			OrderItemID:   v.ID,
			ProductID:     v.ProductID,
			OrderCode:     v.Order.OrderCode,
			Product_name:  v.Product.Name,
			Product_Price: v.Product.Price,
			OrderQuantity: v.Quantity,
			//TotalPrice:    v.SubTotal,
			// CouponDiscount:     int(coupon.Discount),
			// TotalAfterDiscount: int(v.Order.OrderAmount),
			Order_Date:    formatdate,
			Order_Status:  v.OrderStatus,
			OfferDiscount: int(offer),
		}
		List = append(List, show)
	}
	final := grandTotal - int(offer)
	ctx.JSON(200, gin.H{
		"status":             "success",
		"totalDiscount":      int(offer) + couponOffer,
		"totalAmount":        grandTotal - cancelAmount,
		"totalAfterDiscount": (final - couponOffer) - cancelAmount,
		"Orders":             List,
	})
}

func CancelOrder(ctx *gin.Context) {

	var orderitem models.OrderItems

	orderID := ctx.Param("ID")
	convorderid, _ := strconv.ParseUint(orderID, 10, 64)

	if err := initializers.DB.Where("id=?", uint(convorderid)).First(&orderitem); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error":  "Order not Exist",
			"status": 401,
		})
	} else {
		var product models.Product
		if err := initializers.DB.First(&product, orderitem.ProductID).Error; err != nil {
			ctx.JSON(401, gin.H{
				"error": "failed to fetch the product to return quantity",
			})
			return
		}
		beforeCancellationQuantity, _ := strconv.Atoi(product.Quantity)
		fmt.Println("before quantity------------------------------>", beforeCancellationQuantity)

		if orderitem.OrderStatus == "Cancelled" {
			ctx.JSON(200, gin.H{
				"message": "Order aready Cancelled",
				"status":  200,
			})
			return
		}
		var order models.Order
		if err := initializers.DB.Where("id=?", orderitem.OrderID).First(&order).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error": "failed to find order code!!",
			})
		}

		var paymentid models.Payment
		if err := initializers.DB.Where("receipt=?", order.OrderCode).First(&paymentid).Error; err != nil {
			ctx.JSON(404, gin.H{
				"error": "Failed to find payment information",
			})
			return
		}

		cancelAmount := orderitem.SubTotal - float64(orderitem.OfferPercentage)
		fmt.Println("-------------------------->", cancelAmount)
		fmt.Println("payedpaisa-------------------------->", paymentid.PaymentAmount)

		//begin transaction
		tx := initializers.DB.Begin()
		if err := tx.Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to begin transaction",
			})
			return
		}

		if err := initializers.DB.Model(&orderitem).Updates(&models.OrderItems{
			OrderStatus: "Cancelled",
		}); err.Error != nil {
			tx.Rollback()
			ctx.JSON(401, gin.H{
				"error":  "order not cancelled",
				"status": 401,
			})
		} else {

			beforeCancellationQuantity += orderitem.Quantity
			fmt.Println("after quantity--------------------->", beforeCancellationQuantity)
			convQuantity := strconv.Itoa(beforeCancellationQuantity)

			product.Quantity = convQuantity
			if err := initializers.DB.Save(&product).Error; err != nil {
				tx.Rollback()
				ctx.JSON(500, gin.H{
					"error": "Failed to update product Quantity",
				})
				return
			}

			//payment table updation
			if err := tx.Model(&paymentid).Update("payment_amount", paymentid.PaymentAmount-int(cancelAmount)).Error; err != nil {
				tx.Rollback()
				ctx.JSON(500, gin.H{
					"error": "Failed to update payment method",
				})
				return
			}

			userid := ctx.GetUint("userid")
			if err := tx.Create(&models.Wallet{
				Balance: float64(cancelAmount),
				UserID:  userid,
			}).Error; err != nil {
				tx.Rollback()
				ctx.JSON(500, gin.H{
					"error": "Failed to update wallet",
				})
				return
			}

			if err := tx.Commit().Error; err != nil {
				ctx.JSON(500, gin.H{
					"error": "Failed to commit transaction",
				})
				return
			}

			ctx.JSON(200, gin.H{
				"message": "Order Cancelled Succesfully",
			})

		}
	}
}

// func CancelingOrder(ctx *gin.Context) {
// 	var orderitems models.OrderItems

// 	orderItemID := ctx.Param("ID")
// 	convID, _ := strconv.ParseUint(orderItemID, 10, 32)

// 	if err := initializers.DB.Where("id=?", uint(convID)).First(&orderitems).Error; err != nil {
// 		ctx.JSON(400, gin.H{
// 			"error": "failed to find the orderitem",
// 		})
// 		return
// 	}

// 	if orderitems.OrderStatus == "Cancelled" {
// 		ctx.JSON(400, gin.H{
// 			"error": "the order already cancelled",
// 		})
// 		return
// 	}

// 	QtyBeforeCancellation := strconv.Atoi(orderitems.Quantity)

// 	var order models.Order
// 	if err := initializers.DB.Where("id=?", orderitems.OrderID).First(&order).Error; err != nil {
// 		ctx.JSON(400, gin.H{
// 			"error": "failed to find order code!!",
// 		})
// 	}

// 	var paymentid models.Payment
// 	initializers.DB.Where("receipt=?", order.OrderCode).First(&paymentid)

// 	cancelAmount := paymentid.PaymentAmount

// 	if err := initializers.DB.Model(&orderitems).Updates(&models.OrderItems{
// 		OrderStatus: "Cancelled",
// 	}); err.Error != nil {
// 		ctx.JSON(401, gin.H{
// 			"error":  "order not cancelled",
// 			"status": 401,
// 		})
// 	} else {
// 		ctx.JSON(200, gin.H{
// 			"message": "Order Cancelled Succesfully",
// 		})
// 		QtyBeforeCancellation += orderitems.Quantity
// 		fmt.Println("after quantity--------------------->", QtyBeforeCancellation)
// 		convQuantity := strconv.Itoa(QtyBeforeCancellation)

// 		var product models.Product
// 		initializers.DB.Where("id=?", orderitems.ProductID).First(&product)
// 		product.Quantity = convQuantity
// 		if err := initializers.DB.Save(&product).Error; err != nil {
// 			log.Fatalf("Failed to save product: %v", err)
// 		}

// 		userid := ctx.GetUint("userid")
// 		initializers.DB.Create(&models.Wallet{
// 			Balance: float64(cancelAmount),
// 			UserID:  userid,
// 		})

// 	}
// }
