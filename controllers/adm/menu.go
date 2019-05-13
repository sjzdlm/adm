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

//MenuController 控制器
type MenuController struct {
	beego.Controller
}

//List 列表页面
func (c *MenuController) List() {
	var toplist = db.Query("select * from  adm_menu where state=1 and pid=1")
	var qmenu = ""
	for _, v := range toplist {
		qmenu += "<option value='" + v["id"] + "'  title='" + v["label"] + "'>" + v["title"] + "-" + v["id"] + "</option>"
		var sublist = db.Query("select * from adm_menu where state=1 and pid=?", v["id"])
		for _, vv := range sublist {
			qmenu += "<option value='" + vv["id"] + "'   &nbsp;title='" + vv["label"] + "'> --" + vv["title"] + vv["id"] + "</option>"
		}
	}
	fmt.Println("qmenu", qmenu)
	c.Data["qmenu"] = qmenu

	//fmt.Println(qmenu)
	//c.TplName="adm/menu/list.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_menu_list")
	adm_menu_list = strings.Replace(adm_menu_list, "{{str2html .qmenu}}", qmenu, 1)
	tpl.Parse(adm_menu_list)
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

//获取菜单列表
func (c *MenuController) ListJson() {
	var page, _ = c.GetInt("page", 1)
	var pageSize, _ = c.GetInt("rows", 20)
	var qtxt = c.GetString("qtxt")
	var pid, _ = c.GetInt("pid", 1)
	var where = " where pid=" + strconv.Itoa(pid) + " "
	if pid == 0 {
		where = " where 1=1 "
	}

	qtxt = strings.TrimSpace(string(qtxt))
	if qtxt != "" {
		where += " and title like '%" + qtxt + "%'"
	}
	//排序
	var sort = c.GetString("sort")
	var order = c.GetString("order")
	if sort != "" && order != "" {
		where += " order by " + sort + " " + order
	}
	var rst = db.Pager(page, pageSize, "select  *  from adm_menu   "+where)

	c.Data["json"] = rst
	c.ServeJSON()
}

//Edit 列表页面
func (c *MenuController) Edit() {
	var id, _ = c.GetInt("id", 0)
	c.Data["id"] = id

	//顶级菜单
	var toplist = db.Query("select * from  adm_menu where state=1 and pid=1")
	var qmenu = ""
	for _, v := range toplist {
		qmenu += "<option value='" + v["id"] + "'>" + v["title"] + "</option>"
		var sublist = db.Query("select * from adm_menu where state=1 and pid=?", v["id"])
		for _, vv := range sublist {
			qmenu += "<option value='" + vv["id"] + "'> --" + vv["title"] + "</option>"
		}
	}
	c.Data["qmenu"] = qmenu
	//当前菜单
	if id > 0 {
		var m = db.First("select * from adm_menu where id=?", id)
		if m == nil {
			c.Ctx.WriteString("暂无数据!")
			return
		}
		c.Data["m"] = m
	}

	//c.TplName="adm/menu/edit.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("adm_menu_edit")
	adm_menu_edit = strings.Replace(adm_menu_edit, "{{str2html .qmenu}}", qmenu, 1)
	tpl.Parse(adm_menu_edit)
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

//Edit 数据保存
func (c *MenuController) EditPost() {
	var id, _ = c.GetInt("id", 0)
	// if id<1{
	// 	c.Ctx.WriteString("参数错误!")
	// 	return
	// }

	var pid, _ = c.GetInt("parentid", 0)
	var title = c.GetString("title")
	var orders, _ = c.GetInt("orders", 0)
	var image = c.GetString("image")
	var icon = c.GetString("icon")
	var url = c.GetString("url")
	var memo = c.GetString("memo")
	var label = c.GetString("label")
	var state = c.GetString("state")
	if state == "on" {
		state = "1"
	} else {
		state = "0"
	}

	if title == "" || url == "" {
		c.Ctx.WriteString("请填写完整信息!")
		return
	}
	var sql = ""
	if id > 0 {
		sql = `update adm_menu set
		nid=?,
		pid=?,
		title=?,
		orders=?,
		image=?,
		icon=?,
		url=?,
		label=?,
		memo=?,
		state=?
		where id=?
		`
		var i = db.Exec(sql,
			id,
			pid,
			title,
			orders,
			image,
			icon,
			url,
			label,
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
		sql = `
		insert into adm_menu(
			pid,
			title,
			orders,
			image,
			icon,
			url,
			label,
			memo,
			state
		)values(
			?,?,?,?,?,?,?,?,?
		)
		`
		var i = db.Exec(sql,
			pid,
			title,
			orders,
			image,
			icon,
			url,
			label,
			memo,
			state,
		)
		if i > 0 {
			//更新nid
			db.Exec("update adm_menu set nid=id where nid is null")
			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("0")
			return
		}
	}
}

//Remove 菜单删除
func (c *MenuController) Remove() {
	var id, _ = c.GetInt("id", 0)
	if id < 1 {
		c.Ctx.WriteString("参数错误!")
		return
	}
	//校验是否有子菜单
	var list = db.Query("select * from adm_menu where pid=?", id)
	if len(list) > 0 {
		fmt.Println("不能直接删除有子菜单的菜单!")
		c.Ctx.WriteString("不能直接删除有子菜单的菜单!")
		return
	}
	var i = db.Exec("delete from adm_menu where id=?", id)
	if i > 0 {
		c.Ctx.WriteString("1")
		return
	} else {
		c.Ctx.WriteString("0")
		return
	}
}

var adm_menu_list = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title></title>
	<meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">

    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/metro/easyui.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
	<link rel="stylesheet" type="text/css" href="/fonts/iconfont.css">
	<link rel="stylesheet" href="/assets/adminlte/css/font-awesome.min.css">
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
            pid: $('#pid').combobox("getValue"),
			qtxt:$('#qtxt').val()
        });
    }
