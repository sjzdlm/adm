package adm

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	//"strconv"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//UserController 控制器
type UserController struct {
	beego.Controller
}

//输出个人session中的信息
func (c *UserController) Me() {
	var _uid = c.GetSession("_uid")
	c.Data["_uid"] = _uid
	var _mch_id = c.GetSession("_mch_id")
	c.Data["_mch_id"] = _mch_id
	var _pid = c.GetSession("_pid")
	c.Data["_pid"] = _pid
	var _pids = c.GetSession("_pids")
	c.Data["_pids"] = _pids
	var _roles = c.GetSession("_roles")
	c.Data["_roles"] = _roles
	var _username = c.GetSession("_username")
	c.Data["_username"] = _username
	var _sproot = c.GetSession("_sproot")
	c.Data["_sproot"] = _sproot
	var _is_manager = c.GetSession("_is_manager")
	c.Data["_is_manager"] = _is_manager
	var _company = c.GetSession("_company")
	c.Data["_company"] = _company
	var _company_id = c.GetSession("_company_id")
	c.Data["_company_id"] = _company_id
	var _company_pid = c.GetSession("_company_pid")
	c.Data["_company_pid"] = _company_pid
	var _userlevel = c.GetSession("_userlevel")
	c.Data["_userlevel"] = _userlevel

	var html = `
	var _uid={{._uid}} </br>
	var _mch_id={{._mch_id}}</br>
	var _pid={{._pid}}</br>
	var _pids={{._pids}}</br>
	var _roles={{._roles}}</br>
	var _username={{._username}}</br>
	var _sproot={{._sproot}}</br>
	var _sproot={{._sproot}}</br>
	var _is_manager={{._is_manager}}</br>
	var _company={{._company}}</br>
    var _company_id={{._company_id}}</br>
    var _company_pid={{._company_pid}}</br>
	var _userlevel={{._userlevel}}</br>
	`

	var tpl = template.New("")
	tpl.Parse(html)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)
	if e != nil {
		fmt.Println("template 执行错误:", e.Error())
		c.Ctx.WriteString("{}")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))
}

//用户信息页面
func (c *UserController) UInfo() {
	var _uid = c.GetSession("_uid")
	if _uid == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}
	var m = db.First("select * from adm_user where id=?", _uid)
	if m == nil {
		c.Ctx.WriteString("error")
		return
	}
	c.Data["m"] = m

	//c.TplName="adm/user/uinfo.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_uinfo)
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

//修改用户信息
func (c *UserController) UInfoPost() {
	var _uid = c.GetSession("_uid")
	if _uid == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}
	var m = db.First("select * from adm_user where id=?", _uid)
	if m == nil {
		c.Ctx.WriteString("error")
		return
	}

	var realname = c.GetString("realname")
	var headimg = c.GetString("headimg")
	if realname == "" {
		c.Ctx.WriteString("请输入您的姓名!")
		return
	}

	var sql = `
	update adm_user set
	realname=?,
	headimg=?
	where id=?
	`
	var i = db.Exec(sql, realname, headimg, _uid)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

//修改密码页面
func (c *UserController) Pwd() {
	var _uid = c.GetSession("_uid")
	if _uid == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}
	var m = db.First("select * from adm_user where id=?", _uid)
	if m == nil {
		c.Ctx.WriteString("error")
		return
	}
	c.Data["m"] = m
	//c.TplName="adm/user/pwd.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_pwd)
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

