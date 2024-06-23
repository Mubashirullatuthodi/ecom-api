package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func GetFilteredOrders(ctx *gin.Context) {
	filter := ctx.Query("filter")
	var orders []models.Order

	switch filter {
	case "yearly":
		initializers.DB.Where("EXTRACT(YEAR FROM created_at)=?", time.Now().Year()).Find(&orders)
	case "monthly":
		initializers.DB.Where("EXTRACT(MONTH FROM created_at)=? AND EXTRACT(YEAR FROM created_at)=?", time.Now().Month(), time.Now().Year()).Find(&orders)
	default:
		initializers.DB.Find(&orders)
	}

	ctx.JSON(200, orders)
}
