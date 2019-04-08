package adm

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//TbController 控制器
type TbController struct {
	beego.Controller
}

//测试用
func (c *TbController) Get() {
	// var list,pages,rows=db.Page(1,20,"select * from adm_menu where id>?",0)
	// fmt.Println("----------------------------------")
	// fmt.Println(pages,rows,list)
	// c.Data["json"]=&list

	var rst = db.Pager(1, 20, "select * from adm_menu where id>?  ", 0)
	fmt.Println("page json:", rst)
	c.Data["json"] = &rst

	c.ServeJSON()
}

//List 列表页面
func (c *TbController) List() {
	var list = db.Query("select * from tb_table")
	c.Data["list"] = list
	//下拉框列替换赋值
	var xlist = db.Query("select * from tb_table_proj where   state=1  ")
	c.Data["xlist"] = xlist
	var rst = ""
	rst += `$('div .proj_id').each(function(){ `
	for kk, vv := range xlist {
		if kk > 0 {
			rst += " "
		}
		rst += `if($(this).text()=='` + vv["id"] + `'){`
		rst += `	$(this).text('` + vv["proj_name"] + `');`
		rst += `}`
	}
	rst += `});`
	c.Data["jsval"] = rst

	//c.TplName="adm/tb/list.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_list)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//JSON格式数据
func (c *TbController) JsonList() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var proj_id = c.GetString("proj_id")
	var qtxt = c.GetString("qtxt")
	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		qtxt = " where title like '%" + qtxt + "%'"
		if proj_id != "" {
			qtxt += " and proj_id=" + proj_id
		}
	} else {
		if proj_id != "" {
			qtxt += " where proj_id=" + proj_id
		}
	}
	qtxt += " order by id desc "
	var rst = db.Pager(page, pageSize, "select * from tb_table "+qtxt)

	c.Data["json"] = rst

	c.ServeJSON()
}

//Edit 编辑、新增页面
func (c *TbController) Edit() {
	//获取用户商户ID
	var _mch_id = c.GetSession("_mch_id")
	if _mch_id == nil {
		_mch_id = "0"
	}
	var _uid = c.Ctx.Input.Session("_uid")
	if _uid == nil {
		_uid = "0"
	}
	//var _uname=c.Ctx.Input.Session("_username")

	//随机字符串
	var code = "z" + db.RandomString(9)
	var id, _ = c.GetInt("id", 0)
	//如果ID大于0，则修改此数据
	c.Data["is_export"] = "false"
	if id > 0 {
		var m = db.First("select * from tb_table where mch_id=? and id=?", _mch_id, id)
		if m == nil {
			c.Ctx.WriteString("参数错误!")
			return
		}
		//fmt.Println("table:\r\n",m)
		c.Data["m"] = m
		if m["is_import"] == "1" {
			c.Data["is_import"] = "true"
		} else {
			c.Data["is_import"] = "false"
		}
		if m["is_export"] == "1" {
			c.Data["is_export"] = "true"
		} else {
			c.Data["is_export"] = "false"
		}

		code = m["code"]
		c.Data["edit_width"] = m["edit_width"]
		c.Data["edit_height"] = m["edit_height"]
	} else {
		c.Data["edit_width"] = 420
		c.Data["edit_height"] = 380
	}
	c.Data["code"] = code
	//数据库链接
	var dblist = db.Query("select * from adm_conn where state=1")
	if dblist != nil {
		c.Data["dblist"] = dblist
	}
	//项目列表数据
	var projlist = db.Query("select * from tb_table_proj where state=1")
	if projlist != nil {
		c.Data["projlist"] = projlist
	}
	//c.TplName="adm/tb/edit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_edit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//Edit 编辑、新增保存页面
func (c *TbController) EditPost() {
	//获取用户商户ID
	var _mch_id = c.GetSession("_mch_id")
	if _mch_id == nil {
		_mch_id = "0"
	}
	var _uid = c.Ctx.Input.Session("_uid")
	if _uid == nil {
		_uid = "0"
	}
	//ID
	var id, _ = c.GetInt("id", 0)
	//名称
	var title = c.GetString("title")
	title = strings.TrimSpace(string(title))
	if title == "" {
		c.Ctx.WriteString("请输入名称!")
		return
	}
	//项目链接
	var proj_id = c.GetString("proj_id")
	proj_id = strings.TrimSpace(string(proj_id))
	if proj_id == "" {
		c.Ctx.WriteString("请选择项目!")
		return
	}
	//代号
	var code = c.GetString("code")
	code = strings.TrimSpace(string(code))
	if code == "" {
		c.Ctx.WriteString("请输入代号!")
		return
	}
	//数据库链接
	var conn = c.GetString("conn")
	conn = strings.TrimSpace(string(conn))
	if conn == "" {
		c.Ctx.WriteString("请选择数据库!")
		return
	}
	//表名
	var table = c.GetString("table")
	table = strings.TrimSpace(string(table))
	if code == "" {
		c.Ctx.WriteString("请输入表名!")
		return
	}
	//条件
	var where_str = c.GetString("where_str")
	where_str = strings.TrimSpace(string(where_str))

	//数据类型 表、视图、SQL
	var data_type = c.GetString("data_type")
	data_type = strings.TrimSpace(string(data_type))
	if data_type == "" {
		c.Ctx.WriteString("请选择数据类型!")
		return
	}
	//主键字段
	var pri_key = c.GetString("pri_key")
	pri_key = strings.TrimSpace(string(pri_key))
	//排序字段
	var sort_key = c.GetString("sort_key")
	sort_key = strings.TrimSpace(string(sort_key))
	if sort_key == "" {
		c.Ctx.WriteString("请输入排序字段!")
		return
	}
	//新建功能
	var is_create = c.GetString("is_create")
	is_create = strings.TrimSpace(string(is_create))
	if is_create == "on" {
		is_create = "1"
	} else {
		is_create = "0"
	}
	//编辑功能
	var is_edit = c.GetString("is_edit")
	is_edit = strings.TrimSpace(string(is_edit))
	if is_edit == "on" {
		is_edit = "1"
	} else {
		is_edit = "0"
	}
	//详情功能
	var is_detail = c.GetString("is_detail")
	is_detail = strings.TrimSpace(string(is_detail))
	if is_detail == "on" {
		is_detail = "1"
	} else {
		is_detail = "0"
	}
	//删除功能
	var is_del = c.GetString("is_del")
	is_del = strings.TrimSpace(string(is_del))
	if is_del == "on" {
		is_del = "1"
	} else {
		is_del = "0"
	}
	//展示方式 弹窗 新页
	var edit_style = c.GetString("edit_style")
	edit_style = strings.TrimSpace(string(edit_style))
	if edit_style == "" {
		c.Ctx.WriteString("请选择展示方式!")
		return
	}
	//宽度
	var edit_width = c.GetString("edit_width")
	edit_width = strings.TrimSpace(string(edit_width))
	if edit_width == "" {
		c.Ctx.WriteString("请输入宽度!")
		return
	}
	//高度
	var edit_height = c.GetString("edit_height")
	edit_height = strings.TrimSpace(string(edit_height))
	if edit_height == "" {
		c.Ctx.WriteString("请输入高度!")
		return
	}
	//导出功能
	var is_export = c.GetString("is_export")
	is_export = strings.TrimSpace(string(is_export))
	if is_export == "on" {
		is_export = "1"
	} else {
		is_export = "0"
	}
	//导入功能
	var is_import = c.GetString("is_import")
	is_import = strings.TrimSpace(string(is_import))
	if is_import == "on" {
		is_import = "1"
	} else {
		is_import = "0"
	}
	//合计功能
	var is_total = c.GetString("is_total")
	is_total = strings.TrimSpace(string(is_total))
	if is_total == "on" {
		is_total = "1"
	} else {
		is_total = "0"
	}
	//操作扩展
	var ex_html_operate = c.GetString("ex_html_operate")
	ex_html_operate = strings.TrimSpace(string(ex_html_operate))
	if ex_html_operate == "" {

	}
	//JS扩展
	var ex_javascript = c.GetString("ex_javascript")
	//fmt.Println(ex_javascript)
	//新增SQL扩展
	var ex_sql_add = c.GetString("ex_sql_add")
	//编辑SQL扩展
	var ex_sql_edit = c.GetString("ex_sql_edit")
	//删除SQL扩展
	var ex_sql_del = c.GetString("ex_sql_del")

	//工具栏扩展
	var ex_linkbutton = c.GetString("ex_linkbutton")
	fmt.Println(ex_linkbutton)
	//备注
	var memo = c.GetString("memo")
	fmt.Println(memo)
	//状态
	var state = c.GetString("state")
	fmt.Println(state)
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ""
	//如果ID大于0，则修改此数据
	if id > 0 {
		var m = db.First("select * from tb_table where mch_id=? and id=?", _mch_id, id)
		if m == nil {
			c.Ctx.WriteString("参数错误!")
			return
		}
		sql = `update tb_table
		set title=?,
			code=?,
			proj_id=?,
			conn_str=?,`
		sql += "`table`=?,"
		sql += `where_str=?,
			data_type=?,
			pri_key=?,
			sort_key=?,
			is_create=?,
			is_edit=?,
			is_detail=?,
			is_del=?,
			edit_style=?,
			edit_width=?,
			edit_height=?,
			is_export=?,
			is_import=?,
			is_total=?,
			ex_html_operate=?,
			ex_javascript=?,
			ex_linkbutton=?,
			memo=?,
			state=?,
			mch_id=?,
			user_id=?,
			ex_sql_add=?,
			ex_sql_edit=?,
			ex_sql_del=?
		where id=?			
		`
		var i = db.Exec(sql,
			title,
			code,
			proj_id,
			conn,
			table,
			where_str,
			data_type,
			pri_key,
			sort_key,
			is_create,
			is_edit,
			is_detail,
			is_del,
			edit_style,
			edit_width,
			edit_height,
			is_export,
			is_import,
			is_total,
			ex_html_operate,
			ex_javascript,
			ex_linkbutton,
			memo,
			state,
			_mch_id,
			_uid,
			ex_sql_add,
			ex_sql_edit,
			ex_sql_del,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	} else {
		sql = `insert into tb_table(
		title,
		code,
		proj_id,
		conn_str,`
		sql += "`table`,"
		sql += `where_str,
		data_type,
		pri_key,
		sort_key,
		is_create,
		is_edit,
		is_detail,
		is_del,
		edit_style,
		edit_width,
		edit_height,
		is_export,
		is_import,
		ex_html_operate,
		ex_javascript,
		ex_linkbutton,
		memo,
		state,mch_id,user_id,
		ex_sql_add,
		ex_sql_edit,
		ex_sql_del
		)
		values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`
		var i = db.Exec(sql,
			title,
			code,
			proj_id,
			conn,
			table,
			where_str,
			data_type,
			pri_key,
			sort_key,
			is_create,
			is_edit,
			is_detail,
			is_del,
			edit_style,
			edit_width,
			edit_height,
			is_export,
			is_import,
			ex_html_operate,
			ex_javascript,
			ex_linkbutton,
			memo,
			state,
			_mch_id,
			_uid,
			ex_sql_add,
			ex_sql_edit,
			ex_sql_del,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		}
	}
	c.Ctx.WriteString("0")
}

//字段选择页面
func (c *TbController) Fields() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	var tb = db.First("select * from tb_table where id=?", id)
	if tb != nil {
		var sql = "select * from adm_conn where conn='" + tb["conn_str"] + "' "
		var cnntb = db.FirstOrNil(sql)
		if cnntb == nil {
			c.Ctx.WriteString("链接错误!")
			return
		}
		sql = "select  TABLE_NAME,COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,CHARACTER_MAXIMUM_LENGTH AS MAXLENGTH,COLUMN_COMMENT"
		sql += " from Information_schema.columns  where table_schema='" + cnntb["dbname"] + "' and table_Name = '" + tb["table"] + "'"

		//-------------------------------------------
		var dbtype = cnntb["dbtype"]
		c.Data["dbtype"] = dbtype
		if dbtype == "sqlite" {
			sql = " PRAGMA table_info('" + tb["table"] + "') "
		} else if dbtype == "mysql" {
			//读取字段信息
			sql = "select  ORDINAL_POSITION as cid ,COLUMN_DEFAULT as dflt_value,TABLE_NAME as tbname,COLUMN_NAME as name,IS_NULLABLE as notnull,DATA_TYPE as type,COLUMN_KEY as pk,CHARACTER_MAXIMUM_LENGTH AS MAXLENGTH,COLUMN_COMMENT as memo"
			sql += " from Information_schema.columns  where table_Name = '" + tb["table"] + "'"
		} else if dbtype == "mssql" {
			sql = `
			SELECT  row_number() over(order by syscolumns.name) as cid,
                 Sysobjects.name AS [tbname], 
                 syscolumns.name AS [name],  
                 cast(sys.extended_properties.[value] as varchar(200)) AS [memo],  
                 systypes.name AS [type],  
                 syscolumns.length AS [length],
                 CASE syscolumns.isnullable WHEN '1' THEN 'Y' ELSE 'N' END AS [nullable],  
                 syscomments.text AS [dflt_value],
                 COLUMNPROPERTY(syscolumns.id, syscolumns.name, 'IsIdentity') AS [pk] ,  
                 CASE WHEN EXISTS (SELECT 1 FROM sysobjects WHERE xtype = 'PK' AND name IN  
                 (SELECT name  
                 FROM sysindexes  
                 WHERE indid IN  
                 (SELECT indid  
                 FROM sysindexkeys  
                 WHERE id = syscolumns.id AND colid = syscolumns.colid)))  
                 THEN 1 ELSE 0 END AS [iskey]  
                 FROM syscolumns  
                 INNER JOIN systypes  
                 ON syscolumns.xtype = systypes.xtype  
                 LEFT JOIN sysobjects ON syscolumns.id = sysobjects.id  
                 LEFT OUTER JOIN sys.extended_properties ON  
                 ( sys.extended_properties.minor_id = syscolumns.colid  
                 AND sys.extended_properties.major_id = syscolumns.id)  
                 LEFT OUTER JOIN syscomments ON syscolumns.cdefault = syscomments.id  
                 WHERE (systypes.name <> 'sysname')  
                 AND syscolumns.id IN (SELECT id FROM SYSOBJECTS WHERE  NAME = '` + tb["table"] + `') 
			`
		} else if dbtype == "mssql2k" {
			sql = `
			SELECT           1 as cid,
					 sysobjects.name AS [tbname],   
					 syscolumns.name AS [name],
					 cast(properties.[value] as varchar(200)) AS [memo],  
					 systypes.name AS [type],  
					 syscolumns.length AS [length], 
					 CASE syscolumns.isnullable WHEN '1' THEN 'Y' ELSE 'N' END AS [nullable],  
					 CASE WHEN syscomments.text IS NULL THEN '' ELSE syscomments.text END AS [dflt_value],  
					 CASE WHEN COLUMNPROPERTY(syscolumns.id, syscolumns.name, 'IsIdentity') = 1 THEN 1 ELSE 0 END AS [pk],  
					 CASE WHEN EXISTS (SELECT 1 FROM sysobjects WHERE xtype = 'PK' AND name IN  
					 (SELECT name  
					 FROM sysindexes  
					 WHERE indid IN  
					 (SELECT indid  
					 FROM sysindexkeys  
					 WHERE id = syscolumns.id AND colid = syscolumns.colid)))  
					 THEN 1 ELSE 0 END AS [iskey]  
					 FROM syscolumns INNER JOIN  
					 sysobjects ON sysobjects.id = syscolumns.id INNER JOIN  
					 systypes ON syscolumns.xtype = systypes.xtype LEFT OUTER JOIN  
					 sysproperties properties ON syscolumns.id = properties.id AND  
					 syscolumns.colid = properties.smallid LEFT OUTER JOIN  
					 sysproperties ON sysobjects.id = sysproperties.id AND  
					 sysproperties.smallid = 0 LEFT OUTER JOIN  
					 syscomments ON syscolumns.cdefault = syscomments.id  
					 WHERE (sysobjects.xtype = 'u')  
					 AND sysobjects.NAME='` + tb["table"] + `' 
					 ORDER BY [tbname]  
			`
		}
		//--------------------------------------------
		var xx = db.NewDb(tb["conn_str"])
		var list = db.Query2(xx, sql)
		//主键
		for _, v := range list {
			if v["COLUMN_KEY"] == "PRI" || v["COLUMN_KEY"] == "1" {
				db.Exec2(xx, "update tb_table set pri_key=? where id=?", v["name"], id)
			}
		}

		c.Data["list"] = list
		c.Data["tb"] = tb["table"]

		var flist = db.Query("select * from tb_field where tbid=?", tb["id"])
		c.Data["flist"] = flist
	}
	//c.TplName="adm/tb/fields.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_fields)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//字段选择保存页面  -1数据库链接失败 0保存失败 1保存成功 否则属于提示信息
func (c *TbController) FieldsPost() {
	var id, _ = c.GetInt("id", 0)
	var tb = c.GetString("tb")
	var fs = c.GetStrings("fields")
	//fmt.Println("fs",fs)
	if len(fs) < 1 {
		c.Ctx.WriteString("请选中至少一个字段!")
		return
	}

	//将字段记录到字典中,以便循环判断
	var dict = make(map[string]string)
	for _, v := range fs {
		dict[v] = v
	}
	fmt.Println("dict:", dict)
	//读取表记录
	var table = db.First(" select * from tb_table where id=?", id)
	if table == nil {
		c.Ctx.WriteString("0")
		return
	}

	var sql = "select * from adm_conn where conn='" + table["conn_str"] + "' "
	var cnntb = db.FirstOrNil(sql)
	if cnntb == nil {
		c.Ctx.WriteString("链接错误!")
		return
	}

	fmt.Println("table:", table)
	//读取字段信息
	sql = ""
	//-------------------------------------------
	var dbtype = cnntb["dbtype"]
	c.Data["dbtype"] = dbtype
	fmt.Println("dbtype:", dbtype)
	if dbtype == "sqlite" {
		sql = " PRAGMA table_info('" + tb + "') "
	} else if dbtype == "mysql" {
		//读取字段信息
		sql = "select  ORDINAL_POSITION as cid ,COLUMN_DEFAULT as dflt_value,TABLE_NAME as tbname,COLUMN_NAME as name,IS_NULLABLE as notnull,DATA_TYPE as type,COLUMN_KEY as pk,CHARACTER_MAXIMUM_LENGTH AS MAXLENGTH,COLUMN_COMMENT as memo"
		sql += " from Information_schema.columns  where table_Name = '" + tb + "'"
	} else if dbtype == "mssql" {
		sql = `
		SELECT  row_number() over(order by syscolumns.name) as cid,
                 Sysobjects.name AS [tbname], 
                 syscolumns.name AS [name],  
                 cast(sys.extended_properties.[value] as varchar(200)) AS [memo],  
                 systypes.name AS [type],  
                 syscolumns.length AS [length],
                 CASE syscolumns.isnullable WHEN '1' THEN 'Y' ELSE 'N' END AS [nullable],  
                 syscomments.text AS [dflt_value],
                 COLUMNPROPERTY(syscolumns.id, syscolumns.name, 'IsIdentity') AS [pk] ,  
                 CASE WHEN EXISTS (SELECT 1 FROM sysobjects WHERE xtype = 'PK' AND name IN  
                 (SELECT name  
                 FROM sysindexes  
                 WHERE indid IN  
                 (SELECT indid  
                 FROM sysindexkeys  
                 WHERE id = syscolumns.id AND colid = syscolumns.colid)))  
                 THEN 1 ELSE 0 END AS [iskey]  
                 FROM syscolumns  
                 INNER JOIN systypes  
                 ON syscolumns.xtype = systypes.xtype  
                 LEFT JOIN sysobjects ON syscolumns.id = sysobjects.id  
                 LEFT OUTER JOIN sys.extended_properties ON  
                 ( sys.extended_properties.minor_id = syscolumns.colid  
                 AND sys.extended_properties.major_id = syscolumns.id)  
                 LEFT OUTER JOIN syscomments ON syscolumns.cdefault = syscomments.id  
                 WHERE (systypes.name <> 'sysname')  
                 AND syscolumns.id IN (SELECT id FROM SYSOBJECTS WHERE  NAME = '` + tb + `') 
		`
	} else if dbtype == "mssql2k" {
		sql = `
		SELECT           1 as cid,
                 sysobjects.name AS [tbname],   
                 syscolumns.name AS [name],  
                 cast(properties.[value] as varchar(200)) AS [memo],   
                 systypes.name AS [type],  
                 syscolumns.length AS [length], 
                 CASE syscolumns.isnullable WHEN '1' THEN 'Y' ELSE 'N' END AS [nullable],  
                 CASE WHEN syscomments.text IS NULL THEN '' ELSE syscomments.text END AS [dflt_value],  
                 CASE WHEN COLUMNPROPERTY(syscolumns.id, syscolumns.name, 'IsIdentity') = 1 THEN 1 ELSE 0 END AS [pk],  
                 CASE WHEN EXISTS (SELECT 1 FROM sysobjects WHERE xtype = 'PK' AND name IN  
                 (SELECT name  
                 FROM sysindexes  
                 WHERE indid IN  
                 (SELECT indid  
                 FROM sysindexkeys  
                 WHERE id = syscolumns.id AND colid = syscolumns.colid)))  
                 THEN 1 ELSE 0 END AS [iskey]  
                 FROM syscolumns INNER JOIN  
                 sysobjects ON sysobjects.id = syscolumns.id INNER JOIN  
                 systypes ON syscolumns.xtype = systypes.xtype LEFT OUTER JOIN  
                 sysproperties properties ON syscolumns.id = properties.id AND  
                 syscolumns.colid = properties.smallid LEFT OUTER JOIN  
                 sysproperties ON sysobjects.id = sysproperties.id AND  
                 sysproperties.smallid = 0 LEFT OUTER JOIN  
                 syscomments ON syscolumns.cdefault = syscomments.id  
				 WHERE (sysobjects.xtype = 'u')  
                 AND sysobjects.NAME='` + tb + `' 
                 ORDER BY [tbname]   
		`
	}
	//--------------------------------------------
	fmt.Println("sql:", sql)
	if sql == "" {
		c.Ctx.WriteString("系统错误")
		return
	}

	var xx = db.NewDb(table["conn_str"])
	if xx == nil {
		c.Ctx.WriteString("-1")
		return
	}
	var list = db.Query2(xx, sql)
	if list == nil || len(list) < 1 {
		c.Ctx.WriteString("0")
		return
	}
	//fmt.Println("list:",list)
	//遍历字段记录
	var rst = 0
	for k, v := range list {
		if _, ok := dict[v["name"]]; ok {
			//fmt.Println("选中字段:",v["column_name"])

			var MAXLENGTH = "25"
			if dbtype != "sqlite" {
				MAXLENGTH = strings.TrimSpace(string(v["maxlength"]))
			}
			if MAXLENGTH == "" {
				MAXLENGTH = "20"
			}
			var r = db.First("select * from tb_field where tbid=? and tbcode=? and field_code=?", id, tb, v["name"])

			//插入新记录
			if r == nil {
				sql = `
				insert into tb_field(
					tbid,
					tbcode,
					field_name,
					field_type,
					field_code,
					field_defval,
					field_length,
					form_length,
					form_tip,
					form_type,
					form_value,
					memo,
					state,
					form_sort,
					field_prikey,
					view_list
				)values(
					?,?,?,?,?,?,?,?,
					?,?,?,?,?,?,?,?
				)
				`
				var iskey = 0
				if v["pk"] == "pri" || v["pk"] == "1" {
					iskey = 1
				}
				var i = db.Exec(sql,
					id,
					tb,
					v["name"],
					v["type"],
					v["name"],
					v["dflt_value"],
					MAXLENGTH,
					MAXLENGTH,
					v["memo"],
					"文本框",
					"",
					v["memo"],
					1,
					k,
					iskey,
					0,
				)
				if i > 0 {
					rst++
				}
			} else { //更新记录
				sql = `update tb_field set
				field_type=?,
				field_code=?,
				field_length=?,
				field_prikey=?
				where id=?
				`
				var iskey = 0
				if v["pk"] == "pri" || v["pk"] == "1" {
					iskey = 1
				}
				var i = db.Exec(sql,
					v["type"],
					v["name"],
					MAXLENGTH,
					iskey,
					id,
				)
				if i > 0 {
					rst++
				}
			}
		} else { //没有选择的将被删除
			db.Exec("delete from tb_field where tbid=? and tbcode=? and field_code=?", table["id"], tb, v["name"])
		}
	}
	if rst > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

//表单字段页面
func (c *TbController) FieldList() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	if id < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	//读取tb_table表
	var tb = db.First("select * from tb_table where id=?", id)
	c.Data["tb"] = tb

	//c.TplName="adm/tb/fieldlist.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_fieldlist)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//表单字段数据
func (c *TbController) FieldListJson() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	if id < 1 {
		c.Ctx.WriteString("{}")
		return
	}
	//读取tb_table表
	var tb = db.First("select * from tb_table where id=?", id)
	c.Data["tb"] = tb
	//读取表单元素列表
	var list = db.Query("select * from tb_field where tbid=? order by form_sort", id)
	if list == nil || len(list) < 1 {
		c.Ctx.WriteString("{}")
		return
	}
	c.Data["list"] = list

	c.Data["json"] = list
	c.ServeJSON()
}

//表单字段数据
func (c *TbController) FieldEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id
	var m = db.First("select * from tb_field where id=?", id)
	if m == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}
	c.Data["m"] = m
	//c.TplName="adm/tb/fieldedit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_fieldedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//表单字段数据保存
func (c *TbController) FieldEditPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id
	var m = db.First("select * from tb_field where id=?", id)
	if m == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}

	//获取参数
	var field_name = c.GetString("field_name")
	var field_code = c.GetString("field_code")
	var form_type = c.GetString("form_type")
	var form_length, _ = c.GetInt("form_length", 0)
	var form_value = c.GetString("form_value")
	var form_tip = c.GetString("form_tip")
	var field_defval = c.GetString("field_defval")
	var form_sort, _ = c.GetInt("form_sort", 0)
	//列表显示
	var view_list = c.GetString("view_list")
	if view_list == "on" {
		view_list = "1"
	} else {
		view_list = "0"
	}
	//列表显示色值
	var view_list_color = c.GetString("view_list_color")
	//列表编辑
	var view_edit = c.GetString("view_edit")
	if view_edit == "on" {
		view_edit = "1"
	} else {
		view_edit = "0"
	}
	//表单显示
	var view_form = c.GetString("view_form")
	if view_form == "on" {
		view_form = "1"
	} else {
		view_form = "0"
	}
	//详情显示
	var view_detail = c.GetString("view_detail")
	if view_detail == "on" {
		view_detail = "1"
	} else {
		view_detail = "0"
	}

	var is_search = c.GetString("is_search")
	if is_search == "on" {
		is_search = "1"
	} else {
		is_search = "0"
	}

	var search_require = c.GetString("search_require")
	if search_require == "on" {
		search_require = "1"
	} else {
		search_require = "0"
	}

	var is_sort = c.GetString("is_sort")
	if is_sort == "on" {
		is_sort = "1"
	} else {
		is_sort = "0"
	}

	var is_total = c.GetString("is_total")
	if is_total == "on" {
		is_total = "1"
	} else {
		is_total = "0"
	}

	var is_export = c.GetString("is_export")
	if is_export == "on" {
		is_export = "1"
	} else {
		is_export = "0"
	}
	var is_import = c.GetString("is_import")
	if is_import == "on" {
		is_import = "1"
	} else {
		is_import = "0"
	}
	var is_import_unique = c.GetString("is_import_unique")
	if is_import_unique == "on" {
		is_import_unique = "1"
	} else {
		is_import_unique = "0"
	}
	var is_navtree = c.GetString("is_navtree")
	if is_navtree == "on" {
		is_navtree = "1"
	} else {
		is_navtree = "0"
	}
	var navtree_sql = c.GetString("navtree_sql")

	var memo = c.GetString("memo")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}
	//校验必填数据
	if field_name == "" || field_code == "" || form_type == "" {
		c.Ctx.WriteString("请填写必要数据!")
		return
	}

	//准备保存
	var sql = `
	update tb_field set
	field_name=?,
	form_type=?,
	form_length=?,
	form_value=?,
	form_tip=?,
	field_defval=?,
	form_sort=?,
	view_list=?,
	view_list_color=?,
	view_edit=?,
	view_form=?,
	view_detail=?,
	is_search=?,
	is_total=?,
	is_navtree=?,
	navtree_sql=?,
	search_require=?,
	is_sort=?,
	is_export=?,
	is_import=?,
	is_import_unique=?,
	memo=?
	where id=?
	`
	var i = db.Exec(sql,
		field_name,
		form_type,
		form_length,
		form_value,
		form_tip,
		field_defval,
		form_sort,
		view_list,
		view_list_color,
		view_edit,
		view_form,
		view_detail,
		is_search,
		is_total,
		is_navtree,
		navtree_sql,
		search_require,
		is_sort,
		is_export,
		is_import,
		is_import_unique,
		memo,
		id,
	)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

