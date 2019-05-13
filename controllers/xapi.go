package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//非权限api数据调用专用
type XApiController struct {
	beego.Controller
}

//返回一个ID
func (c *XApiController) ID() {
	//var id = fmt.Sprintf("T%v", time.Now().UnixNano()/1e6)
	var t = ""
	var companyid = c.GetSession("_company_id")
	if companyid != nil {
		t = companyid.(string)
	}

	for i := len(t); i < 5; i++ {
		t = "0" + t
	}
	//var id = "T" + t + time.Now().Format("060102150405.000")
	var id = "T" + t + time.Now().Format("060102150405")

	c.Ctx.WriteString(id)
}

//页面导航列表-移动端专用
func (c *XApiController) NavListJson() {
	var pageid = c.GetString("pageid")
	if pageid == "" {
		pageid = "0"
	}
	//fmt.Println("page_id:", pageid)

	var sql = `select * from tbm_nav_list where page_id=? order by sort`
	var list = db.Query(sql, pageid)
	c.Data["json"] = list
	c.ServeJSON()
}

//ListJson数据
func (c *XApiController) ListJson() {

	var code = c.GetString("code")
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("")
		return
	}
	//根据参数读取数据
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""
	var where_date = ""
	var totalField = ""

	var data_type = tb["data_type"] //数据类型  0表 1视图  2语句
	var data_table = ""             //如果是语句则需要别名
	if data_type == "2" {
		data_table = "_tb"
	}
	//搜索条件
	var fields = db.Query("select * from tb_field where tbid=? and  is_search=1 ", tb["id"])
	qtxt = strings.Replace(strings.TrimSpace(string(qtxt)), "'", "''", -1)

	if fields != nil && len(fields) > 0 {
		where = ""
		for k, v := range fields {
			if k > 0 && where != "" {
				where += " or "
			}
			if v["form_type"] == "日期选择" && where_date == "" { //对于日期只取第一个日期搜索
				var start = c.GetString(v["field_code"] + "_start")
				var end = c.GetString(v["field_code"] + "_end")
				where_date = v["field_code"] + ` >= '` + start + `' and ` + v["field_code"] + `<='` + end + ` 23:59:59' `
			} else {
				if qtxt != "" {
					if v["form_type"] != "文本框" {
						where += v["field_code"] + ` like '` + qtxt + `' `
					} else {
						where += v["field_code"] + ` like '%` + qtxt + `%' `
					}

				}
			}

		}
		if where != "" {
			where = " where (" + where + ") "
		} else {
			where = " where 1=1 "
		}
		if where_date != "" {
			where += " and " + where_date
		}
	}

	// //合计字段
	// var hjzd = db.First("select * from tb_field where tbid=? and is_total=1 and view_list=1", tb["id"])
	// if len(hjzd) > 0 {
	// 	totalField = hjzd["field_code"]
	// }
	//其他搜索条件
	if fields != nil && len(fields) > 0 {
		for _, v := range fields {
			var q = c.GetString(v["field_code"])
			q = strings.Replace(strings.TrimSpace(string(q)), "'", "''", -1)
			if q != "" || v["search_require"] == "1" { ////必搜字段即使为空也要增加此字段条件
				if where != "" {
					where += " and "
				} else {
					where += " where "
				}
				if v["form_type"] != "文本框" {
					where += v["field_code"] + ` like '` + q + `' `
				} else {
					where += v["field_code"] + ` like '%` + q + `%' `
				}
			}

		}
	}

	//全局条件
	if tb["where_str"] != "" {
		if where == "" {
			where = " where " + tb["where_str"]
		} else {
			where += " and " + tb["where_str"]
		}
	}

	//排序
	var sort = c.GetString("sort")
	var order = c.GetString("order")
	var orderstr = ""
	if sort != "" && order != "" {
		orderstr = " order by " + sort + " " + order
	} else {
		if tb["pri_key"] != "" {
			orderstr = " order by " + tb["pri_key"] + " desc "
		}
	}

	//搜索资源,手机端修改字段别名
	var fs = db.Query("select * from tb_field where tbid=? and  view_list=1 ", tb["id"])
	var fstr = ""
	for _, v := range fs {
		if fstr != "" {
			fstr += ","
		}
		if v["view_list_place"] != "" {
			fstr += " " + v["field_code"] + " as " + v["view_list_place"]
		} else {
			fstr += " " + v["field_code"]
		}
	}
	if fstr == "" {
		fstr = " * "
	}

	var xx = db.NewDb(tb["conn_str"])
	var sql = "select " + fstr + " from " + tb["table"] + " " + where
	if data_type == "2" { //如果是语句的话
		sql = "select " + fstr + " from (" + tb["table"] + ") as" + data_table + " " + where
	}

	//判断数据库类型
	var con = db.First("select * from adm_conn where conn='" + tb["conn_str"] + "' limit 1")
	if strings.Contains(con["dbtype"], "mssql") {
		var rst = db.Pager2MsSql(xx, page, pageSize, tb["table"], sql, where, orderstr, sql+orderstr)
		c.Data["json"] = rst
	} else {
		fmt.Println(sql + orderstr)
		var rst = db.Pager2(xx, page, pageSize, sql+orderstr)
		rst.Extra = tb["title"] //功能标题
		if totalField != "" {
			var heji = db.Query2(xx, "select sum("+totalField+") as "+totalField+" from "+tb["table"]+" "+where+" "+orderstr)
			if len(heji) > 0 {
				rst.Footer = heji
			}
		}
		c.Data["json"] = rst

	}

	c.ServeJSON()
}

