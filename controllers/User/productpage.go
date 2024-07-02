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
		Price float64
		Stock string
	}

	var List []productlist

	if err := initializers.DB.Find(&Product).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  "Cant find products",
			"code":   500,
		})
		return
	}

	for _, value := range Product {
		list := productlist{
			ID:    int(value.ID),
			Name:  value.Name,
			Price: value.Price,
			Stock: value.Quantity,
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
	var listProduct []models.Product

	type list struct {
		ID                  int      `json:"id"`
		Name                string   `json:"name"`
		Image               []string `json:"images"`
		Description         string   `json:"description"`
		Price               float64  `json:"price"`
		Quantity            string   `json:"quantity"`
		CategoryName        string   `json:"categoryName"`
		CategoryDescription string   `json:"categoryDescription"`
	}

	var List []list

	id := ctx.Param("ID")

	if err := initializers.DB.Preload("Category").Where("id=?", id).Find(&listProduct).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "fail",
			"error":  "failed to list products",
			"code":   500,
		})
		return
	}

	for _, value := range listProduct {
		fmt.Println("image", value.ImagePath)
		listproduct := list{
			ID:                  int(value.ID),
			Image:               value.ImagePath,
			Name:                value.Name,
			Description:         value.Description,
			Price:               value.Price,
			Quantity:            value.Quantity,
			CategoryName:        value.Category.Name,
			CategoryDescription: value.Category.Description,
		}
		List = append(List, listproduct)
	}
	fmt.Println("list roducts: ", List)

	ctx.JSON(200, gin.H{
		"status":   "success",
		"Products": List,
	})

}
