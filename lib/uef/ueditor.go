package uef

import (
	"net/http"
)

const ACTION_CONFIG = "config"
const ACTION_LIST_IMAGE = "listimage"
const ACTION_LIST_FILE = "listfile"
const ACTION_UPLOAD_IMAGE = "uploadimage"
const ACTION_UPLOAD_FILE = "uploadfile"

func UEditor(response http.ResponseWriter, request *http.Request) {
	var action = request.FormValue("action")
	if action == ACTION_CONFIG {
		ConfigData(response, request)
	} else if action == ACTION_LIST_IMAGE {
		ListImage(response, request)
	} else if action == ACTION_LIST_FILE {
		ListFile(response, request)
	} else if action == ACTION_UPLOAD_IMAGE {
		UploadImage(response, request)
	} else if action == ACTION_UPLOAD_FILE {
		UploadFile(response, request)
	}
}