//FieldJson数据
func (c *XApiController) FieldJson() {
	var code = c.GetString("code")
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("{}")
		return
	}
	var fields = db.Query(`select 
	id,tbid,tbcode,field_name,field_code,field_type,field_length,
	field_defval,field_prikey,form_type,form_length,form_sort,form_tip,
	form_value,form_required,view_list ,view_form ,is_sort,is_search,
	is_total ,is_editable,is_readonly ,memo,state
	from tb_field where tbid=? and  view_form=1 `, tb["id"])

	//获取值
	var data_type = tb["data_type"] //数据类型  0表 1视图  2语句
	var data_table = ""             //如果是语句则需要别名
	if data_type == "2" {
		data_table = "_tb"
	}

	var id, _ = c.GetInt("id", 0)
	var where = " where id=?"
	var xx = db.NewDb(tb["conn_str"])
	var sql = "select * from " + tb["table"] + " " + where
	if data_type == "2" { //如果是语句的话
		sql = "select * from (" + tb["table"] + ") as" + data_table + " " + where
	}
	var m = db.First2(xx, sql, id)
	if len(m) > 0 {
		for _, row := range fields {
			row["_value"] = m[row["field_code"]]
		}
	}

	c.Data["json"] = fields
	c.ServeJSON()
}

//DataJson
func (c *XApiController) DataJson() {
	var code = c.GetString("code")
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("{}")
		return
	}
	var data_type = tb["data_type"] //数据类型  0表 1视图  2语句
	var data_table = ""             //如果是语句则需要别名
	if data_type == "2" {
		data_table = "_tb"
	}

	var id, _ = c.GetInt("id", 0)
	var where = " where id=?"
	var xx = db.NewDb(tb["conn_str"])
	var sql = "select * from " + tb["table"] + " " + where
	if data_type == "2" { //如果是语句的话
		sql = "select * from (" + tb["table"] + ") as" + data_table + " " + where
	}
	var m = db.First2(xx, sql, id)
	c.Data["json"] = m
	if len(m) < 1 {
		c.Data["json"] = "{}"
	}
	c.ServeJSON()
}

//专为复选框、单选框、下拉框返回自定义JSON
func (c *XApiController) ItemJson() {
	var id, _ = c.GetInt("id", 0)

	var f = db.First("select * from tb_field where id=?", id)
	if len(f) < 1 {
		c.Ctx.WriteString("{}")
		return
	}

	//获取表信息
	var tb = db.First("select * from tb_table where id=?", f["tbid"])
	if len(tb) < 1 {
		c.Ctx.WriteString("{}")
		return
	}
	var sql = f["form_value"]
	if sql == "" {
		c.Ctx.WriteString("{}")
		return
	}

	var jsonstr = ""
	if strings.Contains(f["form_value"], "select") == true {
		var xx = db.NewDb(tb["conn_str"])
		var list = db.Query2(xx, sql)
		// c.Data["json"] = list
		// c.ServeJSON()
		//jsonstr = `[{"id":"","val":"请选择..."}`
		jsonstr = `[`
		for kk, vv := range list {
			if kk > 0 {
				jsonstr += ","
			}
			jsonstr += `{"id":"` + vv["id"] + `","val":"` + vv["val"] + `"}`
		}
		jsonstr += `]`
	} else {
		var ls = strings.Split(f["form_value"], ";")
		if len(ls) > 0 {
			//jsonstr = `[{"id":"","val":"请选择..."}`
			jsonstr = `[`
			for kk, vv := range ls {
				if kk > 0 {
					jsonstr += ","
				}
				var lsb = strings.Split(vv, ",")
				if len(lsb) > 1 {
					jsonstr += `{"id":"` + lsb[0] + `",`
					jsonstr += `"val":"` + lsb[1] + `"}`
				} else {
					jsonstr += `"id":"` + vv + `",`
					jsonstr += `"val":"` + vv + `"`
				}
			}
			jsonstr += `]`
		}
	}
	if jsonstr == "" {
		c.Ctx.WriteString("{}")
	} else {
		c.Ctx.WriteString(jsonstr)
	}
}