//单个表单字段数据保存
func (c *TbController) FieldSet() {
	var id, _ = c.GetInt("id", 0)

	var m = db.First("select * from tb_field where id=?", id)
	if m == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}
	var f = c.GetString("f")
	var v = c.GetString("v")

	var rst = db.Exec("update tb_field set "+f+"=? where id=?", v, id)
	if rst > 0 {
		c.Ctx.WriteString("<font color='yellow'>修改成功!</font>")
	} else {
		c.Ctx.WriteString("保存失败,请稍后重试!")
	}
}

//页面工具栏设置
func (c *TbController) Btns() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	if id < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	//读取tb_table表
	var tb = db.First("select * from tb_table where id=?", id)
	c.Data["tb"] = tb

	//c.TplName="adm/tb/btns.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_btns)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//工具栏字段数据
func (c *TbController) BtnsJson() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	if id < 1 {
		c.Ctx.WriteString("{}")
		return
	}
	//读取tb_table表
	var tb = db.First("select * from tb_table where id=?", id)
	c.Data["tb"] = tb
	//读取按钮列表
	var list = db.Query("select * from tb_linkbtn where tbid=?", id)
	if list == nil || len(list) < 1 {
		c.Ctx.WriteString("{}")
		return
	}
	c.Data["list"] = list

	c.Data["json"] = list
	c.ServeJSON()
}

//按钮编辑
func (c *TbController) BtnEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	var tbid, _ = c.GetInt("tbid", 0)
	c.Data["tbid"] = tbid

	var m = db.First("select * from tb_linkbtn where id=?", id)
	if m != nil {
		c.Data["m"] = m
		c.Data["tbid"] = m["tbid"]
	}

	//c.TplName="adm/tb/btnedit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_btnedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//按钮编辑保存
