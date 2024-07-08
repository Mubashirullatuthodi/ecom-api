package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/models"
)

func HandleError(ctx *gin.Context, StatusCode int, message string) {
	responseData := models.ErrorHandler{
		Status:     "fail",
		StatusCode: StatusCode,
		Error:      message,
	}
	ctx.JSON(StatusCode, responseData)
}
