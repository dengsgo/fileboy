## 项目说明  

fileboy，文件变更监听通知系统，使用 GO 编写。  

适用于 Hot Reload （典型的如开发go项目，无需每次手动执行 go build；又比如前端 node 打包） 或者 系统监控的场景。  

## 编译环境    

go version >=1.10   

## 更新日志  

[CHANGELOG](CHANGELOG.md)  


## 运行    

### 下载二进制文件   

Github: [正式版 v1.1](https://github.com/dengsgo/fileboy/releases)  
Gitee [正式版 v1.1](https://gitee.com/dengsgo/fileboy/releases)  

直接下载已经编译好的对应平台二进制文件，加入系统 Path 中即可。 

### 源码编译   

clone 该项目，进入主目录，运行命令:    

```shell
## 安装依赖
go get -u gopkg.in/fsnotify/fsnotify.v1
go get -u gopkg.in/yaml.v2
## 编译
go build
## 运行
./fileboy
```

## 使用

fileboy 的正常运行依赖于 `filegirl.yaml` 配置，所以首次在项目中使用需要初始化 `filegirl.yaml`。  

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

# 命令
command:
    # 监听的文件有更改会执行的命令
    # 可以有多条命令，会依次执行
    # 如有多条命令，每条命令都会等待上一条命令执行完毕后才会执行
    # 如遇交互式命令，允许外部获取输入
    # 支持变量占位符,运行命令时会替换成实际值：
    #    {{file}}    文件名(如 a.txt 、test/test2/a.go)
    #    {{ext}}     文件后缀(如 .go)
    #    {{changed}} 文件更新的本地时间戳(纳秒,如 1537326690523046400)
    # 变量占位符使用示例：cp {{file}} /root/sync -rf  、 myCommand --{{ext}} {{changed}}
    exec:
        - go version
        - go env
        - echo {{file}}
```

## QA

#### 很多框架都自带了 hot reload 的功能，为什么还要单独写个 fileboy 呢？  

这个是一款通用的 hot reload 的软件，理论上适用于任何需要 hot reload 的场景，并不局限于语言层面上。只要灵活的配置 `filegirl.yaml`文件就行了。  

#### fileboy 可以应用在那些具体的场景？  

在开发中，我们很需要一款可以帮助我们自动打包编译的工具，那 fileboy 就非常适合这样的场景。比如 go 项目的热编译，让我们可以边修改代码边运行得到反馈。又比如 PHP Swoole 框架，由于常驻进程的原因，无法更改代码立即reload，使用 fileboy 就可以辅助做到传统 PHP 开发的体验。  

对于一些需要监控文件日志或者配置变动的场景， fileboy 同样适合。你可以事先编写好相应的通知报警脚本，然后定义`filegirl.yaml`中的`command`命令，交由 fileboy 自动运行监控报警。  

#### idea 下更改文件，为什么会执行两次或者多次 command ?

由于 idea 系列软件特殊的文件保存策略，他会自动创建一些临时文件，并且在需要时多次重写文件，所以有时反映在文件上就是有多次的更改，所以会出现这种情况。这个后续会做优化.  

#### filegirl.yaml 里面的 command 不支持复杂的命令吗？  

对于“很复杂的命令”这种说法很难去定义，比如 `echo "hello world"`并不复杂，但是对于 fileboy 来讲，目前无法解析这种命令。  

fileboy 目前支持 `命令 + 参数`这种形式的 command，而且 参数中不能有""符号或者有空格。如：  

`go build`:支持；  

`go env`:支持;  

`php swoole start --daemon`:支持  

`cat a.txt | grep "q" | wc -l`:不支持  

对于不支持的命令，可以把它写到一个文件里，然后在 command 中执行这个文件来解决。  

#### 为什么起名为 fileboy，又把配置名叫做 filegirl ？

因为爱情~~

#### 听说有彩蛋？

(◡ᴗ◡✿)
