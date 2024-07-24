go实现linux机器批量自动化安装filebeat采集功能

打包： make all

# 1、安装filebeat
### 1)本地安装
     ./ed
## 2)远程安装
    ./ed -a reinstall -m remote -h 服务器地址(多个以,分隔) -u root -p 服务器密码 -a install

# 2、卸载filebeat
    ./ed -a uninstall -m remote -h 服务器地址(多个以,分隔) -u root -p 服务器密码 -a install


# 适用系统：
    centos7，ubuntu，centos6，centos8,suse