func (c *TbController) BtnEditPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数
	var tbid, _ = c.GetInt("tbid", 0)
	var title = c.GetString("title")
	var types = c.GetString("type")
	var style = c.GetString("style")
	var is_blank, _ = c.GetInt("is_blank", 0)
	var icon = c.GetString("icon")
	var url = c.GetString("url")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var jstr = c.GetString("jstr")

	//校验必填数据
	if title == "" || types == "" || style == "" || tbid == 0 {
		c.Ctx.WriteString("请填写必要数据!")
		return
	}
	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update tb_linkbtn set
		tbid=?,
		title=?,
		type=?,
		style=?,
		url=?,
		icon=?,
		is_blank=?,
		state=?,
		jstr=?
		where id=?
		`
		var i = db.Exec(sql,
			tbid,
			title,
			types,
			style,
			url,
			icon,
			is_blank,
			state,
			jstr,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into tb_linkbtn(`
		sql += `
		tbid,
		title,
		type,
		style,
		url,
		icon,
		is_blank,
		state,
		jstr)values(
			?,?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			tbid,
			title,
			types,
			style,
			url,
			icon,
			is_blank,
			state,
			jstr,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除数据
func (c *TbController) BtnDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from tb_linkbtn where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//删除tb_table表数据
func (c *TbController) TbDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from tb_table where id=?", id)
	if i > 0 {
		db.Exec("delete from tb_field where tbid=?", id)
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//删除tb_field表数据
func (c *TbController) FieldDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from tb_field where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//数据库管理-列出table表数据
func (c *TbController) Tbs() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	c.Data["id"] = id
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb
	//
	//c.TplName="tbs.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_tbs)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//数据库管理-列出数据库所有表列表
func (c *TbController) TbsJson() {
	var id, _ = c.GetInt("id", 0)
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb

	var sql = ""
	if len(tb) > 0 {
		var dbtype = tb["dbtype"]
		if dbtype == "sqlite" {
			sql = " select * from sqlite_master WHERE `name`!='sqlite_sequence' and type = 'table' or type='view'  "
		} else if dbtype == "mysql" {
			//读取表信息
			sql = "select 0 as rootpage,table_name as name ,table_type as type   from information_schema.tables  where table_schema='" + tb["dbname"] + "'"
		} else if dbtype == "mssql" {
			sql = "Select row_number() over(order by SysObjects.name) as rootpage, Name From SysObjects Where XType='U' order By Name"
		} else if dbtype == "mssql2k" {
			sql = "Select 1 rootpage, Name From SysObjects Where XType='U' order By Name"
		}
		fmt.Println("conn:", tb["conn"])

		var xx = db.NewDb(tb["conn"])
		if xx == nil {
			c.Ctx.WriteString("-1")
			return
		}
		var list = db.Query2(xx, sql)
		if list == nil || len(list) < 1 {
			c.Ctx.WriteString("0")
			return
		}
		fmt.Println("表数据:", list)
		c.Data["json"] = list
		c.ServeJSON()
		return
	}
	c.Data["json"] = "{}"
	c.ServeJSON()
}

//数据库管理-增加表页面
func (c *TbController) TbsAdd() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	c.Data["id"] = id
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb
	//
	//c.TplName="tbsadd.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_tbsadd)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//数据库管理-增加表页面
func (c *TbController) TbsAddPost() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var name = c.GetString("name")
	if name == "" {
		c.Ctx.WriteString("-1")
		return
	}
	c.Data["id"] = id
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb
	//
	var sql = `
	CREATE TABLE "` + name + `" (
		"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		"mch_id"  INTEGER,
		"user_id"  INTEGER
	);
	`
	//sql server 数据库
	if tb["dbtype"] == "mssql" || tb["dbtype"] == "mssql2k" {
		sql = `
		CREATE TABLE  [` + name + `](
			[id] [int] IDENTITY(1000,1) NOT NULL,
			[mch_id] [int] NULL,
			[user_id] [int] NULL,
		 CONSTRAINT [PK_` + name + `] PRIMARY KEY CLUSTERED 
		 (
			[id] ASC
		 )
		) 
		`
	}
	var xx = db.NewDb(tb["conn"])
	if xx == nil {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec2(xx, sql)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	}
	c.Ctx.WriteString("0")
}

//数据库管理-字段列表
func (c *TbController) Fs() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	c.Data["id"] = id
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb

	var tname = c.GetString("tname")
	c.Data["tname"] = tname
	//
	//c.TplName="fs.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_tb_fs")
	tpl.Parse(adm_tb_fs)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//数据库管理-字段结构
func (c *TbController) FsJson() {
	var id, _ = c.GetInt("id", 0)
	var tname = c.GetString("tname")

	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb

	var xx = db.NewDb(tb["conn"])
	if xx == nil {
		c.Ctx.WriteString("-1")
		return
	}

	var sql = ""
	if len(tb) > 0 {
		var dbtype = tb["dbtype"]
		if dbtype == "sqlite" {
			sql = " PRAGMA table_info('" + tname + "') "
			// var l = db.Query2(xx,sql)
			// if l ==nil || len(l)<1{
			// 	c.Ctx.WriteString("0")
			// 	return
			// }

		} else if dbtype == "mysql" {
			//读取字段信息
			sql = "select  ORDINAL_POSITION as cid ,COLUMN_DEFAULT as dflt_value,TABLE_NAME,COLUMN_NAME as name,IS_NULLABLE as notnull,DATA_TYPE as type,COLUMN_KEY as pk,CHARACTER_MAXIMUM_LENGTH AS MAXLENGTH,COLUMN_COMMENT"
			sql += " from Information_schema.columns  where table_Name = '" + tname + "'"
		} else if dbtype == "mssql" {
			sql = `
			SELECT  row_number() over(order by syscolumns.name) as cid,
                 Sysobjects.name AS 'tbname', 
                 syscolumns.name AS 'name',  
                 cast(sys.extended_properties.[value] as varchar(200)) AS memo,  
                 systypes.name AS 'type',  
                 syscolumns.length AS [length],
                 CASE syscolumns.isnullable WHEN '1' THEN 'Y' ELSE 'N' END AS 'nullable',  
                 syscomments.text AS [dflt_value],
                 COLUMNPROPERTY(syscolumns.id, syscolumns.name, 'IsIdentity') AS [pk] ,  
                 CASE WHEN EXISTS (SELECT 1 FROM sysobjects WHERE xtype = 'PK' AND name IN  
                 (SELECT name  
                 FROM sysindexes  
                 WHERE indid IN  
                 (SELECT indid  
                 FROM sysindexkeys  
                 WHERE id = syscolumns.id AND colid = syscolumns.colid)))  
                 THEN 1 ELSE 0 END AS [iskey]  
                 FROM syscolumns  
                 INNER JOIN systypes  
                 ON syscolumns.xtype = systypes.xtype  
                 LEFT JOIN sysobjects ON syscolumns.id = sysobjects.id  
                 LEFT OUTER JOIN sys.extended_properties ON  
                 ( sys.extended_properties.minor_id = syscolumns.colid  
                 AND sys.extended_properties.major_id = syscolumns.id)  
                 LEFT OUTER JOIN syscomments ON syscolumns.cdefault = syscomments.id  
                 WHERE (systypes.name <> 'sysname')  
                 AND syscolumns.id IN (SELECT id FROM SYSOBJECTS WHERE  NAME = '` + tname + `')  
			`
		} else if dbtype == "mssql2k" {
			sql = `
			SELECT           1 as cid,
                 sysobjects.name AS [tbname],   
                 syscolumns.name AS [name],  
                 cast(properties.[value] as varchar(200)) AS [memo],    
                 systypes.name AS [type],  
                 syscolumns.length AS [length], 
                 CASE syscolumns.isnullable WHEN '1' THEN 'Y' ELSE 'N' END AS [nullable],  
                 CASE WHEN syscomments.text IS NULL THEN '' ELSE syscomments.text END AS [dflt_value],  
                 CASE WHEN COLUMNPROPERTY(syscolumns.id, syscolumns.name, 'IsIdentity') = 1 THEN 1 ELSE 0 END AS [pk],  
                 CASE WHEN EXISTS (SELECT 1 FROM sysobjects WHERE xtype = 'PK' AND name IN  
                 (SELECT name  
                 FROM sysindexes  
                 WHERE indid IN  
                 (SELECT indid  
                 FROM sysindexkeys  
                 WHERE id = syscolumns.id AND colid = syscolumns.colid)))  
                 THEN 1 ELSE 0 END AS [iskey]  
                 FROM syscolumns INNER JOIN  
                 sysobjects ON sysobjects.id = syscolumns.id INNER JOIN  
                 systypes ON syscolumns.xtype = systypes.xtype LEFT OUTER JOIN  
                 sysproperties properties ON syscolumns.id = properties.id AND  
                 syscolumns.colid = properties.smallid LEFT OUTER JOIN  
                 sysproperties ON sysobjects.id = sysproperties.id AND  
                 sysproperties.smallid = 0 LEFT OUTER JOIN  
                 syscomments ON syscolumns.cdefault = syscomments.id  
				 WHERE (sysobjects.xtype = 'u')  
                 AND sysobjects.NAME='` + tname + `' 
                 ORDER BY [tbname]   

			`
		}
		fmt.Println("sql", sql)

		var list = db.Query2(xx, sql)
		if list == nil || len(list) < 1 {
			c.Ctx.WriteString("0")
			return
		}
		c.Data["json"] = list
		c.ServeJSON()
		return
	}
	c.Data["json"] = "{}"
	c.ServeJSON()
}

//数据库管理-字段添加
func (c *TbController) FsAdd() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	c.Data["id"] = id
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb
	var tname = c.GetString("tname")
	c.Data["tname"] = tname
	//
	//c.TplName="fsadd.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_tb_fsadd")
	tpl.Parse(adm_tb_fsadd)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//数据库管理-增加表字段
func (c *TbController) FsAddPost() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var tname = c.GetString("tname")
	if tname == "" {
		c.Ctx.WriteString("-1")
		return
	}

	var name = c.GetString("name")
	if name == "" {
		c.Ctx.WriteString("-1")
		return
	}
	var ftype = c.GetString("ftype")
	if ftype == "" {
		c.Ctx.WriteString("-1")
		return
	}
	var defval = c.GetString("defval")
	var length = c.GetString("length")

	c.Data["id"] = id
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb

	//sqlite语句
	var sql = `ALTER TABLE "` + tname + `" ADD ` + name + ` ` + ftype
	if defval != "" {
		sql += " DEFAULT '" + defval + "' "
	}
	//mysql语句和mssql语句
	if tb["dbtype"] != "sqlite" {
		sql = `ALTER TABLE ` + tname + ` ADD ` + name + ` `
		if ftype != "int" {
			sql += ftype + "" + length
		} else {
			sql += ftype + " "
		}
		if defval != "" {
			sql += " DEFAULT '" + defval + "' "
		}
	}
	var xx = db.NewDb(tb["conn"])
	if xx == nil {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec2(xx, sql)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	}
	c.Ctx.WriteString("0")
}
func (c *TbController) Conn() {
	//c.TplName="adm/tb/conn.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_conn)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) ConnJson() {
	var rst = db.Query("select * from adm_conn ")
	c.Data["json"] = rst
	c.ServeJSON()
}
func (c *TbController) ConnEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	// if id<1{
	// 	c.Ctx.WriteString("参数错误!")
	// 	return
	// }
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	c.Data["m"] = tb

	//c.TplName="adm/tb/connedit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_tb_connedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) ConnDel() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	//读取adm_conn表
	var tb = db.First("select * from adm_conn where id=?", id)
	if len(tb) < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var fields = db.Query("select id from tb_table where conn_str=?", tb["conn"])
	if len(fields) > 0 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from adm_conn where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	}
	c.Ctx.WriteString("0")
	return
}
func (c *TbController) Clone() {
	var id, err = c.GetInt("id")
	if err != nil {
		c.Ctx.WriteString("-")
		return
	}
	var tb = db.First("select * from tb_table where id=?", id)
	if len(tb) < 1 {
		c.Ctx.WriteString("-")
		return
	}

	var tbid = "0"
	var tbcode = "z" + db.RandomString(9)

	var sql = `
	insert into tb_table(
		mch_id,
		user_id,
		proj_id,
		title,` + "`code`" + `,
		conn_id,
		conn_str,` + "`table`" + `,
		where_str,
		pri_key,
		sort_key,
		data_type,
		is_total,
		is_create,
		is_edit,
		is_detail,
		is_del,
		is_import,
		is_export,
		edit_style,
		edit_width,
		edit_height,
		ex_javascript,
		ex_linkbutton,
		ex_html_operate,
		memo,
		state) 
		select 
		mch_id,
		user_id,
		proj_id,
		concat(title,'-clone'),` + "`code`," + `
		conn_id,
		conn_str,` + "`table`" + `,
		where_str,
		pri_key,
		sort_key,
		data_type,
		is_total,
		is_create,
		is_edit,
		is_detail,
		is_del,
		is_import,
		is_export,
		edit_style,
		edit_width,
		edit_height,
		ex_javascript,
		ex_linkbutton,
		ex_html_operate,
		memo,
		state 
		from tb_table where id=?
	`
	var i = db.Exec(sql, id)
	if i > 0 {
		var au = db.First("SELECT max(id) as id from tb_table")
		tbid = au["id"]

		sql = `
			insert into tb_field(
			tbid,
			tbcode,
			field_name,
			field_code,
			field_type,
			field_length,
			field_defval,
			field_prikey,
			form_type,
			form_length,
			form_sort,
			form_tip,
			form_value,
			form_required,
			view_list,
			view_list_color,
			view_form,
			view_edit,
			view_detail,
			is_sort,
			is_search,
			search_require,
			is_import,
			is_import_unique,
			is_export,
			is_total,
			is_navtree,
			navtree_sql,
			memo,
			state
			)
			select 
			` + tbid + ` as tbid,
			'` + tbcode + `' as tbcode,
			field_name,
			field_code,
			field_type,
			field_length,
			field_defval,
			field_prikey,
			form_type,
			form_length,
			form_sort,
			form_tip,
			form_value,
			form_required,
			view_list,
			view_list_color,
			view_form,
			view_edit,
			view_detail,
			is_sort,
			is_search,
			search_require,
			is_import,
			is_import_unique,
			is_export,
			is_total,
			is_navtree,
			navtree_sql,
			memo,
			state
			from tb_field where tbid=?	
			`
		var j = db.Exec(sql, id)
		if j < 1 {
			db.Exec("delete from tb_table where id=?", tbid)
			c.Ctx.WriteString("0")
			return
		}
	}
	c.Ctx.WriteString(strconv.FormatInt(i, 10))
}
func (c *TbController) ConnPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数
	var title = c.GetString("title")
	var server = c.GetString("server")
	var conn = c.GetString("conn")
	var dbtype = c.GetString("dbtype")
	var port = c.GetString("port")
	var dbname = c.GetString("dbname")
	var uid = c.GetString("uid")
	var pwd = c.GetString("pwd")
	var memo = c.GetString("memo")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}
	//校验必填数据
	if dbtype != "sqlite" && (title == "" || server == "" || dbtype == "") {
		c.Ctx.WriteString("请填写必要数据!")
		return
	}
	if dbtype == "sqlite" && (title == "" || dbname == "" || dbtype == "") {
		c.Ctx.WriteString("请填写必要数据!")
		return
	}
	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update adm_conn set
		title=?,
		dbtype=?,
		conn=?,
		dbname=?,
		server=?,
		port=?,
		uid=?,
		pwd=?,
		memo=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			title,
			dbtype,
			conn,
			dbname,
			server,
			port,
			uid,
			pwd,
			memo,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into adm_conn(`
		sql += `
		title,
			dbtype,
			conn,
			dbname,
			server,
			port,
			uid,
			pwd,
			state,
			memo)values(
			?,?,?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			title,
			dbtype,
			conn,
			dbname,
			server,
			port,
			uid,
			pwd,
			state,
			memo,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//-----------------------------------------------------------------------------
//页面模板编辑
func (c *TbController) Page() {
	var tpl = template.New("")
	tpl.Parse(adm_tb_page)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) PageJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		qtxt = " where title like '%" + qtxt + "%'"
	}
	qtxt += " order by id desc "
	var rst = db.Pager(page, pageSize, "select * from page_list "+qtxt)

	//var rst = db.Query("select * from page_list order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}
func (c *TbController) PageEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//读取page_list表
	var tb = db.First("select * from page_list where id=?", id)
	c.Data["m"] = tb
	//数据库连接信息
	var dblist = db.Query("select * from adm_conn")
	c.Data["dblist"] = dblist

	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	//tpl.Parse(adm_tb_pageedit)
	tpl.Parse(adm_tb_pageedit_tab)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) PagePost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数
	var title = c.GetString("title")
	var code = c.GetString("code")
	var conn_str = c.GetString("conn_str")
	var description = c.GetString("description")
	var module = c.GetString("module")
	var token = c.GetString("token")
	var template = c.GetString("template")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update page_list set
		title=?,
		code=?,
		conn_str=?,
		description=?,
		module=?,
		token=?,
		template=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			title,
			code,
			conn_str,
			description,
			module,
			token,
			template,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into page_list(`
		sql += `
		title,
		code,
		conn_str,
		description,
		module,
		token,
		template,
		state
		)values(
			?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			title,
			code,
			conn_str,
			description,
			module,
			token,
			template,
			state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//-------------------------------------------------------------------------------------------
//API模板编辑
func (c *TbController) API() {
	var tpl = template.New("")
	tpl.Parse(adm_tb_api)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) ApiJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		qtxt = " where title like '%" + qtxt + "%'"
	}
	qtxt += " order by id desc "
	var rst = db.Pager(page, pageSize, "select * from api_list "+qtxt)

	//var rst = db.Query("select * from api_list order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}
func (c *TbController) ApiEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//读取api_list表
	var tb = db.First("select * from api_list where id=?", id)
	c.Data["m"] = tb
	//数据库连接信息
	var dblist = db.Query("select * from adm_conn")
	c.Data["dblist"] = dblist
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	//tpl.Parse(adm_tb_apiedit)
	tpl.Parse(adm_tb_apiedit_tab)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) ApiPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数
	var title = c.GetString("title")
	var api_code = c.GetString("code")
	var conn_str = c.GetString("conn_str")
	var description = c.GetString("description")
	var module = c.GetString("module")
	var token = c.GetString("token")
	var api_template = c.GetString("template")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update api_list set
		title=?,
		api_code=?,
		conn_str=?,
		description=?,
		module=?,
		token=?,
		api_template=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			title,
			api_code,
			conn_str,
			description,
			module,
			token,
			api_template,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into api_list(`
		sql += `
		title,
		api_code,
		conn_str,
		description,
		module,
		token,
		api_template,
		state
		)values(
			?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			title,
			api_code,
			conn_str,
			description,
			module,
			token,
			api_template,
			state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//-------------------------------------------------------------------------------------------
//报表模板编辑
func (c *TbController) Rpt() {
	var tpl = template.New("")
	tpl.Parse(adm_tb_rpt)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) RptJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		qtxt = " where title like '%" + qtxt + "%'"
	}
	qtxt += " order by id desc "
	var rst = db.Pager(page, pageSize, "select * from rpt_list "+qtxt)

	//var rst = db.Query("select * from rpt_list order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}
func (c *TbController) RptEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//读取rpt_list表
	var tb = db.First("select * from rpt_list where id=?", id)
	c.Data["m"] = tb
	//数据库连接信息
	var dblist = db.Query("select * from adm_conn")
	c.Data["dblist"] = dblist

	var tpl = template.New("")
	//tpl.Parse(adm_tb_rptedit)
	tpl.Parse(adm_tb_rptedit_tab)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *TbController) RptPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数
	var title = c.GetString("title")
	var code = c.GetString("code")
	var conn_str = c.GetString("conn_str")
	var description = c.GetString("description")
	var module = c.GetString("module")
	var token = c.GetString("token")
	var template = c.GetString("template")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update rpt_list set
		title=?,
		code=?,
		conn_str=?,
		description=?,
		module=?,
		token=?,
		template=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			title,
			code,
			conn_str,
			description,
			module,
			token,
			template,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into rpt_list(`
		sql += `
		title,
		code,
		conn_str,
		description,
		module,
		token,
		template,
		state
		)values(
			?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			title,
			code,
			conn_str,
			description,
			module,
			token,
			template,
			state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

func (c *TbController) RptDesign() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//读取rpt_list表
	var tb = db.First("select * from rpt_list where id=?", id)
	c.Data["m"] = tb
	//数据库连接信息
	var dblist = db.Query("select * from adm_conn")
	c.Data["dblist"] = dblist

	var tpl = template.New("")

	tpl.Parse(adm_tb_rptdesign)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

///保存报表设计
func (c *TbController) RptDesignSave() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id
	var codes = c.GetString("codes")

	var i = db.Exec("update rpt_list set template=? where id=?", codes, id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

var adm_tb_rptdesign = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.m.title}}-报表设计</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />
	

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

	<script language="javascript" src="/js/lodop/LodopFuncs.js"></script>
	<object id="LODOP1" classid="clsid:2105C259-1E0C-4534-8141-A753534CB4CA" width=0 height=0> 
		<embed id="LODOP_EM1" TYPE="application/x-print-lodop" width=0 height=0 PLUGINSPAGE="install_lodop32.exe"></embed>
	</object> 
    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
    $(function(){

    })
	function SaveDesign(){
		getProgram();
		jQuery.post('/adm/tb/rptdesignsave',{'id':{{.id}},'codes':$('#codes').val()},function(data){
			alert(data);
		})
	}
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">
    <div id="tb" style="padding:5px;height:auto">

	<div class="easyui-panel" title="报表设计" style="width:850px;padding:10px;" data-options="tools:'#tt'">
	
	&nbsp;&nbsp;&nbsp;                                        
一般当窗口弹出显示时，各种ADD或SET语句就无法继续执行，但如果窗口是内嵌的，则可以执行。<br>
下面首先打开显示<a href="javascript:DisplayDesign()">打印设计</a>或<a href="javascript:DisplaySetup()">打印维护</a>窗口,然后点击后面的ADD或SET语句，看看追加效果：<br>

<object id="LODOP2" classid="clsid:2105C259-1E0C-4534-8141-A753534CB4CA" width=810 height=407> 
  <param name="Caption" value="内嵌显示区域">
  <param name="Border" value="1">
  <param name="Color" value="#C0C0C0">
  <embed id="LODOP_EM2" TYPE="application/x-print-lodop" width=810 height=407 PLUGINSPAGE="install_lodop.exe">
</object> 


<textarea id="codes" name="codes" rows="6" id="S1" cols="107">返回的结果值</textarea>
	</div>
	<div id="tt">
	<a href="javascript:void(0)" class="icon-38" title="保存" onclick="SaveDesign();return false;"></a> 
	
	</div>

    </div>

	<script language="javascript" type="text/javascript"> 
	var LODOP; //声明为全局变量 
	function Moditify(item){
		LODOP=getLodop(document.getElementById('LODOP2'),document.getElementById('LODOP_EM2')); 	
        	if ((!LODOP.GET_VALUE("ItemIsAdded",item.name))&&(item.checked)){
		LODOP.ADD_PRINT_TEXTA(item.name,56,32,175,30,item.value); } else {
		LODOP.SET_PRINT_STYLEA(item.name,'Deleted',!item.checked);}
	}	
	function CreatePage(){
		LODOP=getLodop(document.getElementById('LODOP2'),document.getElementById('LODOP_EM2')); 
		LODOP.PRINT_INITA(0,0,760,321,"打印控件功能演示_Lodop功能_在线编辑获得程序代码");
		LODOP.ADD_PRINT_TEXT(10,50,175,30,"先加的内容");
	};	
	function DisplayDesign() {		
		CreatePage();
		LODOP.SET_SHOW_MODE("DESIGN_IN_BROWSE",1);
		LODOP.SET_SHOW_MODE("SETUP_ENABLESS","11111111000000");//隐藏关闭(叉)按钮
		LODOP.SET_SHOW_MODE("HIDE_GROUND_LOCK",true);//隐藏纸钉按钮
		LODOP.PRINT_DESIGN();		
	};
	function DisplaySetup() {		
		CreatePage();
		LODOP.SET_SHOW_MODE("SETUP_IN_BROWSE",1);
		LODOP.SET_SHOW_MODE("MESSAGE_NOSET_PROPERTY",'不能设置属性，请用打印设计(本提示可修改)！');
		LODOP.PRINT_SETUP();		
	};
	function Addhtm() {	
		LODOP.ADD_PRINT_HTM(45,494,288,88,"<table border='1'>\n<tr>\n<td>表格11</td>\n<td>表格12</td>\n<td>表格13</td>\n</tr>\n<tr>\n<td>表格21</td>\n<td>表格22</td>\n<td>表格23</td>\n</tr>\n</table>");
	};
	function SetBKIMG() {
		LODOP=getLodop(document.getElementById('LODOP2'),document.getElementById('LODOP_EM2')); 
                LODOP.ADD_PRINT_SETUP_BKIMG("<img border='0' src='http://s1.sinaimg.cn/middle/721e77e5t99431b026bd0&690'>");	

	};
	function getProgram() {	
		LODOP=getLodop(document.getElementById('LODOP2'),document.getElementById('LODOP_EM2')); 
		if (LODOP.CVERSION) LODOP.On_Return=function(TaskID,Value){document.getElementById('codes').value=Value;};	
		document.getElementById('codes').value=LODOP.GET_VALUE("ProgramCodes",0);
		//document.getElementById('button01').style.display=""; 	
	};	
	function prn_Preview() {		
		LODOP=getLodop(document.getElementById('LODOP1'),document.getElementById('LODOP_EM1')); 
		eval(document.getElementById('S1').value); 
		LODOP.PREVIEW();
		LODOP=getLodop(document.getElementById('LODOP2'),document.getElementById('LODOP_EM2')); 
	};	
	function getMyValue(strType,oResultOB){
		var LODOP=getLodop(document.getElementById('LODOP_X'),document.getElementById('LODOP_EM')); 
		if (LODOP.CVERSION) CLODOP.On_Return=function(TaskID,Value){if (oResultOB) oResultOB.value=Value;}; 
		var stResult=LODOP.GET_VALUE(strType,"0");
		if (!LODOP.CVERSION) oResultOB.value=stResult; 
	};
</script> 

</body>
</html>
`

