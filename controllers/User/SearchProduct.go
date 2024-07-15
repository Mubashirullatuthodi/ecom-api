package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

const (
	PriceLowToHigh = "price_low_to_high"
	PriceHighToLow = "price_high_to_low"
	NewArrivals    = "new_arrivals"
	AToZ           = "a_to_z"
	ZToA           = "z_to_a"
	Popularity     = "popularity"
)

func FetchProducts(orderQuery string) ([]gin.H, error) {
	var products []models.Product
	result := initializers.DB.Order(orderQuery).Joins("Category").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	var prices []gin.H
	for _, v := range products {
		prices = append(prices, gin.H{
			"name":     v.Name,
			"price":    v.Price,
			"category": v.Category.Name,
			"ID":       v.ID,
		})
	}
	return prices, nil
}

func FetchPopularProducts() ([]gin.H, error) {
	var products []models.Product
	query := `SELECT * FROM products
		JOIN ( 
			SELECT product_id, SUM(quantity) AS total_quantity
			FROM order_items
			GROUP BY product_id
			ORDER BY total_quantity DESC
			LIMIT 10
		) AS o ON products.id = o.product_id
		WHERE products.deleted_at IS NULL
		ORDER BY o.total_quantity DESC`
	initializers.DB.Raw(query).Scan(&products)

	var prices []gin.H
	for _, v := range products {
		prices = append(prices, gin.H{
			"name":     v.Name,
			"price":    v.Price,
			"category": v.Category.Name,
			"ID":       v.ID,
		})
	}

	return prices, nil
}

func SearchProduct(ctx *gin.Context) {
	search := ctx.Query("search")

	if search == "" {
		utils.HandleError(ctx, http.StatusBadRequest, "Please enter a search term.")
		return
	}

	var prices []gin.H
	var err error

	switch search {

	case PriceLowToHigh:
		prices, err = FetchProducts("price ASC")

	case PriceHighToLow:
		prices, err = FetchProducts("price DESC")

	case NewArrivals:

		prices, err = FetchProducts("created_at DESC")

	case AToZ:
		prices, err = FetchProducts("name ASC")

	case ZToA:
		prices, err = FetchProducts("name DESC")

	case Popularity:
		prices, err = FetchPopularProducts()

	default:
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid search type.")
		return

	}
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "Product not found.")
		return
	}

	ctx.JSON(200, prices)
}
