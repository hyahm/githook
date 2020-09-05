# githook

githook

这个东西就是git 的钩子，每次push 会自动拉取代码, 暂时只做了gitlab

项目设置集成中， 填写url
Secret Token： 默认123456，  配置文件可以修改

json 文件如下， 默认在项目的json 目录下， 配置文件可以修改
```json
{
    "user": "nginx",
    "command": "git pull",
    "dir": "/home/app",
    "shell": "/bin/bash"
}
```