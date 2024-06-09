package controllers

import (
	"fmt"
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
	OrderDate       string
	TotalAmount     float64
	CouponDeduction int
	Discount        int
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
		Preload("Order.Coupons").
		Preload("Order.User").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN users ON orders.user_id = users.id").
		Joins("JOIN coupons ON orders.coupon_code = coupons.code").
		Where("orders.order_status = ?", "Delivered").
		Where("order_items.product_order_status = ?", "Pending").
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

	//var grandTotal float64
	for _, details := range sales {
		formatDate := details.Order.CreatedAt.Format("2006-01-02 15:04:05")
		new := ReportRequest{
			OrderID:         details.OrderID,
			CustomerName:    details.Order.User.FirstName,
			ProductName:     details.Product.Name,
			OrderDate:       formatDate,
			TotalAmount:     details.Product.Price*float64(details.Quantity) - details.Order.Coupons.Discount,
			OrderStatus:     details.Order.OrderStatus,
			CouponDeduction: int(details.Order.Coupons.Discount),
			PaymentMethod:   details.Order.PaymentMethod,
		}
		// if details.Order.CouponCode != "" {
		// 	new.CouponDeduction = int(details.Order.Coupons.Discount)
		// }
		newsales = append(newsales, new)
		//grandTotal += details.Order.TotalAmount
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

	GeneratePDF(recentSales, grandTotal, ctx)
}

func GeneratePDF(newsales []ReportRequest, grandTotal float64, ctx *gin.Context) {
	//generate pdf
	pdf := gofpdf.New("P", "mm", "A3", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	//Title
	pdf.Cell(40, 10, "Sales Report")
	pdf.Ln(20)

	//table header
	pdf.SetFont("Arial", "B", 12)
	tableWidth := 200.0
	pageWidth := 210.0
	margin := (pageWidth - tableWidth) / 2

	pdf.SetX(margin)

	columnWidth := []float64{20, 40, 40, 40, 30, 40, 40,30}

	pdf.CellFormat(columnWidth[0], 10, "Order ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[1], 10, "Customer Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[2], 10, "Product Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[3], 10, "Order Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[4], 10, "Total Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[5], 10, "Coupon discount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[6], 10, "Payment Method", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[7], 10, "Order Status", "1", 1, "C", false, 0, "")

	//Table Body
	pdf.SetFont("Arial", "", 12)
	for _, sale := range newsales {
		pdf.SetX(margin)
		pdf.CellFormat(columnWidth[0], 10, strconv.Itoa(int(sale.OrderID)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[1], 10, sale.CustomerName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[2], 10, sale.ProductName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[3], 10, sale.OrderDate, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[4], 10, fmt.Sprintf("%.2f", sale.TotalAmount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[5], 10, strconv.Itoa(sale.CouponDeduction), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[6], 10, sale.PaymentMethod, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[7], 10, sale.OrderStatus, "1", 1, "C", false, 0, "")
	}

	//grand total
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2]+columnWidth[3]+columnWidth[4]+columnWidth[5]+columnWidth[6], 10, "Grand Total", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", grandTotal), "1", 1, "C", false, 0, "")
	//pdf.CellFormat(30, 10, "", "1", 1, "C", false, 0, "")

	//write pdf to file
	err := pdf.OutputFileAndClose("sales_report.pdf")
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to Generate PDF",
		})
		return
	}

	//send PDF response
	ctx.FileAttachment("sales_report.pdf", "sales_report.pdf")
}
