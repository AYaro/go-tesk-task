all:
	GOOS=linux GOARCH=amd64 go build -o parts cmd/server/main.go