//修改密码
func (c *UserController) PwdPost() {
	var _uid = c.GetSession("_uid")
	if _uid == nil {
		fmt.Println("_uid 为空,修改密码失败!")
		c.Ctx.WriteString("参数错误!")
		return
	}

	var pwd1 = c.GetString("pwd1")
	var pwd2 = c.GetString("pwd2")
	if pwd1 == "" {
		c.Ctx.WriteString("请输入一个密码!")
		return
	}
	if pwd1 != pwd2 {
		c.Ctx.WriteString("两次密码不一致!")
		return
	}
	var sql = `
	update adm_user set 
	password=?
	where id=?
	`
	var i = db.Exec(sql, pwd1, _uid)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

//头像数组
func (c *UserController) HeadImg() {
	var hi = make([]map[string]string, 0)
	for i := 0; i < 10; i++ {
		var row = make(map[string]string)
		row["id"] = strconv.Itoa(i)
		row["text"] = strconv.Itoa(i)
		row["icon"] = "/images/headimg/" + strconv.Itoa(i) + ".jpg"
		hi = append(hi, row)
	}
	c.Data["json"] = hi
	c.ServeJSON()
}

//List 列表页面
func (c *UserController) List() {
	c.Data["_username"] = c.GetSession("_username").(string)
	//账号类型信息
	var sql = `select * from adm_usertype `
	if c.GetSession("_username").(string) != "root" {
		sql += " where level >= " + c.GetSession("_usertype").(string)
	}
	c.Data["usertype_list"] = db.Query(sql)

	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_list)
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

//获取用户列表
func (c *UserController) ListJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where (`username` like '%" + qtxt + "%' or realname like '%" + qtxt + "%' or company like '%" + qtxt + "%') "
	}

	//---------------------------------------------------------------------------------------
	// //根据自己的权限进行过滤
	// if c.GetSession("_username").(string) != "root" {
	// 	var pid = c.GetSession("_uid").(string)
	// 	var w = ChildIds(pid)
	// 	w = strings.Replace(w, ",,", ",", -1)
	// 	w = strings.Replace(w, ",,", ",", -1)
	// 	if where != "" {
	// 		where += " and id in(" + w + ")"
	// 	} else {
	// 		where = "where id in(" + w + ") "
	// 	}
	// }

	//根据单位级别进行数据过滤
	if c.GetSession("_username").(string) != "root" {
		var _mch_id = c.GetSession("_mch_id").(string)
		var _uid = c.GetSession("_uid").(string)
		var _is_manager = c.GetSession("_is_manager").(string)
		var _usertype = c.GetSession("_usertype").(string)
		var _company_id = c.GetSession("_company_id").(string)
		//如果是管理员,展示自己级别及以下信息,不包括自己
		var ids = ChildIds2(_mch_id, _company_id)
		if _is_manager == "1" {
			ids = _company_id + "," + ids
			if where != "" {
				where += " and usertype >=" + _usertype + " and id !=" + _uid + " and company_id in(" + ids + ") and sproot !=1 "
			} else {
				where = "where usertype >=" + _usertype + " and id !=" + _uid + " and company_id in(" + ids + ") and sproot !=1 "
			}
		} else {
			if where != "" {
				where += " and usertype >=" + _usertype + " and id !=" + _uid + " and company_id in(" + ids + ") and sproot !=1 "
			} else {
				where = "where usertype >=" + _usertype + " and id !=" + _uid + " and company_id in(" + ids + ") and sproot !=1 "
			}
		}
	}
	//------------------------------------------------------------------------------------------

	var usertype = c.GetString("usertype")
	if usertype != "0" && usertype != "" {
		if where != "" {
			where += " and usertype='" + usertype + "' "
		} else {
			where += " where  usertype='" + usertype + "' "
		}

	}

	//排序
	var sort = c.GetString("sort")
	var order = c.GetString("order")
	if sort != "" && order != "" {
		where += " order by " + sort + " " + order
	} else {
		where += " order by usertype,company_pid,company_id "
	}

	//fmt.Println("where:", where)
	var rst = db.Pager(page, pageSize, "select *  from adm_user "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//Edit 用户编辑页面
func (c *UserController) Edit() {
	c.Data["_username"] = c.GetSession("_username").(string) //输出账号到模板
	// var usertype = c.GetSession("_usertype").(string)        //账号类型
	// c.Data["_usertype"] = usertype
	// var is_sq = `0` //是否可以修改pid 和pids  //直接设置style 进行隐藏和显示
	// if usertype == "0" ||
	// 	usertype == "1" ||
	// 	usertype == "2" ||
	// 	usertype == "3" ||
	// 	usertype == "7" ||
	// 	usertype == "4" {
	// 	is_sq = "1"
	// }

	var _company_id = c.GetSession("_company_id").(string) //单位ID
	c.Data["_company_id"] = _company_id

	var _userlevel = c.GetSession("_userlevel").(string) //账号类型
	c.Data["_userlevel"] = _userlevel
	var is_sq = `0` //是否可以修改pid 和pids  //直接设置style 进行隐藏和显示
	// if _userlevel == "0" ||
	// 	_userlevel == "1" ||
	// 	_userlevel == "2" ||
	// 	_userlevel == "3" ||
	// 	_userlevel == "4" {
	// 	is_sq = "1"
	// }
	c.Data["is_sq"] = is_sq

	var id, _ = c.GetInt("id", 0)
	var m = db.First("select * from adm_user where id=?", id)
	c.Data["m"] = m
	if len(m) > 0 {
		c.Data["uname"] = m["username"]
		c.Data["pwd"] = m["password"]
	} else {
		c.Data["uname"] = ""
		c.Data["pwd"] = "358719"
	}
	//公司列表
	var mchlist = db.Query("select * from adm_mch")
	c.Data["mchlist"] = mchlist
	//角色列表
	var where = ""
	if c.GetSession("_sproot").(string) != "1" {
		var roles = c.GetSession("_roles").(string)
		if roles != "" {
			where = " where id in(" + roles + ") "
		}
	}
	var roles = db.Query("select * from adm_role" + where)
	c.Data["roles"] = roles
	var roleids = ""
	for i, v := range roles {
		if i > 0 {
			roleids += ","
		}
		roleids += v["id"]
	}
	c.Data["roleids"] = roleids
	//账号类型列表 根据级别过滤
	var utypelist = db.Query("select * from adm_usertype where level >=? order by orders", _userlevel)
	c.Data["utypelist"] = utypelist
	//根据信息选择已有权限
	var jstr = ""
	if m != nil {
		var r = strings.Split(m["roles"], ",")
		for i := 0; i < len(r); i++ {
			jstr += `$('#role` + r[i] + `').prop('checked',true);`
		}
	}
	c.Data["jstr"] = template.JS(jstr)

	//c.TplName="adm/user/edit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_user_edit")
	tpl.Funcs(template.FuncMap{"str2html": beego.Str2html})
	tpl.Parse(adm_user_edit)
	//fmt.Println(adm_user_edit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!" + e.Error())
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}
func (c *UserController) JsonUType() {
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

	//账号类型对应的单位
	var rcount = 0
	jsonstr += `var jsoncompany_id={ `
	for _, row := range list {
		var m = row
		//单位绑定 绑定字段为 id val  从数据库中读取
		if m["conn_str"] != "" && m["bindapi"] != "" {
			var xdb = db.NewDb(m["conn_str"])
			var list1 = db.Query2(xdb, m["bindapi"])
			for _, vv := range list1 {
				if rcount > 0 {
					jsonstr += ","
				}
				jsonstr += `"key` + vv["id"] + `":"` + vv["val"] + `","pkey` + vv["id"] + `":"` + vv["pname"] + `"`
				rcount++
			}
		}
	}
	jsonstr += `};
	`

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Body([]byte(jsonstr))
}

//用户信息编辑
func (c *UserController) EditPost() {
	var id, _ = c.GetInt("id", 0)
	var username = c.GetString("username")
	var realname = c.GetString("realname")
	var mobile = c.GetString("mobile")
	var usertype, _ = c.GetInt("usertype", 2)
	var level, _ = c.GetInt("level", 1)
	var password = c.GetString("password")
	var defpage = c.GetString("defpage")
	var state = c.GetString("state")
	var pid, _ = c.GetInt("pid", 0)
	var pids = strings.Join(c.GetStrings("pids"), ",")
	var company_id = c.GetString("company_id")
	var company_pid = ChildIdPid2(c.GetSession("_mch_id").(string), company_id) //上级单位id
	var company = c.GetString("company")
	var is_manager = c.GetString("is_manager")
	if is_manager == "on" || is_manager == "1" {
		is_manager = "1"
	} else {
		is_manager = "0"
	}

	if company_id == "" {
		c.Ctx.WriteString("账号所属单位不能为空!")
		return
	}

	var is_sq = "1"
	if pid == 0 && pids == "" {
		is_sq = "0"
	}

	//检查是否有重复账号
	if id > 0 {
		var jcu = db.First("select * from adm_user where username=? and id!=?", username, id)
		if len(jcu) > 0 {
			c.Ctx.WriteString("账号不能重复!")
			return
		}
	} else {
		var jcu = db.First("select * from adm_user where username=?", username)
		if len(jcu) > 0 {
			c.Ctx.WriteString("账号不能重复!")
			return
		}
	}
	//超管权限
	var sproot = c.GetString("sproot")
	if sproot == "" {
		sproot = "0"
	} else {
		//只有root可以启用关闭超管
		if c.GetSession("_username").(string) == "root" {
			if sproot == "on" {
				sproot = "1"
			} else {
				sproot = "0"
			}
		}
	}

	//默认是自己的企业ID,如果当前账号是root则可以修改
	var mch_id = c.GetSession("_mch_id").(string)
	if c.GetSession("_username").(string) == "root" {
		mch_id = c.GetString("mch_id")
	}

	if pid == 0 {
		pid, _ = strconv.Atoi(c.GetSession("_uid").(string))
	}
	if pids == "" {
		pids = strconv.Itoa(pid)
	}

	if username == "root" {
		pid = 0
		pids = "0"
	}

	var role = c.GetStrings("role")
	var roles = ""
	if len(role) < 1 {
		c.Ctx.WriteString("请选择至少一个角色.")
		return
	} else {
		for i := 0; i < len(role); i++ {
			if i > 0 {
				roles += ","
			}
			roles += role[i]
		}
	}
	var memo = c.GetString("memo")
	if id > 0 {
		var m = db.First("select * from adm_user where id=?", id)
		if m == nil {
			c.Ctx.WriteString("参数错误！")
			return
		}

		var sql = ""
		var i int64 = 0

		if is_sq == "0" { //不需要修改pid pids
			sql = `
		update adm_user set 
		sproot=?,
		username=?,
		realname=?,
		mobile=?,
		usertype=?,
		level=?,
		password=?,
		state=?,
		roles=?,
		defpage=?,
		mch_id=?,
		company=?,
		company_id=?,
		company_pid=?,
		is_manager=?,
		memo=?
		where id=?
		`
			i = db.Exec(sql,
				sproot,
				username,
				realname,
				mobile,
				usertype,
				level,
				password,
				state,
				roles,
				defpage,
				mch_id,
				company,
				company_id,
				company_pid,
				is_manager,
				memo,
				id,
			)
		} else {
			sql = `
		update adm_user set 
		sproot=?,
		pid=?,
		pids=?,
		username=?,
		realname=?,
		mobile=?,
		usertype=?,
		level=?,
		password=?,
		state=?,
		roles=?,
		defpage=?,
		mch_id=?,
		company=?,
		company_id=?,
		company_pid=?,
		is_manager=?,
		memo=?
		where id=?
		`
			i = db.Exec(sql,
				sproot,
				pid,
				pids,
				username,
				realname,
				mobile,
				usertype,
				level,
				password,
				state,
				roles,
				defpage,
				mch_id,
				company,
				company_id,
				company_pid,
				is_manager,
				memo,
				id,
			)
		}

		if i > 0 {
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		var sql = `insert into adm_user(
			sproot,
			pid,
			pids,
			username,
			realname,
			mobile,
			usertype,
			level,
			password,
			state,
			roles,
			defpage,
			mch_id,
			company,
			company_id,
			company_pid,
			headimg,
			regtime,
			is_manager,
			memo
		)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`
		var i = db.Exec(sql,
			sproot,
			pid,
			pids,
			username,
			realname,
			mobile,
			usertype,
			level,
			password,
			state,
			roles,
			defpage,
			mch_id,
			company,
			company_id,
			company_pid,
			"/images/headimg/2.jpg",
			time.Now().Format("2006-01-02 15:04:05"),
			is_manager,
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

//删除用户
func (c *UserController) Remove() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	var i = db.Exec("delete from adm_user where id=? and id>1", id)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

//Role 列表页面
func (c *UserController) Role() {
	var sysid = c.GetSession("_sysid")
	if sysid == nil {
		c.Ctx.WriteString("0")
		return
	}
	fmt.Println("-------------------sysid:", c.GetSession("_sysid"))
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_role)
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

//获取角色列表
func (c *UserController) RoleJson() {
	var sysid = c.GetSession("_sysid")
	if sysid == nil {
		c.Ctx.WriteString("0")
		return
	}
	fmt.Println("-------------------sysid:", c.GetSession("_sysid"))

	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")

	var where = " where sysid=" + sysid.(string)

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += "  and `name` like '%" + qtxt + "%'"
	}

	//排序
	var sort = c.GetString("sort")
	var order = c.GetString("order")
	if sort != "" && order != "" {
		where += " order by " + sort + " " + order
	}

	var rst = db.Pager(page, pageSize, "select * from adm_role "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//角色添加/编辑
func (c *UserController) RoleEdit() {
	var id, _ = c.GetInt("id", 0)
	if id > 0 {
		var m = db.First("select * from adm_role where id=?", id)
		c.Data["m"] = m
	}

	//读取账号类型
	var usertypelist = db.Query("select * from adm_usertype order by level")
	c.Data["usertypelist"] = usertypelist
	//c.TplName="adm/user/roleedit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_roleedit)
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

//角色保存提交
func (c *UserController) RoleEditPost() {
	var sysid = c.GetSession("_sysid")
	if sysid == nil {
		c.Ctx.WriteString("0")
		return
	}
	fmt.Println("-------------------sysid:", c.GetSession("_sysid"))

	var id, _ = c.GetInt("id", 0)
	var name = c.GetString("name")
	var r = c.GetStrings("rights")
	if len(r) < 1 {
		c.Ctx.WriteString("请选择角色权限!")
		return
	}
	var rights = ""
	for i := 0; i < len(r); i++ {
		if rights != "" {
			rights += ","
		}
		rights += r[i]
	}
	//找到上级ID
	var sqlstr = "select GROUP_CONCAT(pid) as pid from adm_menu where id in(" + rights + ") "
	var a1 = db.First(sqlstr)
	sqlstr = "select GROUP_CONCAT(pid) as pid from adm_menu where id in(" + a1["pid"] + ")"
	var a2 = db.First(sqlstr)
	rights += "," + a1["pid"] + "," + a2["pid"]

	var info = c.GetString("info")
	var memo = c.GetString("memo")
	var state = c.GetString("state")
	var level = c.GetString("level")

	var sql = ""
	if id > 0 {
		sql = `
		update adm_role set 
		name=?,
		rights=?,
		info=?,
		memo=?,
		level=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			name,
			rights,
			info,
			memo,
			level,
			state,
			id,
		)
		if i > 0 {
			//同步更新角色权限表adm_role_auth
			sql = `
			INSERT into adm_role_auth(roleid,rolename,menuid,menuname,ac_add,ac_del,ac_query,ac_update,ac_import,ac_export,ac_edit,ac_info,ac_print)
			select ` + strconv.Itoa(id) + ` as roleid,'` + name + `' as rolename,id as menuid,title as menuname,'1','1','1','1','1','1','1','1','1' from adm_menu where id not in (
			select menuid from adm_role_auth where roleid=` + strconv.Itoa(id) + `
			) 
			and id in (` + rights + `)
			`
			db.Exec(sql)
			sql = `update adm_role_auth set state=0 where roleid=` + strconv.Itoa(id)
			db.Exec(sql)
			sql = `update adm_role_auth set state=1 where roleid=` + strconv.Itoa(id) + ` and menuid in(` + rights + `)`
			db.Exec(sql)
			//删除多余的权限
			sql = `delete from adm_role_auth where state=0 and roleid=` + strconv.Itoa(id)
			db.Exec(sql)

			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `
		insert into adm_role(
			sysid,
			name,
			rights,
			info,
			memo,
			level,
			state
		)values(
			?,?,?,?,?,?,?
		)
		`
		var i = db.Exec(sql, sysid.(string),
			name, rights,
			info, memo, level, state,
		)
		if i > 0 {
			var m = db.First("select * from adm_role where name=? and rights=? order by id desc limit 1", name, rights)
			if len(m) > 0 {
				//同步更新角色权限表adm_role_auth
				sql = `
				INSERT into adm_role_auth(roleid,rolename,menuid,menuname,ac_add,ac_del,ac_query,ac_update,ac_import,ac_export,ac_edit,ac_info,ac_print)
				select ` + m["id"] + ` as roleid,'` + name + `' as rolename,id as menuid,title as menuname,'1','1','1','1','1','1','1','1','1' from adm_menu where id not in (
				select menuid from adm_role_auth where roleid=` + m["id"] + `
				) 
				and id in (` + rights + `)
				`
				db.Exec(sql)
				sql = `update adm_role_auth set state=0 where roleid=` + m["id"]
				db.Exec(sql)
				sql = `update adm_role_auth set state=1 where roleid=` + m["id"] + ` and menuid in(` + rights + `)`
				db.Exec(sql)
			}
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//DelRole 角色删除
func (c *UserController) DelRole() {
	var id, _ = c.GetInt("id", 0)
	var i = db.Exec("delete from adm_role where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

//早期分级检索方式,使用id  pid方式
// func ChildIds(pid string) string {
// 	//根据like语法读取
// 	var sql = "select GROUP_CONCAT(id) as ids from adm_user where pids ='" + pid + "' or pids like '%," + pid + "' or pids like '" + pid + ",%' or pids like '%," + pid + ",%' "
// 	//fmt.Println("pids sql:", sql)
// 	var p = db.First(sql)
// 	if len(p) < 1 {
// 		return "0"
// 	}
// 	var pids = p["ids"]
// 	// if rst == "" {
// 	// 	return "0"
// 	// }
// 	// return rst

// 	//读取根节点
// 	var list = db.Query("select * from adm_user where pid =" + pid + "")
// 	//第一层节点
// 	var rst = ""
// 	for k, v := range list {
// 		if k > 0 {
// 			rst += ","
// 		}
// 		rst += v["id"]
// 		//第二层节点
// 		var list1 = db.Query("select * from adm_user where pid=?", v["id"])
// 		rst1 := ""
// 		for kk, vv := range list1 {
// 			if kk > 0 {
// 				rst1 += ","
// 			}
// 			rst1 += vv["id"]
// 			//第三层节点
// 			var list2 = db.Query("select * from adm_user where pid=?", vv["id"])
// 			rst2 := ""
// 			for kkk, vvv := range list2 {
// 				if kkk > 0 {
// 					rst2 += ","
// 				}
// 				rst2 += vvv["id"]

// 				//第四层
// 				var list3 = db.Query("select * from adm_user where pid=?", vvv["id"])
// 				rst3 := ""
// 				for kkkk, vvvv := range list3 {
// 					if kkkk > 0 {
// 						rst3 += ","
// 					}
// 					rst3 += vvvv["id"]

// 					//第五层
// 					var list4 = db.Query("select * from adm_user where pid=?", vvvv["id"])
// 					rst4 := ""
// 					for kkkkk, vvvvv := range list4 {
// 						if kkkkk > 0 {
// 							rst4 += ","
// 						}
// 						rst4 += vvvvv["id"]
// 					}
// 					if rst4 != "" {
// 						rst3 += `,` + rst4
// 					}

// 				}
// 				if rst3 != "" {
// 					rst2 += `,` + rst3
// 				}

// 			}
// 			if rst2 != "" {
// 				rst1 += `,` + rst2
// 			}

// 		}
// 		if rst1 != "" {
// 			rst += `,` + rst1
// 		}

// 	}
// 	if rst == "" {
// 		rst = pids
// 	} else {
// 		rst += "," + pids
// 	}
// 	return rst
// }

//根据当前商户配置获取 部门(公司)父节点id 公司表包括 id,name,level,pid,pname
func ChildIdPid2(mch_id string, id string) string {
	//读取分级配置
	// var sys = db.First("select * from adm_system where mch_id=? and user_level=0 and user_id=0 ", mch_id)
	// var dbname = sys["company_db"]
	// var tbname = sys["company_tb"]

	var sys = db.First("select * from adm_mch where id=?", mch_id)
	var dbname = sys["company_db"]
	var tbname = sys["company_tb"]
	var field_id = sys["company_idfield"]
	var field_pid = sys["company_pidfield"]
	//var field_name = sys["company_namefield"]

	if dbname == "" || tbname == "" {
		if tbname != "" && dbname == "" {
			return ""
		}
		tbname = "adm_dept" //默认为系统部门表
		field_id = "id"
		field_pid = "pid"
		//field_name = "name"
	}
	var xdb = db.NewDb(dbname)
	if xdb == nil {
		xdb = db.NewDb("")
		if xdb == nil {
			return "0"
		}
	}
	var m = db.First2(xdb, "select * from "+tbname+" where "+field_id+"=?", id)
	if len(m) > 0 {
		return m[field_pid]
	}
	return ""
}

//根据当前商户配置获取 部门(公司)父节点名称 公司表包括 id,name,level,pid,pname
func ChildIdPname2(mch_id string, id string) string {
	//读取分级配置
	// var sys = db.First("select * from adm_system where mch_id=? and user_level=0 and user_id=0 ", mch_id)
	// var dbname = sys["company_db"]
	// var tbname = sys["company_tb"]

	var sys = db.First("select * from adm_mch where id=?", mch_id)
	var dbname = sys["company_db"]
	var tbname = sys["company_tb"]
	var field_id = sys["company_idfield"]
	var field_pid = sys["company_pidfield"]
	var field_name = sys["company_namefield"]

	if dbname == "" || tbname == "" {
		if tbname != "" && dbname == "" {
			return ""
		}
		tbname = "adm_dept" //默认为系统部门表
		field_id = "id"
		field_pid = "pid"
		field_name = "name"
	}
	var xdb = db.NewDb(dbname)
	if xdb == nil {
		xdb = db.NewDb("")
		if xdb == nil {
			return "0"
		}
	}
	var m = db.First2(xdb, "select * from "+tbname+" where "+field_id+"=?", id)
	if len(m) > 0 {
		var mm = db.First2(xdb, "select * from "+tbname+" where "+field_id+"=?", m[field_pid])
		if len(mm) > 0 {
			return mm[field_name]
		}
	}
	return ""
}

//根据当前商户配置获取所有部门(公司)子节点信息 公司表包括 id,name,level,pid,pname
func ChildIds2(mch_id string, pid string) string {
	//读取分级配置
	//var sys = db.First("select * from adm_system where mch_id=? and user_level=0 and user_id=0 ", mch_id)
	var sys = db.First("select * from adm_mch where id=?", mch_id)
	var dbname = sys["company_db"]
	var tbname = sys["company_tb"]
	var field_id = sys["company_idfield"]
	var field_pid = sys["company_pidfield"]
	//var field_name = sys["company_namefield"]

	if dbname == "" || tbname == "" {
		if tbname != "" && dbname == "" {
			return ""
		}
		tbname = "adm_dept" //默认为系统部门表
		field_id = "id"
		field_pid = "pid"
		//field_name = "name"
	}
	var xdb = db.NewDb(dbname)
	if xdb == nil {
		xdb = db.NewDb("")
		if xdb == nil {
			return "0"
		}
	}
	//读取根节点
	var list = db.Query2(xdb, "select "+field_id+","+field_pid+" from "+tbname+" where "+field_pid+" ="+pid+"")
	//第一层节点
	var rst = ""
	for k, v := range list {
		if v[field_id] == "" {
			continue
		}
		if k > 0 {
			rst += ","
		}
		rst += v[field_id]
		//第二层节点
		var list1 = db.Query2(xdb, "select "+field_id+","+field_pid+" from "+tbname+"  where "+field_pid+" =?", v[field_pid])
		rst1 := ""
		for kk, vv := range list1 {
			if vv[field_pid] == "" {
				continue
			}
			if kk > 0 {
				rst1 += ","
			}
			rst1 += vv[field_pid]
			//第三层节点
			var list2 = db.Query2(xdb, "select  "+field_id+","+field_pid+" from "+tbname+"  where "+field_pid+"=?", vv[field_id])
			rst2 := ""
			for kkk, vvv := range list2 {
				if vvv[field_id] == "" {
					continue
				}
				if kkk > 0 {
					rst2 += ","
				}
				rst2 += vvv[field_id]

				//第四层
				var list3 = db.Query2(xdb, "select  "+field_id+","+field_pid+" from "+tbname+"  where "+field_pid+"=?", vvv[field_id])
				rst3 := ""
				for kkkk, vvvv := range list3 {
					if vvvv[field_id] == "" {
						continue
					}
					if kkkk > 0 {
						rst3 += ","
					}
					rst3 += vvvv[field_id]

					//第五层
					var list4 = db.Query2(xdb, "select  "+field_id+","+field_pid+" from "+tbname+"  where "+field_pid+"=?", vvvv[field_id])
					rst4 := ""
					for kkkkk, vvvvv := range list4 {
						if vvvvv[field_id] == "" {
							continue
						}
						if kkkkk > 0 {
							rst4 += ","
						}
						rst4 += vvvvv[field_id]
						//第6层
						var list5 = db.Query2(xdb, "select  "+field_id+","+field_pid+" from "+tbname+"  where "+field_pid+"=?", vvvvv[field_id])
						rst5 := ""
						for kkkkkk, vvvvvv := range list5 {
							if vvvvvv[field_id] == "" {
								continue
							}
							if kkkkkk > 0 {
								rst5 += ","
							}
							rst5 += vvvvvv[field_id]

							//第7层
							var list6 = db.Query2(xdb, "select  "+field_id+","+field_pid+" from "+tbname+"  where "+field_pid+"=?", vvvvvv[field_id])
							rst6 := ""
							for kkkkkkk, vvvvvvv := range list6 {
								if vvvvvvv[field_id] == "" {
									continue
								}
								if kkkkkkk > 0 {
									rst6 += ","
								}
								rst6 += vvvvvvv[field_id]
							}
							if rst6 != "" {
								rst5 += `,` + rst6
							}
						}
						if rst5 != "" {
							rst4 += `,` + rst5
						}
					}
					if rst4 != "" {
						rst3 += `,` + rst4
					}

				}
				if rst3 != "" {
					rst2 += `,` + rst3
				}

			}
			if rst2 != "" {
				rst1 += `,` + rst2
			}

		}
		if rst1 != "" {
			rst += `,` + rst1
		}

	}
	if rst == "" {
		rst = "0"
	}
	rst = strings.Replace(rst, ",,", ",", -1)
	rst = strings.Replace(rst, ",,", ",", -1)
	rst = strings.Replace(rst, ",,", ",", -1)
	return rst
}

//UserTreeJson 用户JSON字符串 ids为需要选中的节点
func (c *UserController) UserTreeJson() {
	var id, _ = c.GetInt("id", 0)

	var sysid = c.GetSession("_sysid")
	if sysid == nil {
		c.Ctx.WriteString("0")
		return
	}
	//fmt.Println("-------------------sysid:", c.GetSession("_sysid"))
	var pid = sysid.(string)
	// if c.GetSession("_uid") != nil {
	// 	pid = c.GetSession("_uid").(string)
	// }

	//读取根节点
	var list = db.Query("select * from adm_user where id=?", pid)
	//fmt.Println("test ids:", ChildIds(pid))
	//第一层节点
	var rst = "["
	for k, v := range list {
		if k > 0 {
			rst += ","
		}
		rst += "{"
		rst += `"id":` + v["id"] + ","
		rst += `"text":"` + v["realname"] + "-" + v["id"] + `"`
		//第二层节点
		var list1 = db.Query("select * from adm_user where pid=? and id!=?", v["id"], id)
		rst1 := "["
		for kk, vv := range list1 {
			if kk > 0 {
				rst1 += ","
			}
			rst1 += "{"
			rst1 += `"id":` + vv["id"] + ","
			rst1 += `"text":"` + vv["realname"] + "-" + vv["id"] + `"`
			//第三层节点
			var list2 = db.Query("select * from adm_user where pid=? and id!=? ", vv["id"], id)
			rst2 := "["
			for kkk, vvv := range list2 {
				if kkk > 0 {
					rst2 += ","
				}
				rst2 += "{"
				rst2 += `"id":` + vvv["id"] + ","
				rst2 += `"text":"` + vvv["realname"] + "-" + vvv["id"] + `"`

				//第四层
				var list3 = db.Query("select * from adm_user where pid=? and id!=?", vvv["id"], id)
				rst3 := "["
				for kkkk, vvvv := range list3 {
					if kkkk > 0 {
						rst3 += ","
					}
					rst3 += "{"
					rst3 += `"id":` + vvvv["id"] + ","
					rst3 += `"text":"` + vvvv["realname"] + "-" + vvvv["id"] + `"`

					//第五层
					var list4 = db.Query("select * from adm_user where pid=? and id!=?", vvvv["id"], id)
					rst4 := "["
					for kkkkk, vvvvv := range list4 {
						if kkkkk > 0 {
							rst4 += ","
						}
						rst4 += "{"
						rst4 += `"id":` + vvvvv["id"] + ","
						rst4 += `"text":"` + vvvvv["realname"] + "-" + vvvvv["id"] + `"`
						rst4 += "}"
					}
					rst4 += "]"
					rst3 += `,"children":` + rst4
					rst3 += "}"
				}
				rst3 += "]"
				rst2 += `,"children":` + rst3
				rst2 += "}"
			}
			rst2 += "]"
			rst1 += `,"children":` + rst2
			rst1 += "}"
		}
		rst1 += "]"
		rst += `,"children":` + rst1
		rst += "}"
	}
	rst += "]"
	c.Ctx.WriteString(rst)
}

//TreeJson 菜单JSON字符串 ids为需要选中的节点
func (c *UserController) TreeJson() {
	var ids = c.GetString("ids")
	var idarray = strings.Split(ids, ",")
	var mapid = make(map[string]string)
	for _, v := range idarray {
		mapid[v] = v
	}
	//读取根节点
	var list = db.Query("select * from adm_menu where pid=1")

	//第一层节点
	var rst = "["
	for k, v := range list {
		if k > 0 {
			rst += ","
		}
		rst += "{"
		rst += `"id":` + v["id"] + ","
		rst += `"text":"` + v["title"] + "-" + v["id"] + `"`
		//第二层节点
		var list1 = db.Query("select * from adm_menu where pid=?", v["id"])
		rst1 := "["
		for kk, vv := range list1 {
			if kk > 0 {
				rst1 += ","
			}
			rst1 += "{"
			rst1 += `"id":` + vv["id"] + ","
			rst1 += `"text":"` + vv["title"] + "-" + vv["id"] + `"`
			//第三层节点
			var list2 = db.Query("select * from adm_menu where pid=?", vv["id"])
			rst2 := "["
			for kkk, vvv := range list2 {
				if kkk > 0 {
					rst2 += ","
				}
				rst2 += "{"
				rst2 += `"id":` + vvv["id"] + ","
				//校验是否需要选中
				if _, ok := mapid[vvv["id"]]; ok {
					rst2 += `"checked":true,`
				}
				rst2 += `"text":"` + vvv["title"] + "-" + vvv["id"] + `"`
				rst2 += "}"
			}
			rst2 += "]"
			rst1 += `,"children":` + rst2
			rst1 += "}"
		}
		rst1 += "]"
		rst += `,"children":` + rst1
		rst += "}"
	}
	rst += "]"
	c.Ctx.WriteString(rst)
}

//TreeJson 菜单JSON字符串 ids为需要选中的节点
func (c *UserController) TreeJson2() {
	var ids = c.GetString("ids")
	var idarray = strings.Split(ids, ",")
	var mapid = make(map[string]string)
	for _, v := range idarray {
		mapid[v] = v
	}
	//读取根节点
	var list = db.Query("select * from adm_menu where pid=1")

	//第一层节点
	var rst = "["
	for k, v := range list {
		if k > 0 {
			rst += ","
		}
		rst += "{"
		rst += `"id":` + v["id"] + ","
		rst += `"text":"` + v["label"] + "-" + v["id"] + `"`
		//第二层节点
		var list1 = db.Query("select * from adm_menu where pid=?", v["id"])
		rst1 := "["
		for kk, vv := range list1 {
			if kk > 0 {
				rst1 += ","
			}
			rst1 += "{"
			rst1 += `"id":` + vv["id"] + ","
			rst1 += `"text":"` + vv["label"] + "-" + vv["id"] + `"`
			//第三层节点
			var list2 = db.Query("select * from adm_menu where pid=?", vv["id"])
			rst2 := "["
			for kkk, vvv := range list2 {
				if kkk > 0 {
					rst2 += ","
				}
				rst2 += "{"
				rst2 += `"id":` + vvv["id"] + ","
				//校验是否需要选中
				if _, ok := mapid[vvv["id"]]; ok {
					rst2 += `"checked":true,`
				}
				rst2 += `"text":"` + vvv["label"] + "-" + vvv["id"] + `"`
				rst2 += "}"
			}
			rst2 += "]"
			rst1 += `,"children":` + rst2
			rst1 += "}"
		}
		rst1 += "]"
		rst += `,"children":` + rst1
		rst += "}"
	}
	rst += "]"
	c.Ctx.WriteString(rst)
}

var adm_user_uinfo = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>修改信息</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css"> 
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <!--<script type="text/javascript" src="/js/jquery.form.js"></script>-->
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <link href="/css/www.css" rel="stylesheet" type="text/css" />

    <style>
        body {
            background: #fff;
        }
    </style>
    <script type="text/javascript">
        var jq = jQuery;

        function submitForm() {
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                       jq.messager.alert('成功', "操作成功!", "info");
                    }else if(data=="0"){
                        jq.messager.alert('失败', "保存失败,请稍后重试!", "warning");
                    } else {
                       jq.messager.alert('失败', data, "warning");
                    }
                }
            });
        }
        function clearForm() {
            $('#win').window('close');
        }
        $(function () {
            $('#headimg').combobox({
                formatter: function (row) {
                    var imageFile = row.icon;
                    return '<img class="item-img" style="height:50px;width:50px;" src="' + imageFile + '"/><span class="item-text">' + row.text + '</span>';
                }
            });
        })
    </script>

</head>
<body>
    <div class="easyui-panel" title="个人信息" style="width:100%" fix="true" border="false">
        <div style="padding:10px 60px 20px 60px">
            <form id="form1" action="/adm/user/uinfopost" method="post">
                <table cellpadding="5">
                    <tr>
                        <td>账号：</td>
                        <td><input class="easyui-textbox" type="text" name="username" readonly="readonly" value="{{.m.username}}" style="width:200px;" data-options="required:true,missingMessage:'必填字段'" /></td>
                    </tr>
                    <tr>
                        <td>姓名：</td>
                        <td><input class="easyui-textbox" type="text" name="realname" value="{{.m.realname}}" style="width:200px;" data-options="required:true,missingMessage:'必填字段'" /></td>
                    </tr>
                    <tr>
                        <td>头像：</td>
                        <td>
                            <input id="headimg" style="width:200px" name="headimg"
                                   url="/adm/user/headimg"
                                   valuefield="icon" value="{{.m.headimg}}" textfield="text" editable="false">
                            </input>
                        </td>
                    </tr>
                    <tr>
                        <td></td>
                        <td><a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;</td>
                    </tr>
                </table>
            </form>
            <!--<div style="text-align:left;padding:80px 30px 80px 30px ">

                <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;
                <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
            </div>-->
        </div>
    </div>
</body>

`
var adm_user_pwd = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>修改密码</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css"> 
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <!--<script type="text/javascript" src="/js/jquery.form.js"></script>-->
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    <style>
        body {
            background: #fff;
        }
    </style>
    <script type="text/javascript">
        var jq = jQuery;

        function submitForm() {
            $('#form1').form('submit', {
                success: function (data) {
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                     }else if(data=="0"){
                         jq.messager.alert('失败', "保存失败,请稍后重试!", "warning");
                     } else {
                        jq.messager.alert('失败', data, "warning");
                     }
                }
            });
        }
        function clearForm() {
            $('#win').window('close');
        }
        $(function () {

        })
    </script>

</head>
<body>
    <div class="easyui-panel" title="修改密码" style="width:100%" fix="true" border="false">
        <div style="padding:10px 60px 20px 60px">
            <form id="form1" action="/adm/user/pwdpost" method="post">
                <table cellpadding="5">
                    <tr>
                        <td>用户名：</td>
                        <td><input class="easyui-textbox" type="text" name="username" readonly="readonly" value="{{.m.username}}" style="width:200px;" data-options="required:true,missingMessage:'必填字段'" /></td>
                    </tr>
                    <tr>
                        <td>密码：</td>
                        <td><input class="easyui-textbox" type="password" name="pwd1" value="" style="width:200px;" data-options="required:true,missingMessage:'必填字段'" /></td>
                    </tr>
                    <tr>
                        <td>确认密码:</td>
                        <td><input class="easyui-textbox" type="password" name="pwd2" value="" style="width:200px;" data-options="required:true,missingMessage:'必填字段'" /></td>
                    </tr>
                    <tr>
                        <td></td>
                        
                    </tr>
                    <tr>
                        <td></td>
                        <td><a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;</td>
                    </tr>
                </table>
            </form>
            <!--<div style="text-align:left;padding:80px 30px 80px 30px ">

                <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;
                <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
            </div>-->
        </div>
    </div>
</body>

`

var adm_user_list = `
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
	<script type="text/javascript" src="/adm/user/jsonutype"></script>
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
			usertype: $('#usertype').combobox("getValue"),
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:460,
				height:420,
				top:($(window).height() - 350) * 0.5,   
					left:($(window).width() - 680) * 0.5,
				modal:true,
				title:'{{.tb.title}}'+'[编辑账号]'
			});

            $('#win').window('open');
            $('#win').window('refresh', '/adm/user/edit?id='+row.id);
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
	var row = $('#tt').datagrid('getSelected');
	
	var w=$('#win').window({
		width:460,
		height:420,
		top:($(window).height() - 350) * 0.5,   
            left:($(window).width() - 680) * 0.5,
		modal:true,
		title:'{{.tb.title}}'+'[添加账户]'
	});
    w.window('open');
    w.window('refresh', '/adm/user/edit?id=0');
    $('#ff').form('load', row);
}
function doDel(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('Confirm', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/user/remove', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');
                    }else if(result=="0"){
                        jq.messager.alert('警告','删除失败!','warning');
                    } else {
                        jq.messager.alert('警告',result,'warning');
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
	function doMch(){
		top.addTab("企业管理","/adm/mch/list");
	}

	function rowformater_headimg(value, row, index) {
		//return "<span class=' "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</span>";
		return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
	function rowformater_date(value, row, index) {
       if (value == undefined) {
        return "";
		}
		return value;//dateValue.Format("yyyy-MM-dd hh:mm:ss");
    }
	function rowformater_detail(value, row, index) {
		return "<span ></span>";
	}
	function rowformater_usertype(value, row, index) {
			if(value == undefined){
				return '';
			}
			if(value==''){
				return '';
			}
			var v=value;
			if(jsonutype['key'+value]!=undefined){
			 value= jsonutype['key'+value];
			}
			return value;
	}
	function rowformater_company_id(value, row, index) {
			if(value == undefined){
				return '';
			}
			if(value==''){
				return '';
			}
			var v=value;
			if(jsoncompany_id['key'+value]!=undefined){
			 value= jsoncompany_id['key'+value];
			}
			return value;
	}
	function rowformater_company_pid(value, row, index) {
		if(row.company_id == undefined){
			return '';
		}
		if(row.company_id==''){
			return '';
		}
		var v=row.company_id;
		if(jsoncompany_id['pkey'+row.company_id]!=undefined){
		 value= jsoncompany_id['pkey'+row.company_id];
		}
		return value;
}
	function rowformater_state(value, row, index) {
		if(value=="0"){
			return "禁用";
		}
		if(value=="1"){
			return "<font color='green'>启用</font>";
		}
		if(value=="2"){
			return "<font color='red'>封停</font>";
		}
	}
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           title="用户管理" toolbar="#tb" id="tt" iconcls="icon-man"
           singleselect="true" pagination="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="30" sortable="true">编号</th>
				 <th field="headimg" align="center" width="50" data-options="formatter:rowformater_headimg">头像</th>  
				<th field="usertype" align="center" sortable="true" width="65" data-options="formatter:rowformater_usertype">类型</th>
				<th field="company_pid" align="center" width="90" sortable="true" data-options="formatter:rowformater_company_pid">上级单位</th> 
                <th field="company_id" align="center" width="90" sortable="true" data-options="formatter:rowformater_company_id">单位</th> 				            
                <th field="username" align="right" sortable="true" width="70">用户名</th>
				<th field="realname" align="right" sortable="true" width="90">姓名</th>

                <th field="logintime" width="100" data-options="formatter:rowformater_date">登录时间</th> 
				<th field="memo" width="50">备注</th>
				<th field="state" align="center" width="50" sortable="true" data-options="formatter:rowformater_state">状态</th>
				<th field=" " width="50" data-options="formatter:rowformater_detail">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-56" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-no" plain="true" onclick="doDel();">删除</a>
        </div>
		<div>
		用户类型: 
		<select  id="usertype" name="usertype" style="width:130px;" class="easyui-combobox" editable='false'>
		<option value="">请选择...</option>
		{{range $i,$row:=.usertype_list}}
		<option value="{{$row.id}}">{{$row.name}}</option>
		{{end}}
		</select>
            搜索: <input class="easyui-textbox" id="qtxt" prompt="请输入要检索的账号、姓名、单位等..." style="width:210px">


			<a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>&nbsp;
			{{if eq ._username "root"}}
			<a style="display:none;" href="#" class="easyui-linkbutton" iconcls="icon-43" onclick="doMch();">企业</a>
			{{end}}
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:460px;height:420px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>
<script type="text/javascript">
<!--
	$('#tt').datagrid({
        nowrap: false, 
        striped: true, 
        border: true, 
        collapsible:false,//是否可折叠的 
        fit: true,//自动大小 
        url:'/adm/user/listjson', 
        //sortName: 'usertype', 
        //sortOrder: 'asc', 
        remoteSort:true,  
        idField:'id', 
		pageSize:20,
		pageList:[20,50,100],
        singleSelect:true,//是否单选 
        pagination:true,//分页控件 
        
    }); 
//-->
</script>
</body>
</html>
`
var adm_user_role = `
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
            $('#win').window('refresh', '/adm/user/roleedit?id='+row.id);
			$('#win').window("resize",{top:$(document).scrollTop() + ($(window).height()-250) * 0.5});//居中显示
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
    $('#win').window('open');
    $('#win').window('refresh', '/adm/user/roleedit?id=');
    $('#ff').form('load', row);
}
function doRemove(){
    var jq=jQuery;
    //jq.messager.alert('warning',"操作成功!");
        // jq.messager.show({
        //     title:'温馨提示:',
        //     msg:'你好,我是从右下角弹出的窗体!',
        //     timeout:5000,
        //     showType:'slide'
    // });

    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/user/delrole', { id: row.id }, function (result) {
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
	function doAuth(){
		var row = $('#tt').datagrid('getSelected');
		if (row) {
			top.addTab('角色权限设置','/adm/user/roleauthlist?roleid='+row.id);
		} else {
			jq.messager.alert('警告','请选择一行数据','warning');
		}
		
	}
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/user/rolejson"
           title="角色管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           >
        <thead>
            <tr>
                <th field="id" sortable="true" width="20">ID</th>
				<th field="name" sortable="true" width="70">名称</th>
				<th field="level"  sortable="true" width="70">级别</th>
                <th field="info" sortable="true" align="right" width="70">说明</th>
				<th field="memo" width="50">备注</th>
				<th field="state" sortable="true" width="50">状态</th>
				<th field=" " width="50">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-56" plain="true" onclick="doAdd();">新建</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-15" plain="true" onclick="doAuth();">权限</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-no" plain="true" onclick="doRemove();">删除</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:80px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:460px;height:360px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>

</body>
</html>
`
var adm_user_roleedit = `

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
                        layer.msg('<font color="green">操作失败!</font>');
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
        <form id="form1" action="/adm/user/roleeditpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="name" value="{{.m.name}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
				</tr>
				<tr id="trdb">
                    <td>角色级别:</td>
                    <td>
                        <select id="level" name="level" style="width:180px;" class="easyui-combobox" editable="false">
							<option  value="">请选择级别...</option>
							{{range $k,$v :=.usertypelist}}
                            <option  value="{{$v.level}}">{{$v.name}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#level').combobox({
									onLoadSuccess: function () {
										$('#level').combobox('select','{{.m.level}}');
									}
								});	
                            });
                            
                        </script>
                    </td>
				</tr>
				<tr>
                    <td>权限:</td>
                    <td><input   type="text" name="rights" style="width:180px;" id="rights" data-options="method:'get',labelPosition:'top',multiple:true"></input>
					<script type="text/javascript">
						<!--
						/*
							$('#rights').combotree('loadData', [{
								id: 1,
								text: 'Languages',
								children: [{
									id: 11,
									text: 'Java'
								},{
									id: 12,
									text: 'C++'
								}]
							}]);
							*/
							
							$('#rights').combotree({
								panelWidth: 'auto',
								url: '/adm/user/treejson2?ids={{.m.rights}}',
								onCheck:function (item) {
									//alert(JSON.stringify(item));
								}
							});
						//-->
						</script>
					</td>
                </tr>
                <tr>
                    <td>说明:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="info" value="{{.m.info}}" ></input></td>
                </tr>
				
                <tr>
                    <td>备注:</td>
                    <td>
                        <input class="easyui-textbox" type="text" style="width:180px;" name="memo" value="{{.m.memo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>

                <tr>
                    <td>状态:</td>
                    <td>
                        <select id="state" class="easyui-combobox" name="state" editable="false" style="width:180px;" >
                            <option value="0">禁用</option>
                            <option value="1">启用</option>
                            <option value="2">封停</option>
                        </select>
						<script type="text/javascript">
						 
                            $('#state').combobox({
                                onLoadSuccess: function (data) {
                                    $('#state').combobox('setValue', "{{.m.state}}");
                                }
                            }); 
						</script>
                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`

var adm_user_edit = `

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
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        jq.messager.alert('错误', data, "warning");
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
        <form id="form1" action="/adm/user/editpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>账号:</td>
                    <td><input class="easyui-textbox" type="text" name="username" style="width:160px;" value="{{.uname}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
				</tr>
				{{if eq .is_sq "1"}}
				<tr>
					<td>上级:</td>
					<td><input class="easyui-combotree" name="pid" style="width:160px;" id="pid" data-options="method:'get',labelPosition:'top',multiple:false"></input>
					<script type="text/javascript">
							$('#pid').combotree({
								url: '/adm/user/usertreejson?id={{.m.id}}',
								onCheck:function (item) {
									//alert(JSON.stringify(item));
								},
								onLoadSuccess: function () {
									$('#pid').combotree('setValues',{{.m.pid}});
								}
							});
						</script>
					</td>
				</tr>
				<tr>
					<td>所属:</td>
					<td><input class="easyui-combotree" name="pids" style="width:160px;" id="pids" data-options="method:'get',labelPosition:'top',multiple:true,cascadeCheck:false"></input>
					<script type="text/javascript">
							$('#pids').combotree({
								url: '/adm/user/usertreejson?id={{.m.id}}',
								onCheck:function (item) {
									//alert(JSON.stringify(item));
								},
								onLoadSuccess: function () {
									$('#pids').combotree('setValues',eval('['+{{.m.pids}}+']'));
								}
							});
						</script>
					</td>
				</tr>
				{{end}}
                <tr>
                    <td>姓名:</td>
                    <td><input class="easyui-textbox" type="text" name="realname" style="width:160px;" value="{{.m.realname}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>电话:</td>
                    <td><input class="easyui-textbox" type="text" name="mobile" style="width:160px;" value="{{.m.mobile}}" ></input></td>
				</tr>
				{{if eq ._username "root"}}
                <tr style="display:none;">
                    <td>商户:</td>
                    <td>
                        <select id="mch_id" class="easyui-combobox" name="mch_id" style="width:160px;" data-options="required:true" editable="false">
                            {{range $k,$v:=.mchlist}}
                            <option value="{{$v.id}}">{{$v.mch_name}}</option>
                            {{end}}
                        </select>
						<script type="text/javascript">
						$(function(){
							$('#mch_id').combobox({
                                onLoadSuccess: function () {
								    //$('#mch_id').combobox('select','{{.m.mch_id}}');
							    }
                            });							
						})
						</script>
                    </td>
				</tr>
				{{end}}
                <tr>
                    <td>类型:</td>
                    <td>
                        <select id="usertype" class="easyui-combobox" name="usertype" style="width:160px;" data-options="required:true" editable="false">
						<option value="">请选择...</option>
						{{range $k,$v:=.utypelist}}
						<option value="{{$v.level}}">{{$v.name}}</option>
						{{end}}
                        </select>
						<script type="text/javascript">
							$(function(){
								$('#usertype').combobox({
									onLoadSuccess: function () {
										$('#usertype').combobox('select','{{.m.usertype}}');
									},
									onChange: function (n,o) {
										$('#company_id').combobox({
											url:'/adm/user/usertypecompany?cid={{._company_id}}&id='+$('#usertype').combobox('getValue'),
											valueField:'id',
											textField:'val',
											onLoadSuccess: function () {
												var v='{{.m.company_id}}';
												var ds=$('#company_id').combobox('getData');
												//console.log(ds);
												//console.log('------------------------------');
												for (var i = 0; i < ds.length; i++) {
													//console.log(ds[i]["val"]);
													if(ds[i]["id"]==v){
														$('#company_id').combobox('select','{{.m.company_id}}');
													}
												}
											}
										});
										//load rolelist
										$('#divrole').load('/xapi/rolehtmllist?level='+$('#usertype').combobox('getValue'),function(){
											{{.jstr}}
										});
									}
								});								
							})

						</script>
                    </td>
				</tr>
				<tr>
                    <td>单位:</td>
                    <td>
                        <select id="company_id" class="easyui-combobox" data-options="valueField:'id', textField:'val'" name="company_id" style="width:160px;"  editable="false">
						<option value="0">请选择...</option>
                        </select>
                        <script type="text/javascript">
							$('#company_id').combobox({
                                onLoadSuccess: function () {
									$('#company_id').combobox('select','{{.m.company_id}}');
								},
								onChange: function (n,o) {
									$('#company').val($('#company_id').combobox("getText"));
								}
                            });
						</script>
						<input type="hidden" id="company" name="company" vale="{{.m.company}}"/>
                    </td>
				</tr>
				<tr>
                    <td style="width:55px;">管理员:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="is_manager" title="" name="is_manager" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_manager}}'=='1'){
								$('#is_manager').switchbutton({
									checked: true,
								})
							}else{
								$('#is_manager').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
				<tr style="display:none;">
                    <td>级别:</td>
                    <td>
                        <select id="level" class="easyui-combobox" name="level" style="width:160px;"  editable="false">
                            <option select value="0">免费会员</option>
                            <option value="1">普通会员</option>
                            <option value="2">VIP会员</option>
                            <option value="3">超级VIP</option>
                        </select>
                        <script type="text/javascript">
							$('#level').combobox({
                                onLoadSuccess: function () {
								    $('#level').combobox('select','{{.m.level}}');
							    }
                            });
						</script>
                    </td>
                </tr>
                <tr>
                    <td>密码:</td>
                    <td><input class="easyui-textbox" type="text" name="password" style="width:160px;" value="{{.pwd}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>角色:</td>
					<td> 
						
						<div style="max-width:260px;" id="divrole">
						</div> 

					</td>
				</tr>
				{{if eq ._username "root"}}
				<tr>
                    <td style="width:55px;">超管:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="sproot" title="" name="sproot" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.sproot}}'=='1'){
								$('#sproot').switchbutton({
									checked: true,
								})
							}else{
								$('#sproot').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					 
					
					</td>
                </tr>
				{{end}}
				<tr>
                    <td>默认页:</td>
                    <td><input class="easyui-textbox" type="text" name="defpage" style="width:160px;" value="{{.m.defpage}}"></input></td>
                </tr>
                <tr>
                    <td>备注:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="memo" style="width:160px;" value="{{.m.memo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>
				
                <tr>
                    <td>状态:</td>
                    <td>
                        <select id="state" class="easyui-combobox" name="state" style="width:142px;" editable="false">
                            <option value="0">禁用</option>
                            <option value="1">启用</option>
                            <option value="2">封停</option>
                        </select>
						<script type="text/javascript">
							$('#state').combobox({
                                onLoadSuccess: function () {
								    $('#state').combobox('select','{{.m.state}}');
							    }
                            });
						</script>
                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`

//--------------------------------------------------------------------------------
//角色授权列表页面
func (c *UserController) RoleAuthList() {
	var roleid, _ = c.GetInt("roleid", 0)
	if roleid < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	c.Data["roleid"] = roleid

	var tpl = template.New("")
	tpl.Parse(adm_user_roleauthlist)
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

//获取角色权限列表
func (c *UserController) RoleAuthListJson() {
	var roleid, _ = c.GetInt("roleid", 0)
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = "where roleid=" + strconv.Itoa(roleid) + " "

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " and  `roleid` like '%" + qtxt + "%'"
	}
	//排序
	var sort = c.GetString("sort")
	var order = c.GetString("order")
	if sort != "" && order != "" {
		where += " order by " + sort + " " + order
	}

	var rst = db.Pager(page, pageSize, "select * from adm_role_auth "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//角色授权数据保存
func (c *UserController) RoleAuthFieldSet() {
	var id, _ = c.GetInt("id", 0)

	var m = db.First("select * from adm_role_auth where id=?", id)
	if m == nil {
		c.Ctx.WriteString("参数错误!")
		return
	}
	var f = c.GetString("f")
	var v = c.GetString("v")

	var rst = db.Exec("update adm_role_auth set "+f+"=? where id=?", v, id)
	if rst > 0 {
		c.Ctx.WriteString("<font color='yellow'>修改成功!</font>")
	} else {
		c.Ctx.WriteString("保存失败,请稍后重试!")
	}
}

var adm_user_roleauthlist = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>角色权限设置</title>
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
		return '';
    }
    function rowformater_name(value, row, index) {
		var e = '<div class="canedit" oldval="'+value+'" val="'+value+'" sp="field_name" valid="' + row.id + '" >'+value+'</div> ';
        return value;
    }

    function rowformater_add(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_add\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_add'+row.id+'" id="ac_add'+row.id+'">'
        return e;
    }
    function rowformater_del(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_del\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_del'+row.id+'" id="ac_del'+row.id+'">'
        return e;
	}
	function rowformater_query(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_query\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_query'+row.id+'" id="ac_query'+row.id+'">'
        return e;
	}
	function rowformater_update(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_update\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_update'+row.id+'" id="ac_update'+row.id+'">'
        return e;
	}
	function rowformater_import(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_import\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_import'+row.id+'" id="ac_import'+row.id+'">'
        return e;
	}
	function rowformater_export(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_export\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_export'+row.id+'" id="ac_export'+row.id+'">'
        return e;
	}
	function rowformater_print(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_print\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_print'+row.id+'" id="ac_print'+row.id+'">'
        return e;
    }
	function rowformater_edit(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_edit\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_edit'+row.id+'" id="ac_edit'+row.id+'">'
        return e;
    }
    function rowformater_info(value, row, index) {
        var ischecked=''
        if(value==1){
            ischecked='checked';
        }
        var e='<input class="easyui-checkbox" onchange="funField(\'ac_info\','+row.id+',this.checked)" '+ischecked+' type="checkbox" name="ac_info'+row.id+'" id="ac_info'+row.id+'">'
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
           url="/adm/user/roleauthlistjson?roleid={{.roleid}}"  
           title="角色权限设置" toolbar="#tb" id="tt"
           singleselect="true" pageSize="50" pageList="[20, 50, 100]" pagination="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="15" sortable="true">ID</th>
				<th field="roleid" align="center"  width="15" sortable="true">角色ID</th>
				<th field="rolename"  width="20" sortable="true">角色</th>
				<th field="menuid" align="center" sortable="true" width="15" >权限ID</th>
                <th field="menuname" align="center" sortable="true" data-options="formatter:rowformater_name" width="20" >权限</th>
                <th field="ac_add"  width="10"  data-options="formatter:rowformater_add" >新增</th> 
				<th field="ac_del"  width="10"  data-options="formatter:rowformater_del" >删除</th> 
				<th field="ac_query"  width="10"  data-options="formatter:rowformater_query" >查询</th> 
				<th field="ac_update"  width="10"  data-options="formatter:rowformater_update" >修改</th> 
				<th field="ac_info"  width="10"  data-options="formatter:rowformater_info" >详情</th> 
				<th field="ac_edit"  width="10"  data-options="formatter:rowformater_edit" >编辑</th> 
				<th field="ac_import"  width="10"  data-options="formatter:rowformater_import" >导入</th> 
				<th field="ac_export"  width="10"  data-options="formatter:rowformater_export" >导出</th>
				<th field="ac_print"  width="10"  data-options="formatter:rowformater_print" >打印</th>
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
            jQuery.post('/adm/user/roleauthfieldset', { 'tb': 'adm_role_auth', 'id': id, 'f': f, 'v': v }, function (data) {
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

                                jQuery.post('/adm/user/roleauthfieldset',{'tb':'adm_role_auth','id':valid,'f':field,'v':$(this).val()},function(data){
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

//账号类型列表页面
func (c *UserController) UserTypeList() {
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_usertypelist)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))
}
func (c *UserController) UserTypeListJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " and `name` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select id,level,mch_id,name,orders,sysid,state from adm_usertype where 1=1 "+where)

	c.Data["json"] = rst
	c.ServeJSON()
}

//根据账号类型,显示账号类型的企业列表 支持跨库
func (c *UserController) UserTypeCompany() {
	//var rst = ""
	//var cid, _ = c.GetInt("cid", 0)
	var _usertype = c.GetSession("_usertype").(string)
	var id, _ = c.GetInt("id", 0)
	var m = db.First("select * from adm_usertype where level=?", id)
	if len(m) < 1 {
		// var rst = `var jsoncompay_id={ `
		// rst += `"key0":"---"`
		// rst += `};`

		c.Data["json"] = `
		[
			{
				"id": "0",
				"level": "1",
				"pid": "0",
				"pname": "-",
				"val": "请选择..."
			}
		]
		`
		c.ServeJSON()
		return
	}

	//单位绑定 绑定字段为 id val  从数据库中读取
	if m["conn_str"] != "" && m["bindapi"] != "" {
		var xdb = db.NewDb(m["conn_str"])
		//如果等级一样,只显示自己
		if _usertype == strconv.Itoa(id) && c.GetSession("_username").(string) != "root" {
			m["bindapi"] += " and id=" + c.GetSession("_company_id").(string)
		} else {
			var _mch_id = c.GetSession("_mch_id").(string)
			//如果是管理员,展示自己级别及以下信息,不包括自己
			var ids = ChildIds2(_mch_id, c.GetSession("_company_id").(string))
			ids = strings.Replace(ids, ",,", ",", -1)
			ids = strings.Replace(ids, ",,", ",", -1)
			ids = strings.Replace(ids, ",,", ",", -1)
			m["bindapi"] += " and id in(" + ids + ") "
		}
		var list = db.Query2(xdb, m["bindapi"])
		fmt.Println("--------bindapi:", m["bindapi"])
		c.Data["json"] = list
		c.ServeJSON()
	} else {
		c.Data["json"] = `
		[
			{
				"id": "0",
				"level": "1",
				"pid": "0",
				"pname": "-",
				"val": "请选择..."
			}
		]
		`
		c.ServeJSON()
	}

	// //单位绑定 绑定字段为 id val  从数据库中读取
	// var jsonstr = `var jsoncompay_id={ `
	// if m["conn_str"] != "" && m["bindapi"] != "" {
	// 	var xdb = db.NewDb(m["conn_str"])
	// 	var list = db.Query2(xdb, m["bindapi"])
	// 	for kk, vv := range list {
	// 		if kk > 0 {
	// 			jsonstr += ","
	// 		}
	// 		jsonstr += `"key` + vv["id"] + `":"` + vv["val"] + `"`
	// 	}
	// }
	// jsonstr += `};`

	// rst = jsonstr

	//// if m["bindapi"] == "" {
	//// 	rst = "{}"
	//// 	c.Data["json"] = rst
	//// 	c.ServeJSON()
	//// 	return
	//// }

	//// var url = m["bindapi"]
	//// if strings.Contains(url, "http") == false {
	//// 	if strings.Contains(url, "?") {
	//// 		url += "&cid=" + strconv.Itoa(cid)
	//// 	} else {
	//// 		url += "?cid=" + strconv.Itoa(cid)
	//// 	}
	//// 	url = "http://" + c.Ctx.Input.Domain() + ":" + strconv.Itoa(c.Ctx.Input.Port()) + url
	//// }
	//// fmt.Println("接口url:", url)
	//// rst, _ = httpGet(url)
	//// if rst == "" {
	//// 	rst = "{}"
	//// }

	//c.Ctx.Output.Header("Content-Type", "text/json; charset=utf-8")
	//c.Ctx.Output.Body([]byte(rst))
}
func httpGet(url string) (string, error) {
	postReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("请求失败", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		fmt.Println("client请求失败", err)
		return "", err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(data), err
}

//账号类型编辑
func (c *UserController) UserTypeEdit() {
	var id, _ = c.GetInt("id", 0)
	if id > 0 {
		var m = db.First("select * from adm_usertype where  id=?", id)
		c.Data["m"] = m
	}
	//数据库链接
	var dblist = db.Query("select * from adm_conn where state=1")
	if dblist != nil {
		c.Data["dblist"] = dblist
	}
	c.Data["sproot"] = c.GetSession("_sproot").(string)
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_usertypeedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!")
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))
}
func (c *UserController) UserTypeEditPost() {
	var id, _ = c.GetInt("id", 0)
	var name = c.GetString("name")
	var orders, _ = c.GetInt("orders", 0)
	var state = c.GetString("state")
	var conn_str = c.GetString("conn_str")
	var bindapi = c.GetString("bindapi")
	if state == "on" || state == "1" {
		state = "1"
	} else {
		state = "0"
	}

	var sql = ""
	if id > 0 {
		var m = db.First("select * from adm_usertype where id=?", id)
		if len(m) < 1 {
			c.Ctx.WriteString("0")
			return
		}
		if c.GetSession("_sproot").(string) != "1" { //只有超管才能修改绑定接口
			bindapi = m["bindapi"]
		}
		sql = `
		update adm_usertype set 
		name=?,
		orders=?,
		conn_str=?,
		bindapi=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			name,
			orders,
			conn_str,
			bindapi,
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
		insert into adm_usertype(
			name,
			orders,
			conn_str,
			bindapi,
			state
		)values(
			?,?,?,?,?
		)
		`
		var i = db.Exec(sql,
			name,
			orders,
			conn_str,
			bindapi,
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

//账号类型删除
func (c *UserController) UserTypeDel() {
	var id, _ = c.GetInt("id", 0)
	if id < 2 {
		c.Ctx.WriteString("0")
		return
	}
	var i = db.Exec("delete from adm_usertype where  id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

var adm_user_usertypelist = `
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
			var w=$('#win').window({
					width:480,
					height:380,
					top:($(window).height() - 350) * 0.5,   
						left:($(window).width() - 680) * 0.5,
					modal:true
			});

            w.window('open');
            w.window('refresh', '/adm/user/usertypeedit?id='+row.id);
			w.window("resize",{top:$(document).scrollTop() + ($(window).height()-250) * 0.5});//居中显示
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
	var row = $('#tt').datagrid('getSelected');
	var w=$('#win').window({
		width:480,
		height:380,
		top:($(window).height() - 350) * 0.5,   
			left:($(window).width() - 680) * 0.5,
		modal:true
	});
    w.window('open');
    w.window('refresh', '/adm/user/usertypeedit?id=');
    $('#ff').form('load', row);
}
function doDel(){
    var jq=jQuery;
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/user/usertypedel', { id: row.id }, function (result) {
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
	
	function rowformater_state(value, row, index) {
		if(value=="0"){
			return "禁用";
		}else if(value=="1"){
			return "启用";
		}else{
			return value;
		}
	}
	 
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/user/usertypelistjson"
           title="账号类别管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           >
        <thead>
            <tr>
				<th field="id" width="20">ID</th>
				<th field="name" width="70">名称</th>
				<th field="orders" width="50"  >排序</th>
                
				<th field="state" width="50" formatter="rowformater_state">状态</th>
				<th field=" " width="50">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-56" plain="true" onclick="doAdd();">新建</a>
			<a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-no" plain="true" onclick="doDel();">删除</a>
        </div>
        <div>
            
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:80px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:420px;height:420px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>

</body>
</html>
`

var adm_user_usertypeedit = `

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
                        layer.msg('<font color="green">操作失败!</font>');
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
        <form id="form1" action="/adm/user/usertypeeditpost" method="post">
			<table cellpadding="5">
			<tr>
                    <td>ID:</td>
					<td>
                        {{.m.id}}
                    </td>
                </tr>
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="name" value="{{.m.name}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
				<tr>
                    <td>排序:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="orders" value="{{.m.orders}}"  ></input></td>
				</tr>
				{{if eq .sproot "1"}}
				<tr id="trdb">
                    <td>绑定库:</td>
                    <td>
                        <select id="conn_str" name="conn_str" style="width:180px;" class="easyui-combobox" editable="false">
							<option  value="">请选择数据库...</option>
							{{range $k,$v :=.dblist}}
                            <option  value="{{$v.conn}}">{{$v.title}}</option>
                            {{end}}
                        </select>
                        <script language="javascript">
                            $(function(){
                                $('#conn_str').combobox({
									onLoadSuccess: function () {
										$('#conn_str').combobox('select','{{.m.conn_str}}');
									}
								});	
                            });
                            
                        </script>
                    </td>
				</tr>
				<tr>
                    <td>绑定数据:</td>
					<td><input class="easyui-textbox" type="text" style="width:180px;" name="bindapi" value="{{.m.bindapi}}"></input>
					</br>使用id,val两个字段,只支持sql语句
					</td>
				</tr>
				{{end}}
                <tr>
                    <td>状态:</td>
					<td>
					<input type="hidden" id="id" name="id" value="{{.m.id}}" />
                        <select id="state" class="easyui-combobox" name="state" editable="false" style="width:180px;" >
                            <option value="0">禁用</option>
							<option value="1">启用</option>  
                        </select>
						<script type="text/javascript">
						 
                            $('#state').combobox({
                                onLoadSuccess: function (data) {
                                    $('#state').combobox('setValue', "{{.m.state}}");
                                }
                            }); 
						</script>
                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
`

//-----------------------------------------------------------------------------------------------------------

//List 列表页面
func (c *UserController) UList() {
	c.Data["_username"] = c.GetSession("_username").(string)
	//账号类型信息
	var sql = `select * from adm_usertype `
	c.Data["usertype_list"] = db.Query(sql)

	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_user_ulist)
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

//获取用户列表
func (c *UserController) UListJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where (`username` like '%" + qtxt + "%' or realname like '%" + qtxt + "%' or company like '%" + qtxt + "%') "
	}

	var usertype = c.GetString("usertype")
	if usertype != "0" && usertype != "" {
		if where != "" {
			where += " and usertype='" + usertype + "' "
		} else {
			where += " where  usertype='" + usertype + "' "
		}

	}

	//排序
	var sort = c.GetString("sort")
	var order = c.GetString("order")
	if sort != "" && order != "" {
		where += " order by " + sort + " " + order
	} else {
		where += " order by usertype,company_pid,company_id "
	}

	//fmt.Println("where:", where)
	var rst = db.Pager(page, pageSize, "select *  from adm_user "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//Edit 用户编辑页面
func (c *UserController) UEdit() {
	c.Data["is_sq"] = "0"

	//读取子系统列表
	var syslist = db.Query("select * from adm_menu where sysid=0 and pid=0")
	c.Data["syslist"] = syslist

	var _company_id = "0"
	if c.GetSession("_company_id") != nil {
		_company_id = c.GetSession("_company_id").(string) //单位ID
	}
	c.Data["_company_id"] = _company_id

	var id, _ = c.GetInt("id", 0)
	var m = db.First("select * from adm_user where id=?", id)
	c.Data["m"] = m
	if len(m) > 0 {
		c.Data["uname"] = m["username"]
		c.Data["pwd"] = m["password"]
	} else {
		c.Data["uname"] = ""
		c.Data["pwd"] = "358719"
	}
	//公司列表
	var mchlist = db.Query("select * from adm_mch")
	c.Data["mchlist"] = mchlist
	//角色列表
	var where = ""

	var roles = db.Query("select * from adm_role" + where)
	c.Data["roles"] = roles
	var roleids = ""
	for i, v := range roles {
		if i > 0 {
			roleids += ","
		}
		roleids += v["id"]
	}
	c.Data["roleids"] = roleids
	//账号类型列表 根据级别过滤
	var utypelist = db.Query("select * from adm_usertype  order by orders")
	c.Data["utypelist"] = utypelist
	//根据信息选择已有权限
	var jstr = ""
	if m != nil {
		var r = strings.Split(m["roles"], ",")
		for i := 0; i < len(r); i++ {
			jstr += `$('#role` + r[i] + `').prop('checked',true);`
		}
	}
	c.Data["jstr"] = template.JS(jstr)

	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_user_uedit")
	tpl.Funcs(template.FuncMap{"str2html": beego.Str2html})
	tpl.Parse(adm_user_uedit)
	var buf bytes.Buffer
	var e = tpl.Execute(&buf, c.Data)

	if e != nil {
		fmt.Println("tpl.Execute 错误:", e.Error())
		c.Ctx.WriteString("页面模板错误!" + e.Error())
		return
	}
	var rst = buf.String()

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body([]byte(rst))

	c.Ctx.WriteString(rst)
}

var adm_user_ulist = `
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
	<script type="text/javascript" src="/adm/user/jsonutype"></script>
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
			usertype: $('#q_usertype').combobox("getValue"),
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:460,
				height:420,
				top:($(window).height() - 350) * 0.5,   
					left:($(window).width() - 680) * 0.5,
				modal:true,
				title:'{{.tb.title}}'+'[编辑账号]'
			});

            $('#win').window('open');
            $('#win').window('refresh', '/adm/user/uedit?id='+row.id);
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
	var row = $('#tt').datagrid('getSelected');
	
	var w=$('#win').window({
		width:460,
		height:420,
		top:($(window).height() - 350) * 0.5,   
            left:($(window).width() - 680) * 0.5,
		modal:true,
		title:'{{.tb.title}}'+'[添加账户]'
	});
    w.window('open');
    w.window('refresh', '/adm/user/uedit?id=0');
    $('#ff').form('load', row);
}
function doDel(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('Confirm', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/user/remove', { id: row.id }, function (result) {
                    if (result=="1") {
                        $('#tt').datagrid('reload');
                    }else if(result=="0"){
                        jq.messager.alert('警告','删除失败!','warning');
                    } else {
                        jq.messager.alert('警告',result,'warning');
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
	function doMch(){
		top.addTab("企业管理","/adm/mch/list");
	}

	function rowformater_headimg(value, row, index) {
		//return "<span class=' "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</span>";
		return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
	function rowformater_date(value, row, index) {
       if (value == undefined) {
        return "";
		}
		return value;//dateValue.Format("yyyy-MM-dd hh:mm:ss");
    }
	function rowformater_detail(value, row, index) {
		return "<span ></span>";
	}
	function rowformater_usertype(value, row, index) {
			if(value == undefined){
				return '';
			}
			if(value==''){
				return '';
			}
			var v=value;
			if(jsonutype['key'+value]!=undefined){
			 value= jsonutype['key'+value];
			}
			return value;
	}
	function rowformater_company_id(value, row, index) {
			if(value == undefined){
				return '';
			}
			if(value==''){
				return '';
			}
			var v=value;
			if(jsoncompany_id['key'+value]!=undefined){
			 value= jsoncompany_id['key'+value];
			}
			return value;
	}
	function rowformater_company_pid(value, row, index) {
		if(row.company_id == undefined){
			return '';
		}
		if(row.company_id==''){
			return '';
		}
		var v=row.company_id;
		if(jsoncompany_id['pkey'+row.company_id]!=undefined){
		 value= jsoncompany_id['pkey'+row.company_id];
		}
		return value;
}
	function rowformater_state(value, row, index) {
		if(value=="0"){
			return "禁用";
		}
		if(value=="1"){
			return "<font color='green'>启用</font>";
		}
		if(value=="2"){
			return "<font color='red'>封停</font>";
		}
	}
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           title="用户管理" toolbar="#tb" id="tt" iconcls="icon-man"
           singleselect="true" pagination="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="30" sortable="true">编号</th>
				 <th field="headimg" align="center" width="50" data-options="formatter:rowformater_headimg">头像</th>  
				<th field="usertype" align="center" sortable="true" width="65" data-options="formatter:rowformater_usertype">类型</th>
				<th field="company_pid" align="center" width="90" sortable="true" data-options="formatter:rowformater_company_pid">上级单位</th> 
                <th field="company_id" align="center" width="90" sortable="true" data-options="formatter:rowformater_company_id">单位</th> 				            
                <th field="username" align="right" sortable="true" width="70">用户名</th>
				<th field="realname" align="right" sortable="true" width="90">姓名</th>

                <th field="logintime" width="100" data-options="formatter:rowformater_date">登录时间</th> 
				<th field="memo" width="50">备注</th>
				<th field="state" align="center" width="50" sortable="true" data-options="formatter:rowformater_state">状态</th>
				<th field=" " width="50" data-options="formatter:rowformater_detail">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-56" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-no" plain="true" onclick="doDel();">删除</a>
        </div>
		<div>
		用户类型: 
		<select  id="q_usertype" name="usertype" style="width:130px;" class="easyui-combobox" editable='false'>
		<option value="">请选择...</option>
		{{range $i,$row:=.usertype_list}}
		<option value="{{$row.id}}">{{$row.name}}</option>
		{{end}}
		</select>
            搜索: <input class="easyui-textbox" id="qtxt" prompt="请输入要检索的账号、姓名、单位等..." style="width:210px">


			<a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>&nbsp;
			{{if eq ._username "root"}}
			<a style="display:none;" href="#" class="easyui-linkbutton" iconcls="icon-43" onclick="doMch();">企业</a>
			{{end}}
        </div>
    </div>

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:460px;height:420px;padding:5px;overflow-x: hidden;">
        Some Content.
    </div>
<script type="text/javascript">
<!--
	$('#tt').datagrid({
        nowrap: false, 
        striped: true, 
        border: true, 
        collapsible:false,//是否可折叠的 
        fit: true,//自动大小 
        url:'/adm/user/ulistjson', 
        //sortName: 'usertype', 
        //sortOrder: 'asc', 
        remoteSort:true,  
        idField:'id', 
		pageSize:20,
		pageList:[20,50,100],
        singleSelect:true,//是否单选 
        pagination:true,//分页控件 
        
    }); 
//-->
</script>
</body>
</html>
`

var adm_user_uedit = `

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
                    if (data == "1") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        jq.messager.alert('错误', data, "warning");
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
        <form id="form1" action="/adm/user/editpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>账号:</td>
                    <td><input class="easyui-textbox" type="text" name="username" style="width:160px;" value="{{.uname}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
				</tr>
				<tr>
					<td>系统:</td>
					<td>
						<select id="sysid" class="easyui-combobox" name="sysid" style="width:160px;" data-options="required:true" editable="false">
						<option value="0">全部系统</option>	
						{{range $k,$v:=.syslist}}
                            <option value="{{$v.id}}">{{$v.title}}</option>
                            {{end}}
						</select>
						<script type="text/javascript">
							$('#sysid').combobox({
                                onLoadSuccess: function () {
								    $('#sysid').combobox('select','{{.m.sysid}}');
							    }
                            });
						</script>
					</td>
				</tr>
				{{if eq .is_sq "1"}}
				<tr>
					<td>上级:</td>
					<td><input class="easyui-combotree" name="pid" style="width:160px;" id="pid" data-options="method:'get',labelPosition:'top',multiple:false"></input>
					<script type="text/javascript">
							$('#pid').combotree({
								url: '/adm/user/usertreejson?id={{.m.id}}',
								onCheck:function (item) {
									//alert(JSON.stringify(item));
								},
								onLoadSuccess: function () {
									$('#pid').combotree('setValues',{{.m.pid}});
								}
							});
						</script>
					</td>
				</tr>
				<tr>
					<td>所属:</td>
					<td><input class="easyui-combotree" name="pids" style="width:160px;" id="pids" data-options="method:'get',labelPosition:'top',multiple:true,cascadeCheck:false"></input>
					<script type="text/javascript">
							$('#pids').combotree({
								url: '/adm/user/usertreejson?id={{.m.id}}',
								onCheck:function (item) {
									//alert(JSON.stringify(item));
								},
								onLoadSuccess: function () {
									$('#pids').combotree('setValues',eval('['+{{.m.pids}}+']'));
								}
							});
						</script>
					</td>
				</tr>
				{{end}}
                <tr>
                    <td>姓名:</td>
                    <td><input class="easyui-textbox" type="text" name="realname" style="width:160px;" value="{{.m.realname}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>电话:</td>
                    <td><input class="easyui-textbox" type="text" name="mobile" style="width:160px;" value="{{.m.mobile}}" ></input></td>
				</tr>
				
                <tr style="display:none;">
                    <td>商户:</td>
                    <td>
                        <select id="mch_id" class="easyui-combobox" name="mch_id" style="width:160px;" data-options="required:true" editable="false">
                            {{range $k,$v:=.mchlist}}
                            <option value="{{$v.id}}">{{$v.mch_name}}</option>
                            {{end}}
                        </select>
						<script type="text/javascript">
						$(function(){
							$('#mch_id').combobox({
                                onLoadSuccess: function () {
								    //$('#mch_id').combobox('select','{{.m.mch_id}}');
							    }
                            });							
						})
						</script>
                    </td>
				</tr>
				
                <tr>
                    <td>类型:</td>
                    <td>
                        <select id="usertype" class="easyui-combobox" name="usertype" style="width:160px;" data-options="required:true" editable="false">
						<option value="">请选择...</option>
						{{range $k,$v:=.utypelist}}
						<option value="{{$v.level}}">{{$v.name}}</option>
						{{end}}
                        </select>
						<script type="text/javascript">
							$(function(){
								$('#usertype').combobox({
									onLoadSuccess: function () {
										$('#usertype').combobox('select','{{.m.usertype}}');
									},
									onChange: function (n,o) {
										$('#company_id').combobox({
											url:'/adm/user/usertypecompany?cid={{._company_id}}&id='+$('#usertype').combobox('getValue'),
											valueField:'id',
											textField:'val',
											onLoadSuccess: function () {
												var v='{{.m.company_id}}';
												var ds=$('#company_id').combobox('getData');
												//console.log(ds);
												//console.log('------------------------------');
												for (var i = 0; i < ds.length; i++) {
													//console.log(ds[i]["val"]);
													if(ds[i]["id"]==v){
														$('#company_id').combobox('select','{{.m.company_id}}');
													}
												}
											}
										});
										//load rolelist
										$('#divrole').load('/xapi/rolehtmllist?level='+$('#usertype').combobox('getValue'),function(){
											{{.jstr}}
										});
									}
								});	
								$('#usertype').combobox('select','{{.m.usertype}}');							
							})

						</script>
                    </td>
				</tr>
				<tr>
                    <td>单位:</td>
                    <td>
                        <select id="company_id" class="easyui-combobox" data-options="valueField:'id', textField:'val'" name="company_id" style="width:160px;"  editable="false">
						<option value="0">请选择...</option>
                        </select>
                        <script type="text/javascript">
							$('#company_id').combobox({
                                onLoadSuccess: function () {
									$('#company_id').combobox('select','{{.m.company_id}}');
								},
								onChange: function (n,o) {
									$('#company').val($('#company_id').combobox("getText"));
								}
                            });
						</script>
						<input type="hidden" id="company" name="company" vale="{{.m.company}}"/>
                    </td>
				</tr>
				<tr>
                    <td style="width:55px;">管理员:</td>
                    <td>
					
					<input class="easyui-switchbutton" id="is_manager" title="" name="is_manager" style="vertical-align:middle;">
					<script type="text/javascript">
						$(function(){
							if('{{.m.is_manager}}'=='1'){
								$('#is_manager').switchbutton({
									checked: true,
								})
							}else{
								$('#is_manager').switchbutton({
									checked: false,
								})
							}
						})
					</script>
					</td>
                </tr>
				<tr style="display:none;">
                    <td>级别:</td>
                    <td>
                        <select id="level" class="easyui-combobox" name="level" style="width:160px;"  editable="false">
                            <option select value="0">免费会员</option>
                            <option value="1">普通会员</option>
                            <option value="2">VIP会员</option>
                            <option value="3">超级VIP</option>
                        </select>
                        <script type="text/javascript">
							$('#level').combobox({
                                onLoadSuccess: function () {
								    $('#level').combobox('select','{{.m.level}}');
							    }
                            });
						</script>
                    </td>
                </tr>
                <tr>
                    <td>密码:</td>
                    <td><input class="easyui-textbox" type="text" name="password" style="width:160px;" value="{{.pwd}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>角色:</td>
					<td> 
						
						<div style="max-width:260px;" id="divrole">
						</div> 

					</td>
				</tr>
				
				<tr>
                    <td>默认页:</td>
                    <td><input class="easyui-textbox" type="text" name="defpage" style="width:160px;" value="{{.m.defpage}}"></input></td>
                </tr>
                <tr>
                    <td>备注:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="memo" style="width:160px;" value="{{.m.memo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.m.id}}" />
                    </td>
                </tr>
				
                <tr>
                    <td>状态:</td>
                    <td>
                        <select id="state" class="easyui-combobox" name="state" style="width:142px;" editable="false">
                            <option value="0">禁用</option>
                            <option value="1">启用</option>
                            <option value="2">封停</option>
                        </select>
						<script type="text/javascript">
							$('#state').combobox({
                                onLoadSuccess: function () {
								    $('#state').combobox('select','{{.m.state}}');
							    }
                            });
						</script>
                    </td>
                </tr>
            </table>
        </form>
        <div style="text-align:center;padding:5px">

            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-ok" id="btnsave" onclick="submitForm()">保 存&nbsp;</a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>


`