//删除page_list表数据
func (c *TbController) PageDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from page_list where id=?", id)
	if i > 0 {
		db.Exec("delete from page_param where page_id=?", id)
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//删除api_list表数据
func (c *TbController) ApiDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from api_list where id=?", id)
	if i > 0 {
		db.Exec("delete from api_param where api_id=?", id)
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//删除rpt_list表数据
func (c *TbController) RptDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from rpt_list where id=?", id)
	if i > 0 {
		db.Exec("delete from rpt_param where rpt_id=?", id)
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//-------------------------------------------------------------------------------------------
//报表参数编辑
func (c *TbController) RptParam() {
	var rptid = c.GetString("rptid")
	if rptid == "" {
		rptid = "0"
	}
	c.Data["rptid"] = rptid

	var tpl = template.New("")
	tpl.Parse(adm_tb_rptparam)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//参数列表
func (c *TbController) RptParamJson() {
	var rptid = c.GetString("rptid")
	if rptid == "" {
		rptid = "0"
	}
	c.Data["rptid"] = rptid
	var rst = db.Query("select * from rpt_param where rpt_id=" + rptid + " order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}

//参数编辑页面
func (c *TbController) RptParamEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id
	var rptid, _ = c.GetInt("rptid", 0)
	if rptid < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	c.Data["rptid"] = rptid
	//读取rpt_param表
	var tb = db.First("select * from rpt_param where id=?", id)
	c.Data["m"] = tb

	var tpl = template.New("")
	tpl.Parse(adm_tb_rptparamedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//参数数据保存
func (c *TbController) RptParamPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数

	var rpt_id = c.GetString("rpt_id")
	var title = c.GetString("title")
	var param_type = c.GetString("param_type")
	var param_name = c.GetString("param_name")
	var max_length = c.GetString("max_length")
	var is_require = c.GetString("is_require")
	if is_require == "on" {
		is_require = "1"
	} else {
		is_require = "0"
	}
	var param_regex = c.GetString("param_regex")
	var param_value = c.GetString("param_value")
	var memo = c.GetString("memo")

	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update rpt_param set
		rpt_id=?,
		title=?,
		param_type=?,
		param_name=?,
		max_length=?,
		is_require=?,
		param_regex=?,
		param_value=?,
		memo=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			rpt_id,
			title,
			param_type,
			param_name,
			max_length,
			is_require,
			param_regex,
			param_value,
			memo,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into rpt_param(`
		sql += `
		rpt_id,
		title,
		param_type,
		param_name,
		max_length,
		is_require,
		param_regex,
		param_value,
		memo ,
		state
		)values(
			?,?,?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			rpt_id,
			title,
			param_type,
			param_name,
			max_length,
			is_require,
			param_regex,
			param_value,
			memo,
			state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除rpt_param表数据
func (c *TbController) RptParamDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from rpt_param where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//-------------------------------------------------------------------------------------------
//页面参数编辑
func (c *TbController) PageParam() {
	var pageid = c.GetString("pageid")
	if pageid == "" {
		pageid = "0"
	}
	c.Data["pageid"] = pageid

	var tpl = template.New("")
	tpl.Parse(adm_tb_pageparam)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//参数列表
func (c *TbController) PageParamJson() {
	var pageid = c.GetString("pageid")
	if pageid == "" {
		pageid = "0"
	}
	c.Data["pageid"] = pageid
	var rst = db.Query("select * from page_param where page_id=" + pageid + " order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}

//参数编辑页面
func (c *TbController) PageParamEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id
	var pageid, _ = c.GetInt("pageid", 0)
	if pageid < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	c.Data["pageid"] = pageid
	//读取page_param表
	var tb = db.First("select * from page_param where id=?", id)
	c.Data["m"] = tb

	var tpl = template.New("")
	tpl.Parse(adm_tb_pageparamedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//参数数据保存
func (c *TbController) PageParamPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数

	var page_id = c.GetString("page_id")
	var title = c.GetString("title")
	var param_type = c.GetString("param_type")
	var param_name = c.GetString("param_name")
	var max_length = c.GetString("max_length")
	var is_require = c.GetString("is_require")
	if is_require == "on" {
		is_require = "1"
	} else {
		is_require = "0"
	}
	var param_regex = c.GetString("param_regex")
	var param_value = c.GetString("param_value")
	var memo = c.GetString("memo")

	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update page_param set
		page_id=?,
		title=?,
		param_type=?,
		param_name=?,
		max_length=?,
		is_require=?,
		param_regex=?,
		param_value=?,
		memo=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			page_id,
			title,
			param_type,
			param_name,
			max_length,
			is_require,
			param_regex,
			param_value,
			memo,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into page_param(`
		sql += `
		page_id,
		title,
		param_type,
		param_name,
		max_length,
		is_require,
		param_regex,
		param_value,
		memo ,
		state
		)values(
			?,?,?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			page_id,
			title,
			param_type,
			param_name,
			max_length,
			is_require,
			param_regex,
			param_value,
			memo,
			state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除page_param表数据
func (c *TbController) PageParamDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from page_param where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

//-------------------------------------------------------------------------------------------
//API接口参数编辑
func (c *TbController) ApiParam() {
	var apiid = c.GetString("apiid")
	if apiid == "" {
		apiid = "0"
	}
	c.Data["apiid"] = apiid

	var tpl = template.New("")
	tpl.Parse(adm_tb_apiparam)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//参数列表
func (c *TbController) ApiParamJson() {
	var apiid = c.GetString("apiid")
	if apiid == "" {
		apiid = "0"
	}
	c.Data["apiid"] = apiid
	var rst = db.Query("select * from api_param where api_id=" + apiid + " order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}

//参数编辑页面
func (c *TbController) ApiParamEdit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id
	var apiid, _ = c.GetInt("apiid", 0)
	if apiid < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	c.Data["apiid"] = apiid
	//读取api_param表
	var tb = db.First("select * from api_param where id=?", id)
	c.Data["m"] = tb

	var tpl = template.New("")
	tpl.Parse(adm_tb_apiparamedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//参数数据保存
func (c *TbController) ApiParamPost() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//获取参数

	var api_id = c.GetString("api_id")
	var title = c.GetString("title")
	var param_type = c.GetString("param_type")
	var param_name = c.GetString("param_name")
	var max_length = c.GetString("max_length")
	var is_require = c.GetString("is_require") //必填
	if is_require == "on" {
		is_require = "1"
	} else {
		is_require = "0"
	}
	var is_unique = c.GetString("is_unique") //唯一
	if is_unique == "on" {
		is_unique = "1"
	} else {
		is_unique = "0"
	}
	//var is_unique_info = c.GetString("is_unique_info") //唯一提示

	var is_checkout = c.GetString("is_checkout") //是否验证
	if is_checkout == "on" {
		is_checkout = "1"
	} else {
		is_checkout = "0"
	}

	var param_regex = c.GetString("param_regex")
	var param_value = c.GetString("param_value")
	var memo = c.GetString("memo")

	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ``
	if id > 0 {
		//准备保存
		sql = `
		update api_param set
		api_id=?,
		title=?,
		param_type=?,
		param_name=?,
		max_length=?,
		is_require=?,
		param_regex=?,
		param_value=?,
		memo=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			api_id,
			title,
			param_type,
			param_name,
			max_length,
			is_require,
			param_regex,
			param_value,
			memo,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `insert into api_param(`
		sql += `
		api_id,
		title,
		param_type,
		param_name,
		max_length,
		is_require,
		param_regex,
		param_value,
		memo ,
		state
		)values(
			?,?,?,?,?,?,?,?,?,?
		);
		`
		var i = db.Exec(sql,
			api_id,
			title,
			param_type,
			param_name,
			max_length,
			is_require,
			param_regex,
			param_value,
			memo,
			state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除api_param表数据
func (c *TbController) ApiParamDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from api_param where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

var adm_tb_list = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
            proj_id:$('#fproj_id').combobox("getValue"),
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
            //$('#win').window('open');
            //$('#win').window('refresh', '/adm/tb/edit?id='+row.id);
            //$('#ff').form('load',row);

			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/edit?id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}
function doClone(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/clone?id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}
function doField(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
            //$('#win').window('open');
            //$('#win').window('refresh', '/adm/tb/fields?id='+row.id);
            //$('#ff').form('load',row);
			var w=$('#win').window({
				width:520,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/fields?id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}
function doFieldSet(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			top.addTab(row.title+'-表单','/adm/tb/fieldlist?id='+row.id);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}
function doBtn(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			top.addTab(row.title+'-表单','/adm/tb/btns?id='+row.id);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}
function doEx(){
	var row = $('#tt').datagrid('getSelected');
	if (row){
		top.addTab(row.title+'-表单','/adm/tb/tbex?tbid='+row.id);
	}else{
		jq.messager.alert('警告','请选择一行数据','warning');
	}
}
function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	w.window('open');
	w.window('refresh', '/adm/tb/edit?id=');
}
function doData(title,url){
	top.addTab(title,url);
}

function doRemove(){

    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/tbdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })
    function rowformater_proj(value, row, index) {
		return '<div class="proj_id">'+value+'</div>';
    }
	function rowformater_field(value, row, index) {
		//var a= "<a href='/adm/tb/fieldlist?id="+row.id+"' >设置</a>";
		//var b="&nbsp;<a href='/adm/mdata/list?code="+row.code+"' >数据</a>";
		var b="&nbsp;<a href='#' onclick='doData(\""+row.title+"\",\"/adm/d/list/"+row.code+"\");' >数据</a>";
		return b;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/jsonlist" data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="模块管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
                <th field="id" width="5">ID</th>
                <th field="proj_id" data-options="formatter:rowformater_proj" width="10">项目</th>
                <th field="title" width="10">模块</th>
                <th field="code" align="center" width="10">代号</th>
                <th field="table" align="center" width="10" >表名</th>
                <th field="pri_key" width="5">主键</th>
				<th field="is_edit" width="5">编辑</th>
				<th field="is_del" width="5">删除</th>
                <th field="memo" width="10"  >备注</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-46" plain="true" onclick="doClone();">克隆</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            
			
			<a href="#" class="easyui-linkbutton" iconcls="icon-0" plain="true" onclick="doField();">字段</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doFieldSet();">表单</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-3" plain="true" onclick="doBtn();">按钮</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-41" plain="true" onclick="doEx();">扩展</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true" onclick="doRemove();">删除</a>
            
        </div>
        <div>
            项目:<select class="easyui-combobox" style="width:120px;"  editable="false" id="fproj_id" name="proj_id">
                <option value="">请选择...</option>
                {{range $k,$v:=.xlist}}
                <option value="{{$v.id}}">{{$v.proj_name}}</option>
                {{end}}
            </select>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


			<a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-0" onclick="doData('项目管理','/adm/tb/proj');">项目管理</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
            eval({{.jsval}});
        }
    });
</script>
</body>
</html>
`
var adm_tb_fields = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', ''+data, "warning");
                    }else{
						jq.messager.alert('错误', "操作失败!", "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<style type="text/css">
html{
overflow-x:hidden;
}
.grid{font:12px arial,helvetica,sans-serif;border:1px solid #8DB2E3}
.grid td{font:100% arial,helvetica,sans-serif;height:24px;padding:5px}
.grid{width:100%;border-collapse:collapse}
.grid th{background:#E7F3FE;height:27px;line-height:27px;border:1px solid #8DB2E3;padding-left:5px}
.grid td{border:1px solid #8DB2E3;padding-left:5px}
</style>
 
    <div style="padding:2px 2px 2px 2px">

		

		<form id="form1" name="form1"   action="fieldspost" method="post">
 
             
            <table class="grid"  >
                <thead>
                    <tr>
                        <th>#</th>
                        <th>字段</th> 
                        <th>类型</th>
						<th>可空</th>
                        <th>默认值</th> 
						<th>主键</th>
						<th>说明</th>
                    </tr>
					
                </thead>
				    {{range $i,$m :=.list}}
					<tr>
						<td><input type='checkbox' value="{{$m.name}}" id='{{$m.name}}' name='fields'></td>
						<td>{{$m.name}}</td>
						<td>{{$m.type}}</td>
						<td>{{$m.notnull}}</td>
						<td>{{$m.dflt_value}}</td>
						<td>{{$m.pk}}</td>
						<td>{{$m.memo}}</td>
					</tr>
					{{end}}
            </table>
<input type="hidden" value="{{.tb}}" name="tb" />
<input type="hidden" value="{{.id}}" name="id" />
<script type="text/javascript">
    <!--
    $(function(){
        {{range $k,$f :=.flist}}
            $('#{{$f.field_code}}').attr('checked',true);
        {{end}}
    
    });
    
    //-->
    </script>
</form>


        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
 
<script type="text/javascript">

</script>

`
var adm_tb_fieldlist = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
    <script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }

function doField(id){
        $('#win').window('open');
        $('#win').window('refresh', '/adm/tb/fieldedit?id='+id);
}
function doFieldDel(id){
	if(confirm('确实要删除吗?')){
		jQuery.post('/adm/tb/fielddel?id='+id,function(data){
			doSearch();
		})
	}
}
    $(function(){
		
    })
	function rowformater_field(value, row, index) {
		var e = '<a   href="#" onclick="doField(' + row.id + ');" >编辑</a> ';
		var f = '<a   href="#" onclick="doFieldDel(' + row.id + ');" >删除</a> ';
        return e+f;
    }
    function rowformater_name(value, row, index) {
		var e = '<div class="canedit" oldval="'+value+'" val="'+value+'" sp="field_name" valid="' + row.id + '" >'+value+'</div> ';
        return e;
    }
    function rowformater_sort(value, row, index) {
		var e = '<div class="canedit" oldval="'+value+'" val="'+value+'" sp="form_sort" valid="' + row.id + '" >'+value+'</div> ';
        return e;
    }
    function rowformater_search(value, row, index) {
		//var e = '<div class="canedit" oldval="'+value+'" val="'+value+'" sp="is_search" valid="' + row.id + '" >'+value+'</div> ';
        //return e;
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'is_search\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="is_search'+row.id+'" id="is_search'+row.id+'">'
        return e;
    }
    function rowformater_list(value, row, index) {
		//var e = '<div class="canedit" oldval="'+value+'" val="'+value+'" sp="view_list" valid="' + row.id + '" >'+value+'</div> ';
        //return e;
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'view_list\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="view_list'+row.id+'" id="view_list'+row.id+'">'
        return e;
	}
	function rowformater_edit(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'view_edit\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="view_edit'+row.id+'" id="view_edit'+row.id+'">'
        return e;
	}
	function rowformater_detail(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'view_detail\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="view_detail'+row.id+'" id="view_detail'+row.id+'">'
        return e;
	}
	function rowformater_import(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'is_import\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="is_import'+row.id+'" id="is_import'+row.id+'">'
        return e;
	}
	function rowformater_export(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'is_export\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="is_export'+row.id+'" id="is_export'+row.id+'">'
        return e;
    }
	function rowformater_issort(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'is_search\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="is_search'+row.id+'" id="is_search'+row.id+'">'
        return e;
    }
    function rowformater_form(value, row, index) {
		//var e = '<div class="canedit" oldval="'+value+'" val="'+value+'" sp="view_form" valid="' + row.id + '" >'+value+'</div> ';
        //return e;
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'view_form\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="view_form'+row.id+'" id="view_form'+row.id+'">'
        return e;
    }
    function rowformater_memo(value, row, index) {
        if(value+''=='')value='&nbsp;';
		var e = '<div style="min-width:20px;"  class="canedit" oldval="'+value+'" val="'+value+'" sp="memo" valid="' + row.id + '" >'+value+'</div> ';
        return e;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;"  fit="true">

    <table class="easyui-datagrid"  
           url="/adm/tb/fieldlistjson?id={{.id}}"  
           title="{{.tb.title}}-表单设置" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
                <th field="id" width="5">ID</th>
                <th field="field_code"  width="20" >字段</th>
                <th field="field_name" align="center" data-options="formatter:rowformater_name" width="20" >名称</th>
                <th field="form_type" align="center"  width="10" >类型</th>
                <th field="field_length"  width="10" >长度</th>
                <th field="view_list"  width="10"  data-options="formatter:rowformater_list" >列表</th> 
				<th field="view_form"  width="10"  data-options="formatter:rowformater_form" >表单</th> 
				<th field="view_edit"  width="10"  data-options="formatter:rowformater_edit" >编辑</th> 
				<th field="view_detail"  width="10"  data-options="formatter:rowformater_detail" >详情</th> 
				<th field="is_sort"  width="10"  data-options="formatter:rowformater_issort" >排序</th> 
				<th field="is_search"  width="10"  data-options="formatter:rowformater_search" >搜索</th>
				<th field="is_import"  width="10"  data-options="formatter:rowformater_import" >导入</th>
				<th field="is_export"  width="10"  data-options="formatter:rowformater_export" >导出</th> 
				<th field="form_sort"  width="10"  data-options="formatter:rowformater_sort" >顺序</th> 
                <th field="memo"  width="10"  data-options="formatter:rowformater_memo">备注</th>
				<th field=" "  data-options="formatter:rowformater_field"  width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:480px;height:390px;padding:5px;">
        Some Content.
    </div>
    <script language="javascript">
        function funField(f, id,v) {
			if(v=="on"){
				v=true;
			}
            if (v == true) {
                    v = '1';
            } else {
                    v = '0';
            }
            jQuery.post('/adm/tb/fieldset', { 'tb': 'tb_field', 'id': id, 'f': f, 'v': v }, function (data) {
                if (data == "1") {
                    layer.msg('保存成功!');
                }
                else {
                    layer.msg(data);
                }
            })
        }
        $(function () {
            $('#tt').datagrid({   
                onLoadSuccess:function(data){

                    $(".canedit").bind("dblclick", function () {
                        var field=$(this).attr('sp');
                        var valid=$(this).attr('valid');
                        var oldval=$(this).attr('val');

                        var input = "<input type='text' id='temp' value=" + $(this).text() + " >";
                        $(this).text("");
                        $(this).append(input);
                        $("input#temp").focus();
                        $("input#temp").blur(function () {
                            if ($(this).val() == "") {
                                $(this).closest("div").text(oldval);
                                $(this).remove();
                            } else {
                                $(this).closest("div").text($(this).val());
                                if($(this).val()==oldval) return;
                                //alert($(this).val()+"-"+oldval);

                                jQuery.post('/adm/tb/fieldset',{'tb':'tb_field','id':valid,'f':field,'v':$(this).val()},function(data){
                                    if(data=="1"){
                                        layer.msg('保存成功!');
                                    }
                                    else{
                                        layer.msg(data);
                                    }
                                })
                            }
                        });
                    })

                } 
            });  
                    
        })
    </script>

</body>
</html>
`
var adm_tb_fieldedit = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            //$('#pid').val('$!m.parentid');
            if ('{{.m.state}}' == '1') {
                $('#state').attr('checked', 'checked');
            }
            $('#images').val('{{.m.images}}');
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', ''+data, "warning");
                    }else{
						jq.messager.alert('错误', "操作失败!", "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:90%;overflow-x:hidden;" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/fieldeditpost" method="post">
            <table cellpadding="5">
                
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" name="field_name" value="{{.m.field_name}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
                <tr>
                    <td>字段:</td>
                    <td><input class="easyui-textbox" type="text" name="field_code" value="{{.m.field_code}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>类型:</td>
                    <td>
						<select id="form_type" name="form_type" style="width:100px;"  class="easyui-combobox" editable="false">
								<option value="标签框">标签框</option>
								<option value="文本框">文本框</option>
								<option value="下拉框">下拉框</option>
								<option value="单选框">单选框</option>
								<option value="复选框">复选框</option>
								<option value="文本域">文本域</option>
                                <option value="编辑框">编辑框</option>
                                <option value="代码框">代码框</option>
                                <option value="日期选择">日期选择</option>
								<option value="开关按钮">开关按钮</option>
								<option value="图片上传">图片上传</option>
								<option value="文件上传">文件上传</option>
						</select>
					</td>
                </tr>
				<tr>
                    <td>长度:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="form_length" id="form_length" value="{{.m.form_length}}" />

                    </td>
                </tr>
				<tr>
                    <td>绑定:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="form_value" id="form_value" value="{{.m.form_value}}"/>
                        <div>sql语句需要两个字段: id val</div>
                    </td>
                </tr>

				<tr>
                    <td>默认:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="field_defval" id="field_defval" value="{{.m.field_defval}}"/>

                    </td>
                </tr>
				<tr>
                    <td>顺序:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="form_sort" id="form_sort" value="{{.m.form_sort}}"/>
                    </td>
                </tr>
				<tr>
                    <td>列表显示:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="view_list" id="view_list" />

                    </td>
                </tr>
                <tr>
                    <td>列表色值:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="view_list_color" id="view_list_color" value="{{.m.view_list_color}}"/>
                    </td>
                </tr>
                <tr>
                    <td>列表编辑:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="view_edit" id="view_edit" />

                    </td>
                </tr>
				<tr>
                    <td>表单显示:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="view_form" id="view_form" />

                    </td>
                </tr>
                <tr>
                    <td>详情显示:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="view_detail" id="view_detail" />

                    </td>
                </tr>
				<tr>
                    <td>搜索字段:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_search" id="is_search" />
                    </td>
                </tr>
                <tr>
                        <td>合计字段:</td>
                        <td>
                            <input class="easyui-checkbox" type="checkbox" name="is_total" id="is_total" />
                        </td>
                    </tr>
                <tr>
                    <td>必搜字段:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="search_require" id="search_require" />
                    </td>
                </tr>
				<tr>
                    <td>排序字段:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_sort" id="is_sort" />
                    </td>
                </tr>
                <tr>
                    <td>是否导出:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_export" id="is_export" />
                    </td>
                </tr>
                <tr>
                    <td>是否导入:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_import" id="is_import" />
                    </td>
                </tr>
                <tr>
                        <td>导入唯一:</td>
                        <td>
                            <input class="easyui-checkbox" type="checkbox" name="is_import_unique" id="is_import_unique" />
                        </td>
                    </tr>
                <tr>
                        <td>树形导航:</td>
                        <td>
                            <input class="easyui-checkbox" type="checkbox" name="is_navtree" id="is_navtree" />
                        </td>
                </tr>
                <tr>
                        <td>树形SQL:</td>
                        <td>
                            <input class="easyui-textbox" type="text" name="navtree_sql" id="navtree_sql" value="{{.m.navtree_sql}}"/>
                        </td>
				</tr>
				<tr>
                    <td>提示:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="form_tip" value="{{.m.form_tip}}"></input>
                    </td>
                </tr>
                <tr>
                    <td>备注:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="memo" value="{{.m.memo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>
                <tr>
                    <td>启用:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="state" id="state" />

                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
<script type="text/javascript">
<!--
	if('{{.m.view_list}}'=='1') $('#view_list').attr('checked',true);
	if('{{.m.view_form}}'=='1') $('#view_form').attr('checked',true);
    if('{{.m.view_edit}}'=='1') $('#view_edit').attr('checked',true);
    if('{{.m.view_detail}}'=='1') $('#view_detail').attr('checked',true);
	if('{{.m.is_sort}}'=='1') $('#is_sort').attr('checked',true);
	if('{{.m.is_search}}'=='1') $('#is_search').attr('checked',true);
    if('{{.m.is_total}}'=='1') $('#is_total').attr('checked',true);
    if('{{.m.is_navtree}}'=='1') $('#is_navtree').attr('checked',true);
    if('{{.m.search_require}}'=='1') $('#search_require').attr('checked',true);
	if ('{{.m.state}}' == '1') $('#state').attr('checked', true);
	if ('{{.m.is_import}}' == '1') $('#is_import').attr('checked', true);
    if ('{{.m.is_import_unique}}' == '1') $('#is_import_unique').attr('checked', true);
	if ('{{.m.is_export}}' == '1') $('#is_export').attr('checked', true);
	$('#form_type').val("{{.m.form_type}}");
//-->
</script>

`
var adm_tb_edit = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            //$('#pid').val('$!m.parentid');
            if ('{{.m.state}}' == '1') {
                $('#state').attr('checked', 'checked');
            }
            $('#images').val('{{.m.images}}');
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', ''+data, "warning");
                    }else{
						jq.messager.alert('错误', "操作失败!", "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>


<div class="easyui-panel"  style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/editpost" method="post">
            <table cellpadding="5">
                
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" name="title" value="{{.m.title}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
                <tr>
                    <td>项目:</td>
                    <td>
                        <select id="proj_id" name="proj_id" style="width:100px;" class="easyui-combobox" editable="false">
                            {{range $k,$v :=.projlist}}
                            <option  value="{{$v.id}}">{{$v.proj_name}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                
                            });
                            
                        </script>
                    </td>
                </tr>
                <tr>
                    <td>代号:</td>
                    <td><input class="easyui-textbox" type="text" name="code" value="{{.code}}" data-options="required:true"></input></td>
                </tr>
                <tr>
                    <td>数据库:</td>
                    <td>
                        <select id="conn" name="conn" style="width:100px;" class="easyui-combobox" editable="false">
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                
                            });
                            
                        </script>
                    </td>
                </tr>
				<tr>
                    <td>表名:</td>
                    <td><input class="easyui-textbox" type="text" name="table" value="{{.m.table}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>条件:</td>
                    <td><input class="easyui-textbox" type="text" name="where_str" value="{{.m.where_str}}" ></input></td>
                </tr>
				<tr>
                    <td>类型:</td>
                    <td>
						<select id="data_type" name="data_type" style="width:100px;"  class="easyui-combobox" editable="false">
								<option value="0">表</option>
								<option value="1">视图</option>
                            <option value="2">SQL</option>
						</select>
					</td>
                </tr>
				<tr>
                    <td>主键:</td>
                    <td><input class="easyui-textbox" type="text" placeholder='视图时填写' name="pri_key" value="{{.m.pri_key}}" data-options="required:true"></input></td>
                </tr>
                <tr>
                    <td>排序字段:</td>
                    <td><input class="easyui-textbox" type="text" placeholder='默认排序字段' name="sort_key" value="{{.m.sort_key}}" data-options="required:true"></input></td>
                </tr>
                <tr>
                    <td>新建:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_create" id="is_create" />
                    </td>
                </tr>
				<tr>
                    <td>编辑:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_edit" id="is_edit" />
                    </td>
                </tr>
                <tr>
                    <td>详情:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_detail" id="is_detail" />
                    </td>
                </tr>
				<tr>
                    <td>删除:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="is_del" id="is_del" />

                    </td>
                </tr>
                <tr>
                        <td>合计:</td>
                        <td>
                            <input class="easyui-checkbox" type="checkbox" name="is_total" id="is_total" />
    
                        </td>
                    </tr>
				<tr>
                    <td>展示:</td>
                    <td>
						<select id="edit_style" name="edit_style" style="width:100px;"  class="easyui-combobox" editable="false">
                                <option value="1">弹窗</option>
                                <option value="2">标签</option>
						</select>
					</td>
                </tr>
				<tr>
                    <td>宽度:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="edit_width" id="edit_width" value="{{.edit_width}}" />
                    </td>
                </tr>
				<tr>
                    <td>高度:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="edit_height" id="edit_height" value="{{.edit_height}}" />
                    </td>
                </tr>
				<tr>
                    <td>导出:</td>
                    <td>
                        <input class="easyui-switchbutton" id="is_export" name="is_export">
                    </td>
                </tr>
                <tr>
                    <td>导入:</td>
                    <td>
                        <input class="easyui-switchbutton" id="is_import" name="is_import">
                    </td>
                </tr>
                
				<tr style="display:none;">
                    <td>操作扩展:</td>
                    <td>
                        
						<textarea  class="easyui-textbox" name="ex_html_operate"  multiline="true" id="ex_html_operate" style="width:100%;height:120px">{{.m.ex_html_operate}}</textarea>
                    </td>
                </tr>
				<tr>
                    <td>JS扩展:</td>
                    <td>
						<textarea  class="easyui-textbox" name="ex_javascript"  multiline="true" id="ex_javascript" style="width:100%;height:120px">{{.m.ex_javascript}}</textarea>
                    </td>
				</tr>
				<tr>
                    <td>编辑SQL扩展:</td>
                    <td>
						<textarea  class="easyui-textbox" name="ex_sql_edit"  multiline="true" id="ex_sql_edit" style="width:100%;height:120px">{{.m.ex_sql_edit}}</textarea>
                    </td>
				</tr>
				<tr>
                    <td>删除SQL扩展:</td>
                    <td>
						<textarea  class="easyui-textbox" name="ex_sql_del"  multiline="true" id="ex_sql_del" style="width:100%;height:120px">{{.m.ex_sql_del}}</textarea>
                    </td>
                </tr>
				<tr style="display:none;">
                    <td>按钮扩展:</td>
                    <td>
                        
						<textarea  class="easyui-textbox" name="ex_linkbutton"  multiline="true" id="ex_linkbutton" style="width:100%;height:120px">{{.m.ex_linkbutton}}</textarea>
                    </td>
                </tr>
                <tr>
                    <td>备注:</td>
                    <td>
                        <textarea  class="easyui-textbox" name="memo"  multiline="true" id="memo" style="width:100%;height:120px">{{.m.memo}}</textarea>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>
                <tr>
                    <td>启用:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="state" id="state" />

                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
<script type="text/javascript">
<!--
    if('{{.m.is_total}}'=='1') $('#is_total').attr('checked',true);
    if('{{.m.is_create}}'=='1') $('#is_create').attr('checked',true);
	if('{{.m.is_edit}}'=='1') $('#is_edit').attr('checked',true);
	if('{{.m.is_detail}}'=='1') $('#is_detail').attr('checked',true);
    if('{{.m.is_del}}'=='1') $('#is_del').attr('checked',true);
	if('{{.m.state}}'=='1') $('#state').attr('checked',true);
	if('{{.m.is_import}}'=='1') $('#is_import').attr('checked',true);
	if('{{.m.is_export}}'=='1') $('#is_export').attr('checked',true);
	$('#conn').combobox({
	    onLoadSuccess : function(data) {
	        $('#conn').combobox('setValue', "{{.m.conn_str}}");
	    }
	}); 
	$('#edit_style').combobox({
	    onLoadSuccess : function(data) {
	        $('#edit_style').combobox('setValue', "{{.m.edit_style}}");
	    }
	}); 
        $('#proj_id').combobox({
            onLoadSuccess: function (data) {
                $('#proj_id').combobox('setValue', "{{.m.proj_id}}");
            }
        }); 

	$(function(){
        $('#is_import').switchbutton({
            checked: eval("{{.is_import}}"),
            onChange: function(checked){
            }
        })
        $('#is_export').switchbutton({
            checked: eval("{{.is_export}}"),
            onChange: function(checked){
            }
        })
    })
//-->
</script>

`
var adm_tb_connedit = `

<script src="/js/jquery.form.js"></script>
<script src="/js/my97datepicker/wdatepicker.js" type="text/javascript"></script>

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/connpost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">链接:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="conn" 
					id="conn" value='{{.m.conn}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">类型:</td>
                    <td>
					
						<div id="divdbtype">
						<input type="radio" id="dbtype0" onchange="sqlsel(0);" title="类型" name="dbtype" value="mysql" style="vertical-align:middle;" >
						<label for="dbtype0">mysql</label>
						<input type="radio" id="dbtype1" onchange="sqlsel(0);" title="类型" name="dbtype" value="mssql" style="vertical-align:middle;" >
						<label for="dbtype1">mssql</label>
						<input type="radio" id="dbtype2" onchange="sqlsel(1);" title="类型" name="dbtype" value="sqlite" style="vertical-align:middle;" >
						<label for="dbtype2">sqlite</label>
						<input type="radio" id="dbtype3" onchange="sqlsel(0);" title="类型" name="dbtype" value="mssql2k" style="vertical-align:middle;" >
						<label for="dbtype3">mssql2k</label>
						</div>
					<script type="text/javascript">
						function sqlsel(i){
							if(i==0){
								$('#trserver').show();
								$('#trport').show();
								$('#truid').show();
								$('#trpwd').show();
							}else{
								$('#trserver').hide();
								$('#trport').hide();
								$('#truid').hide();
								$('#trpwd').hide();
							}
						}
						$("input[name='dbtype'][value='{{.m.dbtype}}']").attr("checked",true); 
						if('{{.m.dbtype}}'=='sqlite'){
							$('#trserver').hide();
							$('#trport').hide();
							$('#truid').hide();
							$('#trpwd').hide();
						}
					</script>
					 
					
					</td>
                </tr>
                
                <tr id="trserver">
                    <td style="width:55px;">服务器:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="server" 
					id="server" value='{{.m.server}}' style="width:135px;" ></input>
					 
					
					</td>
                </tr>
                
                <tr id="trport">
                    <td style="width:55px;">端口:</td>
                    <td>
					
					
					<select id="port" name="port" title="" class="easyui-combobox" editable="true" style="width:130px;">
						<option value="">请选择...</option>
						<option value="3306">3306</option>
<option value="1433">1433</option>
					</select>
					<script type="text/javascript">
						$('#port').combobox({
							onLoadSuccess: function () {
								$('#port').combobox('select','{{.m.port}}');
							}
						});
						
					</script>
					 
					
					</td>
                </tr>
                
                <tr id="truid">
                    <td style="width:55px;">用户名:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="uid" 
					id="uid" value='{{.m.uid}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr id="trpwd">
                    <td style="width:55px;">密码:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="pwd" 
					id="pwd" value='{{.m.pwd}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="dbname" 
					id="dbname" value='{{.m.dbname}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">说明:</td>
                    <td>
					
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="memo" id="memo" value='{{.m.memo}}'></input>
					
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`
var adm_tb_conn = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/connedit?id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	w.window('open');
	w.window('refresh', '/adm/tb/connedit?id=');
}
function doData(title,url){
	top.addTab(title,url);
}

function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/conndel', { id: row.id }, function (result) {
					if(result=="-1"){
						jq.messager.alert('警告','连接正在使用!','warning');
					}else if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        jq.messager.alert('警告','删除失败!','warning');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })
    function rowformater_proj(value, row, index) {
		return '<div class="proj_id">'+value+'</div>';
    }
	function rowformater_field(value, row, index) {
		var a= "<a href='#' onclick='doData(\""+row.title+"-数据表\",\"/adm/tb/tbs?id="+row.id+"\");' >表管理</a>";
		//var b="&nbsp;<a href='/adm/mdata/list?code="+row.code+"' >数据</a>";
		return a;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/connjson" data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="链接管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
                <th field="id" width="5">ID</th>
                <th field="title" width="10">名称</th>
                <th field="conn" align="center" width="10">代号</th>
                <th field="dbtype" align="center" width="10" >类型</th>
                <th field="server" width="5">服务器</th>
				<th field="port" width="5">端口</th>
				<th field="uid" width="5">用户</th>
                <th field="dbname" width="10"  >数据库</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true"  onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true"  onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_btns = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
    <script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }

function doField(id){
        $('#win').window('open');
        $('#win').window('refresh', '/adm/tb/btnedit?id='+id);
}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
	var w=$('#win').window({
		width:520,
		height:480,
		modal:true,
		title:'{{.tb.title}}'+'[添加按钮]'
	});
    w.window('open');
    w.window('refresh', '/adm/tb/btnedit?tbid={{.id}}');
}
function doDel(id){
    jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/btndel', {id:'{{.m.id}}', id: id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        jq.messager.alert('错误','删除失败!','info');
                    }
                });
            }
        });

}
    $(function(){
		
    })
	function rowformater_field(value, row, index) {
		var e = '<a   href="#" onclick="doField(' + row.id + ');" >编辑</a> ';
        var d = '<a   href="#" onclick="doDel(' + row.id + ');" >删除</a> ';
        return e+d;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;"  fit="true">

    <table class="easyui-datagrid"  
           url="/adm/tb/btnsjson?id={{.id}}"  
           title="{{.tb.title}}-按钮设置" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
                <th field="id" width="5">ID</th>
                <th field="title"  width="10" >名称</th>
                <th field="type" align="center"  width="10" >类型</th>
                <th field="style" align="center"  width="10" >样式</th>
                <th field="icon"  width="10" >图标</th>
                <th field="url"  width="10" >链接</th> 
                <th field="is_blank"  width="10"  >新页</th> 
				<th field="state"  width="20" >状态</th>
				<th field=" "  data-options="formatter:rowformater_field"  width="20" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:480px;height:390px;padding:5px;">
        Some Content.
    </div>
    <script language="javascript">
        $(function () {
            $('#tt').datagrid({   
                onLoadSuccess:function(data){

                    $(".canedit").bind("dblclick", function () {
                        var field=$(this).attr('sp');
                        var valid=$(this).attr('valid');
                        var oldval=$(this).attr('val');

                        var input = "<input type='text' id='temp' value=" + $(this).text() + " >";
                        $(this).text("");
                        $(this).append(input);
                        $("input#temp").focus();
                        $("input").blur(function () {
                            if ($(this).val() == "") {
                                $(this).closest("div").text(oldval);
                                $(this).remove();
                            } else {
                                $(this).closest("div").text($(this).val());
                                if($(this).val()==oldval) return;
                                //alert($(this).val()+"-"+oldval);

                                jQuery.post('/adm/tb/fieldset',{'tb':'tb_field','id':valid,'f':field,'v':$(this).val()},function(data){
                                    if(data=="1"){
                                        layer.msg('保存成功!');
                                    }
                                    else{
                                        layer.msg(data);
                                    }
                                })
                            }
                        });
                    })

                } 
            });  
                    
        })
    </script>

</body>
</html>
`
var adm_tb_btnedit = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            if ('{{.m.state}}' == '1') {
                $('#state').attr('checked', 'checked');
            }
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', ''+data, "warning");
                    }else{
						jq.messager.alert('错误', "操作失败!", "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:90%;overflow-x:hidden;" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/btneditpost" method="post">
            <table cellpadding="5">
                
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" name="title" value="{{.m.title}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
                <tr>
                    <td>类型:</td>
                    <td>
                        <select id="type" name="type" style="width:135px;"  class="easyui-combobox" editable="false">
                            <option value="工具栏">工具栏</option>
                            <option value="搜索栏">搜索栏</option>
                            <option value="操作栏">操作栏</option>
                            <option value="连接页">连接页</option>
                    </select>
                    <script type="text/javascript">
                        $('#type').combobox({
							onLoadSuccess: function () {
								$('#type').combobox('setValue','{{.m.type}}');
							}
						});
					</script>
                    </td>
                </tr>
				<tr>
                    <td>样式:</td>
                    <td>
						<select id="style" name="style" style="width:135px;"  class="easyui-combobox" editable="false">
								<option value="按钮">按钮</option>
								<option value="超链接">超链接</option>
                        </select>
                        <script type="text/javascript">
                            $('#style').combobox({
                                onLoadSuccess: function () {
                                    $('#style').combobox('setValue','{{.m.style}}');
                                }
                            });
                        </script>
					</td>
                </tr>
				<tr>
                    <td>图标:</td>
                    <td>
                        <select id="icon" name="icon" style="width:135px;"  class="easyui-combobox" editable="false">
                            <option value="icon-0">icon-0</option>
                            <option value="icon-1">icon-1</option>
                            <option value="icon-2">icon-2</option>
                            <option value="icon-3">icon-3</option>
                            <option value="icon-4">icon-4</option>
                            <option value="icon-5">icon-5</option>
                            <option value="icon-6">icon-6</option>
                            <option value="icon-7">icon-7</option>
                            <option value="icon-8">icon-8</option>
                            <option value="icon-9">icon-9</option>
                            <option value="icon-10">icon-10</option>
                            <option value="icon-11">icon-11</option>
                            <option value="icon-12">icon-12</option>
                            <option value="icon-13">icon-13</option>
                            <option value="icon-14">icon-14</option>
                            <option value="icon-15">icon-15</option>
                            <option value="icon-16">icon-16</option>
                            <option value="icon-17">icon-17</option>
                            <option value="icon-18">icon-18</option>
                            <option value="icon-19">icon-19</option>
                            <option value="icon-20">icon-20</option>
                            <option value="icon-21">icon-21</option>
                            <option value="icon-22">icon-22</option>
                            <option value="icon-23">icon-23</option>
                            <option value="icon-24">icon-24</option>
                            <option value="icon-25">icon-25</option>
                            <option value="icon-26">icon-26</option>
                            <option value="icon-27">icon-27</option>
                            <option value="icon-28">icon-28</option>
                            <option value="icon-29">icon-29</option>
                            <option value="icon-30">icon-30</option>
                            <option value="icon-31">icon-31</option>
                            <option value="icon-32">icon-32</option>
                            <option value="icon-33">icon-33</option>
                            <option value="icon-34">icon-34</option>
                            <option value="icon-35">icon-35</option>
                            <option value="icon-36">icon-36</option>
                            <option value="icon-37">icon-37</option>
                            <option value="icon-38">icon-38</option>
                            <option value="icon-39">icon-39</option>
                        </select>
                        <script type="text/javascript">
                            $('#icon').combobox({
                                formatter: function (row) {
                                    var imageFile = '/js/easyui/themes/icons/'+row.value.replace('icon-','')+'.png';
                                    return '<img class="item-img" style="height:25px;width:25px;" src="' + imageFile + '"/><span class="item-text">' + row.text + '</span>';
                                },
                                onLoadSuccess: function () {
                                    $('#icon').combobox('setValue','{{.m.icon}}');
                                }
                            });
                        </script>
                    </td>
                </tr>
				<tr>
                    <td>链接:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="url" id="url" value="{{.m.url}}"/>

                    </td>
                </tr>
				<tr>
                    <td>新页:</td>
                    <td>
                        <select id="is_blank" name="is_blank" style="width:135px;"  class="easyui-combobox" editable="false">
                            <option value="0">标签页</option>
                            <option value="1">新页面</option>
                            <option value="2">弹出框</option>
                        </select>
                        <script type="text/javascript">
                            $('#is_blank').combobox({
                                onLoadSuccess: function () {
                                    $('#is_blank').combobox('setValue','{{.m.is_blank}}');
                                }
                            });
                        </script>
                    </td>
                </tr>
				
                <tr>
                    <td>JS脚本:</td>
                    <td>
                        <input class="easyui-textbox" name="jstr" data-options="multiline:true" value="{{.m.jstr}}" style="width:200px;height:100px">
                        <input type="hidden" id="id" name="tbid" value="{{.tbid}}" />
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>
                <tr>
                    <td>启用:</td>
                    <td>
                        <input class="easyui-checkbox" type="checkbox" name="state" id="state" />

                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
<script type="text/javascript">
	if ('{{.m.state}}' == '1') $('#state').attr('checked', true); 
	$('#is_blank').val("{{.m.is_blank}}");
</script>

`
var adm_tb_tbs = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
    <script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    <style>
        .entry {
            position: relative;
            margin-left: 2px;
            margin-top: 2px;
            margin-bottom: 5px;
            width: auto;
            background: #FFFFDD;
            padding: 10px;
            padding-left: 0px;
            /*设置圆角*/
            -webkit-border-radius: 5px;
            -moz-border-radius: 5px;
            border-radius: 5px;
        }
    </style>

    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }

function doField(id){
        $('#win').window('open');
        $('#win').window('refresh', '/adm/tb/btnedit?id='+id);
}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
	var w=$('#win').window({
		width:380,
		height:280,
		modal:true,
		title:'{{.tb.title}}'+'[添加表]'
	});
    w.window('open');
    w.window('refresh', '/adm/tb/tbsadd?id={{.id}}');
}
function doData(title,url){
	top.addTab(title,url);
}
    $(function(){
		
    })
	function rowformater_field(value, row, index) {
		var a= "<a href='#' onclick='doData(\""+row.name+"-字段\",\"/adm/tb/fs?id={{.id}}&tname="+row.name+"\");' >字段管理</a>";
        return a;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;"  fit="true">

    <table class="easyui-datagrid"  
           url="/adm/tb/tbsjson?id={{.id}}"  
           title="{{.m.title}}-数据表管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
                <th field="rootpage" width="5">ID</th>
                <th field="name"  width="10" >名称</th>
                <th field="type" align="center"  width="10" >类型</th>
				<th field=" "  data-options="formatter:rowformater_field"  width="20" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
		<div class="entry">
            <div style=""><font color='red'>注意:</font>表名不支持修改功能,请谨慎操作.</div>
        </div>
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:280px;height:190px;padding:5px;">
        Some Content.
    </div>
    <script language="javascript">
        $(function () {
            $('#tt').datagrid({   
                onLoadSuccess:function(data){
					
                } 
            });  
                    
        })
    </script>

</body>
</html>
`
var adm_tb_tbsadd = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data != "0") {
                        layer.msg('<font color="yellow">操作成功!</font>');
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        layer.msg('<font color="green">操作失败!</font>');
                    }
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }
        
</script>


<div class="easyui-panel" title="" style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/tbsaddpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="name" value="{{.m.name}}" data-options="required:true,missingMessage:'必填字段'"></input>
					<input type="hidden" id="id" name="id" value="{{.m.id}}" />
					</td>
                </tr>
				
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`
var adm_tb_fs = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
    <script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    <style>
        .entry {
            position: relative;
            margin-left: 2px;
            margin-top: 2px;
            margin-bottom: 5px;
            width: auto;
            background: #FFFFDD;
            padding: 10px;
            padding-left: 0px;
            /*设置圆角*/
            -webkit-border-radius: 5px;
            -moz-border-radius: 5px;
            border-radius: 5px;
        }
    </style>

    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }

function doField(id){
        $('#win').window('open');
        $('#win').window('refresh', '/adm/tb/btnedit?id='+id);
}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
	var w=$('#win').window({
		width:420,
		height:260,
		modal:true,
		title:'{{.tb.title}}'+'[添加字段]'
	});
    w.window('open');
    w.window('refresh', '/adm/tb/fsadd?id={{.id}}&tname={{.tname}}');
}

    $(function(){
		
    })
	function rowformater_field(value, row, index) {
		var e = '';//'<a   href="#" onclick="doField(' + row.id + ');" >编辑</a> ';
        return e;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;"  fit="true">

    <table class="easyui-datagrid"  
           url="/adm/tb/fsjson?id={{.id}}&tname={{.tname}}"  
           title="{{.m.title}}-字段管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
                <th field="cid" width="5">ID</th>
                <th field="name"  width="10" >名称</th>
                <th field="type" align="center"  width="10" >类型</th>
				<th field="notnull" align="center"  width="10" >非空</th>
				<th field="dflt_value" align="center"  width="10" >默认值</th>
				 <th field="pk" align="center"  width="10" >主键</th>
				<th field=" "  data-options="formatter:rowformater_field"  width="20" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
		<div class="entry">
            <div style=""><font color='red'>注意:</font>字段不支持修改功能,请谨慎操作.</div>
        </div>
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add"  plain="true" onclick="doAdd();">新建</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:280px;height:190px;padding:5px;">
        Some Content.
    </div>
    <script language="javascript">
        $(function () {
            $('#tt').datagrid({   
                onLoadSuccess:function(data){
					
                } 
            });  
                    
        })
    </script>

</body>
</html>
`
var adm_tb_fsadd = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data != "0") {
                        layer.msg('<font color="yellow">操作成功!</font>');
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        layer.msg('<font color="green">操作失败!</font>');
                    }
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }
        
</script>


<div class="easyui-panel" title="" style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/fsaddpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="name" value="{{.m.name}}" data-options="required:true,missingMessage:'必填字段'"></input>
					<input type="hidden" id="id" name="id" value="{{.m.id}}" />
					<input type="hidden" id="tname" name="tname" value="{{.tname}}" /> 
					</td>
                </tr>
				 <tr>
                    <td>类型:</td>
                    <td>
					{{if eq .m.dbtype "sqlite"}}
						<select id="state" class="easyui-combobox" name="ftype" editable="false" style="width:180px;" >
                            <option value="INTEGER">INTEGER</option>
                            <option value="TEXT">TEXT</option> 
                        </select>
					{{end}}
					{{if eq .m.dbtype "mysql"}}
						<select id="state" class="easyui-combobox" name="ftype" editable="false" style="width:180px;" >
                            <option value="int">int</option>
							<option value="varchar">varchar</option>
							<option value="decimal">decimal</option>  
                        </select>
					{{end}}
					{{if eq .m.dbtype "mssql"}}
						<select id="state" class="easyui-combobox" name="ftype" editable="false" style="width:180px;" >
                            <option value="int">int</option>
							<option value="varchar">varchar</option>
							<option value="decimal">decimal</option>  
                        </select>
					{{end}}
					{{if eq .m.dbtype "mssql2k"}}
						<select id="state" class="easyui-combobox" name="ftype" editable="false" style="width:180px;" >
                            <option value="int">int</option>
							<option value="varchar">varchar</option>
							<option value="decimal">decimal</option>  
                        </select>
					{{end}}
					</td>
                </tr>
				<tr>
                    <td>长度:</td>
                    <td>
						<input class="easyui-textbox" type="text" style="width:180px;" name="length" value="(50)" ></input>
					</td>
                </tr>
				<tr>
                    <td>默认值:</td>
                    <td>
						<input class="easyui-textbox" type="text" style="width:180px;" name="defval" value="" ></input>
					</td>
                </tr>
				
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`

var adm_tb_page = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>页面模板</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			//w.window('open');
			//w.window('refresh', '/adm/tb/pageedit?id='+row.id);
			doTab('编辑-'+row.title,'/adm/tb/pageedit?id='+row.id);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	//w.window('open');
	//w.window('refresh', '/adm/tb/pageedit?id=');
	doTab('编辑-新增纪录','/adm/tb/pageedit?id=');
}
function doParam(){
	var row = $('#tt').datagrid('getSelected');
	if (row){
		doData(row.title+'-参数设置','/adm/tb/pageparam?pageid='+row.id);
	}else{
		jq.messager.alert('警告','请选择一行数据','warning');
	}
}
function doData(title,url){
	top.addTab(title,url);
}
function doTab(title,url){
	top.addTab(title,url);
}
function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/pagedel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })
    function rowformater_proj(value, row, index) {
		return '<div class="proj_id">'+value+'</div>';
    }
	function rowformater_field(value, row, index) {
		var a= "<a href='/p/"+row.module+"/"+row.code+"' target='_blank' >预览</a>";
		return a;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/pagejson" data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="页面管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="module" align="center" width="10" >模块</th>
                <th field="code" align="center" width="10">代号</th>
                <th field="title" width="10">名称</th>
				<th field="state" width="5">状态</th>
				<th field="memo" width="10">说明</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-41" plain="true" onclick="doParam();">参数</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_pageedit = `

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/pagepost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">模块:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="module" 
					id="module" value='{{.m.module}}' ></input>
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">代号:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="code" 
					id="code" value='{{.m.code}}' ></input>
					 
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<select id="conn_str" name="conn_str" style="width:142px;" class="easyui-combobox" editable="false">
					<option  value="">请选择...</option>
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess : function(data) {
										$('#conn_str').combobox('setValue', "{{.m.conn_str}}");
									}
								}); 
                            });
                            
                        </script>
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">模板:</td>
                    <td>
					
					<textarea class="easyui-textbox" multiline="true" style="width:280px;height:120px" id="template" name="template">{{.m.template}}</textarea>
					<style type="text/css">
					.CodeMirror {border: 1px solid #ddd; font-size:13px}
					</style>
					<script type="text/javascript">

					</script>	

					</td>
                </tr>		
                <tr>
                    <td style="width:55px;">说明:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="description" id="description" value='{{.m.description}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">token:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="token" 
					id="token" value='{{.m.token}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`
var adm_tb_pageedit_tab = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>页面模板</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
</head>
<body style="padding:2px;margin-bottom:2px;">
<link rel="stylesheet" href="/js/codemirror-5.31.0/lib/codemirror.css"/>
<script src="/js/codemirror-5.31.0/lib/codemirror.js"></script>
<script src="/js/codemirror-5.31.0/clike.js"></script>
<script src="/js/codemirror-5.31.0/mode/xml/xml.js"></script>
<link rel="stylesheet" href="/js/codemirror-5.31.0/theme/idea.css"/>
<link rel="stylesheet" href="/js/codemirror-5.31.0/addon/fold/foldgutter.css"/>
<script src="/js/codemirror-5.31.0/addon/fold/foldcode.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/foldgutter.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/brace-fold.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/comment-fold.js"></script>
<script src="/js/codemirror-5.31.0/addon/edit/matchbrackets.js"></script>
<script src="/js/codemirror-5.31.0/mode/javascript/javascript.js"></script>
<script src="/js/codemirror-5.31.0/addon/selection/active-line.js"></script> 

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info",function(){
							//window.top.closeTabById(window.frameElement.parentElement.getAttribute('id'));
						});
                        //$('#tt').datagrid('reload');
                        //$('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
			//$('#win').window('close');
			window.top.closeTabById(window.frameElement.parentElement.getAttribute('id'));
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/pagepost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">模块:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="module" 
					id="module" value='{{.m.module}}' ></input>
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">代号:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="code" 
					id="code" value='{{.m.code}}' ></input>
					 
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<select id="conn_str" name="conn_str" style="width:142px;" class="easyui-combobox" editable="false">
					<option  value="">请选择...</option>
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess : function(data) {
										$('#conn_str').combobox('setValue', "{{.m.conn_str}}");
									}
								}); 
                            });
                            
                        </script>
					
					</td>
				</tr>		

				<tr>
                    <td style="width:55px;">模板:</td>
                    <td>
					
					<textarea class="CodeMirror" multiline="true" style="width:280px;height:120px" id="template" name="template"></textarea>
					<style type="text/css">
					.CodeMirror {border: 1px solid #ddd; font-size:13px}
					</style>
					<script type="text/javascript">

					var editor = CodeMirror.fromTextArea(document.getElementById("template"), {
						mode: "text/xml",    
						
						lineNumbers: true,	
						theme: "idea",	
						htmlMode:true,
						lineWrapping: false,	
						foldGutter: true,
						gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
						matchBrackets: true,	
						
					});
					editor.setSize('680px', '520px');     
					editor.setValue('{{.m.template}}');
					</script>	

					</td>
				</tr>
				<tr>
					<td style="width:55px;">说明:</td>
					<td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="description" id="description" value='{{.m.description}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">token:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="token" 
					id="token" value='{{.m.token}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:left;margin-left:200px;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
</body>
</html>
`
var adm_tb_api = ` 
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>页面模板</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			//w.window('open');
			//w.window('refresh', '/adm/tb/apiedit?id='+row.id);
			top.addTab('编辑-'+row.title,'/adm/tb/apiedit?id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	//w.window('open');
	//w.window('refresh', '/adm/tb/apiedit?id=');
	top.addTab('新增接口','/adm/tb/apiedit?id=');
}
function doParam(){
	var row = $('#tt').datagrid('getSelected');
	if (row){
		doData(row.title+'-参数设置','/adm/tb/apiparam?apiid='+row.id);
	}else{
		jq.messager.alert('警告','请选择一行数据','warning');
	}
}
function doData(title,url){
	top.addTab(title,url);
}

function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/apidel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })
    function rowformater_proj(value, row, index) {
		return '<div class="proj_id">'+value+'</div>';
    }
	function rowformater_field(value, row, index) {
		var a= "<a href='/api/"+row.module+"/"+row.api_code+"' target='_blank' >预览</a>";
		return a;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/apijson" data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="页面管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="module" align="center" width="10" >模块</th>
                <th field="api_code" align="center" width="10">代号</th>
                <th field="title" width="10">名称</th>                
                <th field="description" width="25">说明</th>
				<th field="state" width="5">状态</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-41" plain="true" onclick="doParam();">参数</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_apiedit = `


<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/apipost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">模块:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="module" 
					id="module" value='{{.m.module}}' ></input>
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">代号:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="code" 
					id="code" value='{{.m.api_code}}' ></input>
					 
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<select id="conn_str" name="conn_str" style="width:142px;" class="easyui-combobox" editable="false">
					<option  value="">请选择...</option>
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess : function(data) {
										$('#conn_str').combobox('setValue', "{{.m.conn_str}}");
									}
								}); 
                            });
                            
                        </script>
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">模板:</td>
                    <td>
					
					<textarea class="easyui-textbox" multiline="true" style="width:280px;height:120px" id="template" name="template">{{.m.api_template}}</textarea>
					<style type="text/css">
					.CodeMirror {border: 1px solid #ddd; font-size:13px}
					</style>
					<script type="text/javascript">

					</script>	

					</td>
                </tr>		
                <tr>
                    <td style="width:55px;">说明:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="description" id="description" value='{{.m.description}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">token:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="token" 
					id="token" value='{{.m.token}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`
var adm_tb_apiedit_tab = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>API模板</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
</head>
<body style="padding:2px;margin-bottom:2px;">

<link rel="stylesheet" href="/js/codemirror-5.31.0/lib/codemirror.css"/>
<script src="/js/codemirror-5.31.0/lib/codemirror.js"></script>
<script src="/js/codemirror-5.31.0/clike.js"></script>
<script src="/js/codemirror-5.31.0/mode/xml/xml.js"></script>
<link rel="stylesheet" href="/js/codemirror-5.31.0/theme/idea.css"/>
<link rel="stylesheet" href="/js/codemirror-5.31.0/addon/fold/foldgutter.css"/>
<script src="/js/codemirror-5.31.0/addon/fold/foldcode.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/foldgutter.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/brace-fold.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/comment-fold.js"></script>
<script src="/js/codemirror-5.31.0/addon/edit/matchbrackets.js"></script>
<script src="/js/codemirror-5.31.0/mode/javascript/javascript.js"></script>
<script src="/js/codemirror-5.31.0/addon/selection/active-line.js"></script> 

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info",function(){
							window.top.closeTabById(window.frameElement.parentElement.getAttribute('id'));
						});
                        //$('#tt').datagrid('reload');
                        //$('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
			//$('#win').window('close');
			window.top.closeTabById(window.frameElement.parentElement.getAttribute('id'));
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/apipost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">模块:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="module" 
					id="module" value='{{.m.module}}' ></input>
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">代号:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="code" 
					id="code" value='{{.m.api_code}}' ></input>
					 
					
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<select id="conn_str" name="conn_str" style="width:142px;" class="easyui-combobox" editable="false">
					<option  value="">请选择...</option>
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess : function(data) {
										$('#conn_str').combobox('setValue', "{{.m.conn_str}}");
									}
								}); 
                            });
                            
                        </script>
					
					</td>
				</tr>				
				
				<tr>
                    <td style="width:55px;">模板:</td>
                    <td>
					
					<textarea class="CodeMirror" multiline="true" style="width:280px;height:120px" id="template" name="template">{{.m.api_template}}</textarea>
					<style type="text/css">
					.CodeMirror {border: 1px solid #ddd; font-size:13px}
					</style>
					<script type="text/javascript">

					var editor = CodeMirror.fromTextArea(document.getElementById("template"), {
						mode: "text/xml",    
						
						lineNumbers: true,	
						theme: "idea",	
						htmlMode:true,
						lineWrapping: false,	
						foldGutter: true,
						gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
						matchBrackets: true,	
						
					});
					editor.setSize('480px', '320px');     
					editor.setValue('{{.m.api_template}}');
					</script>	

					</td>
				</tr>
				
                <tr>
                    <td style="width:55px;">说明:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="description" id="description" value='{{.m.description}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">token:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="token" 
					id="token" value='{{.m.token}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:left;margin-left:200px;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
</body>
</html>
`
var adm_tb_rpt = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>报表模板</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			//w.window('open');
			//w.window('refresh', '/adm/tb/rptedit?id='+row.id);
			top.addTab('编辑'+row.title,'/adm/tb/rptedit?id='+row.id);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	//w.window('open');
	//w.window('refresh', '/adm/tb/rptedit?id=');
	top.addTab('报表-新增纪录','/adm/tb/rptedit?id=');
}
function doParam(){
	var row = $('#tt').datagrid('getSelected');
	if (row){
		doData(row.title+'-参数设置','/adm/tb/rptparam?rptid='+row.id);
	}else{
		jq.messager.alert('警告','请选择一行数据','warning');
	}
}
function doData(title,url){
	top.addTab(title,url);
}
function doPrint(){
	var row = $('#tt').datagrid('getSelected');
	if (row){
		doData(row.title+'-报表设计','/adm/tb/rptdesign?id='+row.id);
	}else{
		jq.messager.alert('警告','请选择一行数据','warning');
	}
}
function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/rptdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })
    function rowformater_proj(value, row, index) {
		return '<div class="proj_id">'+value+'</div>';
    }
	function rowformater_field(value, row, index) {
		var a= "<a href='/rpt/"+row.module+"/"+row.code+"' target='_blank' >预览</a>";
		return a;
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/rptjson" data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="页面管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="module" align="center" width="10" >模块</th>
                <th field="code" align="center" width="10">代号</th>
                <th field="title" width="10">名称</th>                
                <th field="description" width="25">说明</th>
				<th field="state" width="5">状态</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-41" plain="true" onclick="doParam();">参数</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-0" plain="true" onclick="doPrint();">设计</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_rptedit = `

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/rptpost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">模块:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="module" 
					id="module" value='{{.m.module}}' ></input>
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">代号:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="code" 
					id="code" value='{{.m.code}}' ></input>
					 
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<select id="conn_str" name="conn_str" style="width:142px;" class="easyui-combobox" editable="false">
					<option  value="">请选择...</option>
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess : function(data) {
										$('#conn_str').combobox('setValue', "{{.m.conn_str}}");
									}
								}); 
                            });
                            
                        </script>
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">模板:</td>
                    <td>
					
					<textarea class="easyui-textbox" multiline="true" style="width:280px;height:120px" id="template" name="template">{{.m.template}}</textarea>
					<style type="text/css">
					.CodeMirror {border: 1px solid #ddd; font-size:13px}
					</style>
					<script type="text/javascript">

					</script>	

					</td>
                </tr>		
                <tr>
                    <td style="width:55px;">说明:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="description" id="description" value='{{.m.description}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">token:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="token" 
					id="token" value='{{.m.token}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`
var adm_tb_rptedit_tab = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>报表模板</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>

</head>
<body style="padding:2px;margin-bottom:2px;">
<link rel="stylesheet" href="/js/codemirror-5.31.0/lib/codemirror.css"/>
<script src="/js/codemirror-5.31.0/lib/codemirror.js"></script>
<script src="/js/codemirror-5.31.0/clike.js"></script>
<script src="/js/codemirror-5.31.0/mode/xml/xml.js"></script>
<link rel="stylesheet" href="/js/codemirror-5.31.0/theme/idea.css"/>
<link rel="stylesheet" href="/js/codemirror-5.31.0/addon/fold/foldgutter.css"/>
<script src="/js/codemirror-5.31.0/addon/fold/foldcode.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/foldgutter.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/brace-fold.js"></script>
<script src="/js/codemirror-5.31.0/addon/fold/comment-fold.js"></script>
<script src="/js/codemirror-5.31.0/addon/edit/matchbrackets.js"></script>
<script src="/js/codemirror-5.31.0/mode/javascript/javascript.js"></script>
<script src="/js/codemirror-5.31.0/addon/selection/active-line.js"></script> 

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info",function(){
							window.top.closeTabById(window.frameElement.parentElement.getAttribute('id'));
						});
                        //$('#tt').datagrid('reload');
						//$('#win').window('close');
						
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
			//$('#win').window('close');
			window.top.closeTabById(window.frameElement.parentElement.getAttribute('id'));
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/rptpost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">模块:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="module" 
					id="module" value='{{.m.module}}' ></input>
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">名称:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">代号:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="code" 
					id="code" value='{{.m.code}}' ></input>
					 
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">数据库:</td>
                    <td>
					
					<select id="conn_str" name="conn_str" style="width:142px;" class="easyui-combobox" editable="false">
					<option  value="">请选择...</option>
                            {{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess : function(data) {
										$('#conn_str').combobox('setValue', "{{.m.conn_str}}");
									}
								}); 
                            });
                            
                        </script>
					
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">模板:</td>
                    <td>
					
					<textarea class="CodeMirror" multiline="true" style="width:280px;height:120px" id="template" name="template">{{.m.template}}</textarea>
					<style type="text/css">
					.CodeMirror {border: 1px solid #ddd; font-size:13px}
					</style>
					<script type="text/javascript">

					var editor = CodeMirror.fromTextArea(document.getElementById("template"), {
						mode: "text/xml",    
						
						lineNumbers: true,	
						theme: "idea",	
						htmlMode:true,
						lineWrapping: false,	
						foldGutter: true,
						gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
						matchBrackets: true,	
						
					});
					editor.setSize('480px', '320px');     
					editor.setValue('{{.m.template}}');
					</script>	

					</td>
                </tr>		
                <tr>
                    <td style="width:55px;">说明:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="description" id="description" value='{{.m.description}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">token:</td>
                    <td>
					
					<input class="easyui-textbox" title="" type="text" name="token" 
					id="token" value='{{.m.token}}'  style="width:135px;"></input>
					 
					
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
        </form>
        <div style="text-align:left;margin-left:200px;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
</body>
</html>
`
var adm_tb_rptparam = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>报表参数列表</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
		$('#tt').datagrid({url:'/adm/tb/rptparamjson','queryParams':{
			qtxt:$('#qtxt').val(),
			rptid:{{.rptid}}
		}});
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/rptparamedit?rptid={{.rptid}}&id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	w.window('open');
	w.window('refresh', '/adm/tb/rptparamedit?id=&rptid={{.rptid}}');
}
function doData(title,url){
	top.addTab(title,url);
}

function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/rptparamdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}


	function rowformater_field(value, row, index) {
		var a= "";
		return a;
	}
	
    $(function(){
		doSearch();
    })	
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="参数管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="title" align="center" width="10" >标题</th>
                <th field="param_type" align="center" width="10">类型</th>
                <th field="param_name" width="10">参数</th>                
				<th field="max_length" width="10">长度</th>  
				<th field="is_require" width="10">必填</th>  
				<th field="param_value" width="10">参数值</th>  
				<th field="memo" width="10">备注</th>  
				<th field="state" width="5">状态</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_rptparamedit = `


<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/rptparampost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">标题:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">参数:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="param_name" 
					id="param_name" value='{{.m.param_name}}' ></input>
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">类型:</td>
                    <td>
					<select id="param_type" name="param_type" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="参数">参数</option>
							<option  value="变量">变量</option>
							<option  value="cookie">Cookie</option>
							<option  value="session">Session</option>
							<option  value="sql">SQL</option>
							<option  value="sqlexec">SQLExec</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#param_type').combobox({
									onLoadSuccess : function(data) {
										$('#param_type').combobox('setValue', "{{.m.param_type}}");
									}
								});
                            });
                            
						</script>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">长度:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="max_length" 
					id="max_length" value='{{.m.max_length}}' ></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">参数值:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="param_value" id="param_value" value='{{.m.param_value}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">备注:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="memo" 
					id="memo" value='{{.m.memo}}'  style="width:135px;"></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">必填:</td>
                    <td>
					<input class="easyui-switchbutton" id="is_require" title="" name="is_require" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_require}}'=='1'){
								$('#is_require').switchbutton({
									checked: true,
								})
							}else{
								$('#is_require').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.state}}'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
			<input type="hidden" id="rpt_id" name="rpt_id" value="{{.rptid}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`

