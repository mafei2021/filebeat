package main

import (
	"context"
	vr "ed/src/ed/process/var"
	"ed/src/ed/utils"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/signal"
	"syscall"
	"time"

	"ed/src/ed/initialize"
	"ed/src/ed/process"
)

func main() {
	fmt.Println("starting")
	ctx, serverStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer serverStop()

	// 读取命令行参数
	var configPath string

	var model string
	flag.StringVar(&configPath, "c", "", "yaml配置文件所在的位置")
	flag.StringVar(&model, "m", "local", "模式，支持远程部署agent或者直接执行本机执行  local/remote")
	var c process.Config
	flag.StringVar(&c.BasePath, "l", "/opt/ed/", "部署的路径 ， 默认为/opt/ed")
	flag.StringVar(&c.Action, "a", "install", "要执行的操作,默认安装agent , install,uninstall")
	// remote模式用到
	flag.StringVar(&c.Host, "h", "", "要操作的主机IP，多个以,分隔")
	flag.StringVar(&c.Username, "u", "", "要操作的主机用户名")
	flag.StringVar(&c.Password, "p", "", "要操作的主机密码")
	flag.StringVar(&c.SudoPassword, "p2", "", "要操作的主机SUDO密码")
	flag.IntVar(&c.Port, "port", 22, "要操作的主机端口")
	flag.Parse()

	initialize.Init(configPath)
	switch model {
	case vr.ModelRemote:
		fmt.Println("Remote")
		break
	case vr.ModelLocal:
		log.Debug().Msgf("local模式")
		c.IsLocal = true
		c.Host = "127.0.0.1"
		break
	default:
		log.Fatal().Msgf("unknown model please user local/remote model")
	}
	// 打印c的所有参数和对应的值
	utils.PrintStruct(c)
	c.InitProcess(ctx)

	// 处理关闭
	<-ctx.Done()
	serverStop()
	//log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// 这里用来强制断开所有的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//log.Info().Msg("Server exiting")
}
