package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

//ErrorController 控制器
type ErrorController struct {
	beego.Controller
}

//Error404 错误页
func (c *ErrorController) Error404() {
	c.Ctx.WriteString("page not found")
	// c.Data["content"] = "page not found"
	// c.TplName = "404.html"
	c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
	c.Ctx.WriteString("")
}

//Error501 服务器错误
func (c *ErrorController) Error501() {
	c.Ctx.WriteString("server error")
	c.Data["content"] = "server error"
	c.TplName = "501.html"
	c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
}

//ErrorDb 数据库错误
func (c *ErrorController) ErrorDb() {
	c.Ctx.WriteString("database error")
	c.Data["content"] = "database is now down"
	c.TplName = "error.html"
}
