# MapReduce

## 项目介绍

在这个实验中你会建立一个类似 [MapReduce](../lectures/lecture-2/papers/mapreduce.pdf) 论文提到的 MapReduce 系统:

* Master Process 应用在 分配任务 Worker 和 处理失败 Worker。
* Worker Process 应用在处理 Map 和 Reduce 函数、处理文件读写。

## 项目须知

* 除了实验项目提供的代码，你必需独立完成 6.824 实验项目。
* 不允许查看他人的实验项目代码，当然也不允许看之前几年的实验项目代码。
* 允许你与同学讨论项目实现方案，但是不允许查看复制同学的实验项目代码。
* 请相信你可以独立设计并完成 6.824 实验项目的所有编码。
* 请不要以任何方式公开你的实验项目代码，因为可能会被其他学习 6.824 的同学检索到。
* github 项目创建时会默认公开项目仓库，如果要使用 github 请确保创建私有项目仓库。

总而言之就是 “独立完成，可以交流，不可抄袭”。

## 项目软件

### 开发环境

* 该实验项目使用 Go 语言进行开发，Go 官方站点删该有相当多的教程信息可以自行学习，
* 实验项目的测试用例运行环境为 Go 1.13 版本，所以推荐你也使用该版本，可以用 `go version` 命令确认当前版本。

### Mac OS

* 方案一：安装 [Homebrew](https://brew.sh) 后使用 `brew instal go` 进行 Go 相关工具的安装。
* 方案二：使用 [gvm](https://github.com/moovweb/gvm) 版本管理工具来进行 Go 相关工具的安装。

### Linux

* 方案一：根据相应 Linux 的包管理工具安装最新的 Go 相关工具，例如使用 `apt install golang` 来安装。
* 方案二：从 Go 官方相关的站点下载并解药对应版本的二进制文件到 `bin` 目录下，并确认 `bin` 目录有配置在 `PATH`。

### Windows

* 不建议使用 Windows 来进行项目试验，如果真的需要可以参考 [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) 来进行环境搭建。

### 版本管理

* 使用 Git 来对你的项目进行管理，如果需要学习Git相关知识可以查看 [Pro Git book](https://git-scm.com/book/en/v2) 或者 [Git User Manual](https://mirrors.edge.kernel.org/pub/software/scm/git/docs/user-manual.html)。

## 项目初始化

### 拉取项目

```
git clone git://g.csail.mit.edu/6.824-golabs-2020 6.824
```

### 项目结构

* `src/main/mrsequential.go` 实现了一个简单的 MapReduce，这个例子中 maps 和 reduces 是同时运行的。
* `src/mrapps/wc.go` 提供了一个 word-count 工具。
* `src/mrapps/indexer.go` 提供了一个 text indexer 工具。

### 运行工具

运行 word-count 工具，该工具会读取 `src/main/pg-xxx.txt` 文件，并将结果输出到 `mr-out-0` 的文件中

```
cd ~/6.824
cd src/main
go build -buildmode=plugin ../mrapps/wc.go
rm mr-out*
go run mrsequential.go wc.so pg*.txt
more mr-out-0
```

建议看看 `src/mrapps/wc.go` 了解一下 MapReduce 的基础样例。

## 项目任务

### 基础需求

实现一个 MapReduce 分布式系统，真实的开发环境中，该 MapReduce 分布式系统 往往会运行在多个不同的主机上，但是在这个实验中只需要在本机运行即可：

* 需要实现 Master 程序 和 Worker 两个程序。
* MapReduce 分布式系统运行后将有一个 Master 进程 和 一个或者多个 Worker 进程同时运行。
* Master 和 Worker 通过 RPC 调用的方式来进行沟通。
* Worker 主动向 Master 获取任务，并从一个或多个文件中读取数据、执行任务、并将结果输出到一个或多个文件中。
* Master 负责监控 Worker 任务的执行情况，如果没有在规定的时间内完成该任务，则将该任务重新分配到一个新的 Worker 中。

### 项目实现

* 实验项目提供了一个基础的 Master 和 Worker 的结构在 `src/main/mrmaster.go` 和 `src/main/mrworker.go`中。
* 不要直接在两个文件进行编程，而是将项目实现放到 `src/mr/mrmaster.go` 和 `src/mr/mrworker.go`、`src/mr/rpc.go` 中。