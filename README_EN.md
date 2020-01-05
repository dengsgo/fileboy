## FILEBOY  

[![Build Status](https://travis-ci.org/dengsgo/fileboy.svg?branch=master)](https://travis-ci.org/dengsgo/fileboy) [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/fileboy)](https://goreportcard.com/report/github.com/dengsgo/fileboy)

[简体中文](README.md) | [ENGLISH](README_EN.md)

Fileboy, File Change Monitoring Notification System, written with Go.  
For Hot Reload scenarios (typically for developing go projects without having to perform go build manually every time; for example, front-end node packaging) or system monitoring.  

## FEATURES

- Minimalist usage and configuration  
- Support multiple platforms, Windows/Linux/MacOS  
- Supports custom file listening scope, listening for specified folders/not listening for specified folders/specified suffix files  
- Support for custom monitoring events (write / rename / remove / create / chmod)  
- Support for setting up multiple commands  
- Command support variable placeholders  
- Supporting redundant task discarding and customizing redundant task scope  
- Supporting HTTP notifications  
- more...  

## COMPILE

go version 1.13   

## CHANGELOG  

[CHANGELOG](CHANGELOG.md)  


## RUN    

### BINARIES   

Github: [download v1.15](https://github.com/dengsgo/fileboy/releases)  
Gitee:  [dowmload v1.15](https://gitee.com/dengsgo/fileboy/releases)  

Download the compiled binary file of the corresponding platform, rename it `fileboy`, and add it to the system Path.  

### SOURCE   

Clone project, enter the project directory, run the command:  
```shell
## installation dependency
go get-u gopkg.in/fsnotify/fsnotify.v1
go get-u gopkg.in/yaml.v2
## compile
go build
## run
./fileboy
```
## USAGE

The normal operation of fileboy depends on the `filegirl.yaml` configuration, so for the first time in a project, `filegirl.yaml` needs to be initialized.  
- Enter the project home directory where you want hot reload;  
- Running `fileboy init` will generate `filegirl. yaml` files in this directory.  
- View `filegirl. yaml` and modify it to a configuration item suitable for your project.  
- Run `fileboy`.  
  
If you define the `command-> exec` command, you can run the `fileboy exec` command to confirm whether it can be executed properly in advance. The system will try to run your custom command.  
You can use `fileboy help` to see help info.  

## filegirl.yaml

```yaml
core:
    # config version code
    version: 1

# monitor section
monitor:
    # directories to monitor
    # test1       listen for the test1 directory in the project directory
    # test1/test2 listen for the test1/test2 directory in the project directory
    # test1,*     listen for the test1 directory in the project directory and all its subdirectories (recursion)
    # .,*         listen for the project directory and all its subdirectories (recursion)
    includeDirs:
        - .,*

    # Unmonitored directories
    # .idea   ignore listening to .idea directory and all its subdirectories
    exceptDirs:
        - .idea
        - .git
        - .vscode
        - node_modules
        - vendor

    # the suffix of the listener file, which changes the file to execute commands
    # .go   file changes suffixed with .go execute commands
    # .*    all file changes execute commands in the command
    types:
        - .go

    # the type of event to listen to. Only when such an event occurs can the command in command be executed
    # without this configuration, all events will be monitored by default
    # write    write file event
    # rename   rename file event
    # remove   remove remove file event
    # create   create file event
    # chmod    update file permission event (UNIX like)
    events:
        - write
        - rename
        - remove
        - create
        - chmod

command:
    # the files monitored have commands that change to be executed
    # there can be multiple commands that will be executed in turn
    # in case of interactive commands, allow external access to input
    # variable placeholders are supported, and the actual values are replaced when the command is run:
    #    {{file}}    (e.g: a.txt 、test/test2/a.go)
    #    {{ext}}     (e.g: .go)
    #    {{event}}   event(e.g: write)
    #    {{changed}} local timestamp for file updated(nanosecond,e.g 1537326690523046400)
    # variable placeholders e.g：cp {{file}} /root/sync -rf  、 myCommand --{{ext}} {{changed}}
    exec:
        - go version
        - go env

    # the command will not execute until XX milliseconds after the file changes
    # a change event (A) cancels execution if there is a new file change event (B) within the defined delay time (t).
    # B and subsequent events are analogized in turn until event Z does not produce new events within t, and Z executes.
    # reasonable setting of delay time will effectively reduce redundancy and duplicate task execution.
    # If this feature is not required, set to 0
    delayMillSecond: 2000

    # 特殊指令
    # 可以通过预定义的指令来控制 command 的行为，指令可以有多个
    #   exec-when-start    fileboy启动就绪后，自动执行一次 'exec' 定义的命令
    #   should-finish      触发执行 'exec' 时(C)，如果上一次的命令(L)未退出（还在执行），会等待 L 退出，直到有明确 exit code 才会开始执行本次命令。
    #                      在等待 L 退出时，又有新事件触发了命令执行(N)，则 C 执行取消，只会保留最后一次的 N 执行
    #   ignore-stdout      执行 'exec' 产生的 stdout 会被丢弃
    #   ignore-warn        fileboy 自身的 warn 信息会被丢弃
    #   ignore-info        fileboy 自身的 info 信息会被丢弃
    #   ignore-exec-error  执行 'exec' 出错仍继续执行下面的命令而不退出
    instruction:
        #- should-finish
        #- exec-when-start

notifier:
    # file changes send requests to the URL (POST JSON text data)
    # the timing of triggering the request is consistent with executing the command command
    # timeout 15 second
    # POST :
    #    Content-Type: application/json;charset=UTF-8
    #    User-Agent: FileBoy Net Notifier v1.15
    #    Body: {"project_folder":"/project/path","file":"main.go","changed":1576567861913824940,"ext":".go","event":"write"}
    # e.g: http://example.com/notifier/fileboy-listener
    # no notice is enabled. Please leave it blank. ""
    callUrl: ""
```

### CONTRIBUTOR

[@dengsgo](https://www.yoytang.com)  <dengsgo@gmail.com>  

[@itwesley](https://github.com/itwesley)  <wcshen1126@gmail.com>  

[@jason-gao](https://github.com/jason-gao)  <3048789891@qq.com>  

