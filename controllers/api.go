package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//ApiController 控制器
type ApiController struct {
	beego.Controller
}

//Api返回结果结构
func ApiResult(code, msg, extra, result string) string {
	var rst = `{
	"code":` + code + `,
	"msg":"` + msg + `",
	"extra":"` + extra + `",
	"result":[` + result + `]
}`
	return rst
}

//Get 默认页
func (c *ApiController) Get() {
	var rst = ""
	var extra = ""

	var _date = time.Now().Format("2006-01-02")
	c.Data["_date"] = _date
	var _time = time.Now().Format("2006-01-02 15:04:05")
	c.Data["_time"] = _time

	var module = c.Ctx.Input.Param(":module")
	if module == "" {
		rst = ApiResult("404", "没有找到接口信息", "", "")
		c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))
		return
	}
	var api = c.Ctx.Input.Param(":api")
	if api == "" {
		rst = ApiResult("405", "没有找到接口信息", "", "")
		c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))
		return
	}

	//查找API接口信息
	var m = db.First("select * from api_list where module=? and api_code=?", module, api)
	if m == nil {
		rst = ApiResult("104", "没有找到接口信息", "", "")
		c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))
		return
	} else {
		//定义参数map
		var data map[string]interface{} = map[string]interface{}{}
		data["_date"] = _date
		data["_time"] = _time

		//填充登录参数
		var _uid = c.GetSession("_uid")
		var _mch_id = c.GetSession("_mch_id")
		var _username = c.GetSession("_username")
		var _usertype = c.GetSession("_usertype")
		data["_uid"] = _uid
		data["_mch_id"] = _mch_id
		data["_username"] = _username
		data["_usertype"] = _usertype

		//默认所有参数的初始化
		var _ = c.GetString("id") //初始化
		var _m = c.Ctx.Request.Form
		for k, v := range _m {
			data[k] = v
		}
		//从数据库读取接口的参数,并从客户端获取参数值
		var params = db.Query("select * from api_param where api_id=? and state=1", m["id"])

		//fmt.Println("--------11111111111111111")
		//数据库连接
		var XX *xorm.Engine
		var conn = m["conn_str"]
		c.SetSession("_conn", conn) //将连接串保存的session中,便于模板中获取调用
		if conn != "" {
			XX = db.NewDb(conn)
		}
		//fmt.Println("--------db:", XX)
		//普通参数接收
		for _, v := range params {
			if v["param_type"] == "参数" {
				var p = c.GetString(v["param_name"])
				if p == "" && v["param_value"] != "" {
					p = v["param_value"]
				}
				// if p != "" {
				// 	data[v["param_name"]] = p
				// }
				data[v["param_name"]] = p

				//验证必填
				if p == "" && v["is_require"] == "1" {
					rst = ApiResult("104", v["title"]+"参数缺失", "", "")
					c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
					c.Ctx.Output.Body([]byte(rst))
					return
				}
				//验证唯一
				var id = c.GetString("id")
				if id == "" {
					id = "0"
				}
				if conn != "" && p != "" && v["is_unique"] == "1" {
					var l = db.Query2(XX, "select * from "+v["is_unique_tb"]+" where id='"+id+"' "+v["param_name"]+"='"+p+"' limit 1")
					if len(l) > 0 {
						rst = ApiResult("104", v["title"]+"["+p+"]已存在", "", "")
						c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
						c.Ctx.Output.Body([]byte(rst))
						return
					}
				}
				//验证值是否符合要求
				var info = ""
				if v["is_checkout"] == "1" {
					var t = v["is_checkout_type"]
					if t == "=" {
						if p != v["is_checkout_val"] {
							info = v["is_checkout_info"]
						}
					}
					if t == "!=" {
						if p == v["is_checkout_val"] {
							info = v["is_checkout_info"]
						}
					}
					if t == ">" {
						if p <= v["is_checkout_val"] {
							info = v["is_checkout_info"]
						}
					}
					if t == ">=" {
						if p < v["is_checkout_val"] {
							info = v["is_checkout_info"]
						}
					}
					if t == "<" {
						if p >= v["is_checkout_val"] {
							info = v["is_checkout_info"]
						}
					}
					if t == "<=" {
						if p > v["is_checkout_val"] {
							info = v["is_checkout_info"]
						}
					}
					if t == "正则" {
						match, _ := regexp.MatchString("p([a-z]+)ch", "peddach")
						if match == false {
							info = v["is_checkout_info"]
						}
					}
				}
				if info != "" {
					rst = ApiResult("104", info, "", "")
					c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
					c.Ctx.Output.Body([]byte(rst))
					return
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
				// if p != "" {
				// 	data[v["param_name"]] = p
				// }
				data[v["param_name"]] = p
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
				// if p != "" {
				// 	data[v["param_name"]] = p
				// }
				data[v["param_name"]] = p
			}
		}
		//SQL执行结果
		if conn != "" {
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
						rst = ApiResult("500", "参数错误", e.Error(), "")
						c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
						c.Ctx.Output.Body([]byte(rst))
						return
					}
					if m["conn_str"] != "" {
						if strings.Contains(buf.String(), "limit 1") || strings.Contains(buf.String(), "top 1") {
							var p = db.First2(XX, buf.String())
							if p != nil {
								data[v["param_name"]] = p
							}
						} else {
							var p = db.Query2(XX, buf.String())
							if p != nil {
								data[v["param_name"]] = p
							}
						}

					}

				} else if v["param_type"] == "SQLEXEC" || v["param_type"] == "sqlexec" { //如果参数类型是sql
					var tpl = template.New("")
					tpl.Parse(v["param_value"])
					var buf bytes.Buffer
					var e = tpl.Execute(&buf, data)
					if e != nil {
						fmt.Println("tpl.Execute sqlexec错误:", e.Error())
						rst = ApiResult("500", "参数错误", e.Error(), "")
						c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
						c.Ctx.Output.Body([]byte(rst))
						return
					}
					if m["conn_str"] != "" {
						var p = db.Exec2(XX, buf.String())
						data[v["param_name"]] = p
						if v["param_name"] == "extra" {
							extra = strconv.FormatInt(p, 10)
						}
					}
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
				// if p != "" {
				// 	data[v["param_name"]] = p
				// }
				data[v["param_name"]] = p
			}
		}

		var tpl = NewTpl()
		data["ctx"] = c.Ctx //data中需要有ctx参数

		tpl, er := tpl.Parse(m["api_template"])
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
		c.Data["extra"] = extra
		//rst=ApiResult("100","接口调用成功!",extra,buf.String())
		//所有格式全部由模板控制
		rst = buf.String()
	}
	rst = strings.TrimSpace(rst)
	// 去除空格
	rst = strings.Replace(rst, "	}", "}", -1)
	rst = strings.Replace(rst, "	", " ", -1)
	// 去除换行符
	//rst = strings.Replace(rst, "\n", "", -1)
	//--------------------------------------------------------------------
	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	c.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Ctx.Output.Body([]byte(rst))
}

