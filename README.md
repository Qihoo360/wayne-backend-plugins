# Wayne 后端插件仓库

存放[Wayne](https://github.com/Qihoo360/wayne)项目后端插件。

## 新建插件

该插件主要包含三部分

- controller：创建controller文件夹，实现controller函数

- model：用于跟数据库交互

- router：注册路由。commentsRouter_*文件为运行make run-backend时自动生成的文件

新开发的插件根目录下需包含init.go,引入router包。然后再plugins.go中引入init.go
