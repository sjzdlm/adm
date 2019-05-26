package conf

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/sjzdlm/adm/controllers"
	"github.com/sjzdlm/adm/controllers/adm"
	"github.com/sjzdlm/db"
)

//初始化adm路由
func InitRouter() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/upgradeie", &controllers.MainController{}, "get:UpgradeIE")
	//API接口
	beego.Router("/api", &controllers.ApiController{}, "get,post:Get")
	beego.AutoRouter(&controllers.ApiController{})
	beego.Router("/api/:module/:api", &controllers.ApiController{}, "get,post:Get")

	//rpt打印报表
	beego.Router("/rpt", &controllers.RptController{}, "get,post:Get")
	beego.AutoRouter(&controllers.RptController{})
	beego.Router("/rpt/:module/:rpt", &controllers.RptController{}, "get,post:Get")

	//p动态页面
	beego.Router("/p", &controllers.PController{}, "get,post:Get")
	beego.AutoRouter(&controllers.PController{})
	beego.Router("/p/:module/:page", &controllers.PController{}, "get,post:Get")
	//mp动态页面
	beego.Router("/mp", &controllers.MPController{}, "get,post:Get")
	//非权限校验接口
	beego.AutoRouter(&controllers.XApiController{})

	beego.AutoRouter(&controllers.MPController{})
	beego.Router("/mp/:app/:page", &controllers.MPController{}, "get,post:Get")
	beego.Router("/mp/:app", &controllers.MPController{}, "get,post:Get")
	beego.Router("/mp/vue", &controllers.MPController{}, "get,post:Get")
	beego.Router("/mp/css", &controllers.MPController{}, "get,post:Get")

	//单页应用
	beego.Router("/app", &controllers.AppController{}, "get,post:Get")
	beego.AutoRouter(&controllers.AppController{})
	beego.Router("/app/:app/:page", &controllers.AppController{}, "get,post:Get")
	beego.Router("/app/:app", &controllers.AppController{}, "get,post:Get")
	beego.Router("/app/vue", &controllers.AppController{}, "get,post:Get")
	beego.Router("/app/css", &controllers.AppController{}, "get,post:Get")

	//二维码
	beego.Router("/qrcode", &controllers.QrcodeController{})
	beego.AutoRouter(&controllers.QrcodeController{})

	//验证码
	beego.Router("/yzm", &controllers.YzmController{})
	beego.AutoRouter(&controllers.YzmController{})
	//工具类
	beego.Router("/tool", &controllers.ToolController{})
	beego.AutoRouter(&controllers.ToolController{})

	//ueditor上传文件
	beego.Router("/uef", &controllers.UEFController{})
	beego.AutoRouter(&controllers.UEFController{})

	//adm
	ns_adm :=
		beego.NewNamespace("/adm",
			beego.NSRouter("/login", &adm.LoginController{}),
			beego.NSAutoRouter(&adm.LoginController{}),
			beego.NSRouter("/main", &adm.MainController{}),
			beego.NSAutoRouter(&adm.MainController{}),

			beego.NSRouter("/mp", &adm.MPController{}),
			beego.NSAutoRouter(&adm.MPController{}),

			// beego.NSRouter("/tb", &adm.TbController{}),
			// beego.NSAutoRouter(&adm.TbController{}),

			beego.NSAutoRouter(&adm.MchController{}),
			beego.NSAutoRouter(&adm.MenuController{}),
			beego.NSAutoRouter(&adm.UserController{}),

			// //自定义数据页面
			// beego.NSRouter("/d/list/:code", &adm.DController{}, "get:List"),
			// beego.NSRouter("/d/info/:code", &adm.DController{}, "get:Info"),
			// beego.NSRouter("/d/view/:code", &adm.DController{}, "get:View"),
			// beego.NSRouter("/d/fe/:code", &adm.DController{}, "get:FE"),
			// beego.NSAutoRouter(&adm.DController{}),
		)
	beego.AddNamespace(ns_adm)
}

//应用程序根路径
func GetAPPRootPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}

