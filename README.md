## 项目说明  

fileboy，文件变更监听通知系统，使用 Go 编写。  
适用于 Hot Reload （典型的如开发go项目，无需每次手动执行 go build；又比如前端 node 打包） 或者 系统监控的场景。  
___  
Fileboy, File Change Monitoring Notification System, written with Go.  
For Hot Reload scenarios (typically for developing go projects without having to perform go build manually every time; for example, front-end node packaging) or system monitoring.  


## 特性  

- 极简的用法和配置  
- 支持多平台，Windows/Linux/MacOS  
- 支持自定义文件监听范围，监听指定文件夹/不监听指定文件夹/指定后缀文件  
- 支持设置多条命令  
- 命令支持变量占位符  
- 支持冗余任务丢弃，自定义冗余任务范围  
- 支持 http 通知  
- 更多...  
___  
- Minimalist usage and configuration  
- Support multiple platforms, Windows/Linux/MacOS  
- Supports custom file listening scope, listening for specified folders/not listening for specified folders/specified suffix files  
- Support for setting up multiple commands  
- Command support variable placeholders  
- Supporting redundant task discarding and customizing redundant task scope  
- Supporting HTTP notifications  
- more...  


## 编译环境    

go version 1.12   

## 更新日志  

[CHANGELOG](CHANGELOG.md)  


## 运行    

### 下载二进制文件   

Github: [download v1.9](https://github.com/dengsgo/fileboy/releases)  
Gitee:  [dowmload v1.9](https://gitee.com/dengsgo/fileboy/releases)  

下载已经编译好的对应平台二进制文件，重命名为`fileboy`, 加入系统 Path 中即可。 
___  
Download the compiled binary file of the corresponding platform, rename it `fileboy', and add it to the system Path.  

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
___  
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
## 使用

fileboy 的正常运行依赖于 `filegirl.yaml` 配置，所以首次在项目中使用需要初始化 `filegirl.yaml`。  
- 进入你想要 hot reload 的项目主目录下；  
- 运行 `fileboy init`，会在该目录下生成 `filegirl.yaml`文件；  
- 查看 `filegirl.yaml`,修改为适合自己项目的配置项；  
- 运行 `fileboy`即可.  
  
如果你定义了 `command -> exec`命令，想事先确认是否能正常执行，可以运行 `fileboy exec`命令，系统会尝试运行你的自定义命令。  
你可以使用 `fileboy help`查看使用帮助。  
___  
The normal operation of fileboy depends on the `filegirl.yaml` configuration, so for the first time in a project, `filegirl.yaml` needs to be initialized.  
- Enter the project home directory where you want hot reload;  
- Running `fileboy init` will generate `filegirl. yaml` files in this directory.  
- View `filegirl. yaml` and modify it to a configuration item suitable for your project.  
- Run `fileboy`.  
  
If you define the `command-> exec` command, you can run the `fileboy exec` command to confirm whether it can be executed properly in advance. The system will try to run your custom command.  
You can use `fileboy help` to see help info.  

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
    #    User-Agent: FileBoy Net Notifier v1.8
    #    Body: {"project_folder":"/watcher-dirs","file":"test.go","changed":1546421173070433800,"ext":".go"}
    # 例: http://example.com/notifier/fileboy-listener
    # 不启用通知，请留空 ""
    callUrl: ""
```
___  
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

command:
    # the files monitored have commands that change to be executed
    # there can be multiple commands that will be executed in turn
    # in case of interactive commands, allow external access to input
    # variable placeholders are supported, and the actual values are replaced when the command is run:
    #    {{file}}    (e.g: a.txt 、test/test2/a.go)
    #    {{ext}}     (e.g: .go)
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

notifier:
    # file changes send requests to the URL (POST JSON text data)
    # the timing of triggering the request is consistent with executing the command command
    # timeout 15 second
    # POST :
    #    Content-Type: application/json;charset=UTF-8
    #    User-Agent: FileBoy Net Notifier v1.8
    #    Body: {"project_folder":"/watcher-dirs","file":"test.go","changed":1546421173070433800,"ext":".go"}
    # e.g: http://example.com/notifier/fileboy-listener
    # no notice is enabled. Please leave it blank. ""
    callUrl: ""
```

### TODO

- [x] 命令支持变量占位符  
- [x] 支持多命令  
- [x] 支持监听指定文件夹  
- [x] 支持不监听指定文件夹  
- [x] 支持监听指定后缀文件  
- [x] 支持 http 通知  
- [x] 支持冗余任务丢弃  
- [ ] 支持 http 合并任务的通知  
___  
- [x] command supports variable placeholders  
- [x] Supports multiple commands  
- [x] Supports listening for specified folders  
- [x] Supports not listening to specified folders  
- [x] Supports listening for specified suffix files  
- [x] Supports HTTP notifications  
- [x] Supports redundant task discarding  
- [ ] Notification supporting HTTP merge tasks  

## QA

#### 很多框架都自带了 hot reload 的功能，为什么还要单独写个 fileboy 呢？  

这个是一款通用的 hot reload 的软件，理论上适用于任何需要 hot reload 的场景，并不局限于语言层面上。只要灵活的配置 `filegirl.yaml`文件就行了。  

#### fileboy 可以应用在那些具体的场景？  

在开发中，我们很需要一款可以帮助我们自动打包编译的工具，那 fileboy 就非常适合这样的场景。比如 go 项目的热编译，让我们可以边修改代码边运行得到反馈。又比如 PHP Swoole 框架，由于常驻进程的原因，无法更改代码立即reload，使用 fileboy 就可以辅助做到传统 PHP 开发的体验。  
对于一些需要监控文件日志或者配置变动的场景， fileboy 同样适合。你可以事先编写好相应的通知报警脚本，然后定义`filegirl.yaml`中的`command`命令，交由 fileboy 自动运行监控报警。  

#### 通知器在什么时候会发送 http 请求 ?

通知器发送 http 通知的前提是在配置文件中设置了 `callUrl` 参数（不为空即为已设置）。触发请求的时机和执行 command 命令是一致的，`command -> delayMillSecond` 参数对于触发器同样有效。请求超时默认15秒.  

#### idea 下更改文件，为什么会执行两次或者多次 command ?

由于 idea 系列软件特殊的文件保存策略，他会自动创建一些临时文件，并且在需要时多次重写文件，所以有时反映在文件上就是有多次的更改，所以会出现这种情况。1.5之后的版本增加了 `delayMillSecond` 参数，可以解决这个问题。  

#### filegirl.yaml 里面的 command 不支持复杂的命令吗？  

对于“很复杂的命令”这种说法很难去定义，比如 `echo "hello world"`并不复杂，但是对于 fileboy 来讲，目前无法解析这种命令。  
fileboy 目前支持 `命令 + 参数`这种形式的 command，而且 参数中不能有""符号或者有空格。如：  
`go build`:支持；  
`go env`:支持;  
`php swoole start --daemon`:支持  
`cat a.txt | grep "q" | wc -l`:不支持  
对于不支持的命令，可以把它写到一个文件里，然后在 command 中执行这个文件来解决。  

#### 为什么起名为 fileboy，又把配置名叫做 filegirl ？

因为爱情~~ (◡ᴗ◡✿)  



### 贡献者

> 排名不分先后

[@dengsgo](https://www.yoytang.com)  <dengsgo@gmail.com>  

[@itwesley](https://github.com/itwesley)  <wcshen1126@gmail.com>  

[@jason-gao](https://github.com/jason-gao)  <3048789891@qq.com>  

