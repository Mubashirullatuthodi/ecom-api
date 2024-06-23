package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

type report struct {
	Quantity        int
	ItemDescription string
	UnitPrice       float64
	TotalPrice      float64
}

func GenerateInvoice(ctx *gin.Context) {
	var orderitem []models.OrderItems
	orderid := ctx.Param("ID")
	convid, _ := strconv.ParseUint(orderid, 32, 10)
	if err := initializers.DB.
		Preload("Order").
		Preload("Product").
		Preload("Order.User").
		Preload("Order.Address").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN users ON orders.user_id = users.id").
		Joins("JOIN addresses ON orders.address_id = addresses.id").
		Where("order_id=?", uint(convid)).
		Where("order_items.order_status = ?", "Delivered").
		Where("order_items.deleted_at IS NULL").
		Find(&orderitem).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": "failed to fetch",
		})
		return
	}

	var new []report
	var ShippingAddress string
	var orderDate string
	var orderNumber string

	for _, item := range orderitem {
		r := report{
			Quantity:        item.Quantity,
			ItemDescription: item.Product.Name,
			UnitPrice:       item.Product.Price,
			TotalPrice:      float64(item.Quantity) * item.Product.Price,
		}
		new = append(new, r)
		ShippingAddress = fmt.Sprintf(item.Order.Address.Address+"\n"+item.Order.Address.District+"\n"+item.Order.Address.State+"\n"+item.Order.Address.Pincode)
		orderNumber = item.Order.OrderCode
		orderDate = item.Order.CreatedAt.Format("2006-01-02")
	}

	GeneratePDF(new, ShippingAddress, orderNumber, orderDate, ctx)

	ctx.JSON(200, gin.H{
		"invoice": new,
	})
}

func GeneratePDF(new []report, shippingAddress, orderNumber, OrderDate string, ctx *gin.Context) {
	//generate pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	//Title
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(20)

	//default
	companyAddress := "new company\n1234 Street Address\nCity, State, ZIP\nCountry"
	//add address
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(100, 10, companyAddress, "", "E", false)
	pdf.Ln(10)

	// Add shipping address on the right
	pdf.SetY(30) // Set the y position to align with company address
	pdf.SetX(120)
	pdf.MultiCell(70, 10, "Shipping Address:\n"+shippingAddress, "", "L", false)
	pdf.Ln(20)

	//order number
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(100, 10, "Order Number:"+orderNumber, "", "L", false)
	pdf.Ln(10)

	//order number
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(100, 10, "Order Date:"+OrderDate, "", "L", false)
	pdf.Ln(10)

	//table header
	pdf.SetFont("Arial", "B", 12)
	tableWidth := 200.0
	pageWidth := 210.0
	margin := (pageWidth - tableWidth) / 2

	pdf.SetX(margin)

	columnWidth := []float64{20, 50, 30, 40}

	pdf.CellFormat(columnWidth[0], 10, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[1], 10, "Item Description", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[2], 10, "Unit Price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[3], 10, "Total Price", "1", 1, "C", false, 0, "")

	//Table Body
	pdf.SetFont("Arial", "", 12)
	for _, sale := range new {
		pdf.SetX(margin)
		pdf.CellFormat(columnWidth[0], 10, strconv.Itoa(int(sale.Quantity)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[1], 10, sale.ItemDescription, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[2], 10, fmt.Sprintf("%d", int(sale.UnitPrice)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[3], 10, fmt.Sprintf("%d", int(sale.TotalPrice)), "1", 1, "C", false, 0, "")
	}

	//write pdf to file
	err := pdf.OutputFileAndClose("invoice.pdf")
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Generate PDF",
		})
		return
	}

	//send PDF response
	ctx.FileAttachment("invoice.pdf", "invoice.pdf")
}
