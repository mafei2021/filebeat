package test

import (
	"flag"
	"fmt"

	"ed/src/ed/initialize"
	"ed/src/ed/process"
)

func BaseInit() process.Config {
	fmt.Println("starting")
	// 读取命令行参数
	var configPath string

	var c process.Config

	flag.StringVar(&configPath, "c", "/opt/ed/resources/application.yaml", "yaml配置文件所在的位置")
	flag.StringVar(&c.Host, "h", "", "要操作的主机IP，多个以,分隔")
	flag.StringVar(&c.Username, "u", "", "要操作的主机用户名")
	flag.StringVar(&c.Password, "p", "", "要操作的主机密码")
	flag.StringVar(&c.SudoPassword, "p2", "", "要操作的主机SUDO密码")
	flag.IntVar(&c.Port, "port", 22, "要操作的主机端口")
	flag.StringVar(&c.Action, "a", "install", "要执行的操作,默认安装agent , install,uninstall")

	flag.Parse()
	initialize.Init(configPath)
	return c

}
