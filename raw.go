package main

import "strconv"

var exampleFileGirl string = `# 主配置
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
    # 一个变更事件(A)如果在定义的延迟时间(t)内, 又有新的文件变更事件(B), 那么A会取消执行。
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
    #    User-Agent: FileBoy Net Notifier v1.15
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
`

var firstRunHelp = `第一次运行 fileboy ?
你可能需要先执行 fileboy init 生成配置。
更多信息使用 fileboy help 查看帮助
`

var helpStr = `fileboy [option]
Usage of fileboy:
    无参数 
        读取当前目录下的 filegirl.yaml 配置，开始监听并工作
    init 
        初始化 fileboy, 在当前目录生成 filegirl.yaml 配置文件
    exec 
        尝试运行定义的 command 命令
    deamon 
        读取当前目录下的 filegirl.yaml 配置，以守护进程的方式运行在后台
    stop 
        停止守护进程
    version 
        查看当前版本信息
`

var englishSay = []string{
	`      Have you, the darkness is no darkness.`,
	`    Why do the good girls always love bad boys?`,
	`              If love is not madness.`,
	`         This world is so lonely without you.`,
	`         You lie. Silence in front of me.`,
	`    I need him like I need the air to breathe.`,
	`  Happiness is when the desolated soul meets love.`,
	`   What I can lose, but do not want to lose you.`,
	`     The same words, both miss, is also missed.`,
	`  Each bathed in the love of the people is a poet.`,
}

var logo = `
 _______ _____ _       _______ ______   _____  _     _ 
(_______|_____) |     (_______|____  \ / ___ \| |   | |
 _____     _  | |      _____   ____)  ) |   | | |___| |
|  ___)   | | | |     |  ___) |  __  (| |   | |\_____/ 
| |      _| |_| |_____| |_____| |__)  ) |___| |  ___   
|_|     (_____)_______)_______)______/ \_____/  (___)   V1.15
`
var statement = `Dengsgo [dengsgo@gmail.com] Open Source with MIT License`

var versionDesc = `
 Version   fileboy: v1.15    filegirl: v` + strconv.Itoa(Version) + `
Released   2020.01.05
 Licence   MIT
  Author   dengsgo [dengsgo@gmail.com]
 Website   https://github.com/dengsgo/fileboy
    Blog   https://www.yoytang.com
`
