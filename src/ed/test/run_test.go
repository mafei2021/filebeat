package test

import (
	"fmt"
	"strings"
	"testing"

	"ed/src/ed/process"
	"github.com/rs/zerolog/log"
)

//type process.SSHConfig struct {
//	Username string
//	Password string
//	Host     string
//	Port     int
//	Client   *ssh.Client
//}

func TestRun2(t *testing.T) {
	BaseInit()
	//-h 10.0.70.159 -u tophant -p 1 -a reinstall
	//var sc = &process.SSHConfig{
	//	Host:         "10.0.70.159",
	//	Port:         22,
	//	Username:     "tophant",
	//	Password:     "1",
	//	SudoPassword: "1",
	//}
	//sc.SshConn()
	//command := sc.RunCommand(" echo \"1\" | sudo -S ls /opt/ed")
	//log.Info().Msgf(command)
	//
	//command2 := sc.RunCommand(" echo \"1\" | sudo -S tar zxvf /opt/ed/agent.tar.gz -C /opt/ed")
	//log.Info().Msgf(command2)
	// 读取嵌入的配置文件内容
	//content, err := generated.Asset("assets/application.yaml")
	//if err != nil {
	//	fmt.Println("Error reading embedded file:", err)
	//	return
	//}
	//
	//fmt.Println("Configuration content:")
	//fmt.Println(string(content))

}
func TestRun(t *testing.T) {
	BaseInit()

	var sc = process.SSHConfig{
		Username: "root",
		Password: "password",
		Host:     "10.0.81.8",
		Port:     22,
	}

	sc.SshConn()
	//1、创建目录
	command := "mkdir -p /opt/ed/"
	sc.RunCommand(command)

	// 2、上传filebeat文件
	filebeat := "agent.zip"
	ls := sc.RunCommand("ls /opt/ed ")
	if !strings.Contains(ls, "agent") {
		// 上传/opt/ed/agent.gz 到10.0.81.98服务器上
		err := sc.ScpFile("/opt/ed/asset/"+filebeat, "/opt/ed/"+filebeat)
		if err != nil {
			log.Fatal().Msgf("Failed to run ScpFile: %s", err)
			return
		}
		// 解压filebeat
		sc.RunCommand(fmt.Sprintf("cd /opt/ed/ && unzip /opt/ed/%s", filebeat))
	}
	if !strings.Contains(ls, "config.yml") {
		// 上传/opt/ed/filebeat.yml 到10.0.81.98服务器上
		err := sc.ScpFile("/opt/ed/asset/config.yml", "/opt/ed/config.yml")
		if err != nil {
			log.Fatal().Msgf("Failed to run ScpFile: %s", err)
			return
		}
	}
	if !strings.Contains(ls, "tagent.service") {
		// 上传/opt/ed/filebeat.yml 到10.0.81.98服务器上
		err := sc.ScpFile("/opt/ed/asset/tagent.service", "/opt/ed/tagent.service")
		if err != nil {
			log.Fatal().Msgf("Failed to run ScpFile: %s", err)
			return
		}
		sc.RunCommand(
			"rm -rf /lib/systemd/system/tagent.service && rm -rf /etc/systemd/system/multi-user.target.wants/tagent.service && ln -s /opt/ed/tagent.service /lib/systemd/system/tagent.service && systemctl start tagent && systemctl enable tagent",
		)
	}

	////// 创建新的 SSH 会话
	//session, err := client.NewSession()
	//if err != nil {
	//	log.Fatalf("Failed to create session: %s", err)
	//}
	//defer session.Close()
	//
	//// 执行 `mkdir /opt/ed/` 命令
	//output, err := session.CombinedOutput("ls /opt/ed ")
	//if err != nil {
	//	log.Fatalf("Failed to run mkdir /opt/ed/: %s", err)
	//}
	//fmt.Printf("Output of ls: %s\n", output)

	//session, err = client.NewSession()
	//// 执行 `df -ah` 命令
	//output, err = session.CombinedOutput("du -ah /opt/ed")
	//if err != nil {
	//	log.Fatalf("Failed to run du -ah: %s", err)
	//}
	//fmt.Println("Output of df -ah:")

}

//	func (sc *process.SSHConfig) RunCommand(cmd string) string {
//		session, err := sc.Client.NewSession()
//		defer session.Close()
//		if err != nil {
//			log.Fatal().Msgf("Failed to create session: %s", err)
//		}
//		// 执行 `df -ah` 命令
//		log.Info().Msgf("start exec cmd %s", cmd)
//		output, err := session.CombinedOutput(cmd)
//
//		if err != nil {
//			log.Fatal().Msgf("Failed to host : %s run %s: %s", fmt.Sprintf("%s:%d", sc.Host, sc.Port), cmd, err.Error())
//		}
//		log.Debug().Msgf(" command : %s, output: %s", cmd, output)
//		return string(output)
//	}
//
//	func (sc *process.SSHConfig) ScpFile(lPath string, rPath string) error {
//		log.Info().Msgf("start ScpFile %s to %s", lPath, rPath)
//		// 创建 SSH 客户端配置
//		clientConfig, err := auth.PasswordKey(sc.Username, sc.Password, ssh.InsecureIgnoreHostKey())
//		if err != nil {
//			log.Fatal().Msgf("Failed to create SSH client config: %s", err)
//			return err
//		}
//		// 创建 SCP 客户端
//		client := scp.NewClient(fmt.Sprintf("%s:%d", sc.Host, sc.Port), &clientConfig)
//
//		// 连接到服务器
//		err = client.Connect()
//		if err != nil {
//			log.Fatal().Msgf("Failed to connect to server: %s", err)
//			return err
//		}
//		defer client.Close()
//
//		// 打开要上传的文件
//		file, err := os.Open(lPath)
//		if err != nil {
//			log.Fatal().Msgf("Failed to open local file: %s", err)
//			return err
//		}
//		defer file.Close()
//
//		// 上传文件到服务器
//		err = client.CopyFile(context.Background(), file, rPath, "0644")
//		if err != nil {
//			log.Fatal().Msgf("Failed to upload file: %s", err)
//			return err
//		}
//		log.Printf("File uploaded successfully to %s", rPath)
//		client.Close()
//		return nil
//	}
//
//	func (conf *process.SSHConfig) sshConn() {
//		// 配置 SSH 客户端
//		config := &ssh.ClientConfig{
//			User: conf.Username,
//			Auth: []ssh.AuthMethod{
//				ssh.Password(conf.Password),
//			},
//			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//		}
//		// 连接到 SSH 服务器
//		client, err := ssh.Dial("tcp", conf.Host+":"+fmt.Sprintf("%d", conf.Port), config)
//		if err != nil {
//			log.Fatal().Msgf("Failed to dial: %s", err)
//		}
//
//		conf.Client = client
//		//defer client.Close()
//
// }
func TestRun3(t *testing.T) {
	lPath := "/a/b/a/a.txt"
	res := strings.Split(lPath, "/")[len(strings.Split(lPath, "/"))-1]
	fmt.Println(res)
}
