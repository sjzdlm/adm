package adm

import (
	"bytes"
	"html/template"

	//"strconv"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//MchController 控制器
type MchController struct {
	beego.Controller
}

//List 列表页面
func (c *MchController) List() {
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_mch_list)
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

//获取商户列表
func (c *MchController) ListJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var where = ""

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " where `mch_name` like '%" + qtxt + "%'"
	}

	var rst = db.Pager(page, pageSize, "select * from adm_mch "+where)

	c.Data["json"] = rst
	c.ServeJSON()
}

//企业添加/编辑
func (c *MchController) Edit() {
	var id, _ = c.GetInt("id", 0)
	if id > 0 {
		var m = db.First("select * from adm_mch where id=?", id)
		c.Data["m"] = m
	}

	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_mch_edit)
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

//信息保存提交
func (c *MchController) EditPost() {
	var id, _ = c.GetInt("id", 0)
	var mch_name = c.GetString("mch_name")
	var mch_contacts = c.GetString("mch_contacts")
	var mch_company = c.GetString("mch_company")
	var mch_phone = c.GetString("mch_phone")
	var mch_email = c.GetString("mch_email")
	var state = c.GetString("mch_state")

	var sql = ""
	if id > 0 {
		sql = `
		update adm_mch set 
		mch_name=?,
		mch_contacts=?,
		mch_company=?,
		mch_phone=?,
		mch_email=?,
		mch_state=?
		where id=?
		`
		var i = db.Exec(sql,
			mch_name,
			mch_contacts,
			mch_company,
			mch_phone,
			mch_email,
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
		insert into adm_mch(
			mch_name,
		mch_contacts,
		mch_company,
		mch_phone,
		mch_email,
		mch_state
		)values(
			?,?,?,?,?,?
		)
		`
		var i = db.Exec(sql,
			mch_name,
			mch_contacts,
			mch_company,
			mch_phone,
			mch_email,
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

//Del删除
func (c *MchController) Del() {
	var id, _ = c.GetInt("id", 0)
	var i = db.Exec("delete from adm_mch where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

var adm_mch_list = `
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
            $('#win').window('refresh', '/adm/mch/edit?id='+row.id);
			$('#win').window("resize",{top:$(document).scrollTop() + ($(window).height()-250) * 0.5});//居中显示
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
    var row = $('#tt').datagrid('getSelected');
    $('#win').window('open');
    $('#win').window('refresh', '/adm/mch/edit?id=');
    $('#ff').form('load', row);
}
function doDel(){
    var jq=jQuery;
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/mch/del', { id: row.id }, function (result) {
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
           url="/adm/mch/listjson"
           title="企业管理" toolbar="#tb" id="tt"
           singleselect="true" fitcolumns="true" fit="true"
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true"
           >
        <thead>
            <tr>
				<th field="id" width="20">ID</th>
				<th field="mch_name" width="50">名称</th>
                <th field="mch_company" width="70">公司</th>
				<th field="mch_contacts" width="70">联系人</th>
				<th field="mch_phone" width="70">联系电话</th>
				<th field="mch_email" width="50">邮箱</th>
				<th field="state" width="50">状态</th>
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

var adm_mch_edit = `

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
        <form id="form1" action="/adm/mch/editpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="mch_name" value="{{.m.mch_name}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
                </tr>
				<tr>
                    <td>公司:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="mch_company" value="{{.m.mch_company}}"></input></td>
				</tr>
				<tr>
                    <td>联系人:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="mch_contacts" value="{{.m.mch_contacts}}"></input></td>
				</tr>
				<tr>
                    <td>联系电话:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="mch_phone" value="{{.m.mch_phone}}"></input></td>
				</tr>
				<tr>
                    <td>电子邮箱:</td>
                    <td><input class="easyui-textbox" type="text" style="width:180px;" name="mch_email" value="{{.m.mch_email}}"></input></td>
                </tr>

                <tr>
                    <td>状态:</td>
					<td>
					<input type="hidden" id="id" name="id" value="{{.m.id}}" />
                        <select id="mch_state" class="easyui-combobox" name="mch_state" editable="false" style="width:180px;" >
                            <option value="0">禁用</option>
                            <option value="1">启用</option> 
                        </select>
						<script type="text/javascript">
						 
                            $('#mch_state').combobox({
                                onLoadSuccess: function (data) {
                                    $('#mch_state').combobox('setValue', "{{.m.mch_state}}");
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
