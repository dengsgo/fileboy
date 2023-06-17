## FILEBOY  

[![Build Status](https://travis-ci.org/dengsgo/fileboy.svg?branch=master)](https://travis-ci.org/dengsgo/fileboy) [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/fileboy)](https://goreportcard.com/report/github.com/dengsgo/fileboy)

[简体中文](README.md) | [ENGLISH](README_EN.md)

Fileboy, File Change Monitoring Notification Tool Software, written with Go.  
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
- Advanced instruction usage  
- more...  

## COMPILE

Go >= 1.13   

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
export GO111MODULE=on  
## installation dependency
go get-u 
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

notifier:
    # file changes send requests to the URL (POST JSON text data)
    # the timing of triggering the request is consistent with executing the command command
    # timeout 15 second
    # POST :
    #    Content-Type: application/json;charset=UTF-8
    #    User-Agent: FileBoy Net Notifier v1.17
    #    Body: {"project_folder":"/project/path","file":"main.go","changed":1576567861913824940,"ext":".go","event":"write"}
    # e.g: http://example.com/notifier/fileboy-listener
    # no notice is enabled. Please leave it blank. ""
    callUrl: ""

instruction:
    # command behavior can be controlled by special command options. there can be multiple instructions
    # options:
    #   exec-when-start    when fileboy is ready to start, execute the command defined by 'exec' once automatically
    #   should-finish      when the execution of 'exec' is triggered (C), if the last command (L) does not exit (still executing),
    #                      it will wait for L to exit (instead of forcing kill), and the execution of this command will not start until L has an explicit exit code.
    #                      when waiting for L to exit, and a new event triggers command execution (n), C execution is cancelled, and only the last N execution is retained
    #   ignore-stdout      stdout generated by executing 'exec' will be discarded
    #   ignore-warn        the warn information of fileboy itself will be discarded
    #   ignore-info        the info information of fileboy itself will be discarded
    #   ignore-exec-error  error executing 'exec' continue to execute the following command without exiting 

    #- should-finish
    #- exec-when-start
    - ignore-warn
```

## CONTRIBUTOR

|   |   |   |
| ------------ | ------------ | ------------ |
| <a href="https://github.com/dengsgo"><img src="https://avatars1.githubusercontent.com/u/7929002?s=460&v=4" width=64 style="border-radius:45px;" /></a> | <a href="https://github.com/jason-gao"><img src="https://avatars1.githubusercontent.com/u/9896574?s=460&v=4" width=64 style="border-radius:45px;" /></a> | <a href="https://github.com/itwesley"><img src="https://avatars1.githubusercontent.com/u/1928721?s=460&v=4" width=64 style="border-radius:45px;" /></a> |



## THANKS

|   |
| ------------ |
| <a href="https://www.jetbrains.com/?from=fileboy"><img src="./resources/jetbrains.png" width=140></a> |