var adm_tb_pageparam = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>报表参数列表</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
		$('#tt').datagrid({url:'/adm/tb/pageparamjson','queryParams':{
			qtxt:$('#qtxt').val(),
			pageid:{{.pageid}}
		}});
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/pageparamedit?pageid={{.pageid}}&id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	w.window('open');
	w.window('refresh', '/adm/tb/pageparamedit?id=&pageid={{.pageid}}');
}
function doData(title,url){
	top.addTab(title,url);
}

function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/pageparamdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}


	function rowformater_field(value, row, index) {
		var a= "";
		return a;
	}
	
    $(function(){
		doSearch();
    })	
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="参数管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="title" align="center" width="10" >标题</th>
                <th field="param_type" align="center" width="10">类型</th>
                <th field="param_name" width="10">参数</th>                
				<th field="max_length" width="10">长度</th>  
				<th field="is_require" width="10">必填</th>  
				<th field="param_value" width="10">参数值</th>  
				<th field="memo" width="10">备注</th>  
				<th field="state" width="5">状态</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_pageparamedit = `

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/pageparampost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">标题:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">参数:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="param_name" 
					id="param_name" value='{{.m.param_name}}' ></input>
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">类型:</td>
                    <td>
					<select id="param_type" name="param_type" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="参数">参数</option>
							<option  value="变量">变量</option>
							<option  value="cookie">Cookie</option>
							<option  value="session">Session</option>
							<option  value="sql">SQL</option>
							<option  value="sqlexec">SQLExec</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#param_type').combobox({
									onLoadSuccess : function(data) {
										$('#param_type').combobox('setValue', "{{.m.param_type}}");
									}
								});
                            });
                            
						</script>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">长度:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="max_length" 
					id="max_length" value='{{.m.max_length}}' ></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">参数值:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="param_value" id="param_value" value='{{.m.param_value}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">备注:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="memo" 
					id="memo" value='{{.m.memo}}'  style="width:135px;"></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">必填:</td>
                    <td>
					<input class="easyui-switchbutton" id="is_require" title="" name="is_require" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_require}}'=='1'){
								$('#is_require').switchbutton({
									checked: true,
								})
							}else{
								$('#is_require').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.state}}'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
			<input type="hidden" id="page_id" name="page_id" value="{{.pageid}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`