function doEdit(){
        var row = $('#tt').datagrid('getSelected');
        if (row){
			var w=$('#win').window({
				width:"420",
				height:"380",
				top:($(window).height() - 350) * 0.5,   
				left:($(window).width() - 680) * 0.5,
				modal:true,
				title:'菜单编辑-'+row.id
			});
            w.window('open');
            w.window('refresh', '/adm/menu/edit?id='+row.id);
            $('#ff').form('load',row);
        }else{
            jq.messager.alert('警告','请选择一行数据','warning');
        }

}
function doAdd() {
	var row = $('#tt').datagrid('getSelected');
	
	var w=$('#win').window({
		width:"420",
		height:"380",
		top:($(window).height() - 350) * 0.5,   
		left:($(window).width() - 680) * 0.5,
		modal:true,
		title:'菜单添加'
	});
    w.window('open');
    w.window('refresh', '/adm/menu/edit?id=');
    $('#ff').form('load', row);
}
function doRemove(){
    var row = $('#tt').datagrid('getSelected');
    if (row) {
        jq.messager.confirm('警告', '确定要删除吗?', function (r) {
            if (r) {
                jq.post('/adm/menu/remove', { id: row.id }, function (result) {
                    if (result=="1") {
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
	function rowformater_image(value, row, index) {
		return "<i class='fa "+value+"'>&nbsp;&nbsp;&nbsp;&nbsp;</i>";
		//return "<img src='"+value+"' style='width:25px;height:25px;'>";
    }
    </script>
</head>
<body style="padding:2px;margin-bottom:2px;">

    <table class="easyui-datagrid" style="width:600px;height:250px"
           url="/adm/menu/listjson"
           title="菜单管理" toolbar="#tb" id="tt" 
           singleselect="true" fitcolumns="true" fit="true" 
           data-options="fitColumns:true,pageList:[20,50,100],pageSize:20,pagination:true" >
        <thead>
            <tr>
				<th field="id" align="center" sortable="true" width="60">ID</th>
				<th field="label" align="left" sortable="true" width="80">菜单</th>
				<th field="title" align="left" sortable="true" width="80">名称</th>
				<th field="pid" align="center" sortable="true" width="80">上级</th>
				<th field="ptitle" align="center" sortable="true" width="80">上级菜单</th>
                <th field="orders" align="center" sortable="true" align="right" width="70">排序</th>
                <th field="image" sortable="true" align="center" width="50" data-options="formatter:rowformater_image">图标</th>
                <th field="url" sortable="true" width="200">URL地址</th>
                <th field="memo" width="50">备注</th>
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
            一级菜单: 
						<select  id="pid" name="pid" style="width:180px;" class="easyui-combobox1" editable='false'>
						<option value="0">请选择...</option>
						
						</select>
						<script language="javascript">
						$(function(){
							$('#pid').append('{{.qmenu}}');
							$('#pid').combobox({});
						});
						</script>
            查询参数: <input class="easyui-textbox" id="qtxt" style="width:160px">


            <a href="#" class="easyui-linkbutton" iconcls="icon-search" onclick="doSearch();">查 询</a>
        </div>
    </div>

    <div id="win" class="easyui-window" title="菜单编辑" closed="true" collapsible="false" minimizable="false" maximizable="false" style="width:480px;height:420px;padding:5px;">
        Some Content.
    </div>
</body>
</html>
`

var adm_menu_edit = `

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
						//layer.msg('<font color="yellow">操作成功!</font>');
                        $('#tt').datagrid('reload');
                        $('#win').window('close');
                    } else {
                        jq.messager.alert('错误', "操作失败!", "warning");
						//layer.msg('<font color="green">操作失败!</font>');
                    }
                }
            });
        }
        function clearForm(){
            $('#win').window('close');
        }
        $('#image').combobox({
            formatter: function (row) {

                //return '<span class="iconfont ' + row.text + '">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="item-text">' + row.text + '</span>';
				return '<span class="iconfont ' + row.text + '"></span><span class="item-text">' + row.text + '</span>';
            }
        });
</script>


<div class="easyui-panel" title="" style="width:100%" fix="true" border="false">
    <div style="padding:10px 60px 20px 60px">
        <form id="form1" action="/adm/menu/editpost" method="post">
            <table cellpadding="5">
                <tr>
                    <td>父节点:</td>
                    <td>
                        <select  id="parentid" name="parentid" style="width:165px;" class="easyui-combobox1" editable='false'>
						<option value="1">根节点</option>
						</select>
						<script language="javascript">
						$(function(){
							$('#parentid').append('{{.qmenu}}');
							$('#parentid').val({{.m.pid}});
							$('#parentid').combobox({});
						});
						</script>
                    </td>
                </tr>
                <tr>
                    <td>名称:</td>
                    <td><input class="easyui-textbox" type="text" name="title" style="width:165px;" value="{{.m.title}}" data-options="required:true,missingMessage:'必填字段'"></input></td>
				</tr>
				<tr>
                    <td>标签:</td>
                    <td><input class="easyui-textbox" type="text" name="label" style="width:165px;" value="{{.m.label}}"  ></input></td>
                </tr>
                <tr>
                    <td>排序:</td>
                    <td><input class="easyui-textbox" type="text" name="orders" style="width:165px;" value="{{.m.orders}}" data-options="required:true"></input></td>
                </tr>
                <tr>
                    <td>图标:</td>
                    <td>
                        <select id="image" class="easyui-combobox fa" name="image" editable='false' style="width:165px;" data-options="required:true">
						<!--
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
                            <option value="icon-56">icon-56</option>
                            <option value="icon-edit">icon-edit</option>
                            <option value="icon-ok">icon-ok</option>
                            <option value="icon-blank">icon-blank</option>
                            <option value="icon-clear">icon-clear</option>
                            <option value="icon-save">icon-save</option>
                            <option value="icon-cut">icon-cut</option>
                            <option value="icon-no">icon-no</option>
                            <option value="icon-no">icon-no</option>
                            <option value="icon-reload">icon-reload</option>
                            <option value="icon-search">icon-search</option>
                            <option value="icon-print">icon-print</option>
                            <option value="icon-back">icon-back</option>
                            <option value="icon-sum">icon-sum</option>
                            <option value="icon-lock">icon-lock</option>
                            <option value="icon-more">icon-more</option>
							-->
							<!--
							<option value="icon-quanbudingdan">icon-quanbudingdan</option>
							<option value="icon-pdf">icon-pdf</option>
							<option value="icon-title-left">icon-title-left</option>
							<option value="icon-book">icon-book</option>
							<option value="icon-yuanquan">icon-yuanquan</option>
							<option value="icon-shujuku">icon-shujuku</option>
							<option value="icon-youxiang">icon-youxiang</option>
							<option value="icon-weiwancheng">icon-weiwancheng</option>
							<option value="icon-del">icon-del</option>
							<option value="icon-fabu">icon-fabu</option>
							<option value="icon-xia">icon-xia</option>
							<option value="icon-22">icon-22</option>
							<option value="icon-information">icon-information</option>
							<option value="icon-huanyihuan">icon-huanyihuan</option>
							<option value="icon-lingdang">icon-lingdang</option>
							<option value="icon-wenjianjia3">icon-wenjianjia3</option>
							<option value="icon-gongyingshangguanli">icon-gongyingshangguanli</option>
							<option value="icon-shoucanggongyingshang">icon-shoucanggongyingshang</option>
							<option value="icon-xinzenggongyingshang">icon-xinzenggongyingshang</option>
							<option value="icon-icon">icon-icon</option>
							<option value="icon-xiazai">icon-xiazai</option>
							<option value="icon-yanzhengma">icon-yanzhengma</option>
							<option value="icon-wuliaoguanli">icon-wuliaoguanli</option>
							<option value="icon-gangchang2">icon-gangchang2</option>
							<option value="icon-buxiugang4">icon-buxiugang4</option>
							<option value="icon-55yewushenqingguidang">icon-55yewushenqingguidang</option>
							<option value="icon-56gongwenshenqingguidang">icon-56gongwenshenqingguidang</option>
							<option value="icon-57qitashenqingguidang">icon-57qitashenqingguidang</option>
							<option value="icon-redooutline">icon-redooutline</option>
							<option value="icon-fanhui">icon-fanhui</option>
							<option value="icon-01">icon-01</option>
							<option value="icon-yunxingwancheng">icon-yunxingwancheng</option>
							<option value="icon-chexiao">icon-chexiao</option>
							<option value="icon-sousuo">icon-sousuo</option>
							<option value="icon-renyuantianjia">icon-renyuantianjia</option>
							<option value="icon-dingdanhetong">icon-dingdanhetong</option>
							<option value="icon-jilu">icon-jilu</option>
							<option value="icon-renwu">icon-renwu</option>
							<option value="icon-daiban">icon-daiban</option>
							<option value="icon-gongqiuduihua">icon-gongqiuduihua</option>
							<option value="icon-shuaxin">icon-shuaxin</option>
							<option value="icon-list">icon-list</option>
							<option value="icon-del1">icon-del1</option>
							<option value="icon-yonghu">icon-yonghu</option>
							<option value="icon-iconfont05">icon-iconfont05</option>
							<option value="icon-tuichu">icon-tuichu</option>
							<option value="icon-553zuzhijiagou">icon-553zuzhijiagou</option>
							<option value="icon-fanhui1">icon-fanhui1</option>
							<option value="icon-xiangyou">icon-xiangyou</option>
							<option value="icon-renyuanguanli2">icon-renyuanguanli2</option>
							<option value="icon-xiao21">icon-xiao21</option>
							<option value="icon-1">icon-1</option>
							<option value="icon-3">icon-3</option>
							<option value="icon-mdliucheng">icon-mdliucheng</option>
							<option value="icon-iconfontdaochuexcel">icon-iconfontdaochuexcel</option>
							<option value="icon-jinlingyingcaiwangtubiao98">icon-jinlingyingcaiwangtubiao98</option>
							<option value="icon-45354354">icon-45354354</option>
							<option value="icon-720caogaoxiang">icon-720caogaoxiang</option>
							<option value="icon-iconyichu01">icon-iconyichu01</option>
							<option value="icon-icontianjia01">icon-icontianjia01</option>
							<option value="icon-guanbi">icon-guanbi</option>
							<option value="icon-jixiao">icon-jixiao</option>
							<option value="icon-guanwangicon64">icon-guanwangicon64</option>
							<option value="icon-chenggong">icon-chenggong</option>
							<option value="icon-shaixuan">icon-shaixuan</option>
							<option value="icon-iconfonticon">icon-iconfonticon</option>
							<option value="icon-shangwuhezuo">icon-shangwuhezuo</option>
							<option value="icon-2">icon-2</option>
							<option value="icon-4">icon-4</option>
							<option value="icon-5">icon-5</option>
							<option value="icon-gerenshezhi">icon-gerenshezhi</option>
							<option value="icon-shoujianxiang">icon-shoujianxiang</option>
							<option value="icon-duigou">icon-duigou</option>
							<option value="icon-iconqueren">icon-iconqueren</option>
							<option value="icon-daochu3">icon-daochu3</option>
							<option value="icon-shejituzhi">icon-shejituzhi</option>
							<option value="icon-guanbi1">icon-guanbi1</option>
							<option value="icon-xiangshang01">icon-xiangshang01</option>
							<option value="icon-mima">icon-mima</option>
							<option value="icon-comments">icon-comments</option>
							<option value="icon-jingpaidanju">icon-jingpaidanju</option>
							<option value="icon-excel">icon-excel</option>
							<option value="icon-word">icon-word</option>
							<option value="icon-daishouhuo">icon-daishouhuo</option>
							<option value="icon-xiangzuo">icon-xiangzuo</option>
							<option value="icon-loupan">icon-loupan</option>
							<option value="icon-yonghu1">icon-yonghu1</option>
							<option value="icon-yaoqinghaoyoupengyou3">icon-yaoqinghaoyoupengyou3</option>
							<option value="icon-daishouhuo1">icon-daishouhuo1</option>
							<option value="icon-gongyingshangguanli1">icon-gongyingshangguanli1</option>
							<option value="icon-daochu">icon-daochu</option>
							<option value="icon-key">icon-key</option>
							<option value="icon-shapes">icon-shapes</option>
							<option value="icon-xiangyou1">icon-xiangyou1</option>
							<option value="icon-caiwutongji">icon-caiwutongji</option>
							<option value="icon-weibiaoti1">icon-weibiaoti1</option>
							<option value="icon-zhaotoubiao">icon-zhaotoubiao</option>
							<option value="icon-zhaobiao">icon-zhaobiao</option>
							<option value="icon-gongyingshang">icon-gongyingshang</option>
							<option value="icon-saoma">icon-saoma</option>
							<option value="icon-ping">icon-ping</option>
							<option value="icon-suoxiao">icon-suoxiao</option>
							<option value="icon-richeng">icon-richeng</option>
							<option value="icon-fuwuzhongxingongchangrenzheng">icon-fuwuzhongxingongchangrenzheng</option>
							<option value="icon-gongyingshangkanban">icon-gongyingshangkanban</option>
							<option value="icon-dingdan">icon-dingdan</option>
							<option value="icon-qiyezhongzhi">icon-qiyezhongzhi</option>
							<option value="icon-gongyingshangguanli2">icon-gongyingshangguanli2</option>
							<option value="icon-xianchangkaohe">icon-xianchangkaohe</option>
							<option value="icon-tupian">icon-tupian</option>
							<option value="icon-ceping">icon-ceping</option>
							<option value="icon-diaochawenjuan">icon-diaochawenjuan</option>
							<option value="icon-96">icon-96</option>
							<option value="icon-zugroup">icon-zugroup</option>
							<option value="icon-shanchu">icon-shanchu</option>
							<option value="icon-shanchu1">icon-shanchu1</option>
							<option value="icon-tongbu">icon-tongbu</option>
							<option value="icon-titlebarcaidan">icon-titlebarcaidan</option>
							<option value="icon-flag">icon-flag</option>
							<option value="icon-iconfontpdf">icon-iconfontpdf</option>
							<option value="icon-zuidahua">icon-zuidahua</option>
							<option value="icon-xunjia">icon-xunjia</option>
							<option value="icon-movedown">icon-movedown</option>
							<option value="icon-moveup">icon-moveup</option>
							<option value="icon-baocun">icon-baocun</option>
							<option value="icon-shoudaoxunjia">icon-shoudaoxunjia</option>
							<option value="icon-dingdanguanli">icon-dingdanguanli</option>
							<option value="icon-fuwushangrenzheng">icon-fuwushangrenzheng</option>
							<option value="icon-baocun1">icon-baocun1</option>
							<option value="icon-pinpaishouquan">icon-pinpaishouquan</option>
							<option value="icon-icon9">icon-icon9</option>
							<option value="icon-daichuli">icon-daichuli</option>
							<option value="icon-xunjia1">icon-xunjia1</option>
							<option value="icon-chuangkou">icon-chuangkou</option>
							<option value="icon-notice">icon-notice</option>
							<option value="icon-doc">icon-doc</option>
							<option value="icon-xls">icon-xls</option>
							<option value="icon-xiayi">icon-xiayi</option>
							<option value="icon-shangyi">icon-shangyi</option>
							<option value="icon-qiyerenzheng">icon-qiyerenzheng</option>
							<option value="icon-multi-line-text">icon-multi-line-text</option>
							<option value="icon-mulu">icon-mulu</option>
							<option value="icon-soliddown">icon-soliddown</option>
							<option value="icon-fabu1">icon-fabu1</option>
							<option value="icon-tubiao16">icon-tubiao16</option>
							<option value="icon-tubiao48">icon-tubiao48</option>
							<option value="icon-zizhirenzheng">icon-zizhirenzheng</option>
							<option value="icon-biaoji">icon-biaoji</option>
							<option value="icon-qunfengfuwushang">icon-qunfengfuwushang</option>
							<option value="icon-xinjiaruqiyeyujing">icon-xinjiaruqiyeyujing</option>
							<option value="icon-bianji">icon-bianji</option>
							<option value="icon-neirong">icon-neirong</option>
							<option value="icon-xitongguanli">icon-xitongguanli</option>
							<option value="icon-jihua01">icon-jihua01</option>
							<option value="icon-jiaoxuepeizhi">icon-jiaoxuepeizhi</option>
							<option value="icon-renwu1">icon-renwu1</option>
							<option value="icon-wenhao">icon-wenhao</option>
							<option value="icon-fazhankaidan">icon-fazhankaidan</option>
							<option value="icon-me-copy">icon-me-copy</option>
							<option value="icon-tuopanxiangzi">icon-tuopanxiangzi</option>
							<option value="icon-fanhui2">icon-fanhui2</option>
							<option value="icon-bangzhu">icon-bangzhu</option>
							<option value="icon-555">icon-555</option>
							<option value="icon-kuanyi20guanrenyaoqingyuangong">icon-kuanyi20guanrenyaoqingyuangong</option>
							<option value="icon-zhiliangjianyan">icon-zhiliangjianyan</option>
							<option value="icon-liucheng">icon-liucheng</option>
							<option value="icon-bangzhu-copy">icon-bangzhu-copy</option>
							<option value="icon-iconset0114">icon-iconset0114</option>
							<option value="icon-iconset0115">icon-iconset0115</option>
							<option value="icon-iconset0186">icon-iconset0186</option>
							<option value="icon-iconset0187">icon-iconset0187</option>
							<option value="icon-iconset01100">icon-iconset01100</option>
							<option value="icon-iconset0256">icon-iconset0256</option>
							<option value="icon-iconset0339">icon-iconset0339</option>
							<option value="icon-iconset0340">icon-iconset0340</option>
							<option value="icon-jingxiaoshangfuwu">icon-jingxiaoshangfuwu</option>
							<option value="icon-quanxianguanli">icon-quanxianguanli</option>
							<option value="icon-jingji">icon-jingji</option>
							<option value="icon-gongyinglian">icon-gongyinglian</option>
							<option value="icon-rizhi">icon-rizhi</option>
							<option value="icon-zaixianfuwu">icon-zaixianfuwu</option>
							<option value="icon-daoru">icon-daoru</option>
							<option value="icon-riqi">icon-riqi</option>
							<option value="icon-yingjigongyingfenbutu">icon-yingjigongyingfenbutu</option>
							<option value="icon-jiagecaiji">icon-jiagecaiji</option>
							<option value="icon-weibiaoti7">icon-weibiaoti7</option>
							<option value="icon-windows">icon-windows</option>
							<option value="icon-tongji">icon-tongji</option>
							<option value="icon-hezuohuobanguanli">icon-hezuohuobanguanli</option>
							<option value="icon-shaixuan1">icon-shaixuan1</option>
							<option value="icon-daiwancheng">icon-daiwancheng</option>
							<option value="icon-zonghesum1">icon-zonghesum1</option>
							<option value="icon-76">icon-76</option>
							<option value="icon-fenlei">icon-fenlei</option>
							<option value="icon-richangyewu">icon-richangyewu</option>
							<option value="icon-zhuanyetuandui">icon-zhuanyetuandui</option>
							<option value="icon-solidup">icon-solidup</option>
							<option value="icon-xiangxia">icon-xiangxia</option>
							<option value="icon-title-right">icon-title-right</option>
							<option value="icon-xinzeng">icon-xinzeng</option>
							<option value="icon-edit">icon-edit</option>
							<option value="icon-xiangyou2">icon-xiangyou2</option>
							<option value="icon-xiangzuo1">icon-xiangzuo1</option>
							<option value="icon-dengpao">icon-dengpao</option>
							<option value="icon-hetong">icon-hetong</option>
							<option value="icon-quanxianguanli1">icon-quanxianguanli1</option>
							<option value="icon-yilianjie">icon-yilianjie</option>
							<option value="icon-duankailianjie">icon-duankailianjie</option>
							<option value="icon-tuihuo1">icon-tuihuo1</option>
							<option value="icon-gongyinglian1">icon-gongyinglian1</option>
							<option value="icon-tubiao0102">icon-tubiao0102</option>
							<option value="icon-wenjiangongxiang">icon-wenjiangongxiang</option>
							<option value="icon-tongguo">icon-tongguo</option>
							<option value="icon-dingshijiancha">icon-dingshijiancha</option>
							<option value="icon-xitongguanli1">icon-xitongguanli1</option>
							<option value="icon-icon-rank">icon-icon-rank</option>
							<option value="icon-scan">icon-scan</option>
							<option value="icon-yunyingyong">icon-yunyingyong</option>
							<option value="icon-shangchuan">icon-shangchuan</option>
							<option value="icon-weibiaoti201">icon-weibiaoti201</option>
							<option value="icon-tuzhang02">icon-tuzhang02</option>
							<option value="icon-trade">icon-trade</option>
							<option value="icon-icon01">icon-icon01</option>
							<option value="icon-gong">icon-gong</option>
							<option value="icon-qianzaiyonghu">icon-qianzaiyonghu</option>
							<option value="icon-zuzhiguanli">icon-zuzhiguanli</option>
							<option value="icon-dayin">icon-dayin</option>
							<option value="icon-cf-c19">icon-cf-c19</option>
							<option value="icon-no">icon-no</option>
							<option value="icon-20160518wangzhanshouyeyindingicon">icon-20160518wangzhanshouyeyindingicon</option>
							<option value="icon-gongyingshangguanli3">icon-gongyingshangguanli3</option>
							<option value="icon-suoxiaochuangkou01">icon-suoxiaochuangkou01</option>
							<option value="icon-duidingdan">icon-duidingdan</option>
							<option value="icon-jishuhezuoqiatanhezuohuitan">icon-jishuhezuoqiatanhezuohuitan</option>
							<option value="icon-shangwuhezuohezuo">icon-shangwuhezuohezuo</option>
							<option value="icon-daochu1">icon-daochu1</option>
							<option value="icon-xinpingpujuecefenxicongyerenyuan">icon-xinpingpujuecefenxicongyerenyuan</option>
							<option value="icon-wuliugongying01">icon-wuliugongying01</option>
							<option value="icon-xinxifuwuyewuyujingcongyerenyuanheimingdanchaxun">icon-xinxifuwuyewuyujingcongyerenyuanheimingdanchaxun</option>
							<option value="icon-chexiao1">icon-chexiao1</option>
							<option value="icon-weiwancheng1">icon-weiwancheng1</option>
							<option value="icon-dingdandiaodu">icon-dingdandiaodu</option>
							<option value="icon-zhongzhimima">icon-zhongzhimima</option>
							<option value="icon-reset">icon-reset</option>
							<option value="icon-shuju">icon-shuju</option>
							<option value="icon-xiangyou3">icon-xiangyou3</option>
							<option value="icon-jianqie">icon-jianqie</option>
							<option value="icon-qunzu">icon-qunzu</option>
							<option value="icon-bianji1">icon-bianji1</option>
							<option value="icon-shizhong">icon-shizhong</option>
							<option value="icon-bangzhu1">icon-bangzhu1</option>
							<option value="icon-cuowutishi">icon-cuowutishi</option>
							<option value="icon-yewurenyuanchushen">icon-yewurenyuanchushen</option>
							<option value="icon-shenqingkuaisucaigoujiandan">icon-shenqingkuaisucaigoujiandan</option>
							<option value="icon-diaochawenjuan1">icon-diaochawenjuan1</option>
							<option value="icon-jihuabiao">icon-jihuabiao</option>
							-->
							

							<option value="fa-address-book">fa fa-address-book</option>
							<option value="fa-address-book-o">fa fa-address-book-o</option>
							<option value="fa-address-card">fa fa-address-card</option>
							<option value="fa-address-card-o">fa fa-address-card-o</option>
							<option value="fa-adjust">fa fa-adjust</option>
							<option value="fa-american-sign-language-interpreting">fa fa-american-sign-language-interpreting</option>
							<option value="fa-anchor">fa fa-anchor</option>
							<option value="fa-archive">fa fa-archive</option>
							<option value="fa-area-chart">fa fa-area-chart</option>
							<option value="fa-arrows">fa fa-arrows</option>
							<option value="fa-arrows-h">fa fa-arrows-h</option>
							<option value="fa-arrows-v">fa fa-arrows-v</option>
							<option value="fa-asl-interpreting (alias)">fa fa-asl-interpreting (alias)</option>
							<option value="fa-assistive-listening-systems">fa fa-assistive-listening-systems</option>
							<option value="fa-asterisk">fa fa-asterisk</option>
							<option value="fa-at">fa fa-at</option>
							<option value="fa-audio-description">fa fa-audio-description</option>
							<option value="fa-automobile (alias)">fa fa-automobile (alias)</option>
							<option value="fa-balance-scale">fa fa-balance-scale</option>
							<option value="fa-ban">fa fa-ban</option>
							<option value="fa-bank (alias)">fa fa-bank (alias)</option>
							<option value="fa-bar-chart">fa fa-bar-chart</option>
							<option value="fa-bar-chart-o (alias)">fa fa-bar-chart-o (alias)</option>
							<option value="fa-barcode">fa fa-barcode</option>
							<option value="fa-bars">fa fa-bars</option>
							<option value="fa-bath">fa fa-bath</option>
							<option value="fa-bathtub (alias)">fa fa-bathtub (alias)</option>
							<option value="fa-battery (alias)">fa fa-battery (alias)</option>
							<option value="fa-battery-0 (alias)">fa fa-battery-0 (alias)</option>
							<option value="fa-battery-1 (alias)">fa fa-battery-1 (alias)</option>
							<option value="fa-battery-2 (alias)">fa fa-battery-2 (alias)</option>
							<option value="fa-battery-3 (alias)">fa fa-battery-3 (alias)</option>
							<option value="fa-battery-4 (alias)">fa fa-battery-4 (alias)</option>
							<option value="fa-battery-empty">fa fa-battery-empty</option>
							<option value="fa-battery-full">fa fa-battery-full</option>
							<option value="fa-battery-half">fa fa-battery-half</option>
							<option value="fa-battery-quarter">fa fa-battery-quarter</option>
							<option value="fa-battery-three-quarters">fa fa-battery-three-quarters</option>
							<option value="fa-bed">fa fa-bed</option>
							<option value="fa-beer">fa fa-beer</option>
							<option value="fa-bell">fa fa-bell</option>
							<option value="fa-bell-o">fa fa-bell-o</option>
							<option value="fa-bell-slash">fa fa-bell-slash</option>
							<option value="fa-bell-slash-o">fa fa-bell-slash-o</option>
							<option value="fa-bicycle">fa fa-bicycle</option>
							<option value="fa-binoculars">fa fa-binoculars</option>
							<option value="fa-birthday-cake">fa fa-birthday-cake</option>
							<option value="fa-blind">fa fa-blind</option>
							<option value="fa-bluetooth">fa fa-bluetooth</option>
							<option value="fa-bluetooth-b">fa fa-bluetooth-b</option>
							<option value="fa-bolt">fa fa-bolt</option>
							<option value="fa-bomb">fa fa-bomb</option>
							<option value="fa-book">fa fa-book</option>
							<option value="fa-bookmark">fa fa-bookmark</option>
							<option value="fa-bookmark-o">fa fa-bookmark-o</option>
							<option value="fa-braille">fa fa-braille</option>
							<option value="fa-briefcase">fa fa-briefcase</option>
							<option value="fa-bug">fa fa-bug</option>
							<option value="fa-building">fa fa-building</option>
							<option value="fa-building-o">fa fa-building-o</option>
							<option value="fa-bullhorn">fa fa-bullhorn</option>
							<option value="fa-bullseye">fa fa-bullseye</option>
							<option value="fa-bus">fa fa-bus</option>
							<option value="fa-cab (alias)">fa fa-cab (alias)</option>
							<option value="fa-calculator">fa fa-calculator</option>
							<option value="fa-calendar">fa fa-calendar</option>
							<option value="fa-calendar-check-o">fa fa-calendar-check-o</option>
							<option value="fa-calendar-minus-o">fa fa-calendar-minus-o</option>
							<option value="fa-calendar-o">fa fa-calendar-o</option>
							<option value="fa-calendar-plus-o">fa fa-calendar-plus-o</option>
							<option value="fa-calendar-times-o">fa fa-calendar-times-o</option>
							<option value="fa-camera">fa fa-camera</option>
							<option value="fa-camera-retro">fa fa-camera-retro</option>
							<option value="fa-car">fa fa-car</option>
							<option value="fa-caret-square-o-down">fa fa-caret-square-o-down</option>
							<option value="fa-caret-square-o-left">fa fa-caret-square-o-left</option>
							<option value="fa-caret-square-o-right">fa fa-caret-square-o-right</option>
							<option value="fa-caret-square-o-up">fa fa-caret-square-o-up</option>
							<option value="fa-cart-arrow-down">fa fa-cart-arrow-down</option>
							<option value="fa-cart-plus">fa fa-cart-plus</option>
							<option value="fa-cc">fa fa-cc</option>
							<option value="fa-certificate">fa fa-certificate</option>
							<option value="fa-check">fa fa-check</option>
							<option value="fa-check-circle">fa fa-check-circle</option>
							<option value="fa-check-circle-o">fa fa-check-circle-o</option>
							<option value="fa-check-square">fa fa-check-square</option>
							<option value="fa-check-square-o">fa fa-check-square-o</option>
							<option value="fa-child">fa fa-child</option>
							<option value="fa-circle">fa fa-circle</option>
							<option value="fa-circle-o">fa fa-circle-o</option>
							<option value="fa-circle-o-notch">fa fa-circle-o-notch</option>
							<option value="fa-circle-thin">fa fa-circle-thin</option>
							<option value="fa-clock-o">fa fa-clock-o</option>
							<option value="fa-clone">fa fa-clone</option>
							<option value="fa-close (alias)">fa fa-close (alias)</option>
							<option value="fa-cloud">fa fa-cloud</option>
							<option value="fa-cloud-download">fa fa-cloud-download</option>
							<option value="fa-cloud-upload">fa fa-cloud-upload</option>
							<option value="fa-code">fa fa-code</option>
							<option value="fa-code-fork">fa fa-code-fork</option>
							<option value="fa-coffee">fa fa-coffee</option>
							<option value="fa-cog">fa fa-cog</option>
							<option value="fa-cogs">fa fa-cogs</option>
							<option value="fa-comment">fa fa-comment</option>
							<option value="fa-comment-o">fa fa-comment-o</option>
							<option value="fa-commenting">fa fa-commenting</option>
							<option value="fa-commenting-o">fa fa-commenting-o</option>
							<option value="fa-comments">fa fa-comments</option>
							<option value="fa-comments-o">fa fa-comments-o</option>
							<option value="fa-compass">fa fa-compass</option>
							<option value="fa-copyright">fa fa-copyright</option>
							<option value="fa-creative-commons">fa fa-creative-commons</option>
							<option value="fa-credit-card">fa fa-credit-card</option>
							<option value="fa-credit-card-alt">fa fa-credit-card-alt</option>
							<option value="fa-crop">fa fa-crop</option>
							<option value="fa-crosshairs">fa fa-crosshairs</option>
							<option value="fa-cube">fa fa-cube</option>
							<option value="fa-cubes">fa fa-cubes</option>
							<option value="fa-cutlery">fa fa-cutlery</option>
							<option value="fa-dashboard (alias)">fa fa-dashboard (alias)</option>
							<option value="fa-database">fa fa-database</option>
							<option value="fa-deaf">fa fa-deaf</option>
							<option value="fa-deafness (alias)">fa fa-deafness (alias)</option>
							<option value="fa-desktop">fa fa-desktop</option>
							<option value="fa-diamond">fa fa-diamond</option>
							<option value="fa-dot-circle-o">fa fa-dot-circle-o</option>
							<option value="fa-download">fa fa-download</option>
							<option value="fa-drivers-license (alias)">fa fa-drivers-license (alias)</option>
							<option value="fa-drivers-license-o (alias)">fa fa-drivers-license-o (alias)</option>
							<option value="fa-edit (alias)">fa fa-edit (alias)</option>
							<option value="fa-ellipsis-h">fa fa-ellipsis-h</option>
							<option value="fa-ellipsis-v">fa fa-ellipsis-v</option>
							<option value="fa-envelope">fa fa-envelope</option>
							<option value="fa-envelope-o">fa fa-envelope-o</option>
							<option value="fa-envelope-open">fa fa-envelope-open</option>
							<option value="fa-envelope-open-o">fa fa-envelope-open-o</option>
							<option value="fa-envelope-square">fa fa-envelope-square</option>
							<option value="fa-eraser">fa fa-eraser</option>
							<option value="fa-exchange">fa fa-exchange</option>
							<option value="fa-exclamation">fa fa-exclamation</option>
							<option value="fa-exclamation-circle">fa fa-exclamation-circle</option>
							<option value="fa-exclamation-triangle">fa fa-exclamation-triangle</option>
							<option value="fa-external-link">fa fa-external-link</option>
							<option value="fa-external-link-square">fa fa-external-link-square</option>
							<option value="fa-eye">fa fa-eye</option>
							<option value="fa-eye-slash">fa fa-eye-slash</option>
							<option value="fa-eyedropper">fa fa-eyedropper</option>
							<option value="fa-fax">fa fa-fax</option>
							<option value="fa-feed (alias)">fa fa-feed (alias)</option>
							<option value="fa-female">fa fa-female</option>
							<option value="fa-fighter-jet">fa fa-fighter-jet</option>
							<option value="fa-file-archive-o">fa fa-file-archive-o</option>
							<option value="fa-file-audio-o">fa fa-file-audio-o</option>
							<option value="fa-file-code-o">fa fa-file-code-o</option>
							<option value="fa-file-excel-o">fa fa-file-excel-o</option>
							<option value="fa-file-image-o">fa fa-file-image-o</option>
							<option value="fa-file-movie-o (alias)">fa fa-file-movie-o (alias)</option>
							<option value="fa-file-pdf-o">fa fa-file-pdf-o</option>
							<option value="fa-file-photo-o (alias)">fa fa-file-photo-o (alias)</option>
							<option value="fa-file-picture-o (alias)">fa fa-file-picture-o (alias)</option>
							<option value="fa-file-powerpoint-o">fa fa-file-powerpoint-o</option>
							<option value="fa-file-sound-o (alias)">fa fa-file-sound-o (alias)</option>
							<option value="fa-file-video-o">fa fa-file-video-o</option>
							<option value="fa-file-word-o">fa fa-file-word-o</option>
							<option value="fa-file-zip-o (alias)">fa fa-file-zip-o (alias)</option>
							<option value="fa-film">fa fa-film</option>
							<option value="fa-filter">fa fa-filter</option>
							<option value="fa-fire">fa fa-fire</option>
							<option value="fa-fire-extinguisher">fa fa-fire-extinguisher</option>
							<option value="fa-flag">fa fa-flag</option>
							<option value="fa-flag-checkered">fa fa-flag-checkered</option>
							<option value="fa-flag-o">fa fa-flag-o</option>
							<option value="fa-flash (alias)">fa fa-flash (alias)</option>
							<option value="fa-flask">fa fa-flask</option>
							<option value="fa-folder">fa fa-folder</option>
							<option value="fa-folder-o">fa fa-folder-o</option>
							<option value="fa-folder-open">fa fa-folder-open</option>
							<option value="fa-folder-open-o">fa fa-folder-open-o</option>
							<option value="fa-frown-o">fa fa-frown-o</option>
							<option value="fa-futbol-o">fa fa-futbol-o</option>
							<option value="fa-gamepad">fa fa-gamepad</option>
							<option value="fa-gavel">fa fa-gavel</option>
							<option value="fa-gear (alias)">fa fa-gear (alias)</option>
							<option value="fa-gears (alias)">fa fa-gears (alias)</option>
							<option value="fa-gift">fa fa-gift</option>
							<option value="fa-glass">fa fa-glass</option>
							<option value="fa-globe">fa fa-globe</option>
							<option value="fa-graduation-cap">fa fa-graduation-cap</option>
							<option value="fa-group (alias)">fa fa-group (alias)</option>
							<option value="fa-hand-grab-o (alias)">fa fa-hand-grab-o (alias)</option>
							<option value="fa-hand-lizard-o">fa fa-hand-lizard-o</option>
							<option value="fa-hand-paper-o">fa fa-hand-paper-o</option>
							<option value="fa-hand-peace-o">fa fa-hand-peace-o</option>
							<option value="fa-hand-pointer-o">fa fa-hand-pointer-o</option>
							<option value="fa-hand-rock-o">fa fa-hand-rock-o</option>
							<option value="fa-hand-scissors-o">fa fa-hand-scissors-o</option>
							<option value="fa-hand-spock-o">fa fa-hand-spock-o</option>
							<option value="fa-hand-stop-o (alias)">fa fa-hand-stop-o (alias)</option>
							<option value="fa-handshake-o">fa fa-handshake-o</option>
							<option value="fa-hard-of-hearing (alias)">fa fa-hard-of-hearing (alias)</option>
							<option value="fa-hashtag">fa fa-hashtag</option>
							<option value="fa-hdd-o">fa fa-hdd-o</option>
							<option value="fa-headphones">fa fa-headphones</option>
							<option value="fa-heart">fa fa-heart</option>
							<option value="fa-heart-o">fa fa-heart-o</option>
							<option value="fa-heartbeat">fa fa-heartbeat</option>
							<option value="fa-history">fa fa-history</option>
							<option value="fa-home">fa fa-home</option>
							<option value="fa-hotel (alias)">fa fa-hotel (alias)</option>
							<option value="fa-hourglass">fa fa-hourglass</option>
							<option value="fa-hourglass-1 (alias)">fa fa-hourglass-1 (alias)</option>
							<option value="fa-hourglass-2 (alias)">fa fa-hourglass-2 (alias)</option>
							<option value="fa-hourglass-3 (alias)">fa fa-hourglass-3 (alias)</option>
							<option value="fa-hourglass-end">fa fa-hourglass-end</option>
							<option value="fa-hourglass-half">fa fa-hourglass-half</option>
							<option value="fa-hourglass-o">fa fa-hourglass-o</option>
							<option value="fa-hourglass-start">fa fa-hourglass-start</option>
							<option value="fa-i-cursor">fa fa-i-cursor</option>
							<option value="fa-id-badge">fa fa-id-badge</option>
							<option value="fa-id-card">fa fa-id-card</option>
							<option value="fa-id-card-o">fa fa-id-card-o</option>
							<option value="fa-image (alias)">fa fa-image (alias)</option>
							<option value="fa-inbox">fa fa-inbox</option>
							<option value="fa-industry">fa fa-industry</option>
							<option value="fa-info">fa fa-info</option>
							<option value="fa-info-circle">fa fa-info-circle</option>
							<option value="fa-institution (alias)">fa fa-institution (alias)</option>
							<option value="fa-key">fa fa-key</option>
							<option value="fa-keyboard-o">fa fa-keyboard-o</option>
							<option value="fa-language">fa fa-language</option>
							<option value="fa-laptop">fa fa-laptop</option>
							<option value="fa-leaf">fa fa-leaf</option>
							<option value="fa-legal (alias)">fa fa-legal (alias)</option>
							<option value="fa-lemon-o">fa fa-lemon-o</option>
							<option value="fa-level-down">fa fa-level-down</option>
							<option value="fa-level-up">fa fa-level-up</option>
							<option value="fa-life-bouy (alias)">fa fa-life-bouy (alias)</option>
							<option value="fa-life-buoy (alias)">fa fa-life-buoy (alias)</option>
							<option value="fa-life-ring">fa fa-life-ring</option>
							<option value="fa-life-saver (alias)">fa fa-life-saver (alias)</option>
							<option value="fa-lightbulb-o">fa fa-lightbulb-o</option>
							<option value="fa-line-chart">fa fa-line-chart</option>
							<option value="fa-location-arrow">fa fa-location-arrow</option>
							<option value="fa-lock">fa fa-lock</option>
							<option value="fa-low-vision">fa fa-low-vision</option>
							<option value="fa-magic">fa fa-magic</option>
							<option value="fa-magnet">fa fa-magnet</option>
							<option value="fa-mail-forward (alias)">fa fa-mail-forward (alias)</option>
							<option value="fa-mail-reply (alias)">fa fa-mail-reply (alias)</option>
							<option value="fa-mail-reply-all (alias)">fa fa-mail-reply-all (alias)</option>
							<option value="fa-male">fa fa-male</option>
							<option value="fa-map">fa fa-map</option>
							<option value="fa-map-marker">fa fa-map-marker</option>
							<option value="fa-map-o">fa fa-map-o</option>
							<option value="fa-map-pin">fa fa-map-pin</option>
							<option value="fa-map-signs">fa fa-map-signs</option>
							<option value="fa-meh-o">fa fa-meh-o</option>
							<option value="fa-microchip">fa fa-microchip</option>
							<option value="fa-microphone">fa fa-microphone</option>
							<option value="fa-microphone-slash">fa fa-microphone-slash</option>
							<option value="fa-minus">fa fa-minus</option>
							<option value="fa-minus-circle">fa fa-minus-circle</option>
							<option value="fa-minus-square">fa fa-minus-square</option>
							<option value="fa-minus-square-o">fa fa-minus-square-o</option>
							<option value="fa-mobile">fa fa-mobile</option>
							<option value="fa-mobile-phone (alias)">fa fa-mobile-phone (alias)</option>
							<option value="fa-money">fa fa-money</option>
							<option value="fa-moon-o">fa fa-moon-o</option>
							<option value="fa-mortar-board (alias)">fa fa-mortar-board (alias)</option>
							<option value="fa-motorcycle">fa fa-motorcycle</option>
							<option value="fa-mouse-pointer">fa fa-mouse-pointer</option>
							<option value="fa-music">fa fa-music</option>
							<option value="fa-navicon (alias)">fa fa-navicon (alias)</option>
							<option value="fa-newspaper-o">fa fa-newspaper-o</option>
							<option value="fa-object-group">fa fa-object-group</option>
							<option value="fa-object-ungroup">fa fa-object-ungroup</option>
							<option value="fa-paint-brush">fa fa-paint-brush</option>
							<option value="fa-paper-plane">fa fa-paper-plane</option>
							<option value="fa-paper-plane-o">fa fa-paper-plane-o</option>
							<option value="fa-paw">fa fa-paw</option>
							<option value="fa-pencil">fa fa-pencil</option>
							<option value="fa-pencil-square">fa fa-pencil-square</option>
							<option value="fa-pencil-square-o">fa fa-pencil-square-o</option>
							<option value="fa-percent">fa fa-percent</option>
							<option value="fa-phone">fa fa-phone</option>
							<option value="fa-phone-square">fa fa-phone-square</option>
							<option value="fa-photo (alias)">fa fa-photo (alias)</option>
							<option value="fa-picture-o">fa fa-picture-o</option>
							<option value="fa-pie-chart">fa fa-pie-chart</option>
							<option value="fa-plane">fa fa-plane</option>
							<option value="fa-plug">fa fa-plug</option>
							<option value="fa-plus">fa fa-plus</option>
							<option value="fa-plus-circle">fa fa-plus-circle</option>
							<option value="fa-plus-square">fa fa-plus-square</option>
							<option value="fa-plus-square-o">fa fa-plus-square-o</option>
							<option value="fa-podcast">fa fa-podcast</option>
							<option value="fa-power-off">fa fa-power-off</option>
							<option value="fa-print">fa fa-print</option>
							<option value="fa-puzzle-piece">fa fa-puzzle-piece</option>
							<option value="fa-qrcode">fa fa-qrcode</option>
							<option value="fa-question">fa fa-question</option>
							<option value="fa-question-circle">fa fa-question-circle</option>
							<option value="fa-question-circle-o">fa fa-question-circle-o</option>
							<option value="fa-quote-left">fa fa-quote-left</option>
							<option value="fa-quote-right">fa fa-quote-right</option>
							<option value="fa-random">fa fa-random</option>
							<option value="fa-recycle">fa fa-recycle</option>
							<option value="fa-refresh">fa fa-refresh</option>
							<option value="fa-registered">fa fa-registered</option>
							<option value="fa-remove (alias)">fa fa-remove (alias)</option>
							<option value="fa-reorder (alias)">fa fa-reorder (alias)</option>
							<option value="fa-reply">fa fa-reply</option>
							<option value="fa-reply-all">fa fa-reply-all</option>
							<option value="fa-retweet">fa fa-retweet</option>
							<option value="fa-road">fa fa-road</option>
							<option value="fa-rocket">fa fa-rocket</option>
							<option value="fa-rss">fa fa-rss</option>
							<option value="fa-rss-square">fa fa-rss-square</option>
							<option value="fa-s15 (alias)">fa fa-s15 (alias)</option>
							<option value="fa-search">fa fa-search</option>
							<option value="fa-search-minus">fa fa-search-minus</option>
							<option value="fa-search-plus">fa fa-search-plus</option>
							<option value="fa-send (alias)">fa fa-send (alias)</option>
							<option value="fa-send-o (alias)">fa fa-send-o (alias)</option>
							<option value="fa-server">fa fa-server</option>
							<option value="fa-share">fa fa-share</option>
							<option value="fa-share-alt">fa fa-share-alt</option>
							<option value="fa-share-alt-square">fa fa-share-alt-square</option>
							<option value="fa-share-square">fa fa-share-square</option>
							<option value="fa-share-square-o">fa fa-share-square-o</option>
							<option value="fa-shield">fa fa-shield</option>
							<option value="fa-ship">fa fa-ship</option>
							<option value="fa-shopping-bag">fa fa-shopping-bag</option>
							<option value="fa-shopping-basket">fa fa-shopping-basket</option>
							<option value="fa-shopping-cart">fa fa-shopping-cart</option>
							<option value="fa-shower">fa fa-shower</option>
							<option value="fa-sign-in">fa fa-sign-in</option>
							<option value="fa-sign-language">fa fa-sign-language</option>
							<option value="fa-sign-out">fa fa-sign-out</option>
							<option value="fa-signal">fa fa-signal</option>
							<option value="fa-signing (alias)">fa fa-signing (alias)</option>
							<option value="fa-sitemap">fa fa-sitemap</option>
							<option value="fa-sliders">fa fa-sliders</option>
							<option value="fa-smile-o">fa fa-smile-o</option>
							<option value="fa-snowflake-o">fa fa-snowflake-o</option>
							<option value="fa-soccer-ball-o (alias)">fa fa-soccer-ball-o (alias)</option>
							<option value="fa-sort">fa fa-sort</option>
							<option value="fa-sort-alpha-asc">fa fa-sort-alpha-asc</option>
							<option value="fa-sort-alpha-desc">fa fa-sort-alpha-desc</option>
							<option value="fa-sort-amount-asc">fa fa-sort-amount-asc</option>
							<option value="fa-sort-amount-desc">fa fa-sort-amount-desc</option>
							<option value="fa-sort-asc">fa fa-sort-asc</option>
							<option value="fa-sort-desc">fa fa-sort-desc</option>
							<option value="fa-sort-down (alias)">fa fa-sort-down (alias)</option>
							<option value="fa-sort-numeric-asc">fa fa-sort-numeric-asc</option>
							<option value="fa-sort-numeric-desc">fa fa-sort-numeric-desc</option>
							<option value="fa-sort-up (alias)">fa fa-sort-up (alias)</option>
							<option value="fa-space-shuttle">fa fa-space-shuttle</option>
							<option value="fa-spinner">fa fa-spinner</option>
							<option value="fa-spoon">fa fa-spoon</option>
							<option value="fa-square">fa fa-square</option>
							<option value="fa-square-o">fa fa-square-o</option>
							<option value="fa-star">fa fa-star</option>
							<option value="fa-star-half">fa fa-star-half</option>
							<option value="fa-star-half-empty (alias)">fa fa-star-half-empty (alias)</option>
							<option value="fa-star-half-full (alias)">fa fa-star-half-full (alias)</option>
							<option value="fa-star-half-o">fa fa-star-half-o</option>
							<option value="fa-star-o">fa fa-star-o</option>
							<option value="fa-sticky-note">fa fa-sticky-note</option>
							<option value="fa-sticky-note-o">fa fa-sticky-note-o</option>
							<option value="fa-street-view">fa fa-street-view</option>
							<option value="fa-suitcase">fa fa-suitcase</option>
							<option value="fa-sun-o">fa fa-sun-o</option>
							<option value="fa-support (alias)">fa fa-support (alias)</option>
							<option value="fa-tablet">fa fa-tablet</option>
							<option value="fa-tachometer">fa fa-tachometer</option>
							<option value="fa-tag">fa fa-tag</option>
							<option value="fa-tags">fa fa-tags</option>
							<option value="fa-tasks">fa fa-tasks</option>
							<option value="fa-taxi">fa fa-taxi</option>
							<option value="fa-television">fa fa-television</option>
							<option value="fa-terminal">fa fa-terminal</option>
							<option value="fa-thermometer (alias)">fa fa-thermometer (alias)</option>
							<option value="fa-thermometer-0 (alias)">fa fa-thermometer-0 (alias)</option>
							<option value="fa-thermometer-1 (alias)">fa fa-thermometer-1 (alias)</option>
							<option value="fa-thermometer-2 (alias)">fa fa-thermometer-2 (alias)</option>
							<option value="fa-thermometer-3 (alias)">fa fa-thermometer-3 (alias)</option>
							<option value="fa-thermometer-4 (alias)">fa fa-thermometer-4 (alias)</option>
							<option value="fa-thermometer-empty">fa fa-thermometer-empty</option>
							<option value="fa-thermometer-full">fa fa-thermometer-full</option>
							<option value="fa-thermometer-half">fa fa-thermometer-half</option>
							<option value="fa-thermometer-quarter">fa fa-thermometer-quarter</option>
							<option value="fa-thermometer-three-quarters">fa fa-thermometer-three-quarters</option>
							<option value="fa-thumb-tack">fa fa-thumb-tack</option>
							<option value="fa-thumbs-down">fa fa-thumbs-down</option>
							<option value="fa-thumbs-o-down">fa fa-thumbs-o-down</option>
							<option value="fa-thumbs-o-up">fa fa-thumbs-o-up</option>
							<option value="fa-thumbs-up">fa fa-thumbs-up</option>
							<option value="fa-ticket">fa fa-ticket</option>
							<option value="fa-times">fa fa-times</option>
							<option value="fa-times-circle">fa fa-times-circle</option>
							<option value="fa-times-circle-o">fa fa-times-circle-o</option>
							<option value="fa-times-rectangle (alias)">fa fa-times-rectangle (alias)</option>
							<option value="fa-times-rectangle-o (alias)">fa fa-times-rectangle-o (alias)</option>
							<option value="fa-tint">fa fa-tint</option>
							<option value="fa-toggle-down (alias)">fa fa-toggle-down (alias)</option>
							<option value="fa-toggle-left (alias)">fa fa-toggle-left (alias)</option>
							<option value="fa-toggle-off">fa fa-toggle-off</option>
							<option value="fa-toggle-on">fa fa-toggle-on</option>
							<option value="fa-toggle-right (alias)">fa fa-toggle-right (alias)</option>
							<option value="fa-toggle-up (alias)">fa fa-toggle-up (alias)</option>
							<option value="fa-trademark">fa fa-trademark</option>
							<option value="fa-trash">fa fa-trash</option>
							<option value="fa-trash-o">fa fa-trash-o</option>
							<option value="fa-tree">fa fa-tree</option>
							<option value="fa-trophy">fa fa-trophy</option>
							<option value="fa-truck">fa fa-truck</option>
							<option value="fa-tty">fa fa-tty</option>
							<option value="fa-tv (alias)">fa fa-tv (alias)</option>
							<option value="fa-umbrella">fa fa-umbrella</option>
							<option value="fa-universal-access">fa fa-universal-access</option>
							<option value="fa-university">fa fa-university</option>
							<option value="fa-unlock">fa fa-unlock</option>
							<option value="fa-unlock-alt">fa fa-unlock-alt</option>
							<option value="fa-unsorted (alias)">fa fa-unsorted (alias)</option>
							<option value="fa-upload">fa fa-upload</option>
							<option value="fa-user">fa fa-user</option>
							<option value="fa-user-circle">fa fa-user-circle</option>
							<option value="fa-user-circle-o">fa fa-user-circle-o</option>
							<option value="fa-user-o">fa fa-user-o</option>
							<option value="fa-user-plus">fa fa-user-plus</option>
							<option value="fa-user-secret">fa fa-user-secret</option>
							<option value="fa-user-times">fa fa-user-times</option>
							<option value="fa-users">fa fa-users</option>
							<option value="fa-vcard (alias)">fa fa-vcard (alias)</option>
							<option value="fa-vcard-o (alias)">fa fa-vcard-o (alias)</option>
							<option value="fa-video-camera">fa fa-video-camera</option>
							<option value="fa-volume-control-phone">fa fa-volume-control-phone</option>
							<option value="fa-volume-down">fa fa-volume-down</option>
							<option value="fa-volume-off">fa fa-volume-off</option>
							<option value="fa-volume-up">fa fa-volume-up</option>
							<option value="fa-warning (alias)">fa fa-warning (alias)</option>
							<option value="fa-wheelchair">fa fa-wheelchair</option>
							<option value="fa-wheelchair-alt">fa fa-wheelchair-alt</option>
							<option value="fa-wifi">fa fa-wifi</option>
							<option value="fa-window-close">fa fa-window-close</option>
							<option value="fa-window-close-o">fa fa-window-close-o</option>
							<option value="fa-window-maximize">fa fa-window-maximize</option>
							<option value="fa-window-minimize">fa fa-window-minimize</option>
							<option value="fa-window-restore">fa fa-window-restore</option>
							<option value="fa-wrench">fa fa-wrench</option>

                        </select>
                    </td>
				</tr>
				<tr>
                    <td>图标2:</td>
                    <td>
                        <select id="icon" class="easyui-combobox fa" name="icon" editable='false' style="width:165px;" data-options="required:true">
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
                            <option value="icon-56">icon-56</option>
                            <option value="icon-edit">icon-edit</option>
                            <option value="icon-ok">icon-ok</option>
                            <option value="icon-blank">icon-blank</option>
                            <option value="icon-clear">icon-clear</option>
                            <option value="icon-save">icon-save</option>
                            <option value="icon-cut">icon-cut</option>
                            <option value="icon-no">icon-no</option>
                            <option value="icon-no">icon-no</option>
                            <option value="icon-reload">icon-reload</option>
                            <option value="icon-search">icon-search</option>
                            <option value="icon-print">icon-print</option>
                            <option value="icon-back">icon-back</option>
                            <option value="icon-sum">icon-sum</option>
                            <option value="icon-lock">icon-lock</option>
                            <option value="icon-more">icon-more</option>
							
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
                    <td>URL:</td>
                    <td><input class="easyui-textbox" type="text" id="url" name="url" style="width:165px;" value="{{.m.url}}" data-options="required:true"></input></td>
				</tr>
				<tr>
                    <td>模块:</td>
					<td>
					<select class="easyui-combogrid" style="width:100%" data-options="
							panelWidth: 500,
							idField: 'code',
							textField: 'title',
							url: '/adm/tb/tblist',
							method: 'get',
							columns: [[
								{field:'id',title:'ID',width:50},
								{field:'proj',title:'分类',width:100},
								{field:'title',title:'模块',width:120,align:'right'},
								{field:'code',title:'代号',width:120,align:'right'},
								{field:'table',title:'表名',width:100},
								{field:'memo',title:'备注',width:120,align:'center'}
							]],
							fitColumns: true,
							label: '',
							labelPosition: 'top',
							onSelect: function (rowIndex, row){
								$('#url').textbox('setValue','/adm/tb/list/'+row.code);
							}
						">
					</select>
					</td>
                </tr>
                <tr>
                    <td>备注:</td>
                    <td>
                        <input class="easyui-textbox" type="text" name="memo" style="width:165px;" value="{{.m.memo}}"></input>
                        <input type="hidden" id="id" name="id" value="{{.id}}" />
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
            <a href="javascript:void(0)" class="easyui-linkbutton" iconcls="icon-no" onclick="clearForm()">取 消&nbsp;</a>
        </div>
    </div>
</div>
<script type="text/javascript">
<!--
	$('#image').combobox({
		onLoadSuccess: function (){
			$('#image').combobox('select','{{.m.image}}');
			$('.combo-panel span').removeClass('iconfont');
		}
	});
	$(function(){
		if('{{.m.url}}'==''){
			$('#url').val('javascript:void(0);');
		}
	});
//-->
</script>

`
