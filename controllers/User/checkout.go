package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func PlaceOrder(ctx *gin.Context) {
	var checkout struct {
		AddressID   uint   `json:"addressID"`
		PaymentType string `json:"paymentType"`
		CouponCode  string `json:"couponCode"`
	}
	if err := ctx.ShouldBind(&checkout); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "failed to bind")
		return
	}

	userid := ctx.GetUint("userid")
	var cart []models.Cart

	initializers.DB.Preload("Product").Where("user_id", userid).Find(&cart)

	//appending the amount of each product with the multiple of quantity
	var Total []int
	var Quantity int
	var discount float64

	for _, v := range cart {
		discount = controllers.OfferCalc(v.Product_ID)
		quantityPrice := (float64(v.Quantity) * v.Product.Price) - (float64(v.Quantity) * discount)
		Quantity += int(v.Quantity)
		Total = append(Total, int(quantityPrice))
	}

	//total of carts amount
	sum := 0
	for _, v := range Total {
		sum += v
	}
	if sum > 39000 {
		utils.HandleError(ctx, http.StatusBadRequest, "it cant be paid by online its above 39000.")
		return
	}

	fmt.Println("total=====================", Total)

	//checking coupon
	useridforcoupon := ctx.GetUint("userid")
	var couponcheck models.Coupons
	totalWithoutDiscount := sum
	//coupondiscount global
	var coupDisc float64 = 0

	if checkout.CouponCode != "" {
		//find coupon as valid or not
		if err := initializers.DB.Where("coupon_code=?", checkout.CouponCode).First(&couponcheck).Error; err != nil {

			fmt.Println("coupon code-------------->", couponcheck.CouponCode)
			utils.HandleError(ctx, http.StatusUnauthorized, "Invalid coupon")
			return
		}
		//find the total above the condition
		if totalWithoutDiscount < couponcheck.Condition {
			sum = totalWithoutDiscount
			ctx.JSON(401, gin.H{
				"Error":       fmt.Sprintf("total Amount needed %d to apply coupon", couponcheck.Condition),
				"TotalAmount": totalWithoutDiscount,
			})
			return
		}

		//check the coupon in the database
		var usageCount int64
		initializers.DB.Model(&models.CouponUsage{}).Where("user_id=? AND coupon_id=?", useridforcoupon, couponcheck.ID).Count(&usageCount)
		if usageCount > 0 {
			utils.HandleError(ctx, http.StatusUnauthorized, "you have already use this coupon")
			return
		}

		//log coupon usage
		CouponUsage := models.CouponUsage{
			UserID:   useridforcoupon,
			CouponID: couponcheck.ID,
		}
		initializers.DB.Create(&CouponUsage)

		fmt.Println("before minus discount-------------------->", sum)
		sum -= int(couponcheck.Discount)
		coupDisc = couponcheck.Discount
		fmt.Println("after minus discount------------------>", sum)
	}

	//adrress checking
	var adrress models.Address
	if err := initializers.DB.Where("user_id = ? AND id = ?", userid, checkout.AddressID).First(&adrress).Error; err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Address not found")
		return
	}

	orderCode := GenerateOrderID(10)

	//transaction
	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//check whether COD is done
	if len(cart) == 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "already done the order")
		return
	}

	//shipping charge
	var shippingCharge float64
	if sum < 15000 {
		shippingCharge = 40
		sum += int(shippingCharge)
	}

	//method checking
	if checkout.PaymentType == "COD" {
		if sum > 10000 {
			utils.HandleError(ctx, http.StatusUnauthorized, "COD not available above 10000 rs")
			return
		}
	}
	var orderPayID string
	//payment gateway
	fmt.Println("orderid------------------->", orderCode, "grand total------------->", sum)
	if checkout.PaymentType == "UPI" {
		orderPaymentID, err := PaymentSubmission(orderCode, sum)
		if err != nil {
			utils.HandleError(ctx, http.StatusUnauthorized, err.Error())
			tx.Rollback()
			return
		}
		ctx.JSON(200, gin.H{
			"message":   "Continue to payment",
			"paymentID": orderPaymentID,
			"status":    200,
		})
		orderPayID = orderPaymentID
		fmt.Println("paymentid-------------------->", orderPaymentID)
		fmt.Println("receipt-------------------->", orderCode)
		if err := tx.Create(&models.Payment{
			OrdID:         orderPaymentID,
			Receipt:       orderCode,
			PaymentStatus: "not done",
			PaymentAmount: int(sum),
		}); err.Error != nil {
			utils.HandleError(ctx, http.StatusUnauthorized, "failed to upload payment")
			fmt.Println("failed to upload payment details: ", err.Error)
			tx.Rollback()
		}
	}
	//order tables
	order := models.Order{
		OrderCode:      orderCode,
		PayOrdID:       orderPayID,
		UserId:         userid,
		CouponCode:     checkout.CouponCode,
		PaymentMethod:  checkout.PaymentType,
		AddressID:      checkout.AddressID,
		TotalQuantity:  Quantity,
		ShippingCharge: shippingCharge,
		OrderAmount:    float64(sum),
		CouponDiscount: int(coupDisc),
		OrderDate:      time.Now(),
	}

	if err := tx.Create(&order); err.Error != nil {
		tx.Rollback()
		utils.HandleError(ctx, http.StatusUnauthorized, "failed to place order")
		return
	}

	for _, v := range cart {
		off := controllers.OfferCalc(v.Product_ID)
		orderitems := models.OrderItems{
			OrderID:         order.ID,
			ProductID:       v.Product_ID,
			Quantity:        int(v.Quantity),
			SubTotal:        v.Product.Price * float64(v.Quantity),
			OfferPercentage: int(off),
		}
		if err := tx.Create(&orderitems); err.Error != nil {
			tx.Rollback()
			utils.HandleError(ctx, http.StatusUnauthorized, "failed to place order")
			fmt.Println("failed to place order items: ", err.Error)
			return
		}

		//stock managing
		convert, _ := strconv.ParseUint(v.Product.Quantity, 10, 32)
		convert -= uint64(v.Quantity)
		v.Product.Quantity = fmt.Sprint(convert)
		if err := initializers.DB.Save(&v.Product); err.Error != nil {
			utils.HandleError(ctx, http.StatusUnauthorized, "failed to update stock")
			return
		}
	}

	if err := initializers.DB.Where("user_id=?", userid).Delete(&models.Cart{}); err.Error != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "failed to delete order")
		return
	}

	if err := tx.Commit(); err.Error != nil {
		tx.Rollback()
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to commit transaction")
		return
	}
	if checkout.PaymentType != "UPI" {
		ctx.JSON(200, gin.H{
			"status":          "success",
			"message":         "Order placed successfully",
			"Grand total":     sum,
			"shipping_Charge": shippingCharge,
		})
	}

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
