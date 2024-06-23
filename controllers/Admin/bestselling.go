package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/initializers"
)

func BestSelling(ctx *gin.Context) {
	type BestProduct struct {
		ProductID     uint    `json:"product_id"`
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		TotalQuantity int     `json:"total_quantity"`
	}

	var bestProducts []BestProduct

	state := ctx.Query("type")

	switch state {
	case "product":
		query := `
		SELECT
			p.id AS product_id,
			p.name,
			p.description,
			p.price,
			SUM(oi.quantity) AS total_quantity
		FROM
			order_items oi
		JOIN
			products p
		ON
			oi.product_id = p.id
		GROUP BY
			p.id
		ORDER BY
			total_quantity DESC
		LIMIT 10
	`

		initializers.DB.Raw(query).Scan(&bestProducts)

	case "category":
		categoryID := ctx.Query("category_id")
		convID, _ := strconv.ParseUint(categoryID, 32, 10)
		fmt.Println("-----", categoryID)
		if categoryID == "" {
			ctx.JSON(400, gin.H{
				"error": "category id is required",
			})
			return
		}
		query := `
			SELECT
				p.id AS product_id,
				p.name,
				p.description,
				p.price,
				SUM(oi.quantity) AS total_quantity
			FROM
				order_items oi
			JOIN
				products p
			ON
				oi.product_id = p.id
			JOIN
				categories c
			ON
				p.category_id = c.id
			WHERE
				c.id = ?
			GROUP BY
				p.id 
			ORDER BY
				total_quantity DESC
			LIMIT 10
		`

		initializers.DB.Raw(query, uint(convID)).Scan(&bestProducts)
	}

	ctx.JSON(http.StatusOK, bestProducts)

}
