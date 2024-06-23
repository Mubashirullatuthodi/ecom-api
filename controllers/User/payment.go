package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/razorpay/razorpay-go"
)

func PaymentSubmission(orderid string, amount int) (string, error) {
	fmt.Println("paymentorderID---------------------->", orderid, "paymentamount-------------->", amount)
	keyid := os.Getenv("RAZORPAY_ID")
	secretkey := os.Getenv("RAZORPAY_SECRET")
	fmt.Println("keyid--------->", keyid, "secretkey------------------->", secretkey)
	Client := razorpay.NewClient(keyid, secretkey)

	paymentDetails := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  orderid,
	}
	order, err := Client.Order.Create(paymentDetails, nil)
	if err != nil {
		return "", err
	}
	razororderID, _ := order["id"].(string)
	fmt.Println("razororderid---------------->", razororderID)
	return razororderID, nil
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	secretKey := os.Getenv("RAZORPAY_SECRET")
	signature := sign
	secret := secretKey
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("PAYMENT FAILED")
	} else {
		return nil
	}
}

func CreatePayment(ctx *gin.Context) {
	fmt.Println("-------------------------its in payment-----------------------")
	var Paymentdetails = make(map[string]string)
	var Payment models.Payment
	if err := ctx.ShouldBindJSON(&Paymentdetails); err != nil {
		ctx.JSON(500, gin.H{"Error": "Invalid Request"})
	}
	fmt.Println("====>", Paymentdetails)
	err := RazorPaymentVerification(Paymentdetails["signatureID"], Paymentdetails["order_Id"], Paymentdetails["paymentID"])
	if err != nil {
		fmt.Println("====>", err)
		return
	}

	fmt.Println("======", Paymentdetails["order_Id"])
	if err := initializers.DB.Where("ord_id = ?", Paymentdetails["order_Id"]).First(&Payment); err.Error != nil {
		ctx.JSON(500, gin.H{"Error": "OrderID not found"})
		return
	}
	fmt.Println("-------", Payment)
	Payment.PaymentID = Paymentdetails["paymentID"]
	Payment.PaymentStatus = "Done"
	if err := initializers.DB.Model(&Payment).Updates(&models.Payment{
		PaymentID:     Payment.PaymentID,
		PaymentStatus: Payment.PaymentStatus,
	}); err.Error != nil {
		ctx.JSON(500, gin.H{"Error": "Failed to update paymentID"})
	} else {
		ctx.JSON(200, gin.H{"Message": "Payment Done"})
	}
}