//InfoPost 数据保存
func (c *XApiController) EditPost() {
	var id, _ = c.GetInt("_id", 0)
	c.Data["id"] = id

	var code = c.GetString("_code")
	if code == "" {
		c.Ctx.WriteString("0")
		return
	}
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("code data not found")
		return
	}
	var log = "" //log日志信息
	//
	var fields = db.Query("select * from tb_field where tbid=? and field_code!='id' and view_form=1 and form_type!='标签框'", tb["id"])
	var sql = ""
	if id > 0 {
		sql = `update  ` + tb["table"] + " set "
		for k, v := range fields {
			if k > 0 {
				sql += ","
			}
			var val = c.GetString(v["field_code"])
			val = strings.Replace(val, "'", "''", -1) //字符替换

			if v["form_type"] == "开关按钮" {
				if val == "on" {
					val = "1"
				} else {
					val = "0"
				}
			}
			if v["form_type"] == "复选框" {
				var vs = c.GetStrings(v["field_code"])
				if val == "on" {
					val = "1"
				} else if val == "" {
					val = "0"
				} else {
					val = ""
					if len(vs) > 0 {
						for kk, vv := range vs {
							if kk > 0 {
								val += ","
							}
							val += vv
						}
					}
				}

			}
			sql += " `" + v["field_code"] + "`='" + val + "' "

			log += v["field_code"] + "=" + val + ";\r\n"

			c.Data[v["field_code"]] = val
		}
		sql += " where " + tb["pri_key"] + "='" + strconv.Itoa(id) + "' "

		// //记录修改日志
		// var _mchid = c.GetSession("_mch_id")
		// var _uid = c.GetSession("_uid")
		// var _username = c.GetSession("_username")

		// var atime = time.Now().Format("2006-01-02 15:04:05")
		// var ip = c.Ctx.Request.RemoteAddr
		// db.Exec("insert into adm_log(mch_id,user_id,username,logtype,opertype,title,content,ip,addtime)values(?,?,?,?,?,?,?,?,?)",
		// 	_mchid.(string), _uid.(string), _username.(string), "操作日志", "修改", _username.(string)+"("+_uid.(string)+")修改["+tb["title"]+"]"+tb["table"]+"表记录,id="+strconv.Itoa(id), log, ip, atime,
		// )
	} else {
		sql = `insert into ` + tb["table"] + `(`
		var tmp = ""
		for k, v := range fields {
			if k > 0 {
				sql += ","
				tmp += ","
			}
			sql += "`" + v["field_code"] + "`"

			var val = c.GetString(v["field_code"])
			val = strings.Replace(val, "'", "''", -1) //字符替换

			if v["form_type"] == "开关按钮" {
				if val == "on" {
					val = "1"
				} else {
					val = "0"
				}
			}
			if v["form_type"] == "复选框" {
				var vs = c.GetStrings(v["field_code"])

				for kk, vv := range vs {
					if kk > 0 {
						val += ","
					}
					val += vv
				}
			}
			tmp += "'" + val + "'"

			c.Data[v["field_code"]] = val
		}
		sql += ")values(" + tmp + ")"
	}

	//fmt.Println("sql:", sql)

	var xx = db.NewDb(tb["conn_str"])
	var i = db.Exec2(xx, sql)
	if i > 0 {
		//保存成功,准备检查额外需执行代码
		var exlist = db.Query("select * from tb_table_ex where extype='sql' and expage='info' and explace='AFTER' and tbid=? and state=1", tb["id"])
		for _, v := range exlist {
			var tpl = template.New("")
			if v["excontent"] == "" {
				continue
			}
			tpl.Parse(v["excontent"])
			var buf bytes.Buffer
			var e = tpl.Execute(&buf, c.Data)
			if e != nil {
				fmt.Println("tb_table_ex执行错误:", e.Error())
				continue
			}
			var p = buf.String()
			db.Exec2(xx, p)
		}

		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}

}

