package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sjzdlm/adm/lib/uef"
)

//UEditor文件、图片上传专用
type UEFController struct {
	beego.Controller
}

const ACTION_CONFIG = "config"
const ACTION_LIST_IMAGE = "listimage"
const ACTION_LIST_FILE = "listfile"
const ACTION_UPLOAD_IMAGE = "uploadimage"
const ACTION_UPLOAD_FILE = "uploadfile"

func (c *UEFController) UEditor() {
	var response = c.Ctx.ResponseWriter
	var request = c.Ctx.Request

	var action = request.FormValue("action")
	if action == ACTION_CONFIG {
		uef.ConfigData(response, request)
	} else if action == ACTION_LIST_IMAGE {
		uef.ListImage(response, request)
	} else if action == ACTION_LIST_FILE {
		uef.ListFile(response, request)
	} else if action == ACTION_UPLOAD_IMAGE {
		uef.UploadImage(response, request)
	} else if action == ACTION_UPLOAD_FILE {
		uef.UploadFile(response, request)
	}
}
