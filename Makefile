build:
	@CGO_ENABLED=1 GOOS="darwin" GOARCH="amd64" go build -o ./bin/main
run: build
	@go run main.go
