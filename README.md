## 项目说明  

![Go](https://github.com/dengsgo/fileboy/workflows/Go/badge.svg?branch=master) [![Build Status](https://travis-ci.org/dengsgo/fileboy.svg?branch=master)](https://travis-ci.org/dengsgo/fileboy) [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/fileboy)](https://goreportcard.com/report/github.com/dengsgo/fileboy)

[简体中文](README.md) | [ENGLISH](README_EN.md)

fileboy，文件变更监听通知工具，使用 Go 编写。  
适用于 Hot Reload （典型的如开发go项目，无需每次手动执行 go build；又比如前端 node 打包） 或者 系统监控的场景。  

## 特性  

- 极简的用法和配置  
- 支持多平台，Windows/Linux/MacOS  
- 支持自定义文件监听范围，监听指定文件夹/不监听指定文件夹/指定后缀文件  
- 支持自定义监控事件（write/rename/remove/create/chmod）  
- 支持设置多条命令  
- 命令支持变量占位符  
- 支持冗余任务丢弃，自定义冗余任务范围  
- 支持 http 通知  
- 高级指令用法  
- 更多...  

## 编译环境    

Go >= 1.13   

## 更新日志  

[CHANGELOG](CHANGELOG.md)  

## 运行    

### 下载二进制文件   

Github: [download v1.15](https://github.com/dengsgo/fileboy/releases)  
Gitee:  [dowmload v1.15](https://gitee.com/dengsgo/fileboy/releases)  

下载已经编译好的对应平台二进制文件，重命名为`fileboy`, 加入系统 Path 中即可。 

### 源码编译   

clone 该项目，进入主目录，运行命令:  
```bash
## 确保本地 Go 启用 modules  
export GO111MODULE=on  
go env -w GOPROXY=https://goproxy.io,direct
## 安装依赖
go get -u 
## 编译
go build
## 运行
./fileboy
```

## 使用

fileboy 的正常运行依赖于 `filegirl.yaml` 配置文件，因此首次在项目中使用需要初始化 `filegirl.yaml`。  
- 进入你想要 hot reload 的项目主目录下；  
- 运行 `fileboy init`，会在该目录下生成 `filegirl.yaml`文件；  
- 查看 `filegirl.yaml`,修改为适合自己项目的配置项；  
- 运行 `fileboy`即可.  
  
如果你定义了 `command -> exec`命令，想事先确认是否能正常执行，可以运行 `fileboy exec`命令，系统会尝试运行你的自定义命令。  
你可以使用 `fileboy help`查看使用帮助。   

## filegirl.yaml 配置文件说明

```yaml
# 主配置
core:
    # 配置版本号
    version: 1

# 监控配置
monitor:
    # 要监听的目录
    # test1       监听当前目录下 test1 目录
    # test1/test2 监听当前目录下 test1/test2 目录
    # test1,*     监听当前目录下 test1 目录及其所有子目录（递归）
    # .,*         监听当前目录及其所有子目录（递归）
    includeDirs:
        - .,*

    # 不监听的目录
    # .idea   忽略.idea目录及其所有子目录的监听
    exceptDirs:
        - .idea
        - .git
        - .vscode
        - node_modules
        - vendor

    # 监听文件的格式，此类文件更改会执行 command 中的命令
    # .go   后缀为 .go 的文件更改，会执行 command 中的命令
    # .*    所有的文件更改都会执行 command 中的命令
    types:
        - .go
    
    # 不监听文件的格式，此类文件更改不会执行 command 中的命令
    # .DS_Store   后缀为 .DS_Store 的文件更改，不会执行 command 中的命令
    types:
        - .DS_Store

    # 监听的事件类型，发生此类事件才执行 command 中的命令
    # 没有该配置默认监听所有事件
    # write   写入文件事件
    # rename  重命名文件事件
    # remove  移除文件事件
    # create  创建文件事件
    # chmod   更新文件权限事件(类unix)
    events:
        - write
        - rename
        - remove
        - create
        - chmod

# 命令
command:
    # 监听的文件有更改会执行的命令
    # 可以有多条命令，会依次执行
    # 如有多条命令，每条命令都会等待上一条命令执行完毕后才会执行
    # 如遇交互式命令，允许外部获取输入
    # 支持变量占位符,运行命令时会替换成实际值：
    #    {{file}}    文件名(如 a.txt 、test/test2/a.go)
    #    {{ext}}     文件后缀(如 .go)
    #    {{event}}   事件(上面的events, 如 write)
    #    {{changed}} 文件更新的本地时间戳(纳秒,如 1537326690523046400)
    # 变量占位符使用示例：cp {{file}} /root/sync -rf  、 myCommand --{{ext}} {{changed}}
    exec:
        - go version
        - go env

    # 文件变更后命令在xx毫秒后才会执行，单位为毫秒
    # 一个变更事件(A)如果在定义的延迟时间(t)内，又有新的文件变更事件(B)，那么A会取消执行。
    # B及以后的事件均依次类推，直到事件Z在t内没有新事件产生，Z 会执行
    # 合理设置延迟时间，将有效减少冗余和重复任务的执行
    # 如果不需要该特性，设置为 0
    delayMillSecond: 2000

# 通知器
notifier:
    # 文件更改会向该 url 发送请求（POST 一段 json 文本数据）
    # 触发请求的时机和执行 command 命令是一致的
    # 请求超时 15 秒
    # POST 格式:
    #    Content-Type: application/json;charset=UTF-8
    #    User-Agent: FileBoy Net Notifier v1.16
    #    Body: {"project_folder":"/project/path","file":"main.go","changed":1576567861913824940,"ext":".go","event":"write"}
    # 例: http://example.com/notifier/fileboy-listener
    # 不启用通知，请留空 ""
    callUrl: ""

# 特殊指令
instruction:
    # 可以通过特殊的指令选项来控制 command 的行为，指令可以有多个
    # 指令选项解释：
    #   exec-when-start    fileboy启动就绪后，自动执行一次 'exec' 定义的命令
    #   should-finish      触发执行 'exec' 时(C)，如果上一次的命令(L)未退出（还在执行），会等待 L 退出（而不是强制 kill ），直到 L 有明确 exit code 才会开始执行本次命令。
    #                      在等待 L 退出时，又有新事件触发了命令执行(N)，则 C 执行取消，只会保留最后一次的 N 执行
    #   ignore-stdout      执行 'exec' 产生的 stdout 会被丢弃
    #   ignore-warn        fileboy 自身的 warn 信息会被丢弃
    #   ignore-info        fileboy 自身的 info 信息会被丢弃
    #   ignore-exec-error  执行 'exec' 出错仍继续执行下面的命令而不退出 

    #- should-finish
    #- exec-when-start
    - ignore-warn
```

## QA

### 很多框架都自带了 hot reload 的功能，为什么还要单独写个 fileboy 呢？  

这个是一款通用的 hot reload 的软件，理论上适用于任何需要 hot reload 的场景，并不局限于语言层面上。只要灵活的配置 `filegirl.yaml`文件就行了。  

### fileboy 可以应用在那些具体的场景？  

在开发中，我们很需要一款可以帮助我们自动打包编译的工具，那 fileboy 就非常适合这样的场景。比如 go 项目的热编译，让我们可以边修改代码边运行得到反馈。又比如 PHP Swoole 框架，由于常驻进程的原因，无法更改代码立即reload，使用 fileboy 就可以辅助做到传统 PHP 开发的体验。  
对于一些需要监控文件日志或者配置变动的场景， fileboy 同样适合。你可以事先编写好相应的通知报警脚本，然后定义`filegirl.yaml`中的`command`命令，交由 fileboy 自动运行监控报警。  

### 通知器在什么时候会发送 http 请求 ?

通知器发送 http 通知的前提是在配置文件中设置了 `callUrl` 参数（不为空即为已设置）。触发请求的时机和执行 command 命令是一致的，`command -> delayMillSecond` 参数对于触发器同样有效。请求超时默认15秒.  

### idea 下更改文件，为什么会执行两次或者多次 command ?

由于 idea 系列软件特殊的文件保存策略，他会自动创建一些临时文件，并且在需要时多次重写文件，所以有时反映在文件上就是有多次的更改，所以会出现这种情况。1.5之后的版本增加了 `delayMillSecond` 参数，可以解决这个问题。  

### filegirl.yaml 里面的 command 如何配置复杂命令？  

fileboy 目前支持 `命令 + 参数`这种形式的 command，而且 参数中不能有""符号或者有空格。如：  
`go build`:支持；  
`go env`:支持;  
`php swoole start --daemon`:支持  
`cat a.txt | grep "q" | wc -l`:不支持  
对于不支持的命令，可以把它写到一个文件里，然后在 command 中执行这个文件来解决。  

### 为什么起名为 fileboy，又把配置名叫做 filegirl ？

因为爱情~~ (◡ᴗ◡✿)  



## 贡献者

|   |   |   |
| ------------ | ------------ | ------------ |
| <a href="https://github.com/dengsgo"><img src="https://avatars1.githubusercontent.com/u/7929002?s=460&v=4" width=64 style="border-radius:45px;" /></a> | <a href="https://github.com/jason-gao"><img src="https://avatars1.githubusercontent.com/u/9896574?s=460&v=4" width=64 style="border-radius:45px;" /></a> | <a href="https://github.com/itwesley"><img src="https://avatars1.githubusercontent.com/u/1928721?s=460&v=4" width=64 style="border-radius:45px;" /></a> |

## 感谢支持  

|   |
| ------------ |
| <a href="https://www.jetbrains.com/?from=fileboy"><img src="./resources/jetbrains.png" width=140 /></a> |