var adm_tb_apiparam = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>报表参数列表</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
	<!--<script type="text/javascript" src="/js/easyui/jquery.easyui.plus.js"></script>-->
	<script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
	
function doSearch(){
		$('#tt').datagrid({url:'/adm/tb/apiparamjson','queryParams':{
			qtxt:$('#qtxt').val(),
			apiid:{{.apiid}}
		}});
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:480,
				height:390,
				modal:true
			});
			w.window('open');
			w.window('refresh', '/adm/tb/apiparamedit?apiid={{.apiid}}&id='+row.id);

        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }
}


function doAdd() {
	var w=$('#win').window({
		width:480,
		height:390,
		modal:true
	});
	w.window('open');
	w.window('refresh', '/adm/tb/apiparamedit?id=&apiid={{.apiid}}');
}
function doData(title,url){
	top.addTab(title,url);
}

function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('确认', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/apiparamdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        alert('删除失败!');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}


	function rowformater_field(value, row, index) {
		var a= "";
		return a;
	}
	
    $(function(){
		doSearch();
    })	
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           title="参数管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="title" align="center" width="10" >标题</th>
                <th field="param_type" align="center" width="10">类型</th>
                <th field="param_name" width="10">参数</th>                
				<th field="max_length" width="10">长度</th>  
				<th field="is_require" width="10">必填</th>  
				<th field="param_value" width="10">参数值</th>  
				<th field="memo" width="10">备注</th>  
				<th field="state" width="5">状态</th>
				<th field=" "  data-options="formatter:rowformater_field" align="center" width="10" >操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true"  onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:490px;height:390px;padding:5px;overflow-x:hidden;">
        Some Content.
    </div>
    <script type="text/javascript">
    $('#tt').datagrid({
        onLoadSuccess: function (data) {
           
        }
    });
