package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func SearchProduct(ctx *gin.Context) {
	search := ctx.Query("search")
	// var products []models.Product

	if search == "" {
		ctx.JSON(400, gin.H{
			"error": "enter any letter",
		})
		return
	}

	var prices []gin.H

	switch search {

	case "price_low_to_high":
		var products []models.Product
		result := initializers.DB.Order("price asc").Joins("Category").Find(&products)
		if result.Error != nil {
			ctx.JSON(500, gin.H{
				"error": "not found",
			})
			return
		}

		for _, v := range products {
			prices = append(prices, gin.H{
				"name":     v.Name,
				"price":    v.Price,
				"category": v.Category.Name,
				"ID":       v.ID,
			})
		}
		fmt.Println("=============================", prices)

	case "price_high_to_low":
		var products []models.Product
		result := initializers.DB.Order("price DESC").Joins("Category").Find(&products)
		if result.Error != nil {
			ctx.JSON(500, gin.H{
				"error": "not found",
			})
			return
		}
		for _, v := range products {
			prices = append(prices, gin.H{
				"name":     v.Name,
				"price":    v.Price,
				"category": v.Category.Name,
				"ID":       v.ID,
			})
		}
		fmt.Println("=============================", prices)

	case "new_arrivals":
		var products []models.Product
		result := initializers.DB.Order("created_at desc").Joins("Category").Find(&products)
		if result.Error != nil {
			ctx.JSON(500, gin.H{
				"error": "not found",
			})
			return
		}
		for _, v := range products {
			prices = append(prices, gin.H{
				"name":     v.Name,
				"price":    v.Price,
				"category": v.Category.Name,
				"ID":       v.ID,
			})
		}
		fmt.Println("=============================", prices)

	case "a_to_z":
		var products []models.Product
		result := initializers.DB.Order("name asc").Joins("Category").Find(&products)
		if result.Error != nil {
			ctx.JSON(500, gin.H{
				"error": "not found",
			})
			return
		}
		for _, v := range products {
			prices = append(prices, gin.H{
				"name":     v.Name,
				"price":    v.Price,
				"category": v.Category.Name,
				"ID":       v.ID,
			})
		}
		fmt.Println("=============================", prices)

	case "z_to_a":
		var products []models.Product
		result := initializers.DB.Order("name desc").Joins("Category").Find(&products)
		if result.Error != nil {
			ctx.JSON(500, gin.H{
				"error": "not found",
			})
			return
		}
		for _, v := range products {
			prices = append(prices, gin.H{
				"name":     v.Name,
				"price":    v.Price,
				"category": v.Category.Name,
				"ID":       v.ID,
			})
		}
		fmt.Println("=============================", prices)
	}

	ctx.JSON(200, prices)
}