// //Post 默认页
// func (c *ApiController) Post() {
// 	var rst=""
// 	var extra="";

// 	var module = c.Ctx.Input.Param(":module")
// 	if module == "" {
// 		rst=ApiResult("404","没有找到接口信息","","")
// 		return
// 	}
// 	var api = c.Ctx.Input.Param(":api")
// 	if api == "" {
// 		rst=ApiResult("405","没有找到接口信息","","")
// 		return
// 	}
// 	//查找API接口信息
// 	var m=db.First("select * from api_list where module=? and api_code=?",module,api)
// 	if m==nil{
// 		rst=ApiResult("104","没有找到接口信息","","")
// 	}else{
// 		//定义参数map
// 		var data map[string]interface{} = map[string]interface{}{}
// 		//填充登录参数
// 		var _uid=c.GetSession("_uid")
// 		var _mch_id=c.GetSession("_mch_id")
// 		var _username=c.GetSession("_username")
// 		var _usertype=c.GetSession("_usertype")
// 		data["_uid"]=_uid
// 		data["_mch_id"]=_mch_id
// 		data["_username"]=_username
// 		data["_usertype"]=_usertype

// 		//从数据库读取接口的参数,并从客户端获取参数值
// 		var params=db.Query("select * from api_param where api_id=? and state=1",m["id"])
// 		//fmt.Println("params:",params)

// 		//普通参数接收
// 		for _,v:=range params{
// 			if v["param_type"]=="参数"{
// 				var p=c.GetString(v["param_name"])
// 				if p=="" && v["param_value"]!=""{
// 					p=v["param_value"]
// 				}
// 				if p!=""{
// 					data[v["param_name"]]=p
// 				}
// 			}
// 		}
// 		//SQL执行结果
// 		for _,v:=range params{
// 			if v["param_type"]=="SQL" || v["param_type"]=="sql"{   //如果参数类型是sql

