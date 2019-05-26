package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/sjzdlm/db"

	//"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) UpgradeIE() {
	c.TplName = "upgradeie.html"
}

//Get 默认首页 1.从配置文件读取域名 demo.sjzkakq.com=/p/m/index 2.根据index和域名查找 3.根据index查找模块和code
func (c *MainController) Get() {
	var rst = ""
	var module = "index"
	var code = "index"
	var domain = c.Ctx.Input.Domain()
	//c.Ctx.WriteString("domain:"+domain)
	//先查找是否有配置文件进行域名跳转
	fmt.Println("domain:", domain)
	var url = beego.AppConfig.String(domain)
	if url != "" {
		c.Redirect(url, 301)
		return
	}
	//查找页面信息
	var m = db.FirstOrNil("select * from page_list where  code=? and token=?", code, domain)
	if m == nil {
		m = db.FirstOrNil("select * from page_list where module=? and  code=?", module, code)
	}
	if m == nil {
		rst = "网站正在建设中..."
		c.Ctx.WriteString(rst)
		return
	} else {
		//定义参数map
		var data map[string]interface{} = map[string]interface{}{}

		//从数据库读取报表的参数,并从客户端获取参数值
		var params = db.Query("select * from page_param where page_id=? and state=1", m["id"])

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

		//SQL执行结果
		var conn = m["conn_str"]
		c.SetSession("_conn", conn) //将连接串保存的session中,便于模板中获取调用
		if conn != "" {
			var XX = db.NewDb(conn)
			for _, v := range params {
				if v["param_type"] == "SQL" || v["param_type"] == "sql" { //如果参数类型是sql

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
						}
					} else {
						var p = db.Query2(XX, sql)
						if p != nil {
							data[v["param_name"]] = p
						}
					}
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

		// fmt.Println("页面param:",data)
		// var tpl=template.New("")

		var tpl = NewTpl()
		data["ctx"] = c.Ctx //data中需要有ctx参数

		tpl.Parse(m["template"])
		var buf bytes.Buffer
		var e = tpl.Execute(&buf, data)

		if e != nil {
			fmt.Println("tpl.Execute 错误:", e.Error())
			rst = "网站正在建设中....."
			return
		}

		rst = buf.String()
	}
	rst = strings.TrimSpace(rst)
	// 去除空格
	rst = strings.Replace(rst, " ", "", -1)
	// 去除换行符
	rst = strings.Replace(rst, "\n", "", -1)
	//--------------------------------------------------------------------
	//c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//-------------------------------------------------------------------------------
func NewTpl() *template.Template {
	var tpl = template.New("tpl")
	//添加自定义函数
	tpl.Funcs(template.FuncMap{"CookieGet": funcCookieGet})
	tpl.Funcs(template.FuncMap{"CookieSet": funcCookieSet})
	tpl.Funcs(template.FuncMap{"SessionGet": funcSessionGet})
	tpl.Funcs(template.FuncMap{"SessionSet": funcSessionSet})

	tpl.Funcs(template.FuncMap{"sqlget": funcSQLGet})
	tpl.Funcs(template.FuncMap{"sqlquery": funcSQLQuery})
	tpl.Funcs(template.FuncMap{"sqlexec": funcSQLExec})

	tpl.Funcs(template.FuncMap{"getstring": funcGetString})
	tpl.Funcs(template.FuncMap{"apijson": funcApiJson})

	tpl.Funcs(template.FuncMap{"str2html": beego.Str2html})

	return tpl
}

//获取cookie值
func funcCookieGet(ctx *context.Context, str string) string {

	var rst = ctx.GetCookie(str)
	fmt.Println("GetCookie:", str, rst)
	return rst
}

//设置cookie值
func funcCookieSet(ctx *context.Context, str string, val string) string {
	ctx.SetCookie(str, val)
	fmt.Println("SetCookie:", str, val)
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

//执行SQL查询,返回单条数据
func funcSQLGet(ctx *context.Context, sql string) map[string]string {
	var c = ctx.Input.CruSession.Get("_conn")
	var conn = ""
	if c == nil {
		return nil
	} else {
		conn = c.(string)
	}
	if conn == "" {
		return nil
	}
	var XX = db.NewDb(conn)

	var tpl = template.New("")
	tpl.Parse(sql)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, ctx.Input.Data)
	if e != nil {
		fmt.Println("funcSQLGet sql错误:", e.Error())
		return nil
	}
	sql = buf.String()
	fmt.Println(sql)
	var rst = db.First2(XX, sql)
	return rst
}

//执行SQL查询,返回多条数据
func funcSQLQuery(ctx *context.Context, sql string) []map[string]string {
	var c = ctx.Input.CruSession.Get("_conn")
	var conn = ""
	if c == nil {
		return nil
	} else {
		conn = c.(string)
	}
	if conn == "" {
		return nil
	}
	var XX = db.NewDb(conn)
	var rst = db.Query2(XX, sql)
	return rst
}

//执行SQL语句
func funcSQLExec(ctx *context.Context, sql string) int64 {
	//var conn=ctx.GetCookie("_conn")
	var c = ctx.Input.CruSession.Get("_conn")
	var conn = ""
	if c == nil {
		return 0
	} else {
		conn = c.(string)
	}
	if conn == "" {
		return 0
	}
	var XX = db.NewDb(conn)
	fmt.Println(sql)
	var i = db.Exec2(XX, sql)
	return i
}

//获取客户端参数值
func funcGetString(ctx *context.Context, key string) string {
	var v = ctx.Input.Query(key)
	return v
}

//返回API JSON格式数据
func funcApiJson(ctx *context.Context, code string, msg string, extra string) string {
	var v = `
	{
		"code":` + code + `,
		"msg":"` + msg + `",
		"extra":"` + extra + `",
		"result":[]
	}
	`
	return v
}
