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

### 本地开发

```shell
make run
```

### 线上部署

```shell
# 编译并发布代码
make publish

# 前往服务器
...

# 拉取仓库代码
...

# 拷贝配置
cp -R config release

# 修改配置数据
...

# 建立日志软链
cd release && ln -s ../storage storage

# 修改权限
sudo chmod +x gin-skeleton gin-cli
sudo chown -R www-data:www-data storage

# 使用PM2启动主程序
cd .. && pm2 start pm2.json
```

### 持续集成

```shell
# 编译并发布代码
make publish

# 前往服务器
...

# 拉取仓库代码
...

# 服务重启
PM2 会监听二进制文件并自动重启服务
```
