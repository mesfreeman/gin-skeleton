# gin-skeleton

一个基于 [GIN](https://github.com/gin-gonic/gin) 框架封装的 WEB 项目骨架，旨在快速开启 WEB 相关的项目开发。

> Go版本依赖最好 >= 1.13.4，为了更好的支持 go mod 包管理。

## 基础功能

* 支持优雅重启
* 支持日志记录
* 支持配置热更新
* 支持路由文件分隔
* 支持 `GORM` 查询
* 支持 `Redis` 查询
* 支持 `jwt`、`sign` 中间件
* 支持 `cobra cli` 命令行脚本

## 部署说明

### 编译

#### WIN 环境下

```shell
# 设置环境变量
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

# 编译主程序与命令行脚本
go build -o release/gin-skeleton main.go
go build -o release/gin-cli cmd/main.go
```

#### MAC 环境下

```shell
# 设置环境变量
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 

# 编译主程序与命令行脚本
go build -o release/gin-skeleton main.go
go build -o release/gin-cli cmd/main.go
```

### 部署

```shell
# 拷贝配置
cp -R config release

# 建立日志软链
cd release && ln -s ../storage storage

# 修改权限
sudo chmod +x gin-skeleton gin-cli
sudo chown -R www-data:www-data storage

# 使用PM2启动主程序
cd .. && pm2 start pm2.json
```
