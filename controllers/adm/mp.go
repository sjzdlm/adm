package adm

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

type MPController struct {
	beego.Controller
}

//首页
func (c *MPController) Index() {
	c.TplName = "adm/mp.html"
}

func (c *MPController) AppsJson() {
	var sql = `select id,title,code,icon,memo,state from tbm_app `
	var list = db.Query(sql)
	c.Data["json"] = list
	c.ServeJSON()
}
func (c *MPController) PagesJson() {
	var id = c.GetString("id")
	if id == "" {
		id = "0"
	}
	var m = db.First("select * from tbm_app where id=?", id)
	if len(m) < 1 {
		c.Data["json"] = "[]"
		c.ServeJSON()
		return
	}
	var sql = `select id,title,module,code,icon,memo,state from tbm_page where module=? `
	var list = db.Query(sql, m["code"])
	c.Data["json"] = list
	c.ServeJSON()
}
func (c *MPController) WidgetsJson() {
	var id = c.GetString("id")
	if id == "" {
		id = "0"
	}

	var sql = `select id,tbmid,field_name,field_code,form_type,form_sort,form_value,form_required,data,data_sql,memo,state from tbm_widget where tbmid=? order by form_sort `
	var list = db.Query(sql, id)
	c.Data["json"] = list
	c.ServeJSON()
}

//某个组件实例
func (c *MPController) WidgetJson() {
	var id = c.GetString("id")
	if id == "" {
		id = "0"
	}

	var sql = `select * from tbm_widget where id=? `
	var list = db.First(sql, id)
	c.Data["json"] = list
	c.ServeJSON()
}

//某个页面实例
func (c *MPController) PageJson() {
	var id = c.GetString("id")
	if id == "" {
		id = "0"
	}

	var sql = `select * from tbm_page where id=? `
	var list = db.First(sql, id)
	c.Data["json"] = list
	c.ServeJSON()
}

//某个应用实例
func (c *MPController) AppJson() {
	var id = c.GetString("id")
	if id == "" {
		id = "0"
	}

	var sql = `select * from tbm_app where id=? `
	var list = db.First(sql, id)
	c.Data["json"] = list
	c.ServeJSON()
}

//图标列表
func (c *MPController) IconsJson() {
	var sql = `select * from tbm_icon `
	var list = db.Query(sql)
	c.Data["json"] = list
	c.ServeJSON()
}

//组件列表
func (c *MPController) WidgetTypesJson() {
	var sql = `select id,title,tplcode,tpltype,state from tbm_widget_type `
	var list = db.Query(sql)
	c.Data["json"] = list
	c.ServeJSON()
}

//导航列表
func (c *MPController) NavListJson() {
	var pageid = c.GetString("page_id")
	if pageid == "" {
		pageid = "0"
	}
	//fmt.Println("page_id:", pageid)

	var sql = `select * from tbm_nav_list where page_id=? order by sort`
	var list = db.Query(sql, pageid)
	c.Data["json"] = list
	c.ServeJSON()
}

//某个组件实例POST
func (c *MPController) WidgetPost() {
	var id, _ = c.GetInt("id", 0)
	var tbmid, _ = c.GetInt("tbmid", 0)
	if tbmid <= 0 {
		c.Ctx.WriteString("-1")
		return
	}

	var field_name = c.GetString("field_name")
	var form_sort, _ = c.GetInt("form_sort", 0)
	var form_type = c.GetString("form_type")
	var view_list = c.GetString("view_list")
	var view_form = c.GetString("view_form")

	var tpl = ""
	var data = c.GetString("data")
	var data_sql = c.GetString("data_sql")

	var w = db.First("select * from tbm_widget where tplcode=?", form_type)
	if len(w) > 0 {
		tpl = w["tpltxt"]
		if data == "" {
			data = w["tpldata"]
		}
	}

	var i int64 = 0
	var sql = ``
	if id <= 0 {
		var data = c.GetString("data")
		sql = `insert into tbm_widget(tbmid,field_name,form_sort,form_type,view_list,view_form,tpl,data,data_sql,state)values(
			?,?,?,?,?,
			?,?,?,?,?)`
		i = db.Exec(sql, tbmid, field_name, form_sort, form_type, view_list, view_form, tpl, data, data_sql, 1)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	} else {
		sql = `update tbm_widget set 
		field_name=?,form_sort=?,form_type=?,view_list=?,view_form=?,data=?,data_sql=?
		where id=?
		`
		i = db.Exec(sql, field_name, form_sort, form_type, view_list, view_form, data, data_sql, id)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	}

	c.Ctx.WriteString("0")
}

