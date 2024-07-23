package process

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ed/gen"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

func (sc *SSHConfig) RunCommand(cmd string) string {
	session, err := sc.Client.NewSession()
	defer session.Close()
	// 执行  命令
	log.Info().Msgf("【%s】start exec cmd： %s", sc.Host, cmd)
	output, err := session.CombinedOutput(cmd)

	if err != nil {
		log.Fatal().
			Msgf("Failed to host : %s run:  【%s】 error:【%s】 output: 【%s】", fmt.Sprintf("%s:%d", sc.Host, sc.Port), cmd, err.Error(), output)
	}
	log.Debug().Msgf("【%s】 command : %s, output: %s", sc.Host, cmd, output)
	return string(output)
}

func (sc *SSHConfig) RunSudoCommand(cmds string) string {
	var output []byte
	var err error
	for _, cmd := range strings.Split(cmds, "&&") {
		suCmd := fmt.Sprintf(" echo \"%s\" | sudo -S %s", sc.SudoPassword, cmd)
		// 执行  命令
		log.Info().Msgf("【%s】start exec cmd： %s", sc.Host, suCmd)
		if sc.IsLocal {
			output, err = exec.Command("sh", "-c", suCmd).Output()
		} else {
			session, err := sc.Client.NewSession()

			output, err = session.CombinedOutput(suCmd)
			defer session.Close()
			if err != nil {
				log.Warn().Msgf("Failed to host : %s run %s: %s", fmt.Sprintf("%s:%d", sc.Host, sc.Port), cmd, err.Error())

				// 针对centos6，会出现 sudo: sorry, you must have a tty to run sudo   , 如果真有centos6，并且不是root，会跑不了
				session, err = sc.Client.NewSession()
				defer session.Close()
				output, err = session.CombinedOutput(cmd)
			}
		}
		if err != nil {
			log.Warn().Msgf("Failed to host : %s run %s: %s", fmt.Sprintf("%s:%d", sc.Host, sc.Port), cmd, err.Error())
		}

		log.Debug().Msgf("【%s】 command : %s, output: %s", sc.Host, cmd, output)
	}
	return string(output)
}

func (sc *SSHConfig) RunCommandIgnoreError(cmd string) string {
	session, err := sc.Client.NewSession()
	defer session.Close()
	// 执行  命令
	log.Info().Msgf("【%s】start exec cmd： %s", sc.Host, cmd)
	output, err := session.CombinedOutput(cmd)

	if err != nil {
		log.Warn().Msgf("Failed to host : %s run %s: %s", fmt.Sprintf("%s:%d", sc.Host, sc.Port), cmd, err.Error())
	}
	log.Debug().Msgf("【%s】 command : %s, output: %s", sc.Host, cmd, output)
	return string(output)
}

func (sc *SSHConfig) ScpFile(lPath string, rPath string) error {
	fName := strings.Split(lPath, "/")[len(strings.Split(lPath, "/"))-1]
	log.Info().Msgf("BasePath : %s , lPath: %s, rPath: %s fName: %s", sc.BasePath, lPath, rPath, fName)
	rPath = sc.BasePath + "/" + fName
	log.Info().Msgf("【%s】start ScpFile：  %s to %s", sc.Host, lPath, rPath)
	var err error
	content, err := gen.Asset(lPath)

	if err != nil {
		log.Fatal().Msgf("【%s】Error reading embedded file: %v, lpath: %s", sc.Host, err, lPath)
	}
	// 将 byte 切片转换为 io.Reader
	reader := bytes.NewReader(content)

	if sc.IsLocal {
		err = os.MkdirAll(sc.BasePath, 0777)
		if err != nil {
			log.Fatal().Msgf("【%s】Error creating directory: %v, rpath: %s", sc.Host, err, rPath)
		}
		log.Debug().Msgf("【%s】BasePath : %v, rpath : %s, lPath: %s , path: %s", sc.Host, sc.BasePath, rPath, lPath, sc.BasePath)
		destFile, err := os.Create(sc.BasePath + "/" + filepath.Base(rPath))
		if err != nil {
			log.Fatal().Msgf("【%s】Error creating destination file: %s", sc.Host, err)
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, reader)
		if err != nil {
			log.Fatal().Msgf("【%s】Error copying file:%s", sc.Host, err)
			return err
		}
	} else {
		err = sc.RScpFile(fName, *reader, rPath)
	}
	return err
}
func (sc *SSHConfig) RScpFile(fName string, reader bytes.Reader, rPath string) error {

	// 创建 SSH 客户端配置
	clientConfig, err := auth.PasswordKey(sc.Username, sc.Password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		log.Fatal().Msgf("【%s】Failed to create SSH client config: %s", sc.Host, err)
		return err
	}
	// 创建 SCP 客户端
	client := scp.NewClient(fmt.Sprintf("%s:%d", sc.Host, sc.Port), &clientConfig)

	// 连接到服务器
	err = client.Connect()
	if err != nil {
		log.Fatal().Msgf("Failed to connect to server: %s", err)
		return err
	}
	defer client.Close()

	// 上传文件到服务器,考虑到有些没有权限，先上传到/tmp下，然后移动到目标路径
	err = client.CopyFile(context.Background(), &reader, "/tmp/"+fName, "0644")
	if err != nil {
		log.Fatal().Msgf("Failed to upload file: %s, tmpPath: %s", err, "/tmp/"+fName)
		return err
	}
	sc.RunSudoCommand(fmt.Sprintf("mv /tmp/%s %s", fName, rPath))
	log.Printf("File uploaded successfully to %s", rPath)
	client.Close()
	return nil
}
