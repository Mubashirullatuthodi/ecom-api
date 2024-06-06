package controllers

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

func InvoicePDF(ctx *gin.Context) {
	orderID := ctx.Param("ID")
	var order models.Order

	if err := initializers.DB.First(&order, "id=?", orderID).Error; err != nil {
		ctx.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	//Generate PDF
	pdf := generateInvoicePDF(order)
	ctx.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	ctx.Header("Content-Type", "aplication/pdf")
	ctx.Data(200, "application/pdf", pdf)
}

func generateInvoicePDF(order models.Order) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Order ID: "+order.OrderCode)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Customer Name: "+order.User.FirstName)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Amount: "+fmt.Sprintf("$%.2f", order.TotalAmount))

	buf := new(bytes.Buffer)
	err := pdf.Output(buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
