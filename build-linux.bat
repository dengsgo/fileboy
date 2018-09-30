@echo off
echo build start...
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o fileboy-linux-amd64.bin
echo build success
pause