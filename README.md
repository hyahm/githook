# githook


githook

这个东西就是git 的钩子，每次push 会自动拉取代码  

# 配置文件
hook.ini
```
listen=:10009
jsondir=json
# 读取超时时间， 默认30s
readtimeout=30s
# gitee WebHook 密码
# github secret
# gitlab Secret Token
token=123456
```


# json文件
项目设置集成中， 填写url
Secret Token： 默认123456，  配置文件可以修改

json 文件如下， 默认在项目的json 目录下， 配置文件可以修改
```json
{
    "user": "nginx",   // windows 不支持， 默认会清空
    "command": "git pull",
    "dir": "/home/app",
    "shell": "/bin/bash"  // windows 默认使用powershell, linux 和mac 默认/bin/bash
}
```