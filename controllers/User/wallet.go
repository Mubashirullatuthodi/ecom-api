package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func GetWallet(ctx *gin.Context) {
	userid := ctx.GetUint("userid")

	var balanceSum float64 = 0

	if err := initializers.DB.Model(&models.Wallet{}).Where("user_id=?", userid).Select("COALESCE(SUM(balance),0)").Row().Scan(&balanceSum); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to retrieve wallet balance")
		return
	}

	ctx.JSON(200, gin.H{
		"wallet": balanceSum,
	})
}

func WalletHistory(ctx *gin.Context) {
	var walletHistory []models.Wallet
	userid := ctx.GetUint("userid")

	if err := initializers.DB.Where("user_id = ?", userid).Find(&walletHistory).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to Find wallet history")
		return
	}

	var History []gin.H

	for _, v := range walletHistory {
		fmt.Println("--------->", v)
		formatted := v.CreatedAt.Format("2006-01-02 15:04:05")
		History = append(History, gin.H{
			"ID":               v.ID,
			"Balance":          v.Balance,
			"userID":           v.UserID,
			"transaction_time": formatted,
		})
	}

	ctx.JSON(200, History)
}
