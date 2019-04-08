package adm

import (
	"bytes"
	"html/template"
	"strconv"

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
	//c.TplName="adm/user/list.html"
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
		where += " where `username` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select * from adm_user "+where)
	//fmt.Println(rst)

	c.Data["json"] = rst
	c.ServeJSON()
}

//Edit 用户编辑页面
func (c *UserController) Edit() {
	var id, _ = c.GetInt("id", 0)
	var m = db.First("select * from adm_user where id=?", id)
	c.Data["m"] = m
	//公司列表
	var mchlist = db.Query("select * from adm_mch")
	c.Data["mchlist"] = mchlist
	//角色列表
	var roles = db.Query("select * from adm_role")
	c.Data["roles"] = roles
	//根据信息选择已有权限
	var jstr = ""
	if m != nil {
		var r = strings.Split(m["roles"], ",")
		for i := 0; i < len(r); i++ {
			jstr += `$('#role` + r[i] + `').attr('checked',true);`
		}
	}
	c.Data["jstr"] = template.JS(jstr)

	//c.TplName="adm/user/edit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_user_edit")
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

//用户信息编辑
func (c *UserController) EditPost() {
	var id, _ = c.GetInt("id", 0)
	var username = c.GetString("username")
	var realname = c.GetString("realname")
	var mobile = c.GetString("mobile")
	var usertype, _ = c.GetInt("usertype", 2)
	var level, _ = c.GetInt("level", 1)
	var password = c.GetString("password")
	var state = c.GetString("state")

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
		var sql = `
		update adm_user set 
		username=?,
		realname=?,
		mobile=?,
		usertype=?,
		level=?,
		password=?,
		state=?,
		roles=?,
		memo=?
		where id=?
		`
		var i = db.Exec(sql,
			username,
			realname,
			mobile,
			usertype,
			level,
			password,
			state,
			roles,
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
		var sql = `insert into adm_user(
			username,
			realname,
			mobile,
			usertype,
			level,
			password,
			state,
			roles,
			memo
		)values(?,?,?,?,?,?,?,?,?)
		`
		var i = db.Exec(sql,
			username,
			realname,
			mobile,
			usertype,
			level,
			password,
			state,
			roles,
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
	//c.TplName="adm/user/role.html"
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
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where `name` like '%" + qtxt + "%'"
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

	var sql = ""
	if id > 0 {
		sql = `
		update adm_role set 
		name=?,
		rights=?,
		info=?,
		memo=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			name,
			rights,
			info,
			memo,
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

			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	} else {
		sql = `
		insert into adm_role(
			name,
			rights,
			info,
			memo,
			state
		)values(
			?,?,?,?,?
		)
		`
		var i = db.Exec(sql,
			name, rights,
			info, memo, state,
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
		rst += `"text":"` + v["title"] + v["id"] + `"`
		//第二层节点
		var list1 = db.Query("select * from adm_menu where pid=?", v["id"])
		rst1 := "["
		for kk, vv := range list1 {
			if kk > 0 {
				rst1 += ","
			}
			rst1 += "{"
			rst1 += `"id":` + vv["id"] + ","
			rst1 += `"text":"` + vv["title"] + vv["id"] + `"`
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
				rst2 += `"text":"` + vvv["title"] + vvv["id"] + `"`
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
function doRemove(){
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


	function rowformater_headimg(value, row, index) {
		//return "<span class=' "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</span>";
		return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
	function rowformater_date(value, row, index) {
       if (value == undefined) {
        return "";
		}
		/*json格式时间转js时间格式*/
		// value = value.substr(1, value.length - 2);
		// var obj = eval('(' + "{Date: new " + value + "}" + ')');
		// var dateValue = obj["Date"];
		// if (dateValue.getFullYear() < 1900) {
		// 	return "";
		// }

		return value;//dateValue.Format("yyyy-MM-dd hh:mm:ss");
    }
	function rowformater_detail(value, row, index) {
		return "<span >详情</span>";
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
                <th field="username" align="right" sortable="true" width="70">用户名</th>
				<th field="realname" align="right" sortable="true" width="70">姓名</th>
                <th field="usertype" sortable="true" width="35">类型</th>
                <th field="level" width="50" sortable="true">级别</th>
				<th field="roles" width="50">角色</th>
				<th field="loginip" width="80">登录IP</th>
                <th field="logintime" width="100" data-options="formatter:rowformater_date">登录时间</th>
				<th field="regtime" width="100" data-options="formatter:rowformater_date">注册时间</th>
				<th field="memo" width="50">备注</th>
				<th field="state" width="50" sortable="true">状态</th>
				<th field=" " width="50" data-options="formatter:rowformater_detail">操作</th>
            </tr>
        </thead>
    </table>

    <div id="tb" style="padding:5px;height:auto">
        <div style="margin-bottom:5px">
            <a href="#" class="easyui-linkbutton" iconcls="icon-56" plain="true" onclick="doAdd();">新建</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-1" plain="true" onclick="doEdit();">编辑</a>
            <a href="#" class="easyui-linkbutton" iconcls="icon-no" plain="true" onclick="doRemove();">删除</a>
        </div>
        <div>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:80px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
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
        sortName: 'id', 
        sortOrder: 'asc', 
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
                <th field="id" width="20">ID</th>
                <th field="name" width="70">名称</th>
                <th field="info" align="right" width="70">说明</th>
				<th field="memo" width="50">备注</th>
				<th field="state" width="50">状态</th>
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

    <div id="win" class="easyui-window" title="编辑信息" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:420px;height:320px;padding:5px;overflow-x: hidden;">
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
				<tr>
                    <td>权限:</td>
                    <td><input   type="text" name="rights" style="width:180px;" id="rights" data-options="method:'get',label:'Select Node:',labelPosition:'top',multiple:true"></input>
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
								url: '/adm/user/treejson?ids={{.m.rights}}',
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
                    if (data != "0") {
                        jq.messager.alert('成功', "操作成功!", "info");
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        jq.messager.alert('错误', "操作失败!", "warning");
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
                    <td><input class="easyui-textbox" type="text" name="username" style="width:160px;" value="{{.m.username}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
                <tr>
                    <td>姓名:</td>
                    <td><input class="easyui-textbox" type="text" name="realname" style="width:160px;" value="{{.m.realname}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>电话:</td>
                    <td><input class="easyui-textbox" type="text" name="mobile" style="width:160px;" value="{{.m.mobile}}" ></input></td>
                </tr>
                <tr>
                    <td>企业:</td>
                    <td>
                        <select id="mch_id" class="easyui-combobox" name="mch_id" style="width:160px;" data-options="required:true" editable="false">
                            {{range $k,$v:=.mchlist}}
                            <option value="{{$v.mch_id}}">{{$v.mch_name}}</option>
                            {{end}}
                        </select>
                        <script type="text/javascript">
							$('#mch').combobox({
                                onLoadSuccess: function () {
								    $('#mch').combobox('select','{{.m.mch_id}}');
							    }
                            });
						</script>
                    </td>
                </tr>
                <tr>
                    <td>类型:</td>
                    <td>
                        <select id="usertype" class="easyui-combobox" name="usertype" style="width:160px;" data-options="required:true" editable="false">
                            <option value="0">管理账号</option>
                            <option value="1">员工账号</option>
                        </select>
                        <script type="text/javascript">
							$('#usertype').combobox({
                                onLoadSuccess: function () {
								    $('#usertype').combobox('select','{{.m.usertype}}');
							    }
                            });
						</script>
                    </td>
                </tr>
				<tr>
                    <td>级别:</td>
                    <td>
                        <select id="level" class="easyui-combobox" name="level" style="width:160px;" data-options="required:true" editable="false">
                            <option value="0">免费会员</option>
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
                    <td><input class="easyui-textbox" type="text" name="password" style="width:160px;" value="{{.m.password}}" data-options="required:true"></input></td>
                </tr>
				<tr>
                    <td>角色:</td>
                    <td> 
						{{range $k,$r:=.roles}}
							<input type="checkbox" name="role" id="role{{$r.id}}" value="{{$r.id}}"   /><label for="role{{$r.id}}">{{$r.name}}</label>
							
						{{end}}
						 
						<script type="text/javascript">
						{{.jstr}}
						</script>
					</td>
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
           singleselect="true" fitcolumns="true" fit="true">
        <thead>
            <tr>
				<th field="id" width="5">ID</th>
				<th field="roleid" align="center"  width="10" >角色ID</th>
				<th field="rolename"  width="20" >角色</th>
				<th field="menuid" align="center"  width="10" >权限ID</th>
                <th field="menuname" align="center" data-options="formatter:rowformater_name" width="20" >权限</th>
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
