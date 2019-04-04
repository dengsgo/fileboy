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