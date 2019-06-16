package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//AppController 控制器  动态页面
type AppController struct {
	beego.Controller
}

///Vue组件脚本
func (c *AppController) Vue() {
	//遍历所有get参数信息放到模板变量--------------------------------
	var paramstr = ""
	var urls = strings.Split(c.Ctx.Input.URI(), "?")
	if len(urls) > 1 {
		var params = strings.Split(urls[1], "&")
		for i := 0; i < len(params); i++ {
			if params[i] == "" || params[i] == "&" {
				continue
			}
			var p = strings.Split(params[i], "=")
			if len(p) < 2 {
				continue
			}
			p[1], _ = url.QueryUnescape(p[1])
			c.Data[p[0]] = p[1]
			if paramstr != "" {
				paramstr += ","
			}
			paramstr += "&" + p[0] + "=" + p[1]
		}
	}
	c.Data["_paramstr"] = paramstr
	//------------------------------------------------------------
	var rst = ""
	var appcode = c.GetString("app")
	if appcode != "" {
		var app = db.First("select * from tbm_app where code=? limit 1", appcode)
		if len(app) > 0 {
			var plist = db.Query("select * from tbm_page where app_id=?", app["id"])
			for _, prow := range plist {
				var wlist = db.Query("select * from tbm_widget where tbmid=?", prow["id"])
				for _, row := range wlist {
					//模板
					rst += `
					var ` + prow["code"] + "_" + row["form_type"] + ` = Vue.extend({
						template: '\
						`
					rst += strings.Replace(strings.Replace(row["tpltxt"], "\r\n", "\\\r\n", -1), "'", "\\'", -1)
					rst += `	',
					`
					//对tpldata进行golang的模板解析
					var data map[string]interface{} = map[string]interface{}{}
					data["app"] = app
					data["page"] = prow
					data["m"] = row
					var tpl = template.New("")
					tpl.Parse(row["tpldata"])
					var buf bytes.Buffer
					//fmt.Println("-----------------------")
					//fmt.Println("tpl:", row["tpldata"])
					if row["tpldata"] != "" {
						var e = tpl.Execute(&buf, data)
						if e == nil {
							rst += buf.String()
						} else {
							rst += `
							data(){
								return{
									msg:"` + e.Error() + `"
								}
							}
							`
						}
					}
					//rst += row["tpldata"]

					rst += `})
					Vue.component('` + prow["code"] + "_" + row["form_type"] + `', ` + prow["code"] + "_" + row["form_type"] + `)
					`
				}
			}
		}
	}
	//进行模板解析
	var tpl = template.New("")
	tpl.Parse(rst)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e != nil {
		fmt.Println("vue 组件 template 执行错误:", e.Error())
	} else {
		rst = buf.String()
	}

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	//c.Ctx.WriteString(rst)
}

