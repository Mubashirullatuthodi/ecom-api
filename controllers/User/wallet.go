package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

type walletResponse struct {
	Balance float64 `json:"balance"`
}

func GetWallet(ctx *gin.Context) {
	var wallet []models.Wallet
	userid := ctx.GetUint("userid")

	if err := initializers.DB.Where("user_id=?", userid).Find(&wallet).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "failed to fetch wallet",
		})
		return
	}
	var balance walletResponse
	s := 0.0
	for _, v := range wallet {
		s += v.Balance
		balance.Balance=s
	}

	ctx.JSON(200, gin.H{
		"wallet": balance,
	})
}
