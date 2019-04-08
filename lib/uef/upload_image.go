package uef

import (
	"net/http"
)

func UploadImage(response http.ResponseWriter, request *http.Request) {
	upload(config.ImagePath, response, request)
}