///Vue组件样式
func (c *AppController) CSS() {
	var rst = ""
	var appcode = c.GetString("app")
	if appcode == "" {
		appcode = "0"
	}
	var app = db.First("select * from tbm_app where code=?", appcode)
	if len(app) < 1 {
		c.Ctx.Output.Header("Content-Type", "text/css; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))

		c.Ctx.WriteString(rst)
	}
	var list = db.Query("select  id, form_type,tplcss  from tbm_widget  where tbmid in (select id from tbm_page where app_id=?) and state=1", app["id"])
	var types = make(map[string]string)
	for _, row := range list {
		// if _, ok := types[row["form_type"]]; ok {
		// 	continue
		// }
		types[row["form_type"]] = row["form_type"]
		rst += row["tplcss"] + "\r\n"
	}

	c.Ctx.Output.Header("Content-Type", "text/css; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//Get 默认页
func (c *AppController) Get() {
	var rst = ""
	var appcode = c.Ctx.Input.Param(":app")
	if appcode == "" {
		rst = "没有找到应用信息"
		c.Ctx.WriteString(rst)
		return
	}
	c.Data["appcode"] = appcode

	//遍历所有get参数信息放到模板变量--------------------------------
	var paramstr = ""
	var urls = strings.Split(c.Ctx.Input.URI(), "?")
	if len(urls) > 1 {
		var params = strings.Split(urls[1], "&")
		for i := 0; i < len(params); i++ {
			if params[i] == "" || params[i] == "&" {
				continue
			}
			var p = strings.Split(params[i], "=")
			if len(p) < 2 {
				continue
			}
			p[1], _ = url.QueryUnescape(p[1])
			c.Data[p[0]] = p[1]
			if paramstr != "" {
				paramstr += ","
			}
			paramstr += "&" + p[0] + "=" + p[1]
		}
	}
	c.Data["_paramstr"] = paramstr
	//------------------------------------------------------------

	var app = db.First("select * from tbm_app where code=?", appcode)
	if len(app) < 1 {
		rst = "没有找到应用信息"
		c.Ctx.WriteString(rst)
		return
	}
	//判断是否需要授权
	if app["oauth2_on"] == "1" {
		var openid = c.GetString("_openid")
		if openid == "" {
			var redirect_uri = "http://" + c.Ctx.Input.Domain() + "/app/" + appcode
			var js = `
			<script>
				window.location='/wxmp/oauth2?wxid=` + app["oauth2_wxid"] + "&scope=snsapi_userinfo&url=" + url.QueryEscape(redirect_uri) + `';
			</script>
			`
			c.Ctx.WriteString(js)
			return
		}
		var nickname = c.GetString("_nickname")
		var headimg = c.GetString("_headimgurl")
		c.Data["_openid"] = openid
		c.Data["_nickname"] = nickname
		c.Data["_headimgurl"] = headimg
	}
	c.Data["app"] = app
	//判断是否需要登录
	c.Data["login_on"] = app["login_on"]
	//应用导航
	var isnav = "0"
	var nav = db.Query("select * from  tbm_nav_menu where app_id=? order by sort ", app["id"])
	if len(nav) < 1 {
		isnav = "0"
		c.Data["isnav"] = isnav
	} else {
		isnav = "1"
		c.Data["isnav"] = isnav
	}
	c.Data["navmenu"] = nav

	//默认页面
	var code = c.Ctx.Input.Param(":page")
	if code == "" {
		if app["mainpage"] != "" {
			code = app["mainpage"]
		} else {
			code = "index"
		}
	}
	c.Data["page"] = code

	var m = db.First("select * from tbm_page where state=1 and app_id=? and code=?", app["id"], code)
	if len(m) < 1 {
		rst = "没有找到页面信息."
		c.Ctx.WriteString(rst)
		return
	}
	c.Data["m"] = m

	// var list = db.Query("select * from tbm_widget where state=1 and tbmid=? order by form_sort", m["id"])
	// c.Data["list"] = list
	// var vue = ""
	// for _, row := range list {
	// 	vue += "<" + row["form_type"] + "></" + row["form_type"] + ">\r\n"
	// }
	// c.Data["vue"] = vue
	//生成页面路由
	var routes = ""
	var pagetpl = ""
	var plist = db.Query("select * from tbm_page where state=1 and app_id=? ", app["id"])
	for i, row := range plist {

		if row["code"] == "index" {
			if routes != "" {
				routes += ","
			}
			routes += "{path: '/', component: " + row["code"] + ",meta:{index:0,keepAlive: true}}"
		} else {
			if routes != "" {
				routes += ","
			}
			routes += "{path: '/" + row["code"] + "', component: " + row["code"] + ",meta:{index:" + strconv.Itoa(i) + ",keepAlive: false}}"
		}

		//
		var wlist = db.Query("select * from tbm_widget where state=1 and tbmid=?  order by form_sort", row["id"])
		if len(wlist) > 0 {
			pagetpl += `
			var ` + row["code"] + `=Vue.extend({
				template: '<div>`
			for _, r := range wlist {
				pagetpl += "<" + row["code"] + "_" + r["form_type"] + "></" + row["code"] + "_" + r["form_type"] + ">"
			}
			if isnav == "1" && row["is_nav"] == "1" {
				pagetpl += "<nav_menu></nav_menu>"
			}
			pagetpl += `</div>',
				props: []
			})
			Vue.component('` + row["code"] + `', ` + row["code"] + `)			
			`
		}

	}

	//进行模板解析
	var tpl = template.New("")
	tpl.Parse(pagetpl)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e != nil {
		fmt.Println("vue pagetpl template 执行错误:", e.Error())
	} else {
		pagetpl = buf.String()
	}

	c.Data["pagetpl"] = pagetpl
	//fmt.Println("pagetpl", pagetpl)

	c.Data["routes"] = routes

	c.TplName = "app.html"
}
