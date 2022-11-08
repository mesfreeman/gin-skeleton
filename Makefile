# 常用命令

# 本地初始化
devinit:
	go mod download
	cp ./config/config.yaml.example ./config/config.yaml && chown -R www:www storage
	@echo "本地初始化完成，请前往[config/config.yaml]文件修改相关配置"

# 本地开发模式运行
devrun:
	go run main.go

# 编译并提交新版本
buildpush:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/gin-skeleton main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/gin-cli cmd/main.go
	git add . && git commit -m "build new version: v1.0.0" && git push

# 线上初始化
proinit:
	cp ./config/config.yaml.example ./config/config.yaml && \
	chown -R www:www storage && ln -s ../../storage release/storage && ln -s ../../config release/config
	@echo "线上初始化完成，请前往[config/config.yaml]文件修改相关配置"

# 线上服务启动
prostart:
	chmod +x ./release/gin-skeleton && pm2 start pm2.json
	@echo "[gin-skeleton]服务已启动"

# 线上服务重启
prorestart:
	pm2 restart gin-skeleton
	@echo "[gin-skeleton]服务已重启"

