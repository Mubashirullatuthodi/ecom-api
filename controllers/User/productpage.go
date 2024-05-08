package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func ProductPage(ctx *gin.Context) {
	var Product []models.Product

	type productlist struct {
		ID    int
		Name  string
		Price string
	}

	var List []productlist

	if err := initializers.DB.Find(&Product).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  err.Error(),
			"code":   500,
		})
		return
	}

	for _, value := range Product {
		list := productlist{
			ID:    int(value.ID),
			Name:  value.Name,
			Price: value.Price,
		}
		List = append(List, list)
	}

	fmt.Println("list", List)

	ctx.JSON(200, gin.H{
		"status":   "success",
		"products": List,
	})
}

func ProductDetail(ctx *gin.Context) {
	var Product models.Product

	id := ctx.Param("ID")

	if err := initializers.DB.Preload("Category").First(&Product, id).Error; err != nil {
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "product not found",
			"code":   404,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "Success",
		"product": Product,
	})
}
