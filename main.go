package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/kardianos/service"
	"github.com/sjzdlm/adm/conf"
	_ "github.com/sjzdlm/adm/routers"
	"github.com/sjzdlm/db"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// 重置模板路径
	beego.SetViewsPath(GetAPPRootPath() + "/views/")
	//初始化数据库连接
	db.InitX()
	//初始化adm配置,静态目录和注入函数
	conf.InitConfig()

	beego.Run()
}

func (p *program) Stop(s service.Service) error {
	return nil
}
func main() {
	var srvname = beego.AppConfig.String("srvname")
	if srvname == "" {
		srvname = "fooapp"
	}
	svcConfig := &service.Config{
		Name:        srvname, //服务显示名称
		DisplayName: srvname, //服务名称
		Description: srvname, //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		beego.Error(err.Error())
	}

	if err != nil {
		beego.Error(err.Error())
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			fmt.Println("服务安装成功")
			return
		}

		if os.Args[1] == "uninstall" {
			s.Uninstall()
			fmt.Println("服务卸载成功")
			return
		}
	}

	err = s.Run()
	if err != nil {
		beego.Error(err)
	}
}

//应用程序根路径
func GetAPPRootPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}
