package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func SearchProduct(ctx *gin.Context) {
	search := ctx.Query("search")

	if search == "" {
		utils.HandleError(ctx, http.StatusBadRequest, "enter any letter")
		return
	}

	var prices []gin.H

	switch search {

	case "price_low_to_high":
		var products []models.Product
		result := initializers.DB.Order("price asc").Joins("Category").Find(&products)
		if result.Error != nil {
			utils.HandleError(ctx, http.StatusNotFound, "not found")
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
			utils.HandleError(ctx, http.StatusNotFound, "not found")
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
			utils.HandleError(ctx, http.StatusNotFound, "not found")
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
			utils.HandleError(ctx, http.StatusNotFound, "not found")
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
			utils.HandleError(ctx, http.StatusNotFound, "not found")
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

	case "popularity":
		var products []models.Product
		query := `SELECT *FROM products
		JOIN ( 
		SELECT product_id,SUM(quantity) as total_quantity
		FROM order_items
		GROUP BY product_id
		ORDER BY total_quantity DESC
		LIMIT 10)
		AS o ON products.id = o.product_id
		WHERE products.deleted_at IS NULL
		ORDER BY o.total_quantity DESC`
		initializers.DB.Raw(query).Scan(&products)

		for _, v := range products {
			prices = append(prices, gin.H{
				"name":     v.Name,
				"Price":    v.Price,
				"category": v.Category.Name,
				"ID":       v.ID,
			})
		}

	}

	ctx.JSON(200, prices)
}
