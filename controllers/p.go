package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//PController 控制器  动态页面
type PController struct {
	beego.Controller
}

//Get 默认页
func (c *PController) Get() {
	var rst = ""
	var module = c.Ctx.Input.Param(":module")
	if module == "" {
		rst = "没有找到页面信息"
		return
	}
	var rpt = c.Ctx.Input.Param(":page")
	if rpt == "" {
		rpt = "index"
	}
	if rpt == "" {
		rst = "没有找到页面信息."
		c.Ctx.WriteString(rst)
		return
	}
	//定义参数map
	//var data map[string]interface{} = map[string]interface{}{}
	var data = c.Data
	//查找页面信息
	var m = db.FirstOrNil("select * from page_list where module=? and  code=?", module, rpt)
	if m == nil {
		// rst = "没有找到页面信息."
		// c.Ctx.WriteString(rst)
		rst = notfound
		c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))
		return
	} else {

		data["__date"] = time.Now().Format("2006-01-02")
		data["__time"] = time.Now().Format("2006-01-02 15:04:05")
		//填充登录参数
		var _uid = c.GetSession("_uid")
		var _mch_id = c.GetSession("_mch_id")
		var _username = c.GetSession("_username")
		var _usertype = c.GetSession("_usertype")
		data["_uid"] = _uid
		data["_mch_id"] = _mch_id
		data["_username"] = _username
		data["_usertype"] = _usertype

		//从数据库读取报表的参数,并从客户端获取参数值
		var params = db.Query("select * from page_param where page_id=? and state=1 order by orders ", m["id"])
		//fmt.Println("params:",params)

		//普通参数接收
		for _, v := range params {
			if v["param_type"] == "参数" {
				var p = c.GetString(v["param_name"])
				if p == "" && v["param_value"] != "" {
					p = v["param_value"]
				}
				if p != "" {
					data[v["param_name"]] = p
				}
			}
		}
		//cookie参数接收
		for _, v := range params {
			if v["param_type"] == "cookie" {
				var p = c.Ctx.GetCookie(v["param_name"])
				if p == "" && v["param_value"] != "" {
					p = v["param_value"]
				}
				if p != "" {
					data[v["param_name"]] = p
				}
			}
		}
		//session参数接收
		for _, v := range params {
			if v["param_type"] == "session" {
				var p1 = c.GetSession(v["param_name"])
				var p = ""
				if p1 != nil {
					p = p1.(string)
				}
				if p == "" && v["param_value"] != "" {
					p = v["param_value"]
				}
				if p != "" {
					data[v["param_name"]] = p
				}
			}
		}
		//SQL执行结果
		var conn = m["conn_str"]
		c.SetSession("_conn", conn) //将连接串保存的session中,便于模板中获取调用
		if conn != "" {
			var XX = db.NewDb(conn)
			for _, v := range params {
				if v["param_type"] == "model" { //model 单条纪录
					var tpl = template.New("")
					tpl.Parse(v["param_value"])
					var buf bytes.Buffer
					var e = tpl.Execute(&buf, data)
					if e != nil {
						fmt.Println("tpl.Execute sql错误:", e.Error())
						rst = ApiResult("500", "参数错误", e.Error(), "")
						c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
						c.Ctx.Output.Body([]byte(rst))
						return
					}
					if m["conn_str"] != "" {
						var p = db.First2(XX, buf.String())
						if p != nil {
							data[v["param_name"]] = p
						}

					}
				} else if v["param_type"] == "list" || v["param_type"] == "SQL" || v["param_type"] == "sql" { //如果参数类型是sql

					var tpl = template.New("")
					tpl.Parse(v["param_value"])
					var buf bytes.Buffer
					var e = tpl.Execute(&buf, data)
					if e != nil {
						fmt.Println("tpl.Execute sql错误:", e.Error())
						c.Ctx.WriteString("页面模板错误500!")
						return
					}
					var sql = buf.String()
					//< >号会被转义,此处替换回来
					sql = strings.Replace(sql, "&lt;", "<", -1)
					sql = strings.Replace(sql, "&gt;", ">", -1)

					if strings.Contains(sql, "limit 1") || strings.Contains(sql, "top 1") {
						var p = db.First2(XX, sql)
						if p != nil {
							data[v["param_name"]] = p
						} else {
							data[v["param_name"]] = ""
						}
					} else {
						var _page, _ = c.GetInt("_page", 0)
						var _pagesize, _ = c.GetInt("_pagesize", 20)
						if _page > 0 {
							var p = db.Pager2(XX, _page, _pagesize, sql)
							data[v["param_name"]] = p
						} else {
							var p = db.Query2(XX, sql)
							if p != nil {
								data[v["param_name"]] = p
							}
						}
					}

					// var p = db.Query2(XX, sql)
					// if p != nil {
					// 	data[v["param_name"]] = p
					// } else {
					// 	data[v["param_name"]] = ""
					// }
				} else if v["param_type"] == "SQLEXEC" || v["param_type"] == "sqlexec" { //如果参数类型是sql
					var tpl = template.New("")
					tpl.Parse(v["param_value"])
					var buf bytes.Buffer
					var e = tpl.Execute(&buf, data)
					if e != nil {
						fmt.Println("tpl.Execute sqlexec错误:", e.Error())
						c.Ctx.WriteString("页面模板错误500!")
						return
					}

					var sql = buf.String()
					//< >号会被转义,此处替换回来
					sql = strings.Replace(sql, "&lt;", "<", -1)
					sql = strings.Replace(sql, "&gt;", ">", -1)

					var p = db.Exec2(XX, sql)
					data[v["param_name"]] = p
				}
			}
		}
		//变量输出到模板
		for _, v := range params {
			if v["param_type"] == "变量" {
				var tpl = template.New("")
				tpl.Parse(v["param_value"])
				var buf bytes.Buffer
				var e = tpl.Execute(&buf, data)
				if e != nil {
					fmt.Println("tpl.Execute 变量错误:", e.Error())
					c.Ctx.WriteString("页面模板错误500!")
					return
				}

				var p = buf.String()
				if p == "" && v["param_value"] != "" {
					p = v["param_value"]
				}
				if p != "" {
					data[v["param_name"]] = p
				}
			}
		}
		data["ctx"] = c.Ctx //data中需要有ctx参数

		//如果没有填写模板文件,则使用模板字符串
		if m["template_file"] == "" {
			var tpl = NewTpl()

			//添加自定义函数
			tpl.Funcs(template.FuncMap{"str2html": beego.Str2html})
			tpl.Funcs(template.FuncMap{"funcLower": funcLower})
			tpl.Funcs(template.FuncMap{"funcUpper": funcUpper})
			tpl.Funcs(template.FuncMap{"funcBR": funcBR})
			tpl.Funcs(template.FuncMap{"funcMod": funcMod})
			tpl.Funcs(template.FuncMap{"funcMap": funcMap})

			tpl, er := tpl.Parse(m["template"])
			if er != nil {
				rst = ApiResult("500", "tpl.Parse代码解析错误", strings.Replace(er.Error(), `"`, " ", -1), "")
				c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
				c.Ctx.Output.Body([]byte(rst))
				return
			}
			var buf bytes.Buffer
			var e = tpl.Execute(&buf, data)

			if e != nil {
				rst = ApiResult("500", "tpl.Execute模板解析错误", strings.Replace(e.Error(), `"`, " ", -1), "")
				c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
				c.Ctx.Output.Body([]byte(rst))
				return
			}
			rst = buf.String()
		}
	}

	//--------------------------------------------------------------------
	if m["template_file"] == "" {
		c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))
	} else {
		c.TplName = m["template_file"]
	}
}

var notfound = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"> 
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no">
<title>404没有找到相关信息</title>
<script src="https://cdn.bootcss.com/jquery/1.12.2/jquery.min.js"></script>
<script src="https://cdn.bootcss.com/jquery.form/3.24/jquery.form.min.js"></script>
<link href="/js/jqweui/lib/weui.min.css" rel="stylesheet"> 
<link href="https://cdn.bootcss.com/jquery-weui/1.2.1/css/jquery-weui.min.css" rel="stylesheet">
<script src="https://cdn.bootcss.com/jquery-weui/1.2.1/js/jquery-weui.min.js"></script>
</head>
<body>
    <div class="weui-msg">
      <div class="weui-msg__icon-area"><i class="weui-icon-warn weui-icon_msg"></i></div>
      <div class="weui-msg__text-area">
        <h2 class="weui-msg__title">没有找到相关信息</h2>
        <p class="weui-msg__desc"></p>
      </div>
      <div class="weui-msg__extra-area">
        <div class="weui-footer">
          <p class="weui-footer__links">
             
          </p>
          <p class="weui-footer__text">Copyright © 2008-2019 sjzapps.com</p>
        </div>
      </div>
    </div>
</body>
</html>
`
