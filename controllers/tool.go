package controllers

import (
	"image/jpeg"
	"net/http"
	"strings"
	"time"

	"golang.org/x/image/draw"

	//"image/color"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	// "net/http"
	// "strings"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/golang/freetype"
	"github.com/sjzdlm/db"
)

func init() {

}

//小工具集合
type ToolController struct {
	beego.Controller
}

//请求URL数据
func HttpGet(url string) (string, error) {

	postReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("请求失败", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		fmt.Println("client请求失败", err)
		return "", err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(data), err
}

//-------------------------------------------------------------------------------
//调用一个网址
func funcHttpGet(urlstr string) string {
	var rst, _ = HttpGet(urlstr)
	return rst
}

//转小写
func funcLower(str string) string {
	return strings.ToLower(str)
}

//转大写
func funcUpper(str string) string {
	return strings.ToUpper(str)
}

//将&gt;&lt;转义回来
func funcUnEscape(str string) string {
	fmt.Println("替换前:", str)
	var str1 = strings.Replace(str, "&lt;", "<", -1)
	fmt.Println("替换后1:", str1)
	var str2 = strings.Replace(str1, "&gt;", ">", -1)
	fmt.Println("替换后2:", str2)
	return str2
}

//根据取余数判断是否输出换行
func funcBR(i int, k int) string {
	if i%k > 0 {
		return "<br/>"
	}
	return ""
}

//取余并返回
func funcMod(i int, k int) int {
	return i % k
}

//根据key从map获取值
func funcMap(key string, val map[string]string) string {
	//fmt.Println("funcMap-val:",val)
	//fmt.Println("funcMap:",key,val[key])
	key = strings.ToLower(key) //转小写
	return val[key]
}

///Vue组件脚本
func (c *ToolController) Vue() {
	var list = db.Query("select * from tbm_widget_type where state=1")
	var rst = ""

	for _, row := range list {
		rst += row["tpltxt"]
	}

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//图片上绘制文字
func (this *ToolController) DrawingTextOnImg() {
	const (
		dx       = 500 // 图片的大小 宽度
		dy       = 500 // 图片的大小 高度
		fontFile = "static/fonts/yahei.ttf"
		fontSize = 16 // 字体尺寸
		fontDPI  = 72 // 屏幕每英寸的分辨率
	)
	var fid = this.GetString("fid")
	if fid == "" {
		fid = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	var iname = this.GetString("iname")
	if iname == "" {
		iname = "static/images/logo.jpg"
	}
	var x1, _ = this.GetInt("x1", -1)
	var y1, _ = this.GetInt("y1", -1)
	var txt1 = this.GetString("txt1")

	var x2, _ = this.GetInt("x2", -1)
	var y2, _ = this.GetInt("y2", -1)
	var txt2 = this.GetString("txt2")

	var x3, _ = this.GetInt("x3", -1)
	var y3, _ = this.GetInt("y3", -1)
	var txt3 = this.GetString("txt3")

	var x4, _ = this.GetInt("x4", -1)
	var y4, _ = this.GetInt("y4", -1)
	var txt4 = this.GetString("txt4")

	var x5, _ = this.GetInt("x5", -1)
	var y5, _ = this.GetInt("y5", -1)
	var txt5 = this.GetString("txt5")

	var x6, _ = this.GetInt("x6", -1)
	var y6, _ = this.GetInt("y6", -1)
	var txt6 = this.GetString("txt6")

	file, err := os.Open("000.jpg")
	if err != nil {
		fmt.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	defer file.Close()

	jpg, err := jpeg.Decode(file) //解码
	if err != nil {
		fmt.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 480, 680))
	draw.Draw(img, jpg.Bounds().Add(image.Pt(0, 0)), jpg, jpg.Bounds().Min, draw.Src) //截取图片的一部分

	// 需要保存的文件
	//imgcounter := 123
	//imgfile, _ := os.Create(fmt.Sprintf("static/images/tmp/%03d.png", imgcounter))
	imgfile, _ := os.Create("static/images/tmp/" + fid + ".png")
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
		this.Ctx.WriteString(err.Error())
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("转换字体样式出错")
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}

	c := freetype.NewContext()
	c.SetDPI(fontDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	//-----------------------------------------------------------------------
	if x1 >= 0 && y1 >= 0 && txt1 != "" {
		pt := freetype.Pt(x1, y1+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

		_, err = c.DrawString(txt1, pt)
		if err != nil {
			log.Println("向图片写字体出错1")
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	}

	if x2 >= 0 && y2 >= 0 && txt2 != "" {
		pt := freetype.Pt(x2, y2+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

		_, err = c.DrawString(txt2, pt)
		if err != nil {
			log.Println("向图片写字体出错2")
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	}

	if x3 >= 0 && y3 >= 0 && txt3 != "" {
		pt := freetype.Pt(x3, y3+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

		_, err = c.DrawString(txt3, pt)
		if err != nil {
			log.Println("向图片写字体出错3")
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	}

	if x4 >= 0 && y4 >= 0 && txt4 != "" {
		pt := freetype.Pt(x4, y4+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

		_, err = c.DrawString(txt4, pt)
		if err != nil {
			log.Println("向图片写字体出错4")
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	}

	if x5 >= 0 && y5 >= 0 && txt5 != "" {
		pt := freetype.Pt(x5, y5+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

		_, err = c.DrawString(txt5, pt)
		if err != nil {
			log.Println("向图片写字体出错5")
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	}

	if x6 >= 0 && y6 >= 0 && txt6 != "" {
		pt := freetype.Pt(x6, y6+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

		_, err = c.DrawString(txt6, pt)
		if err != nil {
			log.Println("向图片写字体出错6")
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	}
	//------------------------------------------------------------------------

	// //重新设置第二行y的位置
	// pt.Y += c.PointToFixed(fontSize)
	// _, err = c.DrawString("13932172487", pt)
	// if err != nil {
	// 	log.Println("向图片写字体出错")
	// 	log.Println(err)
	// 	return
	// }

	// 以PNG格式保存文件
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Println("生成图片出错")
		log.Fatal(err)
		this.Ctx.WriteString(err.Error())
		return
	}

	//this.Ctx.WriteString("/images/tmp/"+fid+".png")
	this.Ctx.Output.Download("static/images/tmp/" + fid + ".png")
}
