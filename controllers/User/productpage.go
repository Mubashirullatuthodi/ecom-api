package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

type ProductList struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock string  `json:"stock"`
}

type productDetails struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Image               []string `json:"images"`
	Description         string   `json:"description"`
	Price               float64  `json:"price"`
	Quantity            string   `json:"quantity"`
	CategoryName        string   `json:"categoryName"`
	CategoryDescription string   `json:"categoryDescription"`
}

func ProductPage(ctx *gin.Context) {
	var Products []models.Product
	var productList []ProductList

	err := initializers.DB.Select("ID,Name,Price,Quantity").Find(&Products).Error
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "can't find products")
		return
	}

	for _, product := range Products {
		productList = append(productList, ProductList{
			ID:    int(product.ID),
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Quantity,
		})
	}

	ctx.JSON(200, gin.H{
		"status":   "success",
		"products": productList,
	})
}

func ProductDetail(ctx *gin.Context) {
	var Product models.Product

	id := ctx.Param("ID")
	categoryID, _ := strconv.ParseUint(id, 10, 32)

	err := initializers.DB.Preload("Category").Where("id=?", uint(categoryID)).Find(&Product).Error
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "failed to list product")
		return
	}
	productDetail := productDetails{
		ID:                  int(Product.ID),
		Name:                Product.Name,
		Image:               Product.ImagePath,
		Description:         Product.Description,
		Price:               Product.Price,
		Quantity:            Product.Quantity,
		CategoryName:        Product.Category.Name,
		CategoryDescription: Product.Category.Description,
	}

	ctx.JSON(200, gin.H{
		"status":   "success",
		"Products": productDetail,
	})

}
