package adm

import (
	"bytes"
	"fmt"
	"strings"

	//"fmt"
	"html/template"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//MainController 控制器
type MainController struct {
	beego.Controller
}

//首页
func (c *MainController) Get() {
	//1.检查是否登录
	var _uid = c.Ctx.Input.Session("_uid")
	var _uname = c.Ctx.Input.Session("_username")
	var _mch_id = c.Ctx.Input.Session("_mch_id")
	var _logintime = c.Ctx.Input.Session("_logintime")
	var _loginip = c.Ctx.Input.Session("_loginip")

	if (_uid == nil || _uid == "") && c.Ctx.Request.RequestURI != "/adm/login" {
		c.Ctx.Redirect(302, "/adm/login")
		return
	}
	c.Data["_logintime"] = _logintime
	c.Data["_loginip"] = _loginip
	//2.获取用户信息
	var m = db.FirstOrNil("select * from adm_user where id=?", _uid)
	if m == nil {
		c.Ctx.Redirect(302, "/adm/login")
		return
	}
	c.Data["_user"] = m
	//3.获取菜单参数
	var nid, _ = c.GetInt("nid", 2) //默认为2,'首页'
	c.Data["nid"] = nid
	//fmt.Println("nid",nid)
	//4.默认欢迎页
	var defpage = "/adm/main/def"
	if m["defpage"] != "" {
		defpage = m["defpage"]
	}
	c.Data["defpage"] = defpage
	//5.角色字符串
	var roles = db.Query("select * from adm_role where id in (" + m["roles"] + ")")
	var rs = ""
	for k, v := range roles {
		//fmt.Println("k v:",k,v)
		if k > 0 {
			rs += ","
		}
		rs += v["rights"]
	}
	//6.获取顶部菜单
	var topmenu = db.Query("select * from adm_menu where pid=1 and nid in (" + rs + ") order by orders ")
	c.Data["_topmenu"] = topmenu

	var _datajson = ""

	//只有当访问参与有su=1参数且用户为root时才放开数据管理功能
	var _un = _uname.(string)
	if c.GetString("su") != "1" {
		_un = ""
	}
	for _, v := range topmenu {
		var jsonstr = menuJson(v["nid"], m["roles"], v["title"], _un)
		_datajson += "_datajson" + v["nid"] + "=" + jsonstr + ";"
	}
	c.Data["_datajson"] = template.JS(_datajson)
	//fmt.Println("_datajson\r\n",_datajson)

	// var js=`
	// <script language='javascript'>
	// alert('hello');
	// </script>
	// `
	//c.Data["_js"]=template.HTML(js)

	//7.获取系统信息
	var sys_logo = ""
	var sys = db.First("select * from adm_system where mch_id=? limit 1", _mch_id)
	if sys == nil {
		sys = db.First("select * from adm_system limit 1")
	}
	if sys != nil {
		c.Data["sysname"] = sys["sys_name"]
		sys_logo = sys["sys_logo"]
	} else {
		c.Data["sysname"] = "掌上软件"
	}
	if sys_logo == "" {
		sys_logo = "/images/logo.png"
	}
	c.Data["sys_logo"] = sys_logo

	//如果用户是root,双击头像跳转新连接,可以进行高级管理
	if _uname == "root" {
		var _js = `
        $("#headimg").dblclick(function(){
            window.location='/adm/main?su=1';
          });
        `
		c.Data["_headimg_js"] = _js
	}

	//c.TplName="adm/main/index.html"

	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	if strings.Index(c.Ctx.Request.UserAgent(), "MSIE") > 0 || strings.Index(c.Ctx.Request.UserAgent(), "Trident") > 0 {
		tpl.Parse(adm_main_get_ie)
	} else {
		tpl.Parse(adm_main_get)
	}

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

//根据菜单ID获取菜单JSON字符串,只有 _uname为 root时才有数据管理功能
func menuJson(nid string, roles string, title string, _uname string) string {
	//db.Exec("SET SSESION group_concat_max_len=102400;")//mysql才有此命令
	var rs = db.FirstOrNil("select GROUP_CONCAT(A.rights) as r from adm_role A   where id in(" + roles + ")")
	var rights = "0"
	if rs != nil {
		rights = rs["r"]
	}
	var list = db.Query("select * from adm_menu where pid=? and state=1 and nid in ("+rights+") order by orders ", nid)
	//fmt.Println("list",list)
	var rst = ""

	for i, v := range list {
		if i == 0 {
			rst += "{ "
			rst += "    \"id\": \"" + nid + "\", "
			rst += "    \"text\": \"" + title + "\", "
			rst += "    \"icon\": \"\", "
			rst += "    \"image\": \"\", "
			rst += "\"url\":\"" + v["url"] + "\","
			rst += "    \"isHeader\": true "
			rst += "}, "

			rst += "{"
			rst += "\"isOpen\":true,"
		}
		if i > 0 {
			rst += ",{"
		}

		rst += "\"id\":\"" + v["nid"] + "\","
		rst += "\"text\":\"" + v["title"] + "\","
		rst += "\"icon\":\"fa " + v["image"] + "\","
		rst += "\"image\":\"fa " + v["icon"] + "\","
		rst += "\"url\":\"" + v["url"] + "\""

		var slist = db.Query("select * from adm_menu where  pid=? and state=1 and nid in ("+rights+") order by orders ", v["nid"])
		if slist != nil && len(slist) > 0 {
			rst += ",\"children\":["
			for j, vv := range slist {
				if j > 0 {
					rst += ","
				}
				rst += "{\"id\":\"" + vv["nid"] + "\","
				rst += "\"text\":\"" + vv["title"] + "\","
				rst += "\"icon\":\"fa " + vv["image"] + "\","
				rst += "\"image\":\"fa " + vv["icon"] + "\","
				rst += "\"targetType\":\"iframe-tab\","
				rst += "\"url\":\"" + vv["url"] + "\"}"
			}
			rst += "]"
		}
		rst += "}"
	}
	if title == "系统" && _uname == "root" {
		rst += `
		,{
			"id": "900",
			"text": "数据管理",
			"icon": "fa fa-database",
			"url": "javascript:void(0);",
			"children": [{
				"id": "9000",
				"text": "系统管理",
                "icon": "fa fa-cog",
                "image": "icon-17",
				"targetType": "iframe-tab",
				"url": "/adm/tb/system"
			},{
				"id": "9001",
				"text": "模块管理",
                "icon": "fa fa-eye",
                "image": "icon-10",
				"targetType": "iframe-tab",
				"url": "/adm/tb/list"
			}, {
				"id": "9002",
				"text": "链接管理",
                "icon": "fa fa-plug",
                "image": "icon-11",
				"targetType": "iframe-tab",
				"url": "/adm/tb/conn"
			}, {
				"id": "9003",
				"text": "接口管理",
                "icon": "fa fa-asterisk",
                "image": "icon-12",
				"targetType": "iframe-tab",
				"url": "/adm/tb/api"
			}, {
				"id": "9004",
				"text": "报表管理",
                "icon": "fa fa-paw",
                "image": "icon-13",
				"targetType": "iframe-tab",
				"url": "/adm/tb/rpt"
			}, {
				"id": "9005",
				"text": "页面管理",
                "icon": "fa fa-adjust",
                "image": "icon-14",
				"targetType": "iframe-tab",
				"url": "/adm/tb/page"
			}, {
				"id": "9006",
				"text": "微站管理",
                "icon": "fa fa-tags",
                "image": "icon-0",
				"targetType": "iframe-tab",
				"url": "/adm/tb/applist"
			}]
		}
		`
	}
	rst = "{\"data\":[" + rst + "]}"
	return rst
}

//默认欢迎页
func (c *MainController) Def() {

	//c.TplName="adm/main/def.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(adm_main_def)
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

var adm_main_get_ie = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>{{.sysname}}</title>
    <link href="/css/default.css" rel="stylesheet" type="text/css" />
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/icon.css">
    <link rel="stylesheet" type="text/css" href="/js/easyui/themes/bootstrap/easyui.css">
    <!-- Font Awesome -->
    <link rel="stylesheet" href="/assets/adminlte/css/font-awesome.min.css">
    <!-- Ionicons -->
    <link rel="stylesheet" href="/assets/adminlte/css/ionicons.min.css">


    <script type="text/javascript" src="/js/easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/js/easyui/jquery.easyui.min.js"></script>
    
    <style>
        body {
            overflow-x: hidden;
            overflow-y: hidden;
        }
    </style>
    <script type="text/javascript">
        var jq=jQuery;
        var _menus = {};

        _menus = {
            "menus": [
                {
                    "menuid": 2,
                    "icon": "icon-41",
                    "menuname": "数据查询",
                    "selected": "true",
                    menus: [
                    ]
                }
                
            ]
        }
        function logout() {
            window.location = '/login';
        }
        $(function () {
            addTab("系统首页", '/main/def');
        })
    </script>
    <script type="text/javascript" src='/js/easyui/init.js'></script>
    <script src="/js/loader.js"></script>
</head>

<body class="easyui-layout">
    <div data-options="region:'north',border:false" split="true" border="false" style="overflow: hidden; height: 50px;
        background:  #7f99be repeat-x repeat-y ;
        background-image:url(/images/bg.gif);background-repeat:norepeat;filter:'progid:DXImageTransform.Microsoft.AlphaImageLoader(sizingMethod=scale)';-moz-background-size:100% 100%; 
        line-height: 50px;color: #000; font-family: Verdana, 微软雅黑,黑体">
        <span  style="margin-top:50px;padding-left:10px;color:#000;height:30px; font-size: 16px; ">
        <img alt="" src="{{.sys_logo}}" width="30px" height="30px" style="display: inline; 
        margin-top:1px;margin-right:1px;margin-bottom: 3px;vertical-align:middle;"><span  style="font-size: 16px; ">{{.sysname}}</span>
        </span>
        <span style="margin-top:50px;padding-left:100px;color:red; font-size: 16px;">

            {{range ._topmenu}}
                <a href="#"  id="nav{{.nid}}" class="easyui-linkbutton {{.icon}}"  onclick="leftMenu({{.nid}});" data-options="plain:true">{{.title}}</a>
              {{end}}
        </span>

        <span style="float:right; padding-right:20px;color:black;text-decoration:none;" class="head">
		欢迎回来,[管理员] ... <a href="/login" id="loginOut" style="color:black;text-decoration:none;">安全退出</a></span>



    </div>
    <div data-options="region:'west',split:true,title:'功能菜单',collapsible:false" style="width:190px;">


        <div class="easyui-accordion" id="accmenu" fit="true" border="false">


        </div>

    </div>

    <div data-options="region:'center'" style="background: #eee; overflow-y:hidden">

        <div id="tabs" class="easyui-tabs" fit="true" border="false">

        </div>

    </div>
<script type="text/javascript">
$(function(){
    leftMenu('{{.nid}}');
})
{{._datajson}}
    function leftMenu(id) {
        //$('.navbar-nav li').removeClass('active');
        //$('#nav'+id).addClass('active');
        //$('.sidebar-menu').html('');
        //$('.sidebar-menu').sidebarMenu(eval('_datajson'+id));

        var pnl = $('#accmenu').accordion('panels');
        //alert(pnl.length);
        for (var i = 0; i < pnl.length; i++) {
             var title = pnl[i].panel("options").title;  
             $('#accmenu').accordion("remove",i);  
             pnl = $('#accmenu').accordion('panels');
             i--;
        }
        
        var menulist = '';
        var flag=0;
        $.each(eval('_datajson'+id+'.data'), function(i, n) {
            if(i>0){
                menulist  = '<div title="'+n.text+'"  data-options="iconCls:\''+n.image+'\'" style="overflow:auto;">';
                menulist += '<ul>';
                if(undefined !=n.children){
                    $.each(n.children, function(j, o) {
                        menulist += '<li><div><a style="position:relative;padding-left:20px;" target="mainFrame" way="' + o.url + '" ><span style="position:absolute;left:0px;margin-top:3px;" class="icon '+o.image+'" ></span>' + o.text + '</a></div></li> ';
                    })               
                }

                menulist += '</ul>';
                menulist += '</div>';
        
                $('.easyui-accordion').accordion('add', {
                    title: n.text,
                    content:menulist,
                    iconCls:n.image,
                    selected: n.selected
                });             
            }

        })
    
        
        $('.easyui-accordion li a').click(function(){
            var tabTitle = $(this).text();
            var url = $(this).attr("way");
            //alert(tabTitle);
            addTab(tabTitle,url);
            $('.easyui-accordion li div').removeClass("selected");
            $(this).parent().addClass("selected");
        }).hover(function(){
            $(this).parent().addClass("hover");
        },function(){
            $(this).parent().removeClass("hover");
        });

    }
</script>

</body>

</html>
`
var adm_main_get = `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">

    <meta http-equiv="X-UA-Compatible" content="IE=edge">

    <title>{{.sysname}}</title>
    <!-- Tell the browser to be responsive to screen width -->
    <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">
    <script type="text/javascript" src="/js/easyui/base_loading.js"></script>
    <!-- Bootstrap 3.3.6 -->
    <link rel="stylesheet" href="/assets/bootstrap/css/bootstrap.min.css">
    <!-- Font Awesome -->
    <link rel="stylesheet" href="/assets/adminlte/css/font-awesome.min.css">
    <!-- Ionicons -->
    <link rel="stylesheet" href="/assets/adminlte/css/ionicons.min.css">
    <!-- Theme style -->
    <link rel="stylesheet" href="/assets/adminlte/css/AdminLTE.css">
    <!-- AdminLTE Skins. Choose a skin from the css/skins
         folder instead of downloading all of them to reduce the load. -->
    <link rel="stylesheet" href="/assets/adminlte/css/skins/all-skins.min.css">

    <style type="text/css">
        html {
            overflow: hidden;
        }
        .skin-red .sidebar-menu>li.header {
            color: #fff;
            background: #2c3b41;
        }
		@media (max-width: 991px){
		.navbar-collapse.pull-left { display: block;/* width: 200px;*/}
		}
        .navbar-collapse.pull-left + .navbar-custom-menu {
			top: auto;
			bottom: 5px;
		 }

		.nav>li{ float: left;}
    </style>
    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="../plugins/ie9/html5shiv.min.js"></script>
    <script src="../plugins/ie9/respond.min.js"></script>
    <![endif]-->
</head>
<body class="hold-transition1 skin-green sidebar-mini ">
<div class="wrapper">

    <header class="main-header">
        <!-- Logo -->
        <a href="/adm/main" class="logo">
            <!-- mini logo for sidebar mini 50x50 pixels -->
            <span class="logo-mini">
                <img alt="" src="{{.sys_logo}}" width="30px" height="30px" style="margin-right:1px;margin-bottom: 3px;margin-right:2px;">
            </span>
            <!-- logo for regular state and mobile devices -->
            <span class="logo-lg">
                <img alt="" src="{{.sys_logo}}" width="30px" height="30px" style="margin-left:-15px;margin-bottom: 3px;margin-right:2px;">{{.sysname}}
            </span>
        </a>
        <!-- Header Navbar: style can be found in header.less -->
        <nav class="navbar navbar-static-top">
            <!-- Sidebar toggle button-->
            <a href="#" class="sidebar-toggle" data-toggle="offcanvas" style="margin-top:1px;" role="button">
                <span class="sr-only">切换导航</span>
                <span class="icon-bar"></span>
            </a>
			
<!-- test-->
<div class="collapse navbar-collapse pull-left" id="navbar-collapse" style="height:50px;">
          <ul class="nav navbar-nav">
              {{range ._topmenu}}
                
                <li id="nav{{.nid}}"><a href="#" onclick="leftMenu({{.nid}});">{{.title}}</a></li>
                
              {{end}}
              <!--
            <li class="dropdown ">
              <a href="#" class="dropdown-toggle" data-toggle="dropdown" aria-expanded="true">更多... <span class="caret"></span></a>
              <ul class="dropdown-menu" role="menu">
                <li><a href="#">Action</a></li>
                <li><a href="#">Another action</a></li>
                <li><a href="#">Something else here</a></li>
                <li class="divider"></li>
                <li><a href="#">Separated link</a></li>
                <li class="divider"></li>
                <li><a href="#">One more separated link</a></li>
              </ul>
            </li>-->
          </ul>
        </div>
<!-- test-->

            <div class="navbar-custom-menu">
                <ul class="nav navbar-nav">
                    <li>
                        <a href="#" onclick="logOut();" title='退出系统'><i class="fa fa-sign-out"></i></a>
                    </li>
                    <!-- Control Sidebar Toggle Button -->
                    <li>
                        <a href="#" data-toggle="control-sidebar"><i class="fa fa-gears"></i></a>
                    </li>
                </ul>
            </div>
        </nav>
    </header>
    <!-- Left side column. contains the logo and sidebar -->
    <aside class="main-sidebar">
        <!-- sidebar: style can be found in sidebar.less -->
        <section class="sidebar">
            <!-- Sidebar user panel-->
            <div class="user-panel">
                <div class="pull-left image">
                    {{if eq ._user.headimg "" }}
                    <img src="/assets/adminlte/img/user2-160x160.jpg" class="img-circle"  >
                    {{else}}
                    <img src="{{._user.headimg}}" class="img-circle" id="headimg" >
                    {{end}}
                </div>
                <div class="pull-left info">
                    <p>{{._user.realname}}</p>
                    <a href="#" title="上次登录时间:{{._logintime}},IP:{{._loginip}}"><i class="fa fa-circle text-success"></i>{{._logintime}}</a>
                </div>
            </div> 
            <!-- search form 
            <form action="#" method="get" class="sidebar-form">
                <div class="input-group">
                    <input type="text" name="q" class="form-control" placeholder="Search...">
                    <span class="input-group-btn">
                <button type="button" name="search" id="search-btn" class="btn btn-flat" onclick="search_menu()"><i class="fa fa-search"></i>
                </button>
              </span>
                </div>
            </form>-->
            <!-- /.search form -->
            <!-- sidebar menu: : style can be found in sidebar.less -->
            <ul class="sidebar-menu">
                
            </ul>
		</section>
		<ul class="user-panel user-footer" style="position: absolute;bottom:0px;color:#5e3e3e;font-size:9px;">
            <div>如需帮助,请联系QQ/微信:65916383</div>
        </ul>   
        <!-- /.sidebar -->
    </aside>

    <!-- Content Wrapper. Contains page content -->
    <div class="content-wrapper" id="content-wrapper" style="min-height: 421px;">
        <!--bootstrap tab风格 多标签页-->
        <div class="content-tabs">
            <button class="roll-nav roll-left tabLeft" onclick="scrollTabLeft()">
                <i class="fa fa-arrow-circle-left"></i>
            </button>
            <nav class="page-tabs menuTabs tab-ui-menu" id="tab-menu">
                <div class="page-tabs-content" style="margin-left: 0px;">

                </div>
            </nav>
            
            <div class="btn-group roll-nav roll-right">
                <button class="dropdown tabClose" data-toggle="dropdown">
                    <i class="fa fa-arrow-circle-down" style="padding-left: 3px;"></i>
                </button>
                <ul class="dropdown-menu dropdown-menu-right" style="min-width: 128px;">
                    <li><a class="tabReload" href="javascript:refreshTab();">刷新当前</a></li>
                    <li><a class="tabCloseCurrent" href="javascript:closeCurrentTab();">关闭当前</a></li>
                    <li><a class="tabCloseAll" href="javascript:closeOtherTabs(true);">全部关闭</a></li>
                    <li><a class="tabCloseOther" href="javascript:closeOtherTabs();">除此之外全部关闭</a></li>
                </ul>
            </div>
            <button class="roll-nav roll-right fullscreen" onclick="App.handleFullScreen()"><i
                    class="fa fa-arrows"></i></button>
        </div>
        <div class="content-iframe " style="background-color: #ffffff; ">
            <div class="tab-content " id="tab-content">

            </div>
        </div>
    </div>
    <!-- /.content-wrapper

    <footer class="main-footer">
        <div class="pull-right hidden-xs">
            <b>Version</b> 2.3.8
        </div>
        <strong>Copyright &copy; 2014-2016 <a href="http://almsaeedstudio.com">Almsaeed Studio</a>.</strong> All rights
        reserved.
    </footer>
 -->
    <!-- Control Sidebar -->
    <aside class="control-sidebar control-sidebar-dark">
        <!-- Create the tabs -->
        <ul class="nav nav-tabs nav-justified control-sidebar-tabs">
            <li><a href="#control-sidebar-home-tab" data-toggle="tab"><i class="fa fa-home"></i></a></li>
        </ul>
        <!-- Tab panes -->
        <div class="tab-content">
            <!-- Home tab content -->
            <div class="tab-pane" id="control-sidebar-home-tab">
                <h3 class="control-sidebar-heading">&nbsp;</h3>
                <ul class="control-sidebar-menu">
                    <li>
                       <div style="height:100%;">&nbsp;</div>
                    </li>
                    
                </ul>
                <!-- /.control-sidebar-menu -->

                <!-- /.control-sidebar-menu -->

            </div>
            <!-- /.tab-pane -->
            <!-- Stats tab content -->
            <div class="tab-pane" id="control-sidebar-stats-tab">Stats Tab Content</div>
            <!-- /.tab-pane -->
            
        </div>
    </aside>
    <!-- /.control-sidebar -->
    <!-- Add the sidebar's background. This div must be placed
         immediately after the control sidebar -->
    <div class="control-sidebar-bg"></div>
</div>
<!-- ./wrapper -->

<!-- jQuery 2.2.3 -->
<script src="/assets/plugins/jQuery/jquery-2.2.3.min.js"></script>

<!-- Bootstrap 3.3.6 -->
    <script src="/assets/bootstrap/js/bootstrap.min.js"></script>
<!-- Slimscroll -->
    <script src="/assets/plugins/slimScroll/jquery.slimscroll.min.js"></script>
<!-- FastClick -->
    <script src="/assets/plugins/fastclick/fastclick.js"></script>
    <script src="/js/layer/layer.js"></script>

<script src="/assets/adminlte/js/superui.js"></script>



<script type="text/javascript">
    function addTab(_title,_url){
		var tabid=Math.floor(Math.random()*1000);
		//alert(tabid);
        addTabs({id:tabid,title:_title,url:_url,close:true});
    }

    function logOut(){
        layer.confirm('确实要退出吗？', {
            btn: ['是的','取消'] //按钮
        }, function(){
            window.location='/adm/login';
        }, function(){
            
        });
    }
    $(function () {

        App.setbasePath("../");
        App.setGlobalImgPath("assets/adminlte/img/");

        addTabs({
            id: '10008',
            title: '欢迎页',
            close: false,
            url: '/adm/main/def',
            urlType: "relative"
        });

        App.fixIframeCotent();

        
        leftMenu('{{.nid}}');

    });
    {{._datajson}}
    function leftMenu(id) {
        $('.navbar-nav li').removeClass('active');
        $('#nav'+id).addClass('active');
        $('.sidebar-menu').html('');
        $('.sidebar-menu').sidebarMenu(eval('_datajson'+id));
        //alert(eval('_datajson'+id).data[0].url);
    }
    var _js={{._headimg_js}}
    if(_js!=""){
        eval(_js);
    }
</script>
{{._js}}
</body>
</html>
`
var adm_main_def = `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title></title>
</head>

<body class=" pace-done">
    
</body>

</html>
`