</script>
</body>
</html>
`
var adm_tb_apiparamedit = `

<script type="text/javascript">
    var jq = jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
        $(function () {

        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    }else if(data=='0'){
                        jq.messager.alert('错误', "操作失败!", "warning");
                    }else{
						jq.messager.alert('错误', data, "warning");
					}
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }

</script>

<div class="easyui-panel"  style="width:99%" fix="true" border="false">
    <div style="padding:10px 20px 20px 20px">
        <form id="form1" action="/adm/tb/apiparampost" method="post">
            <table cellpadding="5">
				
				
				
				
                <tr>
                    <td style="width:55px;">编号:</td>
                    <td>
					
						{{.m.id}}
					 
					
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">标题:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="title" 
					id="title" value='{{.m.title}}' ></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">参数:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="param_name" 
					id="param_name" value='{{.m.param_name}}' ></input>
					</td>
                </tr>
                
                <tr>
                    <td style="width:55px;">类型:</td>
                    <td>
						<select id="param_type" name="param_type" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="参数">参数</option>
							<option  value="变量">变量</option>
							<option  value="cookie">Cookie</option>
							<option  value="session">Session</option>
							<option  value="data">Data</option>
							<option  value="dataset">DataSet</option>
							<option  value="sql">SQL</option>
							<option  value="sqlexec">SQLExec</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#param_type').combobox({
									onLoadSuccess : function(data) {
										$('#param_type').combobox('setValue', "{{.m.param_type}}");
									}
								});
                            });
                            
						</script>
						
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">长度:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="max_length" 
					id="max_length" value='{{.m.max_length}}' ></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">参数值:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="param_value" id="param_value" value='{{.m.param_value}}'></input>
					</td>
				</tr>
                <tr>
                    <td style="width:55px;">备注:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="memo" 
					id="memo" value='{{.m.memo}}'  style="width:135px;"></input>
					</td>
                </tr>
                <tr>
                    <td style="width:55px;">必填:</td>
                    <td>
					<input class="easyui-switchbutton" id="is_require" title="" name="is_require" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_require}}'=='1'){
								$('#is_require').switchbutton({
									checked: true,
								})
							}else{
								$('#is_require').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">必填提示:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="is_require_info" 
					id="is_require_info" value='{{.m.is_require_info}}'  style="width:135px;"></input>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">唯一:</td>
                    <td>
					<input class="easyui-switchbutton" id="is_unique" title="" name="is_unique" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_unique}}'=='1'){
								$('#is_unique').switchbutton({
									checked: true,
								})
							}else{
								$('#is_unique').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">唯一提示:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="is_unique_info" 
					id="is_unique_info" value='{{.m.is_unique_info}}'  style="width:135px;"></input>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">校验:</td>
                    <td>
					<input class="easyui-switchbutton" id="is_checkout" title="" name="is_checkout" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_checkout}}'=='1'){
								$('#is_checkout').switchbutton({
									checked: true,
								})
							}else{
								$('#is_checkout').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">校验类型:</td>
                    <td>
						<select id="is_checkout_type" name="is_checkout_type" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="=">=</option>
							<option  value="!=">!=</option>
							<option  value=">">></option>
							<option  value=">=">>=</option>
							<option  value="<"><</option>
							<option  value="<="><=</option>
							<option  value="正则">正则</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#is_checkout_type').combobox({
									onLoadSuccess : function(data) {
										$('#is_checkout_type').combobox('setValue', "{{.m.is_checkout_type}}");
									}
								});
                            });
                            
						</script>
						
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">校验内容:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="is_checkout_val" 
					id="is_checkout_val" value='{{.m.is_checkout_val}}'  style="width:135px;"></input>
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">校验提示:</td>
                    <td>
					<input class="easyui-textbox" title="" type="text" name="is_checkout_info" 
					id="is_checkout_info" value='{{.m.is_checkout_info}}'  style="width:135px;"></input>
					</td>
				</tr>
				
                <tr>
                    <td style="width:55px;">状态:</td>
                    <td>
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.state}}'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
                
            </table>
			<input type="hidden" id="id" name="id" value="{{.m.id}}" />
			<input type="hidden" id="api_id" name="api_id" value="{{.apiid}}" />
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>

