package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	controllers "github.com/mubashir/e-commerce/controllers/Admin"
	"github.com/mubashir/e-commerce/initializers"
	"github.com/mubashir/e-commerce/models"
	"github.com/mubashir/e-commerce/utils"
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to fetch")
		return
	}

	var new []report
	var ShippingAddress string
	var orderDate string
	var orderNumber string
	var SubTotal int
	var OverallDiscount float64
	var shippingcharge int
	var couponDiscount int

	for _, item := range orderitem {
		couponDiscount = item.Order.CouponDiscount
		shippingcharge += int(item.Order.ShippingCharge)
		SubTotal += item.Quantity * int(item.Product.Price)
		OverallDiscount += controllers.OfferCalc(item.ProductID) * float64(item.Quantity)

		r := report{
			Quantity:        item.Quantity,
			ItemDescription: item.Product.Name,
			UnitPrice:       item.Product.Price,
			TotalPrice:      float64(item.Quantity) * item.Product.Price,
		}
		new = append(new, r)
		ShippingAddress = fmt.Sprintf(item.Order.User.FirstName + "\n" + item.Order.Address.Address + "\n" + item.Order.Address.Town + "," + item.Order.Address.District + "," + item.Order.Address.Pincode + ",\n" + item.Order.Address.State + "\n+91 " + item.Order.User.Phone)
		orderNumber = item.Order.OrderCode
		orderDate = item.Order.CreatedAt.Format("2006-01-02")
	}

	GeneratePDF(new, ShippingAddress, orderNumber, orderDate, SubTotal, couponDiscount, shippingcharge, OverallDiscount, ctx)

	ctx.JSON(200, gin.H{
		"invoice": new,
	})
}

func GeneratePDF(new []report, shippingAddress, orderNumber, OrderDate string, subtotal, coupondiscount, shippingcharge int, overalldiscount float64, ctx *gin.Context) {
	fmt.Println("generating-------------------")
	//generate pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)

	//Title
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(20)

	//default
	companyAddress := "From:\nE-Commerce\necommerce@gmail.com\n1234 JP Nagar\nBangalore, Karnataka, 589625\nIndia"
	//add address
	pdf.SetFont("Arial", "", 16)
	pdf.MultiCell(100, 10, companyAddress, "", "E", false)
	pdf.Ln(10)

	// Add shipping address on the right
	pdf.SetY(30) // Set the y position to align with company address
	pdf.SetX(120)
	pdf.MultiCell(70, 10, "Shipping Address:\n"+shippingAddress+"\nIndia", "", "L", false)
	pdf.Ln(20)

	//order number
	pdf.SetFont("Arial", "", 15)
	pdf.MultiCell(100, 5, "Order Number:"+orderNumber, "", "L", false)
	pdf.Ln(5)

	//order number
	pdf.SetFont("Arial", "", 15)
	pdf.MultiCell(100, 5, "Order Date:"+OrderDate, "", "L", false)
	pdf.Ln(20)

	//table header
	pdf.SetFont("Arial", "B", 16)
	tableWidth := 200.0
	pageWidth := 210.0
	margin := (pageWidth - tableWidth) / 2

	pdf.SetX(margin)

	columnWidth := []float64{20, 60, 60, 60}

	pdf.CellFormat(columnWidth[0], 10, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[1], 10, "Item Description", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[2], 10, "Unit Price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(columnWidth[3], 10, "Total Price", "1", 1, "C", false, 0, "")

	//Table Body
	pdf.SetFont("Arial", "", 16)
	for _, sale := range new {
		pdf.SetX(margin)
		pdf.CellFormat(columnWidth[0], 10, strconv.Itoa(int(sale.Quantity)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[1], 10, sale.ItemDescription, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[2], 10, fmt.Sprintf("%d", int(sale.UnitPrice)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth[3], 10, fmt.Sprintf("%.2f", float64(sale.TotalPrice)), "1", 1, "C", false, 0, "")
	}

	// total
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2], 10, "Total:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("%.2f", float64(subtotal)), "1", 1, "C", false, 0, "")

	//Discount
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2], 10, "Overall Discount:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("-%.2f", float64(overalldiscount)+float64(coupondiscount)), "1", 1, "C", false, 0, "")

	//shiping
	//Discount
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2], 10, "Shipping Charge:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("+%.2f", float64(shippingcharge)), "1", 1, "C", false, 0, "")

	//grand total
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(margin)
	pdf.CellFormat(columnWidth[0]+columnWidth[1]+columnWidth[2], 10, "Grand Total:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("%.2f /-", float64((subtotal-int(overalldiscount)-coupondiscount)+shippingcharge)), "1", 1, "C", false, 0, "")

	//generate pdf file
	path := fmt.Sprintf("C:/Users/shanm/Desktop/pdf/invoice/invoice_%s_%s.pdf", time.Now().Format("20060102_150405"), "invoice")
	if err := pdf.OutputFileAndClose(path); err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "failed to generate PDF")
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment: filename=%s", path))
	ctx.Writer.Header().Set("Content-Type", "application/pdf")
	ctx.File(path)

	//send PDF response
	ctx.FileAttachment("invoice.pdf", "invoice.pdf")
}
