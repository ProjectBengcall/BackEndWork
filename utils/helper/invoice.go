package helper

import (
	"bengcall/config"
	"bengcall/utils/database"
	loo "log"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/gommon/log"
	gomail "gopkg.in/gomail.v2"
)

type DataCore struct {
	ID         uint
	Fullname   string
	Email      string
	Phone      string
	Address    string
	Other      string
	Additional int
	Invoice    int
	Total      int
}

type ServiceCore struct {
	ID            uint
	Name_vehicle  string
	Service_name  string
	TransactionID int
	SubTotal      int
}

// var db gorm.DB
var tgl = time.Now().Format("01-02-2006")

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "PT. Bengcall jaya <klukmanul33@gmail.com>"

func Create(id uint) (string, error) {
	var resQry DataCore
	var prodQry ServiceCore
	cfg := config.NewConfig()
	db := database.InitDB(cfg)
	//var res []string
	if err := db.Table("transactions").Select("transactions.id", "transactions.phone", "transactions.address", "transactions.invoice", "transactions.total", "transactions.other", "transactions.additional", "users.fullname", "users.email").Joins("join users on users.id=transactions.user_id").Where("transactions.invoice = ?", id).Scan(&resQry).Error; err != nil {
		log.Error(err)
		return "", err
	}
	if err := db.Table("details").Select("details.id", "vehicles.name_vehicle", "services.service_name", "details.transaction_id", "details.sub_total").Joins("join vehicles on vehicles.id=details.vehicle_id").Joins("join services on services.id=details.service_id").Where("details.transaction_id = ?", resQry.Invoice).Scan(&prodQry).Error; err != nil {
		log.Error(err)
		return "", err
	}
	loo.Println(resQry)
	loo.Println(prodQry)
	//dat := ToCoreData(resQry)
	a := strconv.Itoa(int(resQry.Invoice))
	pdf, err := newReport(resQry.Fullname, a)
	pdf = header(pdf, []string{"vehicle", "Service", "Price"})

	//a := strconv.Itoa(int(prodQry.Id))
	c := strconv.Itoa(prodQry.SubTotal)
	d := strconv.Itoa(resQry.Additional)
	e := strconv.Itoa(resQry.Total)
	pdf = table(pdf, []string{prodQry.Name_vehicle, prodQry.Service_name, c}, resQry.Other, d, e)
	pdf = image(pdf)
	if pdf.Err() {
		log.Error("Failed creating PDF report: %s\n", pdf.Error())
	}

	err = pdf.OutputFileAndClose("ini-invoice.pdf")
	if err != nil {
		log.Error("close file:", err)
	}
	imageFalse, err := os.Open("ini-invoice.pdf")
	if err != nil {
		log.Error("image open:", err)
	}
	imageFalseCnv := &multipart.FileHeader{
		Filename: imageFalse.Name(),
	}
	res, err := UploadFile(imageFalse, imageFalseCnv)
	log.Print(res)

	sendEmail(res, resQry.Email, resQry.Fullname)

	imageFalse.Close()
	DeleteFile()

	return res, nil
}

// Set bagian atas invoice
func newReport(fullname string, invoice string) (*gofpdf.Fpdf, error) {
	//set page invoice
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 16)
	pdf.Cell(16, 5, "Invoice Bengcall")
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Date :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, tgl) //date
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "No Invoice :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, invoice) //no Invoice
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Penerima :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, fullname) //name penerima
	pdf.Ln(20)

	return pdf, nil
}

func header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		pdf.CellFormat(40, 7, str, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

func table(pdf *gofpdf.Fpdf, data []string, other string, add string, total string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(240, 240, 240)

	// Every column gets aligned according to its contents.
	//align := []string{"C", "L", "C"}
	for _, line := range data { //input.detail
		pdf.CellFormat(40, 7, line, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "other :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, other) //other
	pdf.Ln(8)
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Price add :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, add) //add
	pdf.Ln(8)
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Total :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, total) //total
	pdf.Ln(8)
	//text dibawah tabel
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Terima Kasih Telah Mempercayai Kami")
	pdf.Ln(20)
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "PT. Bengcall Jaya")

	return pdf
}

func image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.ImageOptions("image/bc.png", 170, 10, 25, 25, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func sendEmail(file string, email string, fullname string) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", "lh6468072@gmail.com")
	mailer.SetAddressHeader("Cc", email, fullname)
	mailer.SetHeader("Subject", "Invoice Bengcall")
	mailer.SetBody("text/html", file)
	//mailer.Attach(file)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "klukmanul33@gmail.com", "tdmznylftqkvqcuv")

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Error(err.Error())
	}

	loo.Println("Mail sent!")
}

func DeleteFile() {
	err := os.Remove("ini-invoice.pdf")
	if err != nil {
		log.Print(err)
	}

}
