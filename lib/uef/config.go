package uef

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
)

type configs struct {
	ImagePath string
	FilePath  string
}

var config *configs

func init() {
	var result, err = getConfig()
	if err != nil {
		panic(err)
	}
	if result.FilePath == "" {
		panic(errors.New("附件存放路径未指定"))
	}
	if result.ImagePath == "" {
		panic(errors.New("图片上传路径未指定"))
	}
	config = result
}

func getConfig() (config *configs, err error) {
	var result = new(configs)
	result.FilePath = "static/upload/file"
	result.ImagePath = "static/upload/image"

	file, err := os.Open("config_ue.xml")
	if err != nil {
		return result, nil
		//return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return result, nil
		//return nil, err
	}

	err = xml.Unmarshal(data, result)
	if err != nil {
		return result, nil
		//return nil, err
	}
	return result, nil
}
