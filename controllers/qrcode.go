package controllers

import (
	"github.com/astaxie/beego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image"
	"image/png"
	"log"
	"os"
)

//QrcodeController 控制器
type QrcodeController struct {
	beego.Controller
}

//Get 二维码页
func (c *QrcodeController) Get() {
	base64 := "123456"
	//log.Println("Original data:", base64)
	code, err := qr.Encode(base64, qr.L, qr.Unicode)
	// code, err := code39.Encode(base64)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Encoded data: ", code.Content())

	if base64 != code.Content() {
		log.Fatal("data differs")
	}

	code, err = barcode.Scale(code, 300, 300)
	if err != nil {
		log.Fatal(err)
	}

	writePng("test.png", code)
}
func writePng(filename string, img image.Image) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(file, img)
	// err = jpeg.Encode(file, img, &jpeg.Options{100})      //图像质量值为100，是最好的图像显示  
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	log.Println(file.Name())
}
