# notice
# Make is very picky about spaces vs. tabs.
# Command lines absolutely must be indented with a single tab, and not spaces.
# You may need to adjust your editor to generate tab characters.
# http://blog.chinaunix.net/uid/28458801/sid-171170-list-1.html

# use:
# make build-mac  编译
# make start-mac 启动

# make build-start-mac 编译+启动
init:
	 go get -u

build-mac: ;@echo "编译-mac版";
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags "-s -w" -o ./bin/fileboy-darwin-amd64.bin

build-linux: ;@echo "编译-linux版";
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o ./bin/fileboy-linux-amd64.bin

build-win: ;@echo "编译-windows版";
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags "-s -w" -o ./bin/fileboy-windows-amd64.exe

build-all: build-mac build-linux build-win

start-mac: ;@echo "启动服务";
	./bin/fileboy-darwin-amd64.bin init
	./bin/fileboy-darwin-amd64.bin

start-linux: ;@echo "启动服务";
	./bin/fileboy-linux-amd64.bin init
	./bin/fileboy-linux-amd64.bin

start-win: ;@echo "启动服务";
	./bin/fileboy-windows-amd64.exe init
	./bin/fileboy-windows-amd64.exe

build-start-mac: build-mac start-mac

.PHONY: build-mac build-linux build-win build-all start-linux travis-linux start-mac start-win build-start-mac