//某个页面实例POST
func (c *MPController) PagePost() {
	var id, _ = c.GetInt("id", 0)
	var module = c.GetString("module")

	var title = c.GetString("title")
	var code = c.GetString("code")
	var form_url = c.GetString("form_url")
	var ex_javascript = c.GetString("ex_javascript")
	var tip_msg = c.GetString("tip_msg")
	var memo = c.GetString("memo")
	var state = c.GetString("state")
	if state == "on" || state == "1" {
		state = "1"
	} else {
		state = "0"
	}

	var i int64 = 0
	var sql = ``
	if id <= 0 {
		sql = `insert into tbm_page(title,module,code,form_url,ex_javascript,tip_msg,memo,state)values(
			?,?,?,?,?,
			?,?,?)`
		i = db.Exec(sql, title, module, code, form_url, ex_javascript, tip_msg, memo, state)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	} else {
		sql = `update tbm_page set 
		title=?,module=?,code=?,form_url=?,ex_javascript=?,tip_msg=?,memo=?,state=?
		where id=?
		`
		i = db.Exec(sql, title, module, code, form_url, ex_javascript, tip_msg, memo, state, id)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	}

	c.Ctx.WriteString("0")
}

//某个应用实例POST
func (c *MPController) AppPost() {
	var id, _ = c.GetInt("id", 0)

	var title = c.GetString("title")
	var code = c.GetString("code")
	var icon = c.GetString("icon")
	var memo = c.GetString("memo")
	var state = c.GetString("state")
	if state == "on" || state == "1" {
		state = "1"
	} else {
		state = "0"
	}

	var i int64 = 0
	var sql = ``
	if id <= 0 {
		sql = `insert into tbm_app(title,code,icon,memo,state)values(
			?,?,?,?,?)`
		i = db.Exec(sql, title, code, icon, memo, state)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	} else {
		sql = `update tbm_app set 
		title=?,code=?,icon=?,memo=?,state=?
		where id=?
		`
		i = db.Exec(sql, title, code, icon, memo, state, id)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	}

	c.Ctx.WriteString("0")
}

//Post 登录系统  0参数不全或用户名密码错误  1登录成功  2账号异常
func (c *MPController) LoginPost() {
	var username = c.GetString("username")
	var userpwd = c.GetString("userpwd")
	if username == "" || userpwd == "" {
		c.Ctx.WriteString("0")
		return
	}
	var u = db.First("select * from adm_user where username=? and password=?", username, userpwd)
	if u == nil {
		c.Ctx.WriteString("0")
		return
	} else {
		if u["state"] == "1" {
			c.SetSession("_uid", u["id"])
			c.SetSession("_mch_id", u["mch_id"])
			c.SetSession("_roles", u["roles"])
			c.SetSession("_username", u["username"])
			c.SetSession("_usertype", u["usertype"])

			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("2")
			return
		}
	}

}

func (c *MPController) TbListJson() {
	// var _uid = c.GetSession("_uid")
	// if _uid == nil {
	// 	c.Ctx.WriteString("{}")
	// 	return
	// }
	var code = c.GetString("_code")
	var tb = db.First("select * from tb_table where code=?", code)
	if tb == nil {
		c.Ctx.WriteString("{}")
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
	//var data_table_prefix = ""      //如果是语句,别名的前缀
	if data_type == "2" {
		data_table = "_tb" //
		//data_table_prefix = "_tb."
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

	//合计字段
	var hjzd = db.First("select * from tb_field where tbid=? and is_total=1 and view_list=1", tb["id"])
	if len(hjzd) > 0 {
		totalField = hjzd["field_code"]
	}
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
				//where += v["field_code"] + ` like '%` + q + `%' `
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
		//where+=" order by "+sort+" "+order
		orderstr = " order by " + sort + " " + order
	} else {
		if tb["pri_key"] != "" {
			//where+=" order by "+tb["pri_key"]+" desc "
			orderstr = " order by " + tb["pri_key"] + " desc "
		}
	}

	var xx = db.NewDb(tb["conn_str"])
	var sql = "select * from " + tb["table"] + " " + where
	if data_type == "2" { //如果是语句的话
		sql = "select * from (" + tb["table"] + ") as" + data_table + " " + where
	}

	//判断数据库类型
	var con = db.First("select * from adm_conn where conn='" + tb["conn_str"] + "' limit 1")
	if strings.Contains(con["dbtype"], "mssql") {
		// mssql=`
		// SELECT TOP `+page+` * FROM `+tb["table"]+`
		// WHERE id > (
		// 　　SELECT MAX(id) FROM (
		// 　　　　SELECT TOP `+page+`*(`+pageSize+`-1) * FROM `+tb["table"]+`  `+where+`
		// 　　)
		// )
		// `+orderstr

		//var rst=db.Pager2(xx,page,pageSize,mssql)
		var rst = db.Pager2MsSql(xx, page, pageSize, tb["table"], sql, where, orderstr, sql+orderstr)
		// if totalField !=""{
		// 	var heji=db.Query2(xx,"select sum("+totalField+") as "+totalField+" from "+tb["table"]+" "+where)
		// 	fmt.Println("heji:",heji)
		// 	if len(heji)>0{
		// 		rst.Footer=heji
		// 	}
		// }
		c.Data["json"] = rst
	} else {
		fmt.Println(sql + orderstr)
		var rst = db.Pager2(xx, page, pageSize, sql+orderstr)
		if totalField != "" {
			var heji = db.Query2(xx, "select sum("+totalField+") as "+totalField+" from "+tb["table"]+" "+where+" "+orderstr)
			fmt.Println("heji:", heji)
			if len(heji) > 0 {
				rst.Footer = heji
			}
		}
		c.Data["json"] = rst

	}

	c.ServeJSON()
}
