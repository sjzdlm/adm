package uef

import (
	"net/http"
)

func ListImage(response http.ResponseWriter, request *http.Request) {
	var picPathList, start, total = getList(request, config.ImagePath)
	listResult(picPathList, start, total, response)
}
