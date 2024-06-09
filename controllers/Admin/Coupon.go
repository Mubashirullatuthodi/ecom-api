package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func CreateCoupon(ctx *gin.Context) {
	var coupon models.Coupons

	if err := ctx.ShouldBind(&coupon); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Failed to Bind",
		})
		return
	}

	if err := initializers.DB.Create(&coupon).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  "Coupon already exist",
			"Code":   500,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "Success",
		"message": "Added the coupon",
	})
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
		Code        string
		Discount    float64
		Start_Date  string
		Expiry_date string
	}

	var list []show

	for _, v := range listCoupon {
		//Format Date
		startdate := v.Start_Date.Format("2006-01-02 15:04:05")
		enddate := v.Expiry_date.Format("2006-01-02 15:04:05")
		List := show{
			Code:        v.Code,
			Discount:    v.Discount,
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

}
