package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"

	//canvas
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"github.com/tdewolff/canvas/renderers/pdf"
	//inkscape
)

var fontDejaVu *canvas.FontFamily

var X_PACKAGENO = 10.0
var Y_PACKAGENO = 10.0

var X_MARGIN = 12.0
var Y_MARGIN = 5.0

var QR_SIZE = 33
var X_QR = 127.0
var Y_QR = 15.0

var X_CELL = 80.0
var Y_CELL = 50.0

func GenQRCode(waybill string, destinationCode string, width, height int) image.Image {
	qrCode, _ := qr.Encode(waybill, qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, width, height)

	// create the output file
	file, _ := os.Create("qrcode.png")
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return qrCode
}

func GenBarCode(waybill string, width, height int) image.Image {
	code, err := code128.EncodeWithoutChecksum(waybill)
	if err != nil {
		panic(err)
	}
	code, err = barcode.Scale(code, width, height)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("barcode.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, code)

	return code
}

func renderSFContainerTag(tag string, timestamp string, subWaybillCount int, containerNo, destinationCode string, width, height float64, Orientation int) {
	cvs := canvas.New(width, height)
	ctx := canvas.NewContext(cvs)
	ctx.SetFillColor(canvas.Aqua)
	ctx.DrawPath(0, 0, canvas.Rectangle(width, height))

	fontDejaVu = canvas.NewFontFamily("dejavu")
	if err := fontDejaVu.LoadFontFile("./resources/DejaVu_Sans/DejaVuSans.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	leftMargin := 10.0
	baseMargin := 0.09 * height
	y := baseMargin
	Y_MARGIN = 0.08 * height
	X_MARGIN = 0.08 * width

	fmt.Println("default y:", y, ", Y_MARGIN:", Y_MARGIN)
	textTagFace := fontDejaVu.Face(0.24*height, canvas.Black, canvas.FontBold, canvas.FontNormal)
	textContainerNoFace := fontDejaVu.Face(0.12*height, canvas.Black, canvas.FontBold, canvas.FontNormal)
	textTimestampFace := fontDejaVu.Face(0.1*height, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	flatTagFace := fontDejaVu.Face(0.24*height, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	//ContainerNo (Bottom-Left)
	containerNoText := canvas.NewTextBox(textContainerNoFace, containerNo, 0.0, height*0.12, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(float64(leftMargin), y, containerNoText)

	y += containerNoText.Bounds().H + Y_MARGIN
	fmt.Println("y update:", y)

	//BarCode
	barcode := GenBarCode(containerNo, int(width*0.6), int(height*0.3))
	ctx.DrawImage(leftMargin, y, barcode, canvas.DPMM(1.0))

	//qrcode
	qrCodeX := float64(barcode.Bounds().Max.X) + X_MARGIN
	qrCodeY := y + (0.25-0.2)*height
	QR_SIZE = int(height * 0.2)
	qrcode := GenQRCode(containerNo, destinationCode, QR_SIZE, QR_SIZE)
	ctx.DrawImage(qrCodeX, float64(qrCodeY), qrcode, canvas.DPMM(1.0))

	//flatTag
	flatTagX := qrCodeX
	flatTagW := QR_SIZE
	flatTagH := QR_SIZE / 2
	flatTagText := canvas.NewTextBox(flatTagFace, "扁平", float64(flatTagW), float64(flatTagH), canvas.Left, canvas.Top, 0.0, 0.0)
	flatTagY := qrCodeY - flatTagText.Bounds().H - Y_MARGIN

	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(1.0)
	ctx.DrawPath(flatTagX, flatTagY, canvas.Rectangle(float64(flatTagW), float64(flatTagH)))
	ctx.Fill()

	ctx.DrawText(flatTagX, flatTagY+flatTagText.Bounds().H, flatTagText)

	//timestamp
	fmt.Println("barcode.Bounds.Max:", barcode.Bounds().Max, ", barcode.Bounds.Min:", barcode.Bounds().Min, ", diff:", barcode.Bounds().Max.Y-barcode.Bounds().Min.Y)
	y += float64(barcode.Bounds().Max.Y-barcode.Bounds().Min.Y) + 2*Y_MARGIN
	fmt.Println("y update:", y)

	timestampLocateInfo := fmt.Sprintf("%s (512WE %d 件)", timestamp, subWaybillCount)
	fmt.Println(timestampLocateInfo)
	rt := canvas.NewRichText(textTimestampFace)
	rt.WriteString(timestampLocateInfo)
	timestampText := rt.ToText(width, height*0.1, canvas.Left, canvas.Bottom, 0.0, 0.0)
	ctx.DrawText(leftMargin, y, timestampText)

	//Header
	headerText := canvas.NewTextBox(textTagFace, tag, 0.0, height*0.2, canvas.Left, canvas.Top, 0.0, 0.0)
	headerY := height - headerText.Bounds().H
	ctx.DrawText(float64(leftMargin), headerY, headerText)

	//draw the canvas on the PDF
	renderers.Write(containerNo+".svg", cvs)
	renderers.Write(containerNo+".png", cvs, canvas.DPI(96.0))
	renderers.Write(containerNo+".pdf", cvs, &pdf.Options{
		Compress:      false,
		SubsetFonts:   false,
		ImageEncoding: canvas.Lossless,
	})

	// Options{
	// 	Compress:      true,
	// 	SubsetFonts:   true,
	// 	ImageEncoding: canvas.Lossless,
	// }

	fmt.Println("Done!!")
}

func renderTextCoordinate(tag string, timestamp string, subWaybillCount int, containerNo, destinationCode string, width, height float64, Orientation int) {
	cvs := canvas.New(width, height)
	ctx := canvas.NewContext(cvs)

	ctx.SetFillColor(canvas.Yellow)
	ctx.DrawPath(0, 0, canvas.Rectangle(width, height))

	fontDejaVu = canvas.NewFontFamily("latin")
	if err := fontDejaVu.LoadSystemFont("serif", canvas.FontRegular); err != nil {
		panic(err)
	}

	const spacing = 5
	const margin = 10

	text18Face := fontDejaVu.Face(18.0, canvas.Black, canvas.FontBold, canvas.FontNormal)
	text12Face := fontDejaVu.Face(12.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	text10Face := fontDejaVu.Face(10.0, canvas.Green, canvas.FontRegular, canvas.FontNormal)
	text8Face := fontDejaVu.Face(8.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	headerText := canvas.NewTextBox(text18Face, "(10,10)", 0.0, 30.0, canvas.Left, canvas.Top, 0.0, 0.0)
	face18Height := headerText.Bounds().H
	fmt.Println("face18Height:", face18Height)
	ctx.DrawText(margin, margin, headerText)

	containerNoText := canvas.NewTextBox(text12Face, "(width-10,height-10)", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(width-margin, height-margin, containerNoText)

	cornerText := canvas.NewTextBox(text10Face, "(width-10,10)", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(width-margin, margin, cornerText)

	corner2Text := canvas.NewTextBox(text8Face, "(10,height-10)", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(margin, height-margin, corner2Text)

	//ctx.Scale(600.0/width, 400.0/height)

	renderers.Write(containerNo+".svg", cvs)
	renderers.Write(containerNo+".png", cvs, canvas.DPI(96.0))
	renderers.Write(containerNo+".pdf", cvs)
}

func renderElement() {
	jian := "件"
	bianPing := "扁平"

	width := 80.0
	height := 80.0

	cvs := canvas.New(width, height)
	ctx := canvas.NewContext(cvs)

	fontDejaVu = canvas.NewFontFamily("dejavu")
	if err := fontDejaVu.LoadFontFile("./resources/DejaVu_Sans/DejaVuSans.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	textFace := fontDejaVu.Face(24.0, canvas.Black, canvas.FontBold, canvas.FontNormal)

	headerText := canvas.NewTextBox(textFace, jian, width, height, canvas.Left, canvas.Top, 0.0, 0.0)
	ctx.DrawText(0, 0, headerText)

	renderers.Write("test.svg", cvs)
	renderers.Write("test.png", cvs, canvas.DPI(96.0))

	fmt.Println(bianPing)
}

func main() {
	defaultWidth := 600.0
	defaultHeight := 400.0
	//renderElement()
	renderSFContainerTag("2-13-E-G06", "04/23 05:24", 4, "641572254475", "512W", defaultWidth, defaultHeight, 1)

	//renderTextCoordinate("123434", "04/23 05:24", 4, "11641134343", "512W", defaultWidth, defaultHeight, 1)
}
