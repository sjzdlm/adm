package controllers

import (
	"image/jpeg"

	"github.com/sjzdlm/db"
	"golang.org/x/image/draw"

	//"image/color"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/golang/freetype"
	//registry "github.com/golang/sys/windows/registry"
)

func init() {

}

type TestController struct {
	beego.Controller
}

func (c *TestController) Regedit() {
	//注册表测试代码
	// //k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, registry.ALL_ACCESS)
	// k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Clients\StartMenuInternet\Google Chrome\shell\open\command`, registry.ALL_ACCESS)

	// if err != nil {
	// 	c.Ctx.WriteString(err.Error())
	// 	return
	// }

	// fmt.Println("k",k)

	// keys, _ := k.ReadSubKeyNames(0)
	// fmt.Println(keys)

	// //s, _, err := k.GetStringValue("A")
	// s, _, err := k.GetStringValue("")
	// if err != nil {
	// 	c.Ctx.WriteString("b:"+err.Error())
	// 	return
	// }

	// defer k.Close()
	// c.Ctx.WriteString(s)
	c.Ctx.WriteString("ok")
}

//Default 默认首页
func (c *TestController) Get() {
	fmt.Println("this.Ctx.Request.Host:", c.Ctx.Request.Host)
	fmt.Println("a:", c.Ctx.Request.URL)
	fmt.Println("b:", c.Ctx.Input.Host())
	fmt.Println("b:", c.Ctx.Request.URL.Hostname())
	fmt.Println("b:", c.Ctx.Request.URL.RequestURI())
	c.Ctx.WriteString("ok")
}

func (c *TestController) Test12() {
	var str = "“百年盛膳”>：参/鸡汤"
	rs := []rune(str)
	length := len(rs)
	fmt.Println("A", length, strings.LastIndex(str, "/"))
	fmt.Println("B", len(str))

	fmt.Println(UnicodeIndex(str, "/"))
	fmt.Println(SubString(str, UnicodeIndex(str, "参"), 1))
	c.Ctx.WriteString("ok")
}
func UnicodeIndex(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}

	return result
}
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}
func (c *TestController) Test1() {
	var Url = "http://www.baidu.com/"
	requestLine := strings.Join([]string{Url,
		"s?wd=", "xx"}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("发送get请求获取,错误", err)
		c.Ctx.WriteString("发送get请求获取,错误" + err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("发送get请求获取,读取返回body错误", err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}
	fmt.Println("body", body)
	var rst = fmt.Sprintf("%s", body)
	c.Ctx.WriteString(rst)
}
func (c *TestController) Test2() {
	var _uid = c.GetSession("_uid")
	if _uid == nil {
		c.Ctx.WriteString("nil")
	} else {
		c.Ctx.WriteString(fmt.Sprintf("%s", _uid))
	}
	c.Ctx.WriteString(c.Ctx.Input.Domain())
}
func (c *TestController) Test_Upload() {
	c.TplName = "test/test_upload.html"
}
func (c *TestController) AddPhoto() {
	f, h, _ := c.GetFile("image") //获取上传的文件
	path := "img"                 //c.Input().Get("url")	//存文件的路径
	//path = path[7:]
	path = "./static/ufile/" + path + "/" + h.Filename
	f.Close() // 关闭上传的文件，不然的话会出现临时文件不能清除的情况
	fmt.Println("path:", path)
	c.SaveToFile("image", path) //存文件    WaterMark(path)	//给文件加水印
	// c.Redirect("/test/addphoto", 302)
	c.Ctx.WriteString("1")
}

func (c *TestController) Cookie1() {
	var txt = c.GetString("txt")
	if txt == "" {
		txt = "abc"
	}
	c.Ctx.SetCookie("txt", txt)
	c.Ctx.WriteString(txt)
}
func (c *TestController) Cookie2() {
	var txt = "[" + c.Ctx.GetCookie("txt") + "]"

	c.Ctx.WriteString(txt)
}
func (c *TestController) Url() {
	var url = c.Ctx.Request.RequestURI //获取当前路径
	var a = db.Substring(url, "/", "?")
	if a != "" {
		url = a
	}

	c.Ctx.WriteString(url)
}

func (this *TestController) Font() {
	const (
		dx = 500 // 图片的大小 宽度
		dy = 500 // 图片的大小 高度
		// fontFile = "FZFSK.TTF"
		fontFile = "static/fonts/yahei.ttf"
		fontSize = 20 // 字体尺寸
		fontDPI  = 72 // 屏幕每英寸的分辨率
	)

	file, err := os.Open("000.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	jpg, err := jpeg.Decode(file) //解码
	if err != nil {
		fmt.Println(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 640, 906))
	draw.Draw(img, jpg.Bounds().Add(image.Pt(0, 0)), jpg, jpg.Bounds().Min, draw.Src) //截取图片的一部分

	// 需要保存的文件
	imgcounter := 123
	imgfile, _ := os.Create(fmt.Sprintf("%03d.png", imgcounter))
	defer imgfile.Close()
	// 新建一个 指定大小的 RGBA位图
	////img := image.NewNRGBA(image.Rect(0, 0, dx, dy))
	// 画背景
	/*for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			// 设置某个点的颜色，依次是 RGBA
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 0, 255})
		}
	}*/
	// 读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println("读取字体数据出错")
		log.Println(err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("转换字体样式出错")
		log.Println(err)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(fontDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	pt := freetype.Pt(100, 300+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

	_, err = c.DrawString("杜立敏", pt)
	if err != nil {
		log.Println("向图片写字体出错")
		log.Println(err)
		return
	}
	//重新设置第二行y的位置
	pt.Y += c.PointToFixed(fontSize)
	_, err = c.DrawString("13932172487", pt)
	if err != nil {
		log.Println("向图片写字体出错")
		log.Println(err)
		return
	}

	// 以PNG格式保存文件
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Println("生成图片出错")
		log.Fatal(err)
	}

	this.Ctx.WriteString("ok")
}
