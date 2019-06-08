package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/draw"

	"github.com/golang/freetype"

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
	} else {
		t = c.GetString("_company_id")
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

//ListJson数据-----移动端专用
func (c *XApiController) ListJson() {
	//遍历所有get参数信息放到模板变量--------------------------------
	var paramstr = ""
	var urls = strings.Split(c.Ctx.Input.URI(), "?")
	if len(urls) > 1 {
		var params = strings.Split(urls[1], "&")
		for i := 0; i < len(params); i++ {
			var p = strings.Split(params[i], "=")
			p[1], _ = url.QueryUnescape(p[1])
			c.Data[p[0]] = p[1]
			if paramstr != "" {
				paramstr += ","
			}
			paramstr += "&" + p[0] + "=" + p[1]
		}
	}
	c.Data["_paramstr"] = paramstr
	fmt.Println("paramstr:", paramstr)
	//------------------------------------------------------------

	var code = c.GetString("code")
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil || len(tb) < 1 {
		c.Ctx.WriteString("")
		return
	}
	//数据库连接配置
	var con = db.First("select * from adm_conn where conn='" + tb["conn_str"] + "' limit 1")
	if len(con) < 1 {
		c.Ctx.WriteString("")
		return
	}

	//级联传递值 jlfield jlval
	var jlfield = c.GetString("jlfield")
	var jlval = c.GetString("jlval")
	if jlfield != "" && jlval != "" {
		c.Data[jlfield] = jlval
	}
	//------------------------------------------------------------------
	var _me = c.GetSession("_me")
	c.Data["_me"] = _me
	if c.GetSession("_uid") != nil {
		c.Data["_uid"] = c.GetSession("_uid").(string)
	}
	if c.GetSession("_username") != nil {
		c.Data["_username"] = c.GetSession("_username").(string)
	}

	if c.GetSession("_usertype") != nil {
		c.Data["_usertype"] = c.GetSession("_usertype").(string)
	}

	if c.GetSession("_userlevel") != nil {
		c.Data["_userlevel"] = c.GetSession("_userlevel").(string)
	}

	if c.GetSession("_company") != nil {
		c.Data["_company"] = c.GetSession("_company").(string)
	}

	if c.GetSession("_company_id") != nil {
		c.Data["_company_id"] = c.GetSession("_company_id").(string)
	}

	if c.GetSession("_company_pid") != nil {
		c.Data["_company_pid"] = c.GetSession("_company_pid").(string)
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
						where += v["field_code"] + ` like '%` + qtxt + `%' `
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
		if tb["sort_key"] != "" {
			orderstr = " order by " + tb["sort_key"]
		}
	}

	//搜索资源,手机端修改字段别名

	var fs = db.Query("select * from tb_field where tbid=? and  view_list=1 ", tb["id"])

	//取值并输出到模板
	for _, v := range fs {
		if c.GetString(v["field_code"]) != "" {
			c.Data[v["field_code"]] = c.GetString(v["field_code"])
		}
	}

	var fstr = ""
	if strings.Contains(con["dbtype"], "mysql") {
		var fs = db.Query("select * from tb_field where tbid=? and  view_list=1 and (view_list_place !='' or view_list_isoptionid =1) order by view_list_place", tb["id"])
		//图片
		var fstr_tmp = ""
		for _, v := range fs {
			if v["view_list_place"] != "img" {
				continue
			}
			if fstr_tmp != "" {
				fstr_tmp += ","
			}
			fstr_tmp += v["field_code"]
			break
		}
		if fstr_tmp != "" {
			if fstr != "" {
				fstr += ","
			}
			fstr += "IFNULL(CONCAT(" + fstr_tmp + "),'/images/pic6.png') as img"
		}
		//标题
		fstr_tmp = ""
		for _, v := range fs {
			if v["view_list_place"] != "title" {
				continue
			}
			if fstr_tmp != "" {
				fstr_tmp += ","
			}
			fstr_tmp += v["field_code"]
		}
		if fstr_tmp != "" {
			if fstr != "" {
				fstr += ","
			}
			fstr += "CONCAT(" + fstr_tmp + ") as title"
		}
		//描述
		fstr_tmp = ""
		for _, v := range fs {
			if v["view_list_place"] != "desc" {
				continue
			}
			if fstr_tmp != "" {
				fstr_tmp += ","
			}
			fstr_tmp += " '" + v["field_name"] + ":'," + v["field_code"] + ",';'"
		}
		if fstr_tmp != "" {
			if fstr != "" {
				fstr += ","
			}
			fstr += "CONCAT(" + fstr_tmp + ") as `desc`"
		}
		//右角标
		fstr_tmp = ""
		for _, v := range fs {
			if v["view_list_place"] != "rss" {
				continue
			}
			if fstr_tmp != "" {
				fstr_tmp += ","
			}
			fstr_tmp += v["field_code"]
		}
		if fstr_tmp != "" {
			if fstr != "" {
				fstr += ","
			}
			fstr += "CONCAT(" + fstr_tmp + ") as rss"
		}
		// //选项ID,此ID保存到前端本地
		// fstr_tmp = ""
		// for _, v := range fs {
		// 	if v["view_list_place"] != "chooseid" {
		// 		continue
		// 	}
		// 	if fstr_tmp != "" {
		// 		fstr_tmp += ","
		// 	}
		// 	fstr_tmp += v["field_code"]
		// }
		// if fstr_tmp != "" {
		// 	if fstr != "" {
		// 		fstr += ","
		// 	}
		// 	fstr += "CONCAT(" + fstr_tmp + ") as chooseid"
		// }

		//选项ID,此ID保存到前端本地
		fstr_tmp = ""
		for _, v := range fs {
			//fmt.Println(".......................view_list_isoptionid.............A", v["view_list_isoptionid"])
			if v["view_list_isoptionid"] != "1" {
				continue
			}
			//fmt.Println(".......................view_list_isoptionid.............B")
			if fstr_tmp != "" {
				fstr_tmp += ","
			}
			fstr_tmp += v["field_code"]
		}
		if fstr_tmp != "" {
			if fstr != "" {
				fstr += ","
			}
			fstr += "CONCAT(" + fstr_tmp + ") as _optionid"
		}

		if fstr == "" {
			fstr = " * "
		} else {
			fstr = " id," + fstr
		}
	} else {
		for _, v := range fs {
			if fstr != "" {
				fstr += ","
			}
			if v["view_list_place"] != "" {
				fstr += " " + v["field_code"] + " as `" + v["view_list_place"] + "` "
			} else {
				fstr += " " + v["field_code"]
			}
		}
		if fstr == "" {
			fstr = " * "
		}
	}

	var xx = db.NewDb(tb["conn_str"])
	var sql = "select " + fstr + " from " + tb["table"] + " " + where
	if data_type == "2" { //如果是语句的话
		sql = "select " + fstr + " from (" + tb["table"] + ") as " + data_table + " " + where
	}

	//进行sql模板替换操作
	var tpl = template.New("")
	tpl.Parse(sql)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e != nil {
		fmt.Println("xapi template sql 执行错误A:", e.Error())
		c.Ctx.WriteString("{}")
		return
	}
	sql = buf.String()

	tpl = template.New("")
	tpl.Parse(where)
	e = tpl.Execute(&buf, c.Data)
	if e != nil {
		fmt.Println("xapi template sql 执行错误B:", e.Error())
		c.Ctx.WriteString("{}")
		return
	}
	where = buf.String()
	//判断数据库类型
	if strings.Contains(con["dbtype"], "mssql") {
		var rst = db.Pager2MsSql(xx, page, pageSize, tb["table"], sql, where, orderstr, sql+orderstr)
		c.Data["json"] = rst
	} else {
		fmt.Println(sql + orderstr)
		var rst = db.Pager2(xx, page, pageSize, sql+orderstr)
		rst.Extra = tb["title_list"] //功能标题
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

//获取模块信息
func (c *XApiController) TbJson() {
	var code = c.GetString("code")
	var tb = db.First("select id,code,`table`,is_create,is_search,is_edit,is_del,list_link from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("{}")
		return
	}

	//------------------------------------------------------------------
	var _me = c.GetSession("_me")
	c.Data["_me"] = _me
	if c.GetSession("_uid") != nil {
		c.Data["_uid"] = c.GetSession("_uid").(string)
	}
	if c.GetSession("_username") != nil {
		c.Data["_username"] = c.GetSession("_username").(string)
	}

	if c.GetSession("_usertype") != nil {
		c.Data["_usertype"] = c.GetSession("_usertype").(string)
	}

	if c.GetSession("_userlevel") != nil {
		c.Data["_userlevel"] = c.GetSession("_userlevel").(string)
	}

	if c.GetSession("_company") != nil {
		c.Data["_company"] = c.GetSession("_company").(string)
	}

	if c.GetSession("_company_id") != nil {
		c.Data["_company_id"] = c.GetSession("_company_id").(string)
	}

	if c.GetSession("_company_pid") != nil {
		c.Data["_company_pid"] = c.GetSession("_company_pid").(string)
	}
	//进行模板替换操作
	var tpl = template.New("")
	tpl.Parse(tb["list_link"])
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e == nil {
		tb["list_link"] = buf.String()
	}

	c.Data["json"] = tb
	c.ServeJSON()
}

//获取标题信息
func (c *XApiController) TitleJson() {
	var code = c.GetString("code")
	var tb = db.First("select id,code,title,title_list,title_edit from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("{}")
		return
	}
	c.Data["json"] = tb
	c.ServeJSON()
}

//获取模块扩展信息
func (c *XApiController) ExTable() {
	var code = c.GetString("code")
	var tbid, _ = c.GetInt("tbid", 0)
	var extype = c.GetString("extype")
	var expage = c.GetString("expage")
	var explace = c.GetString("explace")
	var id = c.GetString("id")
	c.Data["id"] = id

	var tb = db.First("select excontent from tb_table_ex where (code=? or tbid=?) and extype=? and expage=? and explace=?", code, tbid, extype, expage, explace)
	if tb == nil {
		c.Ctx.WriteString("")
		return
	}
	var rst = tb["excontent"]

	var tpl = template.New("")
	tpl.Parse(rst)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e == nil {
		rst = buf.String()
	}

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))
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
	form_value,form_required,view_list,view_list_module,view_list_islink,is_hide,view_form ,is_sort,is_search,
	is_total ,is_editable,is_readonly,form_cascade ,memo,state
	from tb_field where tbid=? and  view_form=1 order by form_sort`, tb["id"])

	//获取值
	var data_type = tb["data_type"] //数据类型  0表 1视图  2语句
	var data_table = ""             //如果是语句则需要别名
	if data_type == "2" {
		data_table = "_tb"
	}

	//------------------------------------------------------------------
	var _me = c.GetSession("_me")
	c.Data["_me"] = _me
	if c.GetSession("_uid") != nil {
		c.Data["_uid"] = c.GetSession("_uid").(string)
	}
	if c.GetSession("_username") != nil {
		c.Data["_username"] = c.GetSession("_username").(string)
	}

	if c.GetSession("_usertype") != nil {
		c.Data["_usertype"] = c.GetSession("_usertype").(string)
	}

	if c.GetSession("_userlevel") != nil {
		c.Data["_userlevel"] = c.GetSession("_userlevel").(string)
	}

	if c.GetSession("_company") != nil {
		c.Data["_company"] = c.GetSession("_company").(string)
	}

	if c.GetSession("_company_id") != nil {
		c.Data["_company_id"] = c.GetSession("_company_id").(string)
	}

	if c.GetSession("_company_pid") != nil {
		c.Data["_company_pid"] = c.GetSession("_company_pid").(string)
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
	} else {
		//如果纪录不存在则使用默认值
		var tpl = template.New("")

		for _, row := range fields {
			var defval = row["field_defval"]
			if defval != "" {
				tpl = template.New("")
				tpl.Parse(defval)
				var buf bytes.Buffer
				var e = tpl.Execute(&buf, c.Data)
				if e == nil {
					defval = buf.String()
				}
				//如果是网址则调用获取网页内容
				if defval != "" && len(defval) > 2 && (strings.HasPrefix(defval, "/") || strings.HasPrefix(defval, "http")) {
					var url = "http://" + c.Ctx.Input.Domain()
					if strings.HasPrefix(defval, "/") {
						if strconv.Itoa(c.Ctx.Input.Port()) != "80" {
							url = url + ":" + strconv.Itoa(c.Ctx.Input.Port())
						}
						url += defval
					} else if strings.HasPrefix(defval, "http") {
						url = defval
					}
					rst, e := HttpGet(url)
					if e == nil {
						defval = rst
					}
				}
				row["_value"] = defval
			} else {
				row["_value"] = ""
			}
			defval = ""
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
		c.Ctx.WriteString("[]")
		return
	}

	//------------------------------------------------------------------
	var _me = c.GetSession("_me")
	c.Data["_me"] = _me
	if c.GetSession("_uid") != nil {
		c.Data["_uid"] = c.GetSession("_uid").(string)
	}
	if c.GetSession("_username") != nil {
		c.Data["_username"] = c.GetSession("_username").(string)
	}

	if c.GetSession("_usertype") != nil {
		c.Data["_usertype"] = c.GetSession("_usertype").(string)
	}

	if c.GetSession("_userlevel") != nil {
		c.Data["_userlevel"] = c.GetSession("_userlevel").(string)
	}

	if c.GetSession("_company") != nil {
		c.Data["_company"] = c.GetSession("_company").(string)
	}

	if c.GetSession("_company_id") != nil {
		c.Data["_company_id"] = c.GetSession("_company_id").(string)
	}

	if c.GetSession("_company_pid") != nil {
		c.Data["_company_pid"] = c.GetSession("_company_pid").(string)
	}

	//获取表信息
	var tb = db.First("select * from tb_table where id=?", f["tbid"])
	if len(tb) < 1 {
		c.Ctx.WriteString("[]")
		return
	}
	var sql = f["form_value"]
	if sql == "" {
		c.Ctx.WriteString("[]")
		return
	}

	var tpl = template.New("")
	tpl.Parse(sql)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e == nil {
		sql = buf.String()
	}

	var jsonstr = ""
	if strings.Contains(f["form_value"], "select") == true {
		var xx = db.NewDb(tb["conn_str"])

		var list = db.Query2(xx, sql)
		result, err := json.MarshalIndent(list, "", "	")
		if err != nil {
			c.Ctx.WriteString("[]")
			return
		}
		jsonstr = string(result)
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
	var tbname = tb["table"]
	if tb["data_type"] != "0" && tb["table_master"] == "" {
		c.Ctx.WriteString("没有找到主表信息")
		return
	}
	if tb["data_type"] != "0" && tb["table_master"] != "" {
		tbname = tb["table_master"] //主表操作
	}

	var xx = db.NewDb(tb["conn_str"])
	//纪录唯一性检查
	var uniqlist = db.Query("select * from tb_field where tbid=? and is_unique=1", tb["id"])
	if len(uniqlist) > 0 {
		var uniwhere = ""
		var unifname = ""
		for _, v := range uniqlist {
			var fv = c.GetString(v["field_code"])
			if fv != "" {
				//主表字段名
				var field_code = v["field_code"]
				if tb["data_type"] != "0" && v["field_tb_code"] == "" {
					continue
				}
				if tb["data_type"] != "0" && v["field_tb_code"] != "" {
					field_code = v["field_tb_code"]
				}
				//----------------------------------------------------------
				if uniwhere != "" {
					uniwhere += " and "
				}

				uniwhere += field_code + "='" + c.GetString(v["field_code"]) + "' "

				if unifname != "" {
					unifname += ","
				}
				unifname += field_code
			}
		}
		if id > 0 {
			//var row = db.Query2(xx, `select * from `+tb["table"]+` where id !=? and `+v["field_code"]+`=?`, id, c.GetString(v["field_code"]))
			var row = db.Query2(xx, `select * from `+tb["table"]+" where "+uniwhere+` and id !=? `, id)
			if len(row) > 0 {
				c.Ctx.WriteString(unifname + "值不能重复!")
				return
			}
		} else {
			//var row = db.Query2(xx, `select * from `+tb["table"]+` where   `+v["field_code"]+`=?`, c.GetString(v["field_code"]))
			var row = db.Query2(xx, `select * from `+tb["table"]+" where "+uniwhere+``)
			if len(row) > 0 {
				c.Ctx.WriteString(unifname + "值不能重复!")
				return
			}
		}
	}
	var log = "" //log日志信息
	//
	var fields = db.Query("select * from tb_field where tbid=? and field_code!='id' and view_form=1 and form_type!='标签框'", tb["id"])
	var sql = ""
	var tmp = ""
	if id > 0 {
		sql = `update  ` + tbname + " set "
		for k, v := range fields {
			//主表字段名
			var field_code = v["field_code"]
			if tb["data_type"] != "0" && v["field_tb_code"] == "" {
				continue
			}
			if tb["data_type"] != "0" && v["field_tb_code"] != "" {
				field_code = v["field_tb_code"]
			}
			//----------------------------------------------------------
			if k > 0 && tmp != "" {
				sql += ","
				tmp += ","
			}
			var val = c.GetString(v["field_code"])
			val = strings.Replace(val, "'", "''", -1) //字符替换

			//检查是否必填项
			if v["form_required"] == "1" && val == "" {
				c.Ctx.WriteString(v["field_name"] + "不能为空")
				return
			}
			//字段默认值
			if val == "" {
				val = v["field_defval"]
			}
			//默认值替换---------------------------------------------
			//默认时间
			if val == "CURRENT_TIME" {
				val = time.Now().Format("2006-01-02 15:04:05")
			}
			if val == "CURRENT_DATE" {
				val = time.Now().Format("2006-01-02")
			}

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
			sql += " `" + field_code + "`='" + val + "' "

			log += field_code + "=" + val + ";\r\n"

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
		sql = `insert into ` + tbname + `(`
		var tmp = ""
		for k, v := range fields {
			//主表字段名
			var field_code = v["field_code"]
			if tb["data_type"] != "0" && v["field_tb_code"] == "" {
				continue
			}
			if tb["data_type"] != "0" && v["field_tb_code"] != "" {
				field_code = v["field_tb_code"]
			}
			//----------------------------------------------------------

			if k > 0 && tmp != "" {
				sql += ","
				tmp += ","
			}
			sql += "`" + field_code + "`"

			var val = c.GetString(v["field_code"])
			val = strings.Replace(val, "'", "''", -1) //字符替换

			//检查是否必填项
			if v["form_required"] == "1" && val == "" {
				c.Ctx.WriteString(v["field_name"] + "不能为空")
				return
			}
			//字段默认值
			if val == "" {
				val = v["field_defval"]
			}

			//默认值替换---------------------------------------------
			//默认时间
			if val == "CURRENT_TIME" {
				val = time.Now().Format("2006-01-02 15:04:05")
			}
			if val == "CURRENT_DATE" {
				val = time.Now().Format("2006-01-02")
			}

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

	var i int64 = 0 //db.Insert2(xx, sql, tb["table"]) //mysql msssql 可以返回主键
	if id > 0 {
		i = db.Exec2(xx, sql)
		c.Data["id"] = id
	} else {
		i = db.Insert2(xx, sql, tbname) //mysql msssql 可以返回主键
		c.Data["id"] = i
		id, _ = strconv.Atoi(fmt.Sprintf("%d", i))
	}
	fmt.Println("移动端返回主键:", c.Data["id"])
	if i > 0 {
		//保存成功,准备检查额外需执行代码
		var exlist = db.Query("select * from tb_table_ex where extype like 'sql%' and expage='info' and explace='AFTER' and tbid=? and state=1", tb["id"])
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
			//普通执行语句
			if v["extype"] == "sql" {
				c.Data[v["code"]] = db.Exec2(xx, p)
			}
			if v["extype"] == "sql_model" {
				c.Data[v["code"]] = db.First2(xx, p)
			}
			if v["extype"] == "sql_list" {
				c.Data[v["code"]] = db.Query2(xx, p)
			}
			if v["extype"] == "sql_insert" {
				c.Data[v["code"]] = db.Insert2(xx, p, v["title"]) //临时方案,后期增加单独字段 2019-05-03
			}
			if v["extype"] == "sql_execute" {
				fmt.Println("额外sql_execute:", p)
				c.Data[v["code"]] = db.Exec2(xx, p)
			}
		}

		//c.Ctx.WriteString("1")
		c.Ctx.WriteString(fmt.Sprintf("%d", id))
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

	//文件默认保存路径
	var savePath = "/upload/" + time.Now().Format("2006-01") + "/"
	if c.GetSession("_uid") != nil {
		savePath += c.GetSession("_uid").(string) + "/"
	} else {
		savePath += "0000/"
	}
	os.MkdirAll("static"+savePath, 0755) //目录不存在则创建
	//---------------------------------------------------------------------------------------
	var yin = c.GetString("watermark")

	bbb := bytes.NewBuffer(ddd)
	jpg, _, _ := image.Decode(bbb) // 图片文件解码
	//img := image.NewRGBA(image.Rect(0, 0, 480, 680))
	img := image.NewRGBA(jpg.Bounds())
	draw.Draw(img, jpg.Bounds().Add(image.Pt(0, 0)), jpg, jpg.Bounds().Min, draw.Src) //截取图片的一部分

	const (
		fontFile = "static/fonts/yahei.ttf"
		fontSize = 16 // 字体尺寸
		fontDPI  = 72 // 屏幕每英寸的分辨率
	)
	// 读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err == nil {
		font, err := freetype.ParseFont(fontBytes)
		if err == nil {
			cc := freetype.NewContext()
			cc.SetDPI(fontDPI)
			cc.SetFont(font)
			cc.SetFontSize(fontSize)
			cc.SetClip(img.Bounds())
			cc.SetDst(img)
			cc.SetSrc(image.Black)

			if yin != "0" {
				//输出姓名
				var realname = "-"
				if c.GetSession("_me") != nil {
					realname = c.GetSession("_me").(map[string]string)["realname"]
				}
				var x1 = 20
				var y1 = 30
				var txt1 = realname
				pt := freetype.Pt(x1, y1+int(cc.PointToFixed(fontSize)>>8)) // 字出现的位置

				_, err = cc.DrawString(txt1, pt)
				if err != nil {
					log.Println("向图片写字体出错1")
					log.Println(err)
				}
				//输出时间日期
				x1 = 20
				y1 = 60
				txt1 = time.Now().Format("2006-01-02 15:04:05")
				pt = freetype.Pt(x1, y1+int(cc.PointToFixed(fontSize)>>8)) // 字出现的位置

				_, err = cc.DrawString(txt1, pt)
				if err != nil {
					log.Println("向图片写字体出错1")
					log.Println(err)
				}
				//输出拍摄地点
				x1 = 20
				y1 = 90
				txt1 = c.GetString("_address")
				if txt1 != "" {
					pt = freetype.Pt(x1, y1+int(cc.PointToFixed(fontSize)>>8)) // 字出现的位置

					_, err = cc.DrawString(txt1, pt)
					if err != nil {
						log.Println("向图片写字体出错1")
						log.Println(err)
					}
				}

			}

			// 以PNG格式保存文件
			var fn = fmt.Sprintf(savePath+"%d.jpg", time.Now().UnixNano()/1e6)
			imgfile, err := os.Create("static" + fn)
			if err != nil {
				fmt.Println(err)
			}
			defer imgfile.Close()

			//转成[]byte
			// buf := new(bytes.Buffer)
			// err = jpeg.Encode(buf, img, &jpeg.Options{100})
			// ddd = buf.Bytes()

			// //以下生成图片方式,4,5百k;速度慢,客户端能感受出来
			err = png.Encode(imgfile, img)
			if err != nil {
				log.Println("生成图片出错")
				log.Fatal(err)
			} else {
				c.Ctx.WriteString(fn)
				return
			}

		}

	}

	//-------------------------------------------------------------------------------------
	//以下方式生成速度快,图片质量低
	var fn = fmt.Sprintf(savePath+"%d.jpg", time.Now().UnixNano()/1e6)
	err = ioutil.WriteFile(path+"/static"+fn, ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
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
	var lt = c.GetString("lt")
	var gt = c.GetString("gt")

	var lte = c.GetString("lte")
	var gte = c.GetString("gte")

	where = "where level>=" + strconv.Itoa(level)
	if gt != "" || lt != "" {
		where = "where 1=1 "
	}
	//读取 > >= < <= 参数
	if gt != "" {
		where += " and level >=" + gt
	}
	if lt != "" {
		where += " and level <=" + lt
	}
	if gte != "" {
		where += " and level >" + gte
	}
	if lte != "" {
		where += " and level <" + lte
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

//检查是否已登录
func (c *XApiController) IsLogin() {
	var rst = ``
	var _uid = c.GetSession("_uid")
	if _uid == nil {
		rst = `{
			"code":104,
			"msg":"未登录或会话超时",
			"extra":"0",
			"result":[]
		}`

		c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		c.Ctx.Output.Body([]byte(rst))
	}
	var m = db.First("select id,mch_id,usertype,username,mobile,openid,realname,headimg,company,company_id,company_pid,is_manager,state from adm_user where id=?", _uid)
	result, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		rst = `{
			"code":104,
			"msg":"读取失败,请稍后重试",
			"extra":"0",
			"result":[]
		}`
	} else {
		rst = `{
			"code":100,
			"msg":"已登录",
			"extra":"0",
			"result":[` + string(result) + `]
		}`
	}
	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))
}

//Login 登录系统  0参数不全或用户名密码错误  1登录成功  2账号异常
func (c *XApiController) Login() {
	var username = c.GetString("username")
	var userpwd = c.GetString("userpwd")
	if username == "" || userpwd == "" {
		c.Ctx.WriteString("0")
		return
	}
	fmt.Println("db", db.X)
	var u = db.First("select * from adm_user where username=? and password=?", username, userpwd)
	if u == nil {
		c.Ctx.WriteString("0")
		return
	} else {
		if u["state"] == "1" {
			c.SetSession("_me", u)
			c.SetSession("_uid", u["id"])
			c.SetSession("_mch_id", u["mch_id"])
			c.SetSession("_pid", u["pid"])
			c.SetSession("_pids", u["pids"])
			c.SetSession("_roles", u["roles"])
			c.SetSession("_username", u["username"])
			c.SetSession("_usertype", u["usertype"])
			c.SetSession("_sproot", u["sproot"])
			c.SetSession("_is_manager", u["is_manager"])
			c.SetSession("_company", u["company"])
			c.SetSession("_company_id", u["company_id"])
			c.SetSession("_company_pid", u["company_pid"])

			if u["username"] == "root" {
				c.SetSession("_sproot", "1")
			}
			//用户级别
			var ut = db.First("select * from adm_usertype where level=?", u["usertype"])
			if len(ut) > 0 {
				c.SetSession("_userlevel", ut["level"])
			} else {
				c.Ctx.WriteString("0")
				return
			}
			c.SetSession("_logintime", u["logintime"])
			c.SetSession("_loginip", u["loginip"])

			db.Exec("update adm_user set logintime=?,loginip=? where id=?", time.Now().Format("2006-01-02 15:04:05"), c.Ctx.Request.RemoteAddr, u["id"])

			c.Ctx.WriteString(u["id"])
			return
		} else {
			c.Ctx.WriteString("2")
			return
		}
	}
}

//Login 登录系统  0参数不全或用户名密码错误  1登录成功  2账号异常
func (c *XApiController) LogOut() {
	c.DelSession("_me")
	c.DelSession("_uid")
	c.DelSession("_mch_id")
	c.DelSession("_pid")
	c.DelSession("_pids")
	c.DelSession("_roles")
	c.DelSession("_username")
	c.DelSession("_usertype")
	c.DelSession("_sproot")
	c.DelSession("_is_manager")
	c.DelSession("_company")
	c.DelSession("_company_id")
	c.DelSession("_company_pid")

	c.Ctx.WriteString("1")
}
