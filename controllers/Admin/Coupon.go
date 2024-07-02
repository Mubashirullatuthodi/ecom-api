package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

type newcoupon struct {
	Code        string  `json:"code"`
	Discount    float64 `json:"discount"`
	Condition   int     `json:"condition"`
	Description string  `json:"descrition"`
	MaxUsage    int     `json:"maxUsage"`
	Start_Date  string  `json:"startDate"`
	Expiry_date string  `json:"expiryDate"`
}

func CreateCoupon(ctx *gin.Context) {
	var coupon newcoupon

	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		ctx.JSON(500, "Failed to Bind")
		return
	}

	startDate, err := time.Parse("2006-01-02", coupon.Start_Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid created date format"})
		return
	}
	fmt.Println("----------------------------->", startDate)
	endDate, err := time.Parse("2006-01-02", coupon.Expiry_date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid created date format"})
		return
	}

	if err := initializers.DB.Create(&models.Coupons{
		CouponCode:  coupon.Code,
		Discount:    coupon.Discount,
		Condition:   coupon.Condition,
		Description: coupon.Description,
		MaxUsage:    coupon.MaxUsage,
		Start_Date:  startDate,
		Expiry_date: endDate,
	}); err.Error != nil {
		ctx.JSON(401, gin.H{
			"error":  "Coupon Already exist",
			"status": 401,
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "New coupon added",
		})
	}
}

func ListCoupon(ctx *gin.Context) {
	var listCoupon []models.Coupons

	if err := initializers.DB.Find(&listCoupon).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  "Failed to find coupon details",
			"code":   500,
		})
		return
	}
	type show struct {
		ID          uint
		Code        string 
		Discount    float64 
		Condition   int
		Description string
		MaxUsage    int
		Start_Date  string
		Expiry_date string
	}

	var list []show

	for _, v := range listCoupon {
		//Format Date
		startdate := v.Start_Date.Format("2006-01-02 15:04:05")
		enddate := v.Expiry_date.Format("2006-01-02 15:04:05")
		List := show{
			ID:          v.ID,
			Code:        v.CouponCode,
			Discount:    v.Discount,
			Condition:   v.Condition,
			Description: v.Description,
			MaxUsage:    v.MaxUsage,
			Start_Date:  startdate,
			Expiry_date: enddate,
		}
		list = append(list, List)
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"Coupons": list,
	})
}

func DeleteCoupon(ctx *gin.Context) {
	var coupon models.Coupons
	couponId := ctx.Param("ID")
	if err := initializers.DB.Where("id=?", couponId).Delete(&coupon); err.Error != nil {
		ctx.JSON(400, gin.H{
			"error":  "coupon not found",
			"status": 400,
		})
		return
	}

	ctx.JSON(204, gin.H{
		"message": "Coupon Deleted",
		"status":  204,
	})
}
