package routers

import (
	"github.com/astaxie/beego"
	"github.com/sjzdlm/adm/conf"
	"github.com/sjzdlm/adm/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.MainController{})

	conf.InitRouter()
}
