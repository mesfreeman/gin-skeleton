# gin-skeleton

一个基于 [GIN](https://github.com/gin-gonic/gin) 框架封装的 WEB 项目骨架，旨在快速开启 WEB 相关的项目开发。

> 💡 Go版本依赖最好 >= 1.18.0，为了更好的支持 go mod 包管理及泛型操作。
> 
> 该项目为后端项目，配合前端项目 [vben-skeleton](https://github.com/mesfreeman/vben-skeleton) 项目，直接拥有一个完整的管理后台。

## ✨ 框架特性

* 🍥 支持优雅重启
* 🍤 支持日志记录
* 🍣 支持配置热更新
* 🍔 支持路由文件分隔
* 🍕 支持 `GORM` 查询
* 🌮 支持 `Redis` 查询
* 🍵 支持 `jwt`、`sign` 中间件 
* 🍟 支持 `cobra cli` 命令行脚本
* 🍭 支持 `rabc` 权限模型

## 🌴 目录结构

```text
@todo 待补充
```

## 📖 部署说明

⚠️ Mysql 相关表结构在 `template/sql` 目录下，开发前请自行导入。

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
* 注：该脚本依赖 `pm2` 服务，请自行安装并基于实际情况调整 `deploy.sh` 中的相关配置。

## 🎨 后台截图

前端项目 - [传送门](https://github.com/mesfreeman/vben-skeleton)

|                                   🙅 账号管理                                    |                                   📚 菜单管理                                    |
|:----------------------------------------------------------------------------:|:----------------------------------------------------------------------------:|
| ![账号管理.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.03.55.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) | ![菜单管理.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.04.38.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) |
|                                   👨 角色管理                                    |                                   🗂 文件管理                                    |
| ![角色管理.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.04.53.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) | ![文件管理.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.06.01.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) |
|                                   📨 邮件管理                                    |                                   📝 登录日志                                    |
|      ![邮件管理.png](https://file.dandy.fun/picgo/swap/202406172233282.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x)      | ![登录日志.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.06.36.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) |
|                                   🔍 操作日志                                    |                                   🌓 暗黑模式                                    |
| ![操作日志.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.34.23.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) | ![暗黑模式.png](https://file.dandy.fun/picgo/swap/iShot_2022-11-08_17.09.37.png?imageView2/0/q/70%7Cimageslim%7CimageMogr2/thumbnail/650x) |