// 				var tpl=template.New("")
// 				tpl.Parse(v["param_value"])
// 				var buf bytes.Buffer
// 				var e=tpl.Execute(&buf, data)
// 				if e!=nil{
// 					fmt.Println("tpl.Execute sql错误:",e.Error())
// 					rst=ApiResult("500","参数错误",e.Error(),"")
// 					c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
// 					c.Ctx.Output.Body([]byte(rst))
// 					return
// 				}
// 				if m["conn_str"]!=""{
// 					var xx = db.NewDb(m["conn_str"])

// 					if strings.Contains(buf.String(),"limit 1") ||  strings.Contains(buf.String(),"top 1"){
// 						var p=db.First2(xx,buf.String())
// 						if p!=nil{
// 							data[v["param_name"]]=p
// 						}
// 					}else{
// 						var p=db.Query2(xx,buf.String())
// 						if p!=nil{
// 							data[v["param_name"]]=p
// 						}
// 					}

// 				}

// 			}else if v["param_type"]=="SQLEXEC" || v["param_type"]=="sqlexec"{   //如果参数类型是sql
// 				var tpl=template.New("")
// 				tpl.Parse(v["param_value"])
// 				var buf bytes.Buffer
// 				var e=tpl.Execute(&buf, data)
// 				if e!=nil{
// 					fmt.Println("tpl.Execute sqlexec错误:",e.Error())
// 					rst=ApiResult("500","参数错误",e.Error(),"")
// 					c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
// 					c.Ctx.Output.Body([]byte(rst))
// 					return
// 				}
// 				if m["conn_str"]!=""{
// 					var xx = db.NewDb(m["conn_str"])
// 					var p=db.Exec2(xx,buf.String())
// 					data[v["param_name"]]=p
// 					if v["param_name"]=="extra" {
// 						extra=strconv.FormatInt(p,10)
// 					}
// 				}
// 			}
// 		}

// 		var tpl=template.New("")
// 		//添加自定义函数
// 		tpl.Funcs(template.FuncMap{"CookieGet": funcCookieGet})
// 		tpl.Funcs(template.FuncMap{"CookieSet": funcCookieSet})
// 		tpl.Funcs(template.FuncMap{"SessionGet": funcSessionGet})
// 		tpl.Funcs(template.FuncMap{"SessionSet": funcSessionSet})

// 		tpl.Parse(m["api_template"])
// 		var buf bytes.Buffer
// 		var e=tpl.Execute(&buf, data)

// 		if e!=nil{
// 			fmt.Println("tpl.Execute 错误:",e.Error())
// 			rst=ApiResult("500","模板解析错误",e.Error(),"")
// 			c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
// 			c.Ctx.Output.Body([]byte(rst))
// 			return
// 		}
// 		c.Data["extra"]=extra
// 		//rst=ApiResult("100","接口调用成功!",extra,buf.String())
// 		//所有格式全部由模板控制
// 		rst=buf.String()
// 	}

// 	//--------------------------------------------------------------------
// 	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
// 	c.Ctx.Output.Body([]byte(rst))
// }
// func (c *ApiController) Test() {
// 	c.Ctx.SetCookie("hello", "world22222")
// 	//c.Data["c"]=c
// 	//c.Data["ctx"]=c.Ctx
// 	//beego.AddFuncMap("CookieGet",funcCookieGet)
// 	c.TplName = "test.html"
// }
// func (c *ApiController) Test2() {
// 	c.Ctx.SetCookie("hello", "abcd123456")
// 	var abc = `hello{{CookieGet .ctx "hello"}}
// 	<br/>
// 	hello2{{CookieSet .ctx "a"  "hello00000"}}
// 	`

// 	var tpl = NewTpl()
// 	tpl, er := tpl.Parse(abc)
// 	if er != nil {
// 		c.Ctx.WriteString("error0:" + er.Error())
// 		return
// 	}

// 	var buf bytes.Buffer
// 	var e = tpl.Execute(&buf, c.Data)
// 	if e != nil {
// 		c.Ctx.WriteString("error:" + e.Error())
// 		return
// 	}
// 	c.Ctx.WriteString(buf.String())
// }
