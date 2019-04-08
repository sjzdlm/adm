package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//RptController 控制器
type RptController struct {
	beego.Controller
}

//Get 默认页
func (c *RptController) Get() {
	var rst = ""
	var module = c.Ctx.Input.Param(":module")
	if module == "" {
		rst = "没有找到报表信息"
		c.Ctx.WriteString(rst)
		return
	}
	var rpt = c.Ctx.Input.Param(":rpt")
	if rpt == "" {
		rst = "没有找到报表信息."
		c.Ctx.WriteString(rst)
		return
	}
	//查找rpt报表信息
	var m = db.First("select * from rpt_list where module=? and  code=?", module, rpt)
	if m == nil {
		rst = "没有找到报表信息."
		c.Ctx.WriteString(rst)
		return
	} else {
		//定义参数map
		var data map[string]interface{} = map[string]interface{}{}
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
		var params = db.Query("select * from rpt_param where rpt_id=? and state=1", m["id"])
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
						c.Ctx.WriteString("报表模板错误500!")
						return
					}

					var sql = buf.String()
					//< >号会被转义,此处替换回来
					sql = strings.Replace(sql, "&lt;", "<", -1)
					sql = strings.Replace(sql, "&gt;", ">", -1)

					// if strings.Contains(sql, "limit 1") || strings.Contains(sql, "top 1") {
					// 	var p = db.First2(XX, sql)
					// 	if p != nil {
					// 		data[v["param_name"]] = p
					// 	}
					// } else {
					// 	var p = db.Query2(XX, sql)
					// 	if p != nil {
					// 		data[v["param_name"]] = p
					// 	}
					// }

					var p = db.Query2(XX, sql)
					if p != nil {
						data[v["param_name"]] = p
					} else {
						data[v["param_name"]] = ""
					}
				} else if v["param_type"] == "SQLEXEC" || v["param_type"] == "sqlexec" { //如果参数类型是sql
					var tpl = template.New("")
					tpl.Parse(v["param_value"])
					var buf bytes.Buffer
					var e = tpl.Execute(&buf, data)
					if e != nil {
						fmt.Println("tpl.Execute sqlexec错误:", e.Error())
						c.Ctx.WriteString("报表模板错误500!")
						return
					}
					var p = db.Exec2(XX, buf.String())
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

		var tpl = NewTpl()
		data["ctx"] = c.Ctx //data中需要有ctx参数

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

	//--------------------------------------------------------------------
	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
