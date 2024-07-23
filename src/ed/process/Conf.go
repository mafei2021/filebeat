package process

import "golang.org/x/crypto/ssh"

type Config struct {
	Host         string `yaml:"host"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	SudoPassword string `yaml:"sudo_password"`
	Port         int    `yaml:"port"`
	Action       string `yaml:"action"`
	IsLocal      bool   `yaml:"is_local"`
	BasePath     string `yaml:"base_path"`
}

type SSHConfig struct {
	Username     string
	Password     string
	SudoPassword string
	Host         string
	Port         int
	Client       *ssh.Client
	IsLocal      bool
	BasePath     string
}
