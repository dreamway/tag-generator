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
)

var fontLatin *canvas.FontFamily

var totalHeight = 200.0
var y = totalHeight - 10.0

func GenQRCode(waybill string, width, height int) image.Image {
	// Create the barcode
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

func drawTextAndMoveDown(c *canvas.Context, x float64, text *canvas.Text) {
	const spacing = 5
	y -= text.Bounds().H + spacing
	c.DrawText(x, y, text)
}

func renderSFContainerTag(tag string, timestamp string, subWaybillCount int, containerNo string, width, height float64, Orientation int) {
	cvs := canvas.New(width, height)
	ctx := canvas.NewContext(cvs)

	fontLatin = canvas.NewFontFamily("latin")
	if err := fontLatin.LoadSystemFont("serif", canvas.FontRegular); err != nil {
		panic(err)
	}

	fmt.Println("default y:", y)
	headerFace := fontLatin.Face(24.0, canvas.Black, canvas.FontBold, canvas.FontNormal)
	text12Face := fontLatin.Face(12.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	text10Face := fontLatin.Face(10.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	headerText := canvas.NewTextBox(headerFace, tag, 0.0, 30.0, canvas.Left, canvas.Top, 0.0, 0.0)
	drawTextAndMoveDown(ctx, 10.0, headerText)
	fmt.Println("y update:", y)

	timestamp_and_locate_info := fmt.Sprintf("%s 	(512WE %d ", timestamp, subWaybillCount)
	timestamp_and_locate_info += "件" + ")"
	fmt.Println(timestamp_and_locate_info)
	rt := canvas.NewRichText(text10Face)
	rt.WriteString(timestamp_and_locate_info)
	timestampText := rt.ToText(100, 80, canvas.Left, canvas.Right, 0.0, 0.0)
	drawTextAndMoveDown(ctx, 10.0, timestampText)
	fmt.Println("y update:", y)

	barcode := GenBarCode(containerNo, 200, 60)
	y -= 60
	ctx.DrawImage(10, y, barcode, canvas.DPMM(3.2))
	y -= float64(barcode.Bounds().Size().Y + 5)
	fmt.Println("y update:", y)

	containerNoText := canvas.NewTextBox(text12Face, containerNo, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	drawTextAndMoveDown(ctx, 10.0, containerNoText)

	//qrcode
	qrcode := GenQRCode(containerNo, 30, 30)
	ctx.DrawImage(width-50.0, y, qrcode, canvas.DPMM(3.2))

	//draw the canvas on the PDF
	//cvs.RenderTo(pdf)
	//cvs.RenderTo(canvasWriter)
	//cvs.WriteFile(containerNo+".pdf", renderers.PNG())
	renderers.Write(containerNo+".pdf", cvs)
}

func main() {
	renderSFContainerTag("123434", "04/23 05:24", 4, "641572254475", 300.0, totalHeight, 1)
}