//初始化静态文件目录
func InitConfig() {
	beego.SetStaticPath("/assets", beego.AppPath+"/static/assets")
	beego.SetStaticPath("/images", beego.AppPath+"/static/images")
	beego.SetStaticPath("/common", beego.AppPath+"/static/common")
	beego.SetStaticPath("/ufile", beego.AppPath+"/static/ufile")
	beego.SetStaticPath("/css", beego.AppPath+"/static/css")
	beego.SetStaticPath("/fonts", beego.AppPath+"/static/fonts")
	beego.SetStaticPath("/js", beego.AppPath+"/static/js")
	beego.SetStaticPath("/scripts", beego.AppPath+"/static/scripts")
	beego.SetStaticPath("/upload", beego.AppPath+"/static/upload")
	beego.SetStaticPath("/ueditor", beego.AppPath+"/static/ueditor")
	beego.SetStaticPath("/readme.txt", beego.AppPath+"/static/readme.txt")
	beego.SetStaticPath("/e", beego.AppPath+"/static/e")

	var files = db.Query("select * from adm_file ")
	for _, row := range files {
		beego.SetStaticPath(row["filepath"], beego.AppPath+"/static"+row["filepath"])
		fmt.Println(row["filepath"])
	}
	if len(files) < 1 {
		fmt.Println("......no file")
	}

	//服务配置
	//beego 服务器默认在请求的时候输出 server 为 beego。
	beego.BConfig.ServerName = "foo"
	beego.BConfig.RouterCaseSensitive = false //路由不区分大小写
	//启用会话，并将会话数据记录到文件中
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "_foosessionid"
	beego.BConfig.WebConfig.Session.SessionProvider = "file"
	os.MkdirAll("tmp", 0755)
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
	beego.BConfig.EnableGzip = true //启用压缩

	//日志不输出到终端
	//beego.BeeLogger.DelLogger("console")

	date := time.Now().Format("2006010215")
	os.MkdirAll("logs", 0755)
	beego.SetLogger("file", `{"filename":"logs/log`+date+`.log"}`)

	beego.ErrorController(&controllers.ErrorController{})
	beego.ErrorHandler("404", PageNotFound)
	// beego.Errorhandler("400", PageNotFound)
	// beego.Errorhandler("401", controllers.PageNotFound)
	// beego.Errorhandler("403", controllers.PageNotFound)
	//beego.Errorhandler("404", PageNotFound)
	// beego.Errorhandler("405", controllers.PageNotFound)
	// beego.Errorhandler("500", controllers.ServerError)
	// beego.Errorhandler("502", controllers.ServerError)
	// beego.Errorhandler("503", controllers.ServerError)
	// beego.Errorhandler("504", controllers.ServerError)

	// beego.EnableXSRF = true
	// beego.XSRFKEY = "61oETzKXQAGaYdkL5gEmGeJJFuYh7EQnp2XdTP1o"
	// beego.XSRFExpire = 3600

	var FilterUser = func(ctx *context.Context) {
		var _key = ctx.Input.Query("_key")
		fmt.Println("_key", _key)
		var _uid = ctx.Input.Session("_uid")
		if _key != "sjzapps" {
			if (_uid == nil || _uid == "") && ctx.Request.RequestURI != "/adm/login" && !strings.Contains(ctx.Request.RequestURI, "/adm/mp") {
				fmt.Println("_uid is nil now to go adm/login...")
				ctx.Redirect(302, "/adm/login")
			}
		}

	}
	beego.InsertFilter("/adm/*", beego.BeforeRouter, FilterUser)

	var FilterCtx = func(ctx *context.Context) {
		ctx.Input.SetData("ctx", ctx)
		//客户端标识---------------------------
		var ucid = ctx.GetCookie("_ucid")
		if ucid == "" {
			ucid = db.RandomString(10)
			ctx.SetCookie("_ucid", ucid)
		}
		//-------------------------------------
	}
	beego.InsertFilter("/*", beego.BeforeExec, FilterCtx)

	//添加自定义函数
	beego.AddFuncMap("funcLower", funcLower)
	beego.AddFuncMap("funcUpper", funcUpper)
	beego.AddFuncMap("funcBR", funcBR)
	beego.AddFuncMap("funcMap", funcMap)
	beego.AddFuncMap("funcOption", funcOption)

	beego.AddFuncMap("CookieGet", funcCookieGet)
	beego.AddFuncMap("CookieSet", funcCookieSet)
	beego.AddFuncMap("SessionGet", funcSessionGet)
	beego.AddFuncMap("SessionSet", funcSessionSet)

}

func PageNotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

func ServerError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}

//获取cookie值
func funcCookieGet(ctx *context.Context, str string) string {
	var rst = ctx.GetCookie(str)
	return rst
}

//设置cookie值
func funcCookieSet(ctx *context.Context, str string, val string) string {
	ctx.SetCookie(str, val)
	return ""
}

//获取session值
func funcSessionGet(ctx *context.Context, str string) string {
	if ctx.Input.Session(str) == nil {
		return ""
	}
	var rst = ctx.Input.Session(str).(string)
	return rst
}

