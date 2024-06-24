package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
)

type ReportRequest struct {
	OrderID         uint
	CustomerName    string
	ProductName     string
	ProductQuantity int
	OrderDate       string
	TotalAmount     float64
	CouponDeduction int
	OfferDiscount   int
	OrderStatus     string
	PaymentMethod   string
}

type Search struct {
	Type      string    `json:"type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func SalesReport(ctx *gin.Context) {
	var search Search
	if err := ctx.ShouldBindJSON(&search); err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to bind",
		})
		return
	}

	var sales []models.OrderItems

	if err := initializers.DB.
		Preload("Order").
		Preload("Product").
		Preload("Order.User").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN users ON orders.user_id = users.id").
		Where("order_items.order_status = ?", "Delivered").
		Where("order_items.deleted_at IS NULL").
		Find(&sales).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": "failed to fetch",
		})
		return
	}

	if len(sales) == 0 {
		ctx.JSON(200, gin.H{
			"status": "success",
			"Report": "No delivered orders Found",
		})
		return
	}

	//apending
	var newsales []ReportRequest
	var overallSales float64
	var overallDiscount float64

	//var grandTotal float64
	for _, details := range sales {
		formatDate := details.Order.CreatedAt.Format("2006-01-02 15:04:05")
		new := ReportRequest{
			OrderID:         details.OrderID,
			CustomerName:    details.Order.User.FirstName,
			ProductName:     details.Product.Name,
			ProductQuantity: details.Quantity,
			OrderDate:       formatDate,
			TotalAmount:     details.Product.Price * float64(details.Quantity),
			OrderStatus:     details.OrderStatus,
			CouponDeduction: int(details.Order.CouponDiscount),
			OfferDiscount:   details.OfferPercentage,
			PaymentMethod:   details.Order.PaymentMethod,
		}
		newsales = append(newsales, new)
		overallSales += details.Product.Price * float64(details.Quantity)
		overallDiscount = float64(details.OfferPercentage)*float64(details.Quantity) + float64(details.Order.CouponDiscount)
	}

	var recentSales []ReportRequest
	threshold := time.Now()

	switch search.Type {
	case "Daily":
		threshold = time.Now().Add(-24 * time.Hour)
	case "Weekly":
		threshold = time.Now().Add(-7 * 24 * time.Hour)
	case "Monthly":
		threshold = time.Now().Add(-30 * 24 * time.Hour)
	}

	var grandTotal float64

	for _, sale := range newsales {
		orderTime, err := time.Parse("2006-01-02 15:04:05", sale.OrderDate)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to parse order date",
			})
			return
		}

		if orderTime.After(threshold) {
			recentSales = append(recentSales, sale)
			grandTotal += sale.TotalAmount
		}

	}
	if len(recentSales) == 0 {
		ctx.JSON(200, gin.H{
			"status": "success",
			"Report": "No orders found in the specified period",
		})
		return
	}

	GeneratePDF(recentSales, grandTotal, overallSales, overallDiscount, ctx)
}

func GeneratePDF(newsales []ReportRequest, grandTotal, overallSales, overallDiscount float64, ctx *gin.Context) {
	//generate pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	//Title
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(20)

	//table header
	pdf.SetFont("Arial", "B", 12)
	tableWidth := 200.0
	pageWidth := 210.0
	margin := (pageWidth - tableWidth) / 2

	pdf.SetX(margin)

	columnWidth := []float64{20, 30, 20, 40, 30, 35, 27}

	pdf.CellFormat(columnWidth[0], 10, "Order ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[1], 10, "Product Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[2], 10, "Quantity", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[3], 10, "Order Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[4], 10, "Total Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[5], 10, "Payment Method", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[6], 10, "Order Status", "1", 1, "C", false, 0, "")

	//Table Body
	pdf.SetFont("Arial", "", 12)
	for _, sale := range newsales {
		pdf.SetX(margin)
		pdf.CellFormat(columnWidth[0], 10, strconv.Itoa(int(sale.OrderID)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[1], 10, sale.ProductName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[2], 10, fmt.Sprintf("%d", sale.ProductQuantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[3], 10, sale.OrderDate, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[4], 10, fmt.Sprintf("%.2f", sale.TotalAmount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[5], 10, sale.PaymentMethod, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[6], 10, sale.OrderStatus, "1", 1, "C", false, 0, "")
	}

	// total
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2]+columnWidth[3]+columnWidth[4]+columnWidth[5], 10, "Total:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(27, 10, fmt.Sprintf("%.2f", grandTotal), "1", 1, "C", false, 0, "")

	//Discount
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2]+columnWidth[3]+columnWidth[4]+columnWidth[5], 10, "Overall Discount:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(27, 10, fmt.Sprintf("%.2f", overallDiscount), "1", 1, "C", false, 0, "")

	//grand total
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2]+columnWidth[3]+columnWidth[4]+columnWidth[5], 10, "Grand Total:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(27, 10, fmt.Sprintf("%.2f", overallSales-overallDiscount), "1", 1, "C", false, 0, "")

	//overall sales,order amount
	pdf.Ln(10)
	pdf.SetX(margin)
	pdf.Cell(40, 10, "Overall Sales Summary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.SetX(margin)
	pdf.CellFormat(100, 10, "Overall Sales Amount:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", overallSales), "0", 1, "L", false, 0, "")

	pdf.SetX(margin)
	pdf.CellFormat(100, 10, "Overall Discount:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", overallDiscount), "0", 1, "L", false, 0, "")

	//write pdf to file
	// err := pdf.OutputFileAndClose("sales_report.pdf")
	// if err != nil {
	// 	ctx.JSON(500, gin.H{
	// 		"error": "Failed to Generate PDF",
	// 	})
	// 	return
	// }

	path := fmt.Sprintf("C:/Users/shanm/Desktop/pdf/salesReport_%s_%s.pdf", time.Now().Format("20060102_150405"), "sales")
	if err := pdf.OutputFileAndClose(path); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    401,
			"message": "failed to generate pdf",
		})
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment: filename=%s", path))
	ctx.Writer.Header().Set("Content-Type", "application/pdf")
	ctx.File(path)

	//send PDF response
	ctx.FileAttachment("sales_report.pdf", "sales_report.pdf")
}
