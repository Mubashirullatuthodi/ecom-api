package controllers

import (
	//"io/ioutil"
	"fmt"
	"net/http"

	//"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func CreateCategory(ctx *gin.Context) {
	var category models.Category

	if err := ctx.BindJSON(&category); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind category")
		return
	}

	insert := initializers.DB.Create(&category)
	if insert.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to insert category")
		return
	}
	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "category created successfully",
	})

}

func GetCategory(ctx *gin.Context) {
	var listCategory []models.Category

	type List struct {
		ID          int
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var list []List

	if err := initializers.DB.Find(&listCategory).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to list category")
		return
	}

	for _, value := range listCategory {
		category := List{
			ID:          int(value.ID),
			Name:        value.Name,
			Description: value.Description,
		}

		list = append(list, category)
	}
	fmt.Println("list category: ", list)
	ctx.JSON(http.StatusOK, list)
}

func UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("ID")

	var category models.Category

	if err := initializers.DB.First(&category, id).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "category not found")
		return
	}

	var UpdateCategory models.Category
	if err := ctx.BindJSON(&UpdateCategory); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to bind")
		return
	}

	category.Name = UpdateCategory.Name
	category.Description = UpdateCategory.Description

	if err := initializers.DB.Save(&category).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to update category")
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "category updated successfully",
	})
}

func DeleteCategory(ctx *gin.Context) {
	var DeleteCategory models.Category

	id := ctx.Param("ID")
	err := initializers.DB.First(&DeleteCategory, id)
	if err.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "category not found")
		return
	}
	err = initializers.DB.Delete(&DeleteCategory)
	if err.Error != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to delete category")
		return
	}
	ctx.JSON(200, gin.H{
		"status": "Success",
		"Error":  "Category Deleted Successfully",
	})
}
