package controllers

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

//QrcodeController 控制器
type QrcodeController struct {
	beego.Controller
}

//Get 二维码页
func (c *QrcodeController) Get() {
	var w, _ = c.GetInt("w", 350)
	var h, _ = c.GetInt("h", 350)
	base64 := c.GetString("text")
	if base64 == "" {
		base64 = "hello world!"
	}
	//log.Println("Original data:", base64)
	code, err := qr.Encode(base64, qr.L, qr.Unicode)
	// code, err := code39.Encode(base64)
	if err != nil {
		//log.Fatal(err)
		c.Ctx.WriteString("-")
		return
	}
	log.Println("Encoded data: ", code.Content())

	if base64 != code.Content() {
		//log.Fatal("data differs")
		c.Ctx.WriteString("-")
		return
	}

	code, err = barcode.Scale(code, w, h)
	if err != nil {
		//log.Fatal(err)
		c.Ctx.WriteString("-")
		return
	}
	//生成二维码文件
	//writePng("test.png", code)
	//设置输出头,直接输出二维码图片
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")
	png.Encode(c.Ctx.ResponseWriter, code)
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
