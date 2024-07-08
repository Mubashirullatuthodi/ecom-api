package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
)

func AddProduct(ctx *gin.Context) {
	var Product models.Product
	var category models.Category

	file, _ := ctx.MultipartForm()

	categoryId, _ := strconv.Atoi(ctx.Request.FormValue("categoryID"))
	Product.CategoryID = uint(categoryId)
	if err := initializers.DB.First(&category, Product.CategoryID).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "no category found")
		return
	}

	Product.Name = ctx.Request.FormValue("name")
	Product.Quantity = ctx.Request.FormValue("quantity")
	Product.Description = ctx.Request.FormValue("description")
	Product.Price, _ = strconv.ParseFloat(ctx.Request.FormValue("price"), 64)
	images := file.File["images"]
	for _, img := range images {
		filePath := "./images/" + img.Filename
		if err := ctx.SaveUploadedFile(img, filePath); err != nil {
			utils.HandleError(ctx, http.StatusBadRequest, "failed to save image")
			return
		}
		Product.ImagePath = append(Product.ImagePath, filePath)
	}

	if err := initializers.DB.Create(&Product).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "failed to create product")
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "Success",
		"message": "Product Created succesfully",
	})
}

func ListProducts(ctx *gin.Context) {
	var listProduct []models.Product

	type list struct {
		ID           int      `json:"id"`
		Name         string   `json:"name"`
		Image        []string `json:"images"`
		Description  string   `json:"description"`
		Price        float64  `json:"price"`
		Quantity     string   `json:"quantity"`
		CategoryName string   `json:"categoryName"`
	}

	var List []list

	if err := initializers.DB.Preload("Category").Find(&listProduct).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to list product")
		return
	}

	for _, value := range listProduct {
		fmt.Println("image", value.ImagePath)
		listproduct := list{
			ID:           int(value.ID),
			Image:        value.ImagePath,
			Name:         value.Name,
			Description:  value.Description,
			Price:        value.Price,
			Quantity:     value.Quantity,
			CategoryName: value.Category.Name,
		}
		List = append(List, listproduct)
	}
	fmt.Println("list roducts: ", List)

	ctx.JSON(200, gin.H{
		"status":   "success",
		"Products": List,
	})
}

func EditProduct(ctx *gin.Context) {
	var Product models.Product

	id := ctx.Param("ID")

	if err := initializers.DB.First(&Product, id).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "product not found")
		return
	}

	contentType := ctx.GetHeader("Content-Type")

	switch contentType {
	case "application/json":
		if err := ctx.BindJSON(&Product); err != nil {
			utils.HandleError(ctx, http.StatusBadRequest, "failed to bind json")
			return
		}

		if err := initializers.DB.Model(&Product).Updates(Product).Error; err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "failed to edit product")
			return
		}

	case "multipart/form-data":

	default:
		utils.HandleError(ctx, http.StatusBadRequest, "unsupported content type")
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Product Edited Successfully",
	})
}

func ImageUpdate(ctx *gin.Context) {
	var Product models.Product

	id := ctx.Param("ID")

	if err := initializers.DB.First(&Product, id).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "product not found")
		return
	}

	if err := ctx.Request.ParseMultipartForm(0); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "failed to parse from data")
		return
	}

	images := ctx.Request.MultipartForm.File["images"]
	for _, img := range images {
		filepath := "./images/" + img.Filename
		if err := ctx.SaveUploadedFile(img, filepath); err != nil {
			utils.HandleError(ctx, http.StatusBadRequest, "failed dto save image")
			return
		}
		Product.ImagePath = append(Product.ImagePath, filepath)
		fmt.Println("new: ", Product.ImagePath)
	}
	if err := initializers.DB.Save(&Product).Error; err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to update product details")
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Product Edited Successfully",
	})
}

func DeleteProduct(ctx *gin.Context) {
	var product models.Product

	id := ctx.Param("ID")

	if err := initializers.DB.Where("ID = ?", id).First(&product).Error; err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "user not found")
		return
	} else {
		//soft delete
		if err:=initializers.DB.Delete(&product);err!=nil{
			utils.HandleError(ctx, http.StatusInternalServerError, "failed to delete the product")
			return
		}

		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "user delete succesfully",
		})
	}

}
