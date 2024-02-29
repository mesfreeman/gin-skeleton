# gin-skeleton

一个基于 [GIN](https://github.com/gin-gonic/gin) 框架封装的 WEB 项目骨架，旨在快速开启 WEB 相关的项目开发。

> Go版本依赖最好 >= 1.18.0，为了更好的支持 go mod 包管理及泛型操作。
> 
> 该项目为后端项目，配合前端项目 [vben-skeleton](https://github.com/mesfreeman/vben-skeleton) 项目，直接拥有一个完整的管理后台。

## 基础功能

* 支持优雅重启
* 支持日志记录
* 支持配置热更新
* 支持路由文件分隔
* 支持 `GORM` 查询
* 支持 `Redis` 查询
* 支持 `jwt`、`sign` 中间件 
* 支持 `cobra cli` 命令行脚本
* 支持 `rabc` 权限模型

## 目录结构

```text
@todo 待补充
```

## 部署说明

使用`Makefile`来完成项目初始化及服务的启动、重启等操作，如下：

### 一、本地开发

#### 1. 项目初始化

```shell
go mod tidy
cp ./config/config.yaml.example ./config/config.yaml
chown -R www:www storage
```
注：手动调整配置文件 `config.yaml` 中的相关配置。

#### 2. 服务启动

```shell
go run main.go
```

### 二、线上部署

使用部署脚本 `deploy.sh`，具体使用方法如下：

```shell
./deploy.sh [server] [project_path]
```

说明：

* `server`：服务器地址，默认值：`tank.server.cn`
* `project_path`：项目路径，默认值：`/data/services/projects/gin-skeleton/code`
* 注：请基于实际情况调整 `deploy.sh` 中的相关配置。
