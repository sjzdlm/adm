package uef

import (
	"net/http"
)

func ListFile(response http.ResponseWriter, request *http.Request) {
	var picPathList, start, total = getList(request, config.FilePath)
	listResult(picPathList, start, total, response)
}
