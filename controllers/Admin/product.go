package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func AddProduct(ctx *gin.Context) {
	var Product models.Product
	var category models.Category

	file, _ := ctx.MultipartForm()

	categoryId, _ := strconv.Atoi(ctx.Request.FormValue("categoryID"))
	Product.CategoryID = uint(categoryId)
	if err := initializers.DB.First(&category, Product.CategoryID).Error; err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"error":  "no category found",
			"code":   400,
		})
		return
	}

	Product.Name = ctx.Request.FormValue("name")
	Product.Quantity = ctx.Request.FormValue("quantity")
	Product.Description = ctx.Request.FormValue("description")
	Product.Price = ctx.Request.FormValue("price")
	images := file.File["images"]
	for _, img := range images {
		filePath := "./images/" + img.Filename
		if err := ctx.SaveUploadedFile(img, filePath); err != nil {
			ctx.JSON(400, gin.H{
				"status":  "Fail",
				"message": "failed to save image",
				"code":    400,
			})
		}
		Product.ImagePath = append(Product.ImagePath, filePath)
	}

	if err := initializers.DB.Create(&Product).Error; err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "Failed To Create Products",
			"code":   400,
		})
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
		Price        string   `json:"price"`
		Quantity     string   `json:"quantity"`
		CategoryName string   `json:"category_name"`
	}

	var List []list

	if err := initializers.DB.Preload("Category").Find(&listProduct).Error; err != nil {
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
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "product not found",
			"code":   404,
		})
		return
	}

	contentType := ctx.GetHeader("Content-Type")

	switch contentType {
	case "application/json":
		if err := ctx.BindJSON(&Product); err != nil {
			ctx.JSON(400, gin.H{
				"status": "Fail",
				"Error":  "Failed to bind json",
				"code":   400,
			})
			return
		}

		if err := initializers.DB.Model(&Product).Updates(Product).Error; err != nil {
			ctx.JSON(500, gin.H{
				"status": "Fail",
				"Error":  "Failed To Edit Product",
				"code":   500,
			})
			return
		}

	case "multipart/form-data":

	default:
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "Unsupported content type",
			"code":   400,
		})
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
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "product not found",
			"code":   404,
		})
		return
	}

	if err := ctx.Request.ParseMultipartForm(0); err != nil {
		ctx.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "Failed to parse form data",
			"code":   400,
		})
		return
	}

	// file, err := ctx.MultipartForm()
	// if err != nil {
	// 	ctx.JSON(400, gin.H{
	// 		"status": "fail",
	// 		"Error":  "parsed to multipart form",
	// 		"code":   400,
	// 	})
	// 	return
	// }
	images := ctx.Request.MultipartForm.File["images"]
	for _, img := range images {
		filepath := "./images/" + img.Filename
		if err := ctx.SaveUploadedFile(img, filepath); err != nil {
			ctx.JSON(400, gin.H{
				"status":  "Fail",
				"message": "Failed to save image",
				"code":    400,
			})
		}
		Product.ImagePath = append(Product.ImagePath, filepath)
		fmt.Println("new: ", Product.ImagePath)
	}
	if err := initializers.DB.Save(&Product).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  "Failed to update product details",
			"code":   500,
		})
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
		ctx.JSON(404, gin.H{
			"status": "Fail",
			"Error":  "User not found",
			"code":   404,
		})
	} else {
		//soft delete
		initializers.DB.Delete(&product)

		ctx.JSON(200, gin.H{
			"status":  "success",
			"message": "user delete succesfully",
		})
	}

}
