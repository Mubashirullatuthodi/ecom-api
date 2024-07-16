package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
	"github.com/razorpay/razorpay-go"
)

func PaymentSubmission(orderID string, amount int) (string, error) {
	fmt.Println("paymentorderID:", orderID, "\npaymentamount:", amount)
	keyID := os.Getenv("RAZORPAY_ID")
	secretkey := os.Getenv("RAZORPAY_SECRET")
	fmt.Println("keyid:", keyID, "\nsecretkey:", secretkey)
	Client := razorpay.NewClient(keyID, secretkey)

	paymentDetails := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  orderID,
	}
	order, err := Client.Order.Create(paymentDetails, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}
	razorOrderID, ok := order["id"].(string)
	fmt.Println("razororderid:", razorOrderID)
	if !ok {
		return "", errors.New("failed to retrieve order ID from Razorpay response")
	}
	return razorOrderID, nil
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	secretKey := os.Getenv("RAZORPAY_SECRET")
	signature := sign
	secret := secretKey
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write([]byte(data)); err != nil {
		return fmt.Errorf("failed to write HMAC: %w", err)
	}

	expectedSignature := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(signature)) != 1 {
		return errors.New("payment verification failed")
	}
	return nil

}

func CreatePayment(ctx *gin.Context) {
	fmt.Println("-------------------------payment processing-----------------------")
	var PaymentDetails = make(map[string]string)
	var Payment models.Payment
	if err := ctx.ShouldBindJSON(&PaymentDetails); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid request")
		return
	}
	fmt.Printf("Received Payment Details: %+v\n", PaymentDetails)
	err := RazorPaymentVerification(PaymentDetails["signatureID"], PaymentDetails["order_Id"], PaymentDetails["paymentID"])
	if err != nil {
		fmt.Println("====>", err)
		return
	}

	fmt.Println("======", PaymentDetails["order_Id"])
	if err := initializers.DB.Where("ord_id = ?", PaymentDetails["order_Id"]).First(&Payment); err.Error != nil {
		utils.HandleError(ctx, http.StatusNotFound, "OrderID not found")
		return
	}

	Payment.PaymentID = PaymentDetails["paymentID"]
	Payment.PaymentStatus = "Done"
	if err := initializers.DB.Save(&Payment).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to update payment status")
		return
	}
	ctx.JSON(200, gin.H{"Message": "Payment Done"})
}