//设置session值
func funcSessionSet(ctx *context.Context, str string, val string) string {
	ctx.Input.CruSession.Set(str, val)
	return ""
}

//转小写
func funcLower(str string) string {
	return strings.ToLower(str)
}

//转大写
func funcUpper(str string) string {
	return strings.ToUpper(str)
}

//根据取余数判断是否输出换行
func funcBR(i int, k int) string {
	if i%k > 0 {
		return "<br/>"
	}
	return ""
}

//根据key从map获取值
func funcMap(key string, val map[string]string) string {
	fmt.Println("funcMap:", key, val[key])
	key = strings.ToLower(key) //转小写
	return val[key]
}

//获取表单选项字符串
func funcOption(code, fieldcode, fieldname, formtype, formvalue string) string {
	if code == "" {
		return "-1"
	}
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil {
		return "-1"
	}

	var rst = ""
	//复选框
	if formtype == "复选框" {
		var ls = strings.Split(formvalue, ";")
		if len(ls) > 0 {
			for k, v := range ls {
				if k > 0 {
					rst += "\r\n"
				}

				var lsb = strings.Split(v, ",")
				if len(lsb) > 1 {
					if lsb[0] != "" {
						rst += `<input type="checkbox" id="` + fieldcode + strconv.Itoa(k) + `" title="` + fieldname + `" 
						name="` + fieldcode + `" value="` + lsb[0] + `"  style="vertical-align:middle;">&nbsp;
						<label for='` + fieldcode + strconv.Itoa(k) + `'  >` + lsb[1] + `</label>
						`
					} else {
						rst += `<input type="checkbox" id="` + fieldcode + strconv.Itoa(k) + `" title="` + fieldname + `" 
						name="` + fieldcode + `" style="vertical-align:middle;">&nbsp;
						<label for='` + fieldcode + strconv.Itoa(k) + `'  >` + lsb[1] + `</label>
						`
					}

				} else {
					if v != "" {
						rst += `<input type="checkbox" id="` + fieldcode + strconv.Itoa(k) + `" title="` + fieldname + `" 
						name="` + fieldcode + `" value="` + v + `"  style="vertical-align:middle;">&nbsp;
						<label for='` + fieldcode + strconv.Itoa(k) + `'  >` + v + `</label>
						`
					} else {
						rst += `<input type="checkbox" id="` + fieldcode + strconv.Itoa(k) + `" title="` + fieldname + `" 
						name="` + fieldcode + `" style="vertical-align:middle;">&nbsp;
						<label for='` + fieldcode + strconv.Itoa(k) + `'  >` + v + `</label>
						`
					}

				}
			}
		}
	}
	//单选框
	if formtype == "单选框" {
		var ls = strings.Split(formvalue, ";")
		if len(ls) > 0 {
			for kk, vv := range ls {
				if kk > 0 {
					rst += "\r\n"
				}
				var lsb = strings.Split(vv, ",")
				var kkstr = strconv.Itoa(kk)
				if len(lsb) > 1 {
					rst += `<input type="radio" id="` + fieldcode + kkstr + `" title="` + fieldname + `" 
					name="` + fieldcode + `" value="` + lsb[0] + `"  style="vertical-align:middle;">&nbsp;
					<label for='` + fieldcode + kkstr + `'  >` + lsb[1] + `</label>
					`
				} else {
					rst += `<input type="radio" id="` + fieldcode + kkstr + `" title="` + fieldname + `" 
					name="` + fieldcode + `" value="` + vv + `"  style="vertical-align:middle;">&nbsp;
					<label for='` + fieldcode + kkstr + `'  >` + vv + `</label>
					`
				}

			}
		}
	}
	//下拉框
	if formtype == "下拉框" {
		if strings.Contains(formvalue, "select") {
			var xx = db.NewDb(tb["conn_str"])
			var list = db.Query2(xx, formvalue)
			for kk, vv := range list {
				if kk > 0 {
					rst += "\r\n"
				}
				rst += `<option value="` + vv["id"] + `">` + vv["val"] + `</option>`
			}
		} else {
			var ls = strings.Split(formvalue, ";")
			//fmt.Println("下拉框:",ls)
			if len(ls) > 0 {
				for kk, vv := range ls {
					if kk > 0 {
						rst += "\r\n"
					}
					var lsb = strings.Split(vv, ",")
					if len(lsb) > 1 {
						rst += `<option value="` + lsb[0] + `">` + lsb[1] + `</option>`
					} else {
						rst += `<option value="` + vv + `">` + vv + `</option>`
					}

				}
			}
		}

	}

	return rst
}
