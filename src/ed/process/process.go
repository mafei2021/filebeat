package process

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func (conf *Config) InitProcess(ctx context.Context) {

	if !conf.IsLocal && conf.Host == "" {
		log.Fatal().Msgf("至少传一个IP过来")
	}

	for _, host := range strings.Split(conf.Host, ",") {
		log.Debug().Msgf("start init %s", host)
		var sshConfig = &SSHConfig{
			Host:     host,
			Port:     conf.Port,
			Username: conf.Username,
			Password: conf.Password,
			IsLocal:  conf.IsLocal,
			BasePath: conf.BasePath,
		}
		// 如果没有输入sudo，则共用密码
		if conf.SudoPassword == "" {
			sshConfig.SudoPassword = conf.Password
		}
		switch conf.Action {
		case "reinstall":
			sshConfig.uninstall()
			sshConfig.install(ctx)
		case "install":
			sshConfig.install(ctx)
		case "uninstall":
			sshConfig.uninstall()
		default:
			log.Fatal().Msgf("action %s not support", conf.Action)
		}

	}
	os.Exit(0)
}

func (sc *SSHConfig) Baseinstall(ctx context.Context) {
	if !sc.IsLocal {
		sc.SshConn()
	}
	//1、创建目录
	sc.RunSudoCommand("mkdir -p " + sc.BasePath)

	// 2、上传filebeat文件
	filebeat := "agent.tar.gz"
	ls := sc.RunSudoCommand("ls " + sc.BasePath)
	if !strings.Contains(ls, filebeat) {
		// 上传/opt/ed/agent.gz 到服务器上
		err := sc.ScpFile("asset/"+filebeat, "")
		if err != nil {
			log.Fatal().Msgf("1Failed to run ScpFile: %s", err)
			return
		}
		// 解压filebeat
		sc.RunSudoCommand(
			fmt.Sprintf("cd %s && tar zxvf %s%s -C %s && chmod -R a+wx %sagent", sc.BasePath, sc.BasePath, filebeat, sc.BasePath, sc.BasePath),
		)
	}
	if !strings.Contains(ls, "config.yml") {
		// 上传/opt/ed/filebeat.yml 到10.0.81.98服务器上
		err := sc.ScpFile("asset/config.yml", "config.yml")
		if err != nil {
			log.Fatal().Msgf("2Failed to run ScpFile: %s", err)
			return
		}
		sc.RunSudoCommand(fmt.Sprintf("chmod -R a+wx %sconfig.yml", sc.BasePath))
		// ubuntu 报错： Exiting: error loading config file: config file ("/opt/ed/config.yml") can only be writable
		// by the owner but the permissions are "-rwxrwxrwx" (to fix the permissions use: 'chmod go-w
		// /opt/ed/config.yml')
		sc.RunSudoCommand(fmt.Sprintf("chmod go-w %sconfig.yml", sc.BasePath))
		sc.RunSudoCommand(fmt.Sprintf(" chown root:root %sconfig.yml", sc.BasePath))
	}
}

func (sc *SSHConfig) install(ctx context.Context) {

	sc.Baseinstall(ctx)
	l := sc.RunSudoCommand("cat /etc/os-release || cat /etc/system-release || uname -a")
	log.Info().Msgf("当前主机: %s, 操作系统版本: %s", sc.Host, l)
	lsEd := sc.RunSudoCommand(fmt.Sprintf("ls  %s", sc.BasePath))

	//针对centos7/8版本
	if strings.Contains(l, "Linux 7") || strings.Contains(l, "release 7") || strings.Contains(l, "release 8") || strings.Contains(l, "Ubuntu") {
		if !strings.Contains(lsEd, "tagent.service") {
			// 上传/opt/ed/config.yml 到指定服务器上
			err := sc.ScpFile("asset/tagent.service", "tagent.service")
			if err != nil {
				log.Fatal().Msgf("Failed to run ScpFile: %s", err)
				return
			}
			sc.RunSudoCommand(
				"systemctl daemon-reload && rm -rf /lib/systemd/system/tagent.service && rm -rf /etc/systemd/system/multi-user.target.wants/tagent.service && ln -s " + sc.BasePath + "tagent.service /lib/systemd/system/tagent.service && systemctl start tagent && systemctl enable tagent",
			)
		}
		// 针对centos6版本
	} else if strings.Contains(l, "release 6") {
		if !strings.Contains(lsEd, "tagent") {
			err := sc.ScpFile("asset/tagent", "tagent")
			if err != nil {
				log.Fatal().Msgf("Failed to run ScpFile: %s", err)
				return
			}
		}
		sc.RunSudoCommand("systemctl daemon-reload && rm -rf /etc/init.d/tagent && ln -s " + sc.BasePath + "tagent /etc/init.d/tagent && chmod +x /etc/init.d/tagent && chkconfig --add tagent && chkconfig tagent on && service tagent start")

	} else if strings.Contains(l, "FreeBSD 11") {
		if !strings.Contains(lsEd, "tagent") {
			err := sc.ScpFile("asset/freebsd/tagent", "tagent")
			if err != nil {
				log.Fatal().Msgf("Failed to run ScpFile: %s", err)
				return
			}
		}
		sc.RunSudoCommand("rm -rf /usr/local/etc/rc.d/tagent && ln -s " + sc.BasePath + "tagent /usr/local/etc/rc.d/tagent && chmod +x /usr/local/etc/rc.d/tagent && echo 'tagent_enable=\"YES\"' >> /etc/rc.conf && service tagent start")
	} else {
		log.Fatal().Msgf("not support os version")
	}
	log.Info().Msgf(" install tagent service successfully")
	//os.Exit(0)
}

func (conf *SSHConfig) uninstall() {
	if !conf.IsLocal {
		conf.SshConn()
	}
	// 删除服务文件
	// 删除centos7/8版本
	conf.RunSudoCommand(
		"systemctl disable tagent && systemctl stop tagent && rm -rf /etc/systemd/system/tagent.service && rm -rf /lib/systemd/system/tagent* && rm -rf /etc/systemd/system/multi-user.target.wants/tagent.service",
	)

	// 删除centos6版本
	conf.RunSudoCommand("service tagent stop && rm -rf /etc/init.d/tagent")
	conf.RunSudoCommand("rm -rf /opt/ed")
	log.Info().Msgf("【%s】 uninstalling tagent service successfully", conf.Host)
}

func (conf *SSHConfig) SshConn() {
	// 配置 SSH 客户端
	config := &ssh.ClientConfig{
		User: conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到 SSH 服务器
	log.Info().Msgf(" start ssh connection to %s", conf.Host+":"+fmt.Sprintf("%d", conf.Port))
	client, err := ssh.Dial("tcp", conf.Host+":"+fmt.Sprintf("%d", conf.Port), config)
	if err != nil {
		log.Fatal().Msgf("Failed to dial: %s", err)
	}
	conf.Client = client
	//defer client.Close()
}
