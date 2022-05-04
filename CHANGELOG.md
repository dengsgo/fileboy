### Release v1.16

2022.05.04

- 修复 windows使用-filegirl加载配置出错  
  
2022.02.13

- 优化 net client  
- 增加 -filegirl 参数，允许加载指定路径的配置  
- 优化 一些细节

2020.12.03

- 修改 进程退出清理日志级别 
- 增加 MIT Copyright  

2020.08.23

- 优化 文件扫描性能

2020.07.19

- 增加 pid 文件处理  
- 增加 信息处理  

2020.03.16

- typo deamon->daemon  

### Release v1.15

2020.03.08

- 优化 指令模式  
- 使用 mod 管理依赖  
- go version >= 1.13  
- 优化 一些细节  

2020.01.02

- 增加 指令配置项 `instruction`, 可以通过预定义的指令来控制 command 的行为  
- 增加 `should-finish` 指令  
- 增加 `exec-when-start` 指令   
- 增加 `ignore-warn` 指令   
- 增加 `ignore-info` 指令   
- 增加 `ignore-stdout` 指令   
- 增加 `ignore-exec-error` 指令   


2019.12.28

- 增加 `deamon`命令，支持以守护进程的方式运行在后台 **Unix only**  
- 增加 `stop`命令，用来停止 deamon 进程 **Unix only**  
- 优化 exec stdout  


### Release v1.12

2019.12.18  

- 增加 自定义监控事件（write/rename/remove/create/chmod）  
- 增加 {{event}} 事件名占位符 / event 网络通知字段  
- 增加 文件(夹)变更动态添加/删除监听（beta）  
- 优化 init 命令,如果已有`filegirl.yaml`现在提示错误,不会自动覆盖  
- 优化 日志输出缓冲,可以通过 >> 将fileboy自身输出日志重定向到文件  
- 修复 始终默认监听主目录的问题  
- 升级 底层依赖  
- PR Makefile 的支持 @jason-gao  


### Release v1.10

2019.06.04  

- 优化 log  
- 优化 代码逻辑  
- 修复 gin框架hotReload  


### Release v1.9

2019.04.03  

- 优化 文件夹监听效率，减少大量深层文件夹遍历的时间  
- 优化 代码逻辑  
- 增加 readme 英文说明  
- 修复 偶现监听项目主目录无效的问题  


### Release v1.8

2019.02.27  

- 使用 go1.12 编译  



### Release v1.7

2019.01.24  

- 修复 time 内存  
- 修复 某些情况下cmd异常导致进程挂掉的问题  



### Release v1.6

2019.01.19  

- 修复 http 通知失败导致进程崩溃  
- 增加 includeDirs 参数规则验证  
- 修改 delayMillSecond 默认值,2000  
- 增加 贡献者 @jason-gao  
- 优化 log  



### Release v1.5

2019.01.03  

- 增加 http 通知  
- 增加 callUrl 参数  
- 优化 command 稳定性  



2019.01.02  

- 增加 command -> delayMillSecond 参数  
- 优化 文案  



2018.12.30  

- 增加 在指定时间内堆叠的任务自动丢弃  
- 增加 version 信息  
- 优化 代码逻辑  
- 优化 提示文案  



### Release v1.2

2018.09.30  

- 修复 递归监听会忽略主级目录的bug  
- 增加 `fileboy version`命令  



### Release v1.1

2018.09.19

- command 命令支持变量占位符 `{{file}}`、`{{ext}}`、`{{changed}}`；  

```yaml
# {{file}}    文件名(如 a.txt 、test/test2/a.go)
# {{ext}}     文件后缀(如 .go)
# {{changed}} 文件更新的本地时间戳(纳秒,如 1537326690523046400)
# 变量占位符使用示例：cp {{file}} /root/sync -rf  、 myCommand --{{ext}} {{changed}}
```

- 增加 较深目录递归提示；  
- 优化 文字提示；  
- 修复 command 命令执行时目录不正确的问题； 
- 修复 其他bug; 



### Release v1.0

2018.09.10

- 文件变更监听支持多平台 （windows/linux测试，mac未测试）；  

- 支持灵活配置监听 包含文件夹/排除文件夹/特定文件类型；  

- 支持配置变更要运行命令，可以有多条，会依次执行；  

- 支持 fileboy init 初始化配置，生成 filegirl.yaml 文件；  

- 支持 fileboy exec 直接执行配置的 command 命令；  