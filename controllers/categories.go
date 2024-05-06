package controllers

import (
	//"io/ioutil"
	"fmt"
	"net/http"
	//"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func CreateCategory(ctx *gin.Context) {
	var category models.Category

	err := ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to bind category",
		})
		return
	}

	// var existingcategory models.Category
	// initializers.DB.Where("name=?",category.Name).First(&existingcategory)
	// if existingcategory.Name == category.Name{
	// 	ctx.JSON(http.StatusConflict,gin.H{
	// 		"error":"this category already exist",
	// 	})
	// 	return
	// }

	insert := initializers.DB.Create(&category)
	if insert.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create category",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "category created successfully",
	})

}

func GetCategory(ctx *gin.Context) {
	var listcategory []models.Category

	type List struct {
		ID          int
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var list []List

	if err := initializers.DB.Find(&listcategory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list category",
		})
		return
	}

	for _, value := range listcategory {
		category := List{
			ID:          int(value.ID),
			Name:        value.Name,
			Description: value.Description,
		}

		list = append(list, category)
	}
	fmt.Println("list category: ",list)
	ctx.JSON(http.StatusOK,list)
}

func UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("ID")

	var category models.Category

	if err := initializers.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "category not found",
		})
		return
	}

	var UpdateCategory models.Category
	if err := ctx.BindJSON(&UpdateCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category.Name = UpdateCategory.Name
	category.Description = UpdateCategory.Description

	if err := initializers.DB.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update category",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "category updated successfully",
	})
}

func DeleteCategory(ctx *gin.Context){
	
}

// func CreateProduct(ctx *gin.Context) {
// 	var product models.Product

// 	if err := ctx.BindJSON(&product); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	for _, image := range product.Images {

// 		imagePath := filepath.Join("images", image.Filename)

// 		_, err := ioutil.ReadFile(imagePath)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}
// 		image.URL = imagePath
// 	}
// 	ctx.JSON(http.StatusCreated, product)
// }
