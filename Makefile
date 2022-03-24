.PHONY: run
# run
run:
	cp ./config/config.yaml.example ./config/config.yaml
	go run main.go

.PHONY: publish
# publish
publish:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/gin-skeleton main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/gin-cli cmd/main.go
	git add . && git commit -m "build new version" && git push