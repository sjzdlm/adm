package adm

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/sjzdlm/db"
)

//LoginController 控制器
type LoginController struct {
	beego.Controller
}

//Get 登录页面
func (c *LoginController) Get() {
	//不再支持IE浏览器
	var ie = c.GetString("ie")
	if ie != "1" && (strings.Index(c.Ctx.Request.UserAgent(), "MSIE") > 0 || strings.Index(c.Ctx.Request.UserAgent(), "Trident") > 0) {
		c.Ctx.Redirect(302, "/upgradeie")
		return
	}
	//定义模板参数map
	//var data map[string]interface{} = map[string]interface{}{}
	//获取回调地址
	var backurl = c.GetString("backurl")
	c.Data["backurl"] = backurl

	var m = db.First("select * from adm_system limit 1 ")
	c.Data["m"] = m

	//c.TplName="adm/login/index.html"
	//开始渲染页面---------------------------------------------------------------------------
	var tpl = template.New("")
	tpl.Parse(loginGet)
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

//Post 登录系统  0参数不全或用户名密码错误  1登录成功  2账号异常
func (c *LoginController) Post() {
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
			//首先清楚开发账号模式
			c.SetSession("_root", nil)
			c.SetSession("_root_level_", nil)
			c.DelSession("_root")
			c.DelSession("_root_level_")

			c.SetSession("_sysid", u["sysid"])
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

			c.Ctx.WriteString("1")
			return
		} else {
			c.Ctx.WriteString("2")
			return
		}
	}

}

var loginGet = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>用户登录-{{.m.sys_name}}</title>
  <!-- Tell the browser to be responsive to screen width -->
  <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">
  <!-- Bootstrap 3.3.6 -->
  <link rel="stylesheet" href="/assets/bootstrap/css/bootstrap.min.css">
  <!-- Font Awesome -->
  <link rel="stylesheet" href="/assets/adminlte/css/font-awesome.min.css">
  <!-- Ionicons -->
  <link rel="stylesheet" href="/assets/adminlte/css/ionicons.min.css">
  <!-- Theme style -->
  <link rel="stylesheet" href="/assets/adminlte/css/AdminLTE.min.css">
  <!-- iCheck -->
  <link rel="stylesheet" href="/assets/plugins/iCheck/square/blue.css">

  

  <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
  <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
  <!--[if lt IE 9]>
  <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
  <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
  <![endif]-->
  <style>
	.login-page , .register-page {
		background: #444;
		background-image: url(/images/login_bg.jpg);
		display: flex;
		align-items: center; 
		display: -webkit-box;display: -ms-flexbox;display: flex;display:-webkit-flex;
		background-repeat: no-repeat;
    background-size: cover;
    background-color: transparent;
    background-position: 50% 50%;
	}
  </style>
</head>
<body class="hold-transition login-page">
<!--[if lt IE 10]>
<script> // 如果推荐语使用默认值，可以删除此 script 标签
// IEDIE_HINT = '<p>自定义的提示语</p>';
</script>
<script src="/js/iedie/v1.2/script.min.js"></script>
<![endif]-->

<div class="login-box">
  <div class="login-logo">
     
  </div>
  <!-- /.login-logo -->
  <div class="login-box-body">
  <h3 class="form-title">请登录您的账户</h3>
    <p class="login-box-msg"></p>

    <form id="form1" action="/adm/login" method="post">
      <div class="form-group has-feedback">
        <input type="text" id="username" name="username" class="form-control" placeholder="账号">
        <span class="glyphicon glyphicon-user form-control-feedback"></span>
      </div>
      <div class="form-group has-feedback">
        <input type="password" name="userpwd" id="password" class="form-control" placeholder="密码">
        <span class="glyphicon glyphicon-lock form-control-feedback"></span>
      </div>
      <div class="row">
        <div class="col-xs-8">
          <div class="checkbox icheck">
            <label>
              <input type="checkbox"> 记住密码
            </label>
          </div>
        </div>
        <!-- /.col -->
        <div class="col-xs-4">
            <input type="hidden" name="backurl" id="backurl" value="{{.backurl}}"/>
          <button type="button"  name="btn" id="btn" class="btn btn-primary btn-block btn-flat">登录</button>
        </div>
        <!-- /.col -->
      </div>
    </form>
    
  <div style="height:50px;">&nbsp;</div>
  </div>
  <!-- /.login-box-body -->
</div>
<!-- /.login-box --> 
</div>
<script type="text/javascript" src="/js/jquery-1.11.1.min.js"></script>
<script type="text/javascript" src="/js/jquery.form.js"></script>
<script type="text/javascript" src="/js/layer/layer.js"></script>
<!-- Bootstrap 3.3.6 -->
<script src="/assets/bootstrap/js/bootstrap.min.js"></script>
<!-- iCheck -->
<script src="/assets/plugins/iCheck/icheck.min.js"></script>
<script>
if(top !== self){             top.location.href = location.href;         }
  $(function () {
    $('input').iCheck({
      checkboxClass: 'icheckbox_square-blue',
      radioClass: 'iradio_square-blue',
      increaseArea: '20%' // optional
    });
  });
</script>
<script type="text/javascript">
$(function(){
    $('#btn').click(function () {
        if ($('#username').val() == '' || $('#userpwd').val() == '') {
            layer.msg('请输入账号和密码!');
            return;
        }
        if ($('#backurl').val() == '') {
            $('#form1').ajaxSubmit(function (data) {
                if (data!= '1') {
                  if(data=='0'){
                    layer.msg('用户名或密码错误！');
                  }else{
                    layer.msg('账号异常，请与管理员联系！');
                  }
                } else {
                    $('#btn').attr('disabled',"true");
                    window.location = "/adm/main";
                }
            });
        } else {
            $('#form1').submit();
        }

		return false;
	})
  
})
</script>
 
</body>
</html>

`
