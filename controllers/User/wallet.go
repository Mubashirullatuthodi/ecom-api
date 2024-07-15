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
	userID := ctx.GetUint("userid")

	var balanceSum float64

	err := initializers.DB.Model(&models.Wallet{}).Where("user_id=?", userID).Select("COALESCE(SUM(balance),0)").Row().Scan(&balanceSum)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to retrieve wallet balance")
		return
	}

	ctx.JSON(200, gin.H{
		"wallet": formatCurrency(balanceSum),
	})
}

func WalletHistory(ctx *gin.Context) {
	var walletHistory []models.Wallet
	userID := ctx.GetUint("userid")

	err := initializers.DB.Where("user_id = ?", userID).Find(&walletHistory).Error
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to find wallet history")
		return
	}

	history := make([]gin.H, len(walletHistory))

	for i, v := range walletHistory {
		history[i] = gin.H{
			"ID":              v.ID,
			"Balance":         v.Balance,
			"transactionTime": v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f rs", amount)
}
