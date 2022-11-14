package helper

import (
	"bengcall/config"
	"bengcall/utils/database"
	"context"
	loo "log"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/gommon/log"
	gomail "gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type DataCore struct {
	ID       uint
	Fullname string
	Email    string
	Phone    string
	Addres   string
	Other    string
	Invoice  int
	Total    int
}
type Data struct {
	gorm.Model
	Fullname string
	Email    string
	Phone    string
	Addres   string
	Other    string
	Invoice  int
	Total    int
}
type ServiceCore struct {
	ID            uint
	Vehicle_name  string
	Service_name  string
	TransactionID int
	SubTotal      int
}
type Service struct {
	gorm.Model
	Vehicle_name  string
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
	var resQry Data
	var prodQry Service
	cfg := config.NewConfig()
	db := database.InitDB(cfg)
	//var res []string
	if err := db.Table("transactions").Select("transactions.id", "transactions.phone", "transactions.address", "transactions.invoice", "transactions.total", "transactions.other", "transactions.status", "users.fullname", "users.email").Joins("join users on users.id=transactions.user_id").Where("transactions.id = ?", id).Scan(&resQry).Error; err != nil {
		log.Error(err)
		return "", err
	}
	if err := db.Table("details").Select("vehicles.name_vehicle", "services.service_name", "details.sub_total").Joins("join vehicles on vehicles.id=details.vehicle_id").Joins("join services on services.id=details.service_id").Where("details.transaction_id = ?", resQry.Invoice).Scan(&prodQry).Error; err != nil {
		log.Error(err)
		return "", err
	}
	//dat := ToCoreData(resQry)
	//ser := ToCoreService(prodQry)
	pdf, err := newReport(resQry.Fullname, resQry.Invoice)
	pdf = header(pdf, []string{"vehicle", "Service", "Price"})

	//a := strconv.Itoa(int(prodQry.Id))
	c := strconv.Itoa(prodQry.SubTotal)

	pdf = table(pdf, []string{prodQry.Vehicle_name, prodQry.Service_name, c})
	pdf = image(pdf)
	if pdf.Err() {
		log.Fatalf("Failed creating PDF report: %s\n", pdf.Error())
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
	res, err := uploadFile(imageFalse, imageFalseCnv)
	log.Print(res)

	sendEmail(res, resQry.Email, resQry.Fullname)

	imageFalse.Close()
	DeleteFile()

	return res, nil
}

// Set bagian atas invoice
func newReport(fullname string, invoice int) (*gofpdf.Fpdf, error) {
	//set page invoice
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 16)
	pdf.Cell(16, 5, "Invoice Bengcall")
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Time :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, tgl) //date
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "No Invoice :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, String(invoice)) //no Invoice
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
		pdf.CellFormat(30, 7, str, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

func table(pdf *gofpdf.Fpdf, data []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(240, 240, 240)

	// Every column gets aligned according to its contents.
	//align := []string{"C", "L", "C"}
	for _, line := range data { //input.detail
		pdf.CellFormat(30, 7, line, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
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

func Stringx(length int) string {
	return autoGeneratex(length, charset)
}

func autoGeneratex(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

const charsetx = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRandx *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func uploadFile(file multipart.File, fileheader *multipart.FileHeader) (string, error) {

	randomStr := Stringx(20)

	s3Config := &aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	}
	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)

	input := &s3manager.UploadInput{
		Bucket:      aws.String("bengcallbucket"),                                // bucket's name
		Key:         aws.String("docs/" + randomStr + "-" + fileheader.Filename), // files destination location
		Body:        file,                                                        // content of the file
		ContentType: aws.String("application/pdf"),                               // content type
	}
	res, err := uploader.UploadWithContext(context.Background(), input)

	// RETURN URL LOCATION IN AWS
	return res.Location, err
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
		log.Fatal(err.Error())
	}

	loo.Println("Mail sent!")
}

func DeleteFile() {
	err := os.Remove("ini-invoice.pdf")
	if err != nil {
		log.Print(err)
	}

}

func ToCoreData(d Data) DataCore {
	return DataCore{
		ID:       d.ID,
		Fullname: d.Fullname,
		Email:    d.Email,
		Phone:    d.Phone,
		Addres:   d.Addres,
		Other:    d.Other,
		Invoice:  d.Invoice,
		Total:    d.Total,
	}
}

func ToCoreService(s Service) ServiceCore {
	return ServiceCore{
		ID:            s.ID,
		Vehicle_name:  s.Vehicle_name,
		Service_name:  s.Service_name,
		TransactionID: s.TransactionID,
		SubTotal:      s.SubTotal,
	}
}
