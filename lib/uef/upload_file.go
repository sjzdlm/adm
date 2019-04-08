package uef

import (
	"net/http"
)

func UploadFile(response http.ResponseWriter, request *http.Request) {
	upload(config.FilePath, response, request)
}
