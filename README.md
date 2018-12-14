# fcli

## fcli 简介

fcli 是阿里云函数计算的命令行工具，可以便捷的管理函数计算中的资源。

```
$ fcli
fcli: function compute command line tools

Usage:
  fcli [flags]
  fcli [command]

Available Commands:
  alias           alias related operation
  config          Configure the fcli
  function        function related operation
  help            Help about any command
  service         service related operation
  shell           interactive shell
  sls             sls related operation
  trigger         trigger related operation
  version         fcli version information

Flags:
  -h, --help   help for fcli

Use "fcli [command] --help" for more information about a command.
```

详细的使用手册：[函数计算工具 fcli 帮助文档](https://help.aliyun.com/document_detail/52995.html?spm=5176.10695662.1996646101.searchclickresult.d81a50128SHSpG)

## 如何贡献代码
### 开发环境配置

__1. 安装并配置 Golang 开发环境__
根据 [官方文档](https://golang.org/) 安装并设置环境变量，主要是设置好 `$GOPATH` 环境变量。

__2. Fork Repository__
- 在 [aliyun/fcli](https://github.com/aliyun/fcli) 项目中，点击 `fork`，将项目 fork 到个人仓库
- 在本地 `$GOPATH/src` 目录下，创建 `github.com/aliyun` 目录
- `cd ${GOPATH}/src/github.com/aliyun`
- `git clone https://github.com/个人账号/fcli.git`

__3. 安装 glide 包管理器__
```
$ go get github.com/Masterminds/glide
$ go install github.com/Masterminds/glide
```

__4. 安装依赖__
在项目根目录下，执行 `glide i -v` 进行依赖安装

### 提交 pull request

__1. 将修改 push 到个人账号里的本地仓库__

__2. 发起 pull request 请求__
在 pull request 请求的 comment 中，写明此次修改的内容，并添加此次修改的命令交互示例。

假设此次修改设置到 service list 子命令

- fcli service
```
$ fcli service
service related operation

Usage:
  fcli service [flags]
  fcli service [command]

Aliases:
  service, s

Available Commands:
  create      create service
  delete      Delete service
  get         Get the information of service
  list        List services of the current account
  update      update service
  version     service version related operation

Flags:
      --help   Print Usage (default true)

Use "fcli service [command] --help" for more information about a command.
```

- fcli service list --help
```
william:fcli zechen$ go run main.go service list --help
List services of the current account

Usage:
  fcli service list [option] [flags]

Aliases:
  list, l

Flags:
      --help                list functions
  -l, --limit int32         the max number of the returned services (default 100)
      --name-only           display service name only (default true)
  -t, --next-token string   continue listing the functions from the previous point
  -p, --prefix string       list the services whose names contain the specified prefix
  -k, --start-key string    start key is where you want to start listing from
```

- fcli service list
```
$ fcli service list
{
  "Services": [
    "demo"
  ],
  "NextToken": null
}
```

我们鼓励将此次修改涉及的命令交互，展现的越详细越好。(还可以测试此次修改影响到的命令的各种参数)

## 版本号说明
fcli 在开源后使用 `主版本号.次版本号.修订号` 的版本格式

关于版本格式，可以参考 [Semantic Versioning](https://semver.org/)


## Auto complete under shell

### bash

run command as follows:

```
curl -s https://raw.githubusercontent.com/aliyun/fcli/master/misc/completion/fcli-completion.bash > /usr/local/etc/bash_completion.d/fcli-completion.bash
```

or checkout this repo and run command as follows:

```
source misc/completion/fcli-completion.bash
```

and then relogin

### zsh

run command as follows:

```
curl -s https://raw.githubusercontent.com/aliyun/fcli/master/misc/completion/_fcli > /usr/local/share/zsh/site-functions/_fcli
```

and then relogin

### 