`

//系统参数-----------------------------------------------------------
//system 列表页面
func (c *TbController) System() {
	var tpl = template.New("")
	tpl.Parse(adm_tb_system)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//获取adm_system列表
func (c *TbController) SystemJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where `sys_name` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select * from adm_system "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//系统参数添加/编辑
func (c *TbController) SystemEdit() {
	var id, _ = c.GetInt("id", 0)
	if id > 0 {
		var m = db.First("select * from adm_system where id=?", id)
		c.Data["m"] = m
	}

	var tpl = template.New("")
	tpl.Parse(adm_tb_systemedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//系统参数保存提交
func (c *TbController) SystemEditPost() {
	var id, _ = c.GetInt("id", 0)
	var mch_id, _ = c.GetInt("mch_id", 0)
	var user_id, _ = c.GetInt("user_id", 0)
	var sys_name = c.GetString("sys_name")
	var sys_logo = c.GetString("sys_logo")
	if len(sys_name) < 1 {
		c.Ctx.WriteString("请输入系统名称!")
		return
	}

	var sql = ""
	if id > 0 {
		sql = `
		update adm_system set 
		mch_id=?,
		user_id=?,
		sys_name=?,
		sys_logo=?
		where id=?
		`
		var i = db.Exec(sql,
			mch_id,
			user_id,
			sys_name,
			sys_logo,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `
		insert into adm_sysem(
			mch_id,
			user_id,
			sys_name,
			sys_logo
		)values(
			?,?,?,?
		)
		`
		var i = db.Exec(sql,
			mch_id, user_id,
			sys_name, sys_logo,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除adm_sysem表数据
func (c *TbController) SystemDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from adm_sysem where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

var adm_tb_system = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />
	

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
	<script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
            $('#win').window('open');
            $('#win').window('refresh', '/adm/tb/systemedit?id='+row.id);
			$('#win').window("resize",{top:$(document).scrollTop() + ($(window).height()-250) * 0.5});//居中显示
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
    $('#win').window('open');
    $('#win').window('refresh', '/adm/tb/systemedit?id=');
    $('#ff').form('load', row);
}
function doRemove(){
    var jq=jQuery;

    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/delsys', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        jq.messager.alert('警告','删除失败','warning');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })

	Date.prototype.Format = function (fmt) { //author: meizz   
            var o = {
                "M+": this.getMonth() + 1,                 //月份   
                "d+": this.getDate(),                    //日   
                "h+": this.getHours(),                   //小时   
                "m+": this.getMinutes(),                 //分   
                "s+": this.getSeconds(),                 //秒   
                "q+": Math.floor((this.getMonth() + 3) / 3), //季度   
                "S": this.getMilliseconds()             //毫秒   
            };
            if (/(y+)/.test(fmt))
                fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
            for (var k in o)
                if (new RegExp("(" + k + ")").test(fmt))
                    fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
            return fmt;
    }
	function rowformater_headimg(value, row, index) {
		return "<span class=' "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</span>";
		//return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
	function rowformater_date(value, row, index) {
       if (value == undefined) {
        return "";
		}
		/*json格式时间转js时间格式*/
		value = value.substr(1, value.length - 2);
		var obj = eval('(' + "{Date: new " + value + "}" + ')');
		var dateValue = obj["Date"];
		if (dateValue.getFullYear() < 1900) {
			return "";
		}

		return dateValue.Format("yyyy-MM-dd hh:mm:ss");
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/systemjson"
           title="系统管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           >
        <thead>
            <tr>
                <th field="id" width="20">编号</th>
				<th field="mch_id" width="70">商户</th>
				<th field="user_id" width="70">用户</th>
                <th field="sys_name" align="left" width="70">系统</th>
				<th field="sys_log" width="50">LOGO</th>
				<th field=" " width="50">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true" onclick="doRemove();">删除</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:80px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:420px;height:320px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>

</body>
</html>
`
var adm_tb_systemedit = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            //$('#pid').val('$!m.parentid');
            if ('{{.m.state}}' == '1') {
                $('#state').attr('checked', 'checked');
            }
            $('#images').val('$!m.images');
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data != "0") {
                        layer.msg('<font color="yellow">操作成功!</font>');
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        layer.msg('<font color="red">操作失败!</font>');
                    }
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }
        $('#image').combobox({
            formatter: function (row) {

                return '<span class="' + row.text + '">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="item-text">' + row.text + '</span>';
            }
        });
</script>


<div class="easyui-panel" title="" style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/systemeditpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>商户ID:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="mch_id" value="{{.m.mch_id}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
				</tr>
				<tr>
                    <td>用户ID:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="user_id" value="{{.m.user_id}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
                <tr>
                    <td>系统名称:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="sys_name" value="{{.m.sys_name}}" ></input></td>
                </tr>
				
                <tr>
                    <td>LOGO:</td>
                    <td>
                        <input class="easyui-textbox" type="text" style="width:180px;" name="sys_logo" value="{{.m.sys_logo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>

            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`

//模块项目-----------------------------------------------------------
func (c *TbController) Proj() {
	var tpl = template.New("")
	tpl.Parse(adm_tb_proj)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//获取tb_table_proj列表
func (c *TbController) ProjJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where `proj_name` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select * from tb_table_proj "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//项目添加/编辑
func (c *TbController) ProjEdit() {
	var id, _ = c.GetInt("id", 0)
	if id > 0 {
		var m = db.First("select * from tb_table_proj where id=?", id)
		c.Data["m"] = m
	}

	var tpl = template.New("")
	tpl.Parse(adm_tb_projedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//系统参数保存提交
func (c *TbController) ProjEditPost() {
	var id, _ = c.GetInt("id", 0)
	var proj_name = c.GetString("proj_name")
	var memo = c.GetString("memo")
	if len(proj_name) < 1 {
		c.Ctx.WriteString("请输入项目名称!")
		return
	}

	var sql = ""
	if id > 0 {
		sql = `
		update tb_table_proj set 
		proj_name=?,
		state=1,
		memo=?
		where id=?
		`
		var i = db.Exec(sql,
			proj_name,
			memo,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `
		insert into tb_table_proj(
			proj_name,
			state,
			memo
		)values(
			?,?,?
		)
		`
		var i = db.Exec(sql,
			proj_name, 1, memo,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除tb_table_proj表数据
func (c *TbController) ProjDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from tb_table_proj where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

var adm_tb_proj = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />
	

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
	<script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
            $('#win').window('open');
            $('#win').window('refresh', '/adm/tb/projedit?id='+row.id);
			$('#win').window("resize",{top:$(document).scrollTop() + ($(window).height()-250) * 0.5});//居中显示
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
    $('#win').window('open');
    $('#win').window('refresh', '/adm/tb/projedit?id=');
    $('#ff').form('load', row);
}
function doRemove(){
    var jq=jQuery;

    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/projdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        jq.messager.alert('警告','删除失败','warning');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })

	Date.prototype.Format = function (fmt) { //author: meizz   
            var o = {
                "M+": this.getMonth() + 1,                 //月份   
                "d+": this.getDate(),                    //日   
                "h+": this.getHours(),                   //小时   
                "m+": this.getMinutes(),                 //分   
                "s+": this.getSeconds(),                 //秒   
                "q+": Math.floor((this.getMonth() + 3) / 3), //季度   
                "S": this.getMilliseconds()             //毫秒   
            };
            if (/(y+)/.test(fmt))
                fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
            for (var k in o)
                if (new RegExp("(" + k + ")").test(fmt))
                    fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
            return fmt;
    }
	function rowformater_headimg(value, row, index) {
		return "<span class=' "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</span>";
		//return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
	function rowformater_date(value, row, index) {
       if (value == undefined) {
        return "";
		}
		/*json格式时间转js时间格式*/
		value = value.substr(1, value.length - 2);
		var obj = eval('(' + "{Date: new " + value + "}" + ')');
		var dateValue = obj["Date"];
		if (dateValue.getFullYear() < 1900) {
			return "";
		}

		return dateValue.Format("yyyy-MM-dd hh:mm:ss");
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/projjson"
           title="项目管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           >
        <thead>
            <tr>
                <th field="id" width="20">ID</th>
                <th field="proj_name" align="left" width="30">项目</th>
				<th field="memo" width="30">备注</th>
				<th field=" " width="10">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true" onclick="doRemove();">删除</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:80px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:420px;height:320px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>

</body>
</html>
`
var adm_tb_projedit = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            if ('{{.m.state}}' == '1') {
                $('#state').attr('checked', 'checked');
            }
            $('#images').val('$!m.images');
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data != "0") {
                        layer.msg('<font color="yellow">操作成功!</font>');
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        layer.msg('<font color="red">操作失败!</font>');
                    }
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }
        $('#image').combobox({
            formatter: function (row) {

                return '<span class="' + row.text + '">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="item-text">' + row.text + '</span>';
            }
        });
</script>


<div class="easyui-panel" title="" style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/projeditpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>编号:</td>
                    <td>{{.m.id}}</td>
				</tr>
                <tr>
                    <td>项目:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="proj_name" value="{{.m.proj_name}}" ></input></td>
                </tr>
				
                <tr>
                    <td>备注:</td>
                    <td>
                        <input class="easyui-textbox" type="text" style="width:180px;" name="memo" value="{{.m.memo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>

            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`

//操作日志-----------------------------------------------------------
func (c *TbController) Log() {
	var tpl = template.New("")
	tpl.Parse(adm_tb_log)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//获取adm_log列表
func (c *TbController) LogJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where `title` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select * from adm_log "+where+" order by id desc ")
	c.Data["json"] = rst
	c.ServeJSON()
}

//日志详情
func (c *TbController) LogInfo() {
	var id = c.GetString("id")
	if id == "" {
		id = "0"
	}
	var where = " where id=" + id
	var m = db.First("select * from adm_log " + where + " ")
	c.Data["m"] = m

	var tpl = template.New("")
	tpl.Parse(adm_tb_loginfo)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()
	//rst=strings.Replace(rst,";",";<br/>",-1)

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

var adm_tb_log = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />
	

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
	<script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        //alert($('#pid').combobox("getValue"));
        $('#tt').datagrid('load',{
			qtxt:$('#qtxt').val()
        });
    }

	
function doRemove(){
    var jq=jQuery;

    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/logdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        jq.messager.alert('警告','删除失败','warning');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })

	function rowformater_log(value, row, index) {
		//return "<div title=' "+row.content+"'>value</div>";
		return "<div><a target='_blank' href='/adm/tb/loginfo?id="+row.id+"'>"+value+"</a></div>";
    }

    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/logjson"
           title="操作日志" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[50,100],pageSize:50,pagination:true"
           >
        <thead>
            <tr>
                <th field="id" width="10">ID</th>
				<th field="logtype" width="10">日志类型</th>
				<th field="opertype" width="10">操作类型</th>
				<th field="username" width="10">操作员</th>
				<th field="title" width="50" data-options="formatter:rowformater_log">操作内容</th>
				<th field="ip" width="12">IP</th>
				<th field="addtime" width="10" >时间</th>
				<th field=" " width="5">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <!--<a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true" onclick="doRemove();">删除</a>-->
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:120px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:420px;height:320px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>

</body>
</html>
`

var adm_tb_loginfo = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>日志详情</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />
	

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
    $(function(){

    })

    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">
    <div id="tb" style="padding:5px;height:auto">
		<div>
		{{.m.content}}
        </div>
    </div>


</body>
</html>
`

//扩展功能-----------------------------------------------------------
func (c *TbController) TbEx() {
	var tbid, _ = c.GetInt("tbid", 0)
	c.Data["tbid"] = tbid

	var tpl = template.New("")
	tpl.Parse(adm_tb_ex)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//获取tb_table_ex列表
func (c *TbController) TbExJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where `title` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select * from tb_table_ex "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//扩展添加/编辑
func (c *TbController) TbExEdit() {
	var id, _ = c.GetInt("id", 0)
	var tbid, _ = c.GetInt("tbid", 0)
	c.Data["tbid"] = tbid

	if id > 0 {
		var m = db.First("select * from tb_table_ex where id=?", id)
		c.Data["m"] = m
	}

	var tpl = template.New("")
	tpl.Parse(adm_tb_exedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

//保存提交
func (c *TbController) TbExPost() {
	var id, _ = c.GetInt("id", 0)
	var tbid, _ = c.GetInt("tbid", 0)
	var title = c.GetString("title")
	var extype = c.GetString("extype")
	var expage = c.GetString("expage")
	var explace = c.GetString("explace")
	var excontent = c.GetString("excontent")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}
	if len(title) < 1 {
		c.Ctx.WriteString("请输入名称!")
		return
	}

	var sql = ""
	if id > 0 {
		sql = `
		update tb_table_ex set 
		tbid=?,
		title=?,
		extype=?,
		expage=?,
		explace=?,
		excontent=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			tbid,
			title,
			extype,
			expage,
			explace,
			excontent,
			state,
			id,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `
		insert into tb_table_ex(
			tbid,
			title,
			extype,
			expage,
			explace,
			excontent,
			state
		)values(
			?,?,?,?,?,?,?
		)
		`
		var i = db.Exec(sql,
			tbid, title, extype, expage, explace, excontent, state,
		)
		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//删除tb_table_ex表数据
func (c *TbController) TbExDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("-1")
		return
	}
	var i = db.Exec("delete from tb_table_ex where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
	} else {
		c.Ctx.WriteString("0")
	}
}

var adm_tb_ex = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link href="/css/www.css" rel="stylesheet" type="text/css" />
	

	<script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/jquery.form.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/js/easyui/locale/easyui-lang-zh_CN.js"></script>
	<script type="text/javascript" src="/js/layer/layer.js"></script>

    <style>
        body {
            background: #fff;
        }
    </style>
    </style>
    <script type="text/javascript">
	var jq=jQuery;
	if(jq==undefined){
		jq=jQuery;
	}
function doSearch(){
        $('#tt').datagrid('load',{
			tbid:{{.tbid}},
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
            $('#win').window('open');
            $('#win').window('refresh', '/adm/tb/tbexedit?tbid={{.tbid}}&id='+row.id);
			$('#win').window("resize",{top:$(document).scrollTop() + ($(window).height()-250) * 0.5});//居中显示
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
    $('#win').window('open');
    $('#win').window('refresh', '/adm/tb/tbexedit?tbid={{.tbid}}&id=');
    $('#ff').form('load', row);
}
function doRemove(){
    var jq=jQuery;

    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/tb/tbexdel', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');	// reload the user data
                    } else {
                        jq.messager.alert('警告','删除失败','warning');
                    }
                });
            }
        });
    } else {
        jq.messager.alert('警告','请选择一行数据','warning');
    }

}
    $(function(){

    })

	Date.prototype.Format = function (fmt) { //author: meizz   
            var o = {
                "M+": this.getMonth() + 1,                 //月份   
                "d+": this.getDate(),                    //日   
                "h+": this.getHours(),                   //小时   
                "m+": this.getMinutes(),                 //分   
                "s+": this.getSeconds(),                 //秒   
                "q+": Math.floor((this.getMonth() + 3) / 3), //季度   
                "S": this.getMilliseconds()             //毫秒   
            };
            if (/(y+)/.test(fmt))
                fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
            for (var k in o)
                if (new RegExp("(" + k + ")").test(fmt))
                    fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
            return fmt;
    }
	function rowformater_headimg(value, row, index) {
		return "<span class=' "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</span>";
		//return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
	function rowformater_date(value, row, index) {
       if (value == undefined) {
        return "";
		}
		/*json格式时间转js时间格式*/
		value = value.substr(1, value.length - 2);
		var obj = eval('(' + "{Date: new " + value + "}" + ')');
		var dateValue = obj["Date"];
		if (dateValue.getFullYear() < 1900) {
			return "";
		}

		return dateValue.Format("yyyy-MM-dd hh:mm:ss");
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/tb/tbexjson?tbid={{.tbid}}"
           title="扩展管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           >
        <thead>
            <tr>
                <th field="id" width="10">ID</th>
				<th field="title" align="left" width="20">标题</th>
				<th field="extype" align="left" width="20">类型</th>
				<th field="expage" align="left" width="20">页面</th>
				<th field="explace" align="left" width="20">位置</th>
				<th field="state" width="10">状态</th>
				<th field=" " width="10">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-add" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-cancel" plain="true" onclick="doRemove();">删除</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:80px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:580px;height:470px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>

</body>
</html>
`
var adm_tb_exedit = `

<script type="text/javascript">
    var jq = jQuery;
        $(function () {
            if ('{{.m.state}}' == '1') {
                $('#state').attr('checked', 'checked');
            }
            $('#images').val('$!m.images');
        })
        function submitForm(){
            $('#form1').form('submit', {
                success: function (data) {
                    if (data != "0") {
                        layer.msg('<font color="yellow">操作成功!</font>');
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        layer.msg('<font color="red">操作失败!</font>');
                    }
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }
        $('#image').combobox({
            formatter: function (row) {

                return '<span class="' + row.text + '">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="item-text">' + row.text + '</span>';
            }
        });
</script>


<div class="easyui-panel" title="" style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/tb/tbexpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>编号:</td>
                    <td>{{.m.id}}</td>
				</tr>
                <tr>
                    <td>标题:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="title" value="{{.m.title}}" ></input></td>
                </tr>
				<tr>
                    <td style="width:55px;">类型:</td>
                    <td>
						<select id="extype" name="extype" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="sql">sql</option>
							<option  value="go">变go量</option>
							<option  value="javascript">javascript</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#extype').combobox({
									onLoadSuccess : function(data) {
										$('#extype').combobox('setValue', "{{.m.extype}}");
									}
								});
                            });
                            
						</script>
						
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">页面:</td>
                    <td>
						<select id="expage" name="expage" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="list">列表</option>
							<option  value="info">编辑</option>
							<option  value="view">详情</option>
							<option  value="del">删除</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#expage').combobox({
									onLoadSuccess : function(data) {
										$('#expage').combobox('setValue', "{{.m.expage}}");
									}
								});
                            });
                            
						</script>
						
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">位置:</td>
                    <td>
						<select id="explace" name="explace" style="width:142px;" class="easyui-combobox" editable="false">
							<option  value="LOAD">LOAD</option>
							<option  value="AFTER">AFTER</option>
							<option  value="BEFORE">BEFORE</option>
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#explace').combobox({
									onLoadSuccess : function(data) {
										$('#explace').combobox('setValue', "{{.m.explace}}");
									}
								});
                            });
                            
						</script>
						
					</td>
				</tr>
				<tr>
                    <td style="width:55px;">内容:</td>
                    <td>
					<input class="easyui-textbox" multiline="true" style="width:280px;height:120px" title="" type="text" 
					name="excontent" id="excontext" value='{{.m.excontent}}'></input>
						
					</td>
				</tr>
                <tr>
                    <td>状态:</td>
                    <td>
					<input class="easyui-switchbutton" id="state" title="" name="state" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('1'=='1'){
								$('#state').switchbutton({
									checked: true,
								})
							}else{
								$('#state').switchbutton({
									checked: false,
								})
							}
						})
					</script>
						<input type="hidden" id="id" name="id" value="{{.m.id}}" />
						<input type="hidden" id="tbid" name="tbid" value="{{.tbid}}" />
                    </td>
                </tr>

            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-cancel" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`

//生成代码: 列表 新增 删除 查询 修改 导入 导出 详情 打印----------------------------------------------------------------------------------
func (c *TbController) Codes() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	var tb = db.First("select * from tb_table where id=?", id)
	if tb == nil {
		c.Ctx.WriteString("code  not found")
		return
	}

	var codes = `
	package controllers
	import (
		"github.com/astaxie/beego"
	)
	
	//XController 控制器
	type XController struct {
		beego.Controller
	}
		
	`

	c.Ctx.WriteString(codes)
}