//上传图片base64编码
func (c *XApiController) UpBase64() {
	var b = c.GetString("b")
	b = strings.Replace(b, "data:image/jpeg;base64,", "", -1)
	var path = beego.AppPath
	//fmt.Println("b:", b)
	//fmt.Println("path:", path)
	ddd, er := base64.StdEncoding.DecodeString(b) //成图片文件并把文件写入到buffer
	if er != nil {
		c.Ctx.WriteString("")
		fmt.Println("er:", er.Error())
		return
	}
	var fn = fmt.Sprintf("/upload/%d.jpg", time.Now().UnixNano()/1e6)
	err := ioutil.WriteFile(path+"/static"+fn, ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		c.Ctx.WriteString("")
		fmt.Println("error:", err.Error())
		return
	}
	c.Ctx.WriteString(fn)
}

//返回街道json;code 为区县代码
func (c *XApiController) StreetJson() {
	var code = c.GetString("code")

	var sql = `select * from tb_region where substr(code,1,6)=? and length(code) =9 and pid!=0 order by code `
	var list = db.Query(sql, code)
	c.Data["json"] = list
	c.ServeJSON()
}

//返回区县json;code 为城市代码
func (c *XApiController) CountyJson() {
	var code = c.GetString("code")

	var sql = `select * from tb_region where substr(code,1,4)||'00'=?  and substr(code,5,6)!='00' and length(code) =6 and pid!=0 order by code `
	var list = db.Query(sql, code)
	c.Data["json"] = list
	c.ServeJSON()
}

//返回城市json;
func (c *XApiController) CityJson() {
	var code = c.GetString("code")

	var sql = `select * from tb_region where substr(code,1,2)||'0000'=?  and substr(code,5,6)='00' and length(code) =6 and pid!=0 order by code
	`
	var list = db.Query(sql, code)
	c.Data["json"] = list
	c.ServeJSON()
}

//返回省份json;
func (c *XApiController) ProvJson() {
	var sql = `select * from tb_region where substr(code,3,6)='0000' order by code 
	`
	var list = db.Query(sql)
	c.Data["json"] = list
	c.ServeJSON()
}

//返回街道用户类型JSON
func (c *XApiController) UserTypeJson() {
	var where = ""
	var level, _ = c.GetInt("level", 0)
	var ignore = c.GetString("ignore")

	where = "where level>=" + strconv.Itoa(level)
	//读取 > >= < <= 参数
	if ignore == "" {
		var rt, _ = c.GetInt("rt", 0)
		if rt > 0 {
			where += " and level>" + strconv.Itoa(rt)
		}
		var rte, _ = c.GetInt("rte", 0)
		if rte > 0 {
			where += " and level>=" + strconv.Itoa(rte)
		}
		var lt, _ = c.GetInt("lt", 0)
		if lt > 0 {
			where += " and level<" + strconv.Itoa(lt)
		}
		var lte, _ = c.GetInt("lte", 0)
		if lte > 0 {
			where += " and level<=" + strconv.Itoa(lte)
		}
	}

	if where == "" {
		where = " where state=1 "
	} else {
		where += " and state=1 "
	}
	var sql = `select mch_id,level as id, name as val from adm_usertype ` + where
	fmt.Println("sql usertype:", sql)
	var list = db.Query(sql)
	c.Data["json"] = list
	c.ServeJSON()
}

//用户类型JS对象,主要用于下拉框的列表绑定显示
func (c *XApiController) JsonUserType() {
	//账号类型
	var list = db.Query("select * from adm_usertype where state=1 order by orders ")
	var jsonstr = `var jsonutype={ `
	for kk, vv := range list {
		if kk > 0 {
			jsonstr += ","
		}
		jsonstr += `"key` + vv["level"] + `":"` + vv["name"] + `"`
	}
	jsonstr += `};
	`

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Body([]byte(jsonstr))
}

//根据级别调出角色的复选框html
func (c *XApiController) RoleHtmlList() {
	if c.GetSession("_uid") == nil {
		c.Ctx.WriteString("")
		fmt.Println("xxxxxxxxxxxxxxxxxx---------------------------")
		return
	}
	var _userlevel = c.GetSession("_userlevel").(string)

	var level = c.GetString("level")

	var roles = c.GetSession("_roles").(string)
	var where = " "
	if _userlevel == level {
		where = " where  id in(" + roles + ") "
	} else {
		where = " where level >= " + level
	}
	if c.GetSession("_sproot").(string) == "1" {
		where = ""
	}
	where = "select * from adm_role " + where
	fmt.Println("xapi role:", where)
	var list = db.Query(where)
	var rst = ` `
	for _, vv := range list {
		rst += `
		<input type="checkbox" name="role" id="role` + vv["id"] + `" value="` + vv["id"] + `"   /><label for="role` + vv["id"] + `">` + vv["name"] + `</label>
		`
	}
	rst += `
	`

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))
}
