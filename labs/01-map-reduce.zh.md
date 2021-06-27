# MapReduce

## 项目介绍

在这个实验中你会建立一个类似 [MapReduce](../lectures/lecture-2/papers/mapreduce.pdf) 论文提到的 MapReduce 系统：

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
$ git clone git://g.csail.mit.edu/6.824-golabs-2020 6.824
```

### 项目结构

* `src/main/mrsequential.go` 实现了一个简单的 MapReduce，这个例子中 maps 和 reduces 是同时运行的。
* `src/mrapps/wc.go` 提供了一个 word-count 工具。
* `src/mrapps/indexer.go` 提供了一个 text indexer 工具。

### 运行工具

运行 word-count 工具，该工具会读取 `src/main/pg-xxx.txt` 文件，并将结果输出到 `mr-out-0` 的文件中

```
$ cd ~/6.824
$ cd src/main
$ go build -buildmode=plugin ../mrapps/wc.go
$ rm mr-out*
$ go run mrsequential.go wc.so pg*.txt
$ more mr-out-0
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

### 项目步骤

首先要确定 `work-count` 插件是最新构建的：

```
$ go build -buildmode=plugin ../mrapps/wc.go
$ rm mr-out*
```

其次在 `main` 目录, 运行 Master 程序， `pg-*.txt` 将会作分布式文件输入到 Map 任务中：

```
$ go run mrmaster.go pg-*.txt
```

再者在另外的窗口运行一些 Worker 程序：

```
$ go run mrworker.go wc.so
```

当 Master 和 Worker 完成后将会看到输出文件 `mr-out-*`。你需要对输出文件进行一个简单的排序：

```
$ cat mr-out-* | sort | more
```

### 项目测试

使用通过 `main/test-mr.sh` 运行测试文件 `wc` 和 `indexer` 来测试项目：

* 输出结果是否正确。
* Map 和 Reduce 任务是否并行。
* Worker 程序在崩溃后能否自我恢复。

如果你现在运行测试脚本，脚本会进入挂起状态，因为 Master 程序由永远不会完成：

```
$ cd ~/6.824/src/main
$ sh test-mr.sh
*** Starting wc test.
```

如果你修改 `mr/master.go` 中的 `ret := false` 为 true，这个时候运行测试脚本 Master 程序将会立即完成：

```
$ sh ./test-mr.sh
*** Starting wc test.
sort: No such file or directory
cmp: EOF on mr-wc-all
--- wc output is not the same as mr-correct-wc.txt
--- wc test: FAIL
```

测试脚本会检查每一个输出文件 `mr-out-*`，如果没有实现 `mr/master.go` 和 `mr/worker.go` 将不会有输出文件和执行任何操作，所以测试结果为失败。如果你通过了测试，测试脚本会输出与下列相似的内容：

```
$ sh ./test-mr.sh
*** Starting wc test.
--- wc test: PASS
*** Starting indexer test.
--- indexer test: PASS
*** Starting map parallelism test.
--- map parallelism test: PASS
*** Starting reduce parallelism test.
--- reduce parallelism test: PASS
*** Starting crash test.
--- crash test: PASS
*** PASSED ALL TESTS
```

如果你看到来自 Go RPC 包的一些错误，请忽视这些错误：

```
2019/12/16 13:27:09 rpc.Register: method "Done" has 1 input parameters; needs exactly three
```

### 实验规则

* 在 Map 阶段的时候将中间值存储到 `nReduce` Reduce 任务中，`nReduce` 将会作为参数传递给 `main/mrmaster.go` 中的 `MakeMaster()` 函数。
* Worker 程序的实现需要将输出结果放置到对应的 `mr-out-*` 文件中。
* 每个 `mr-out-*` 应该包含一行 Reduce 输出，这行输出使用 Go 中 `%v %v` 格式生成，分别对应键和值。可以参考 `main/mrsequential.go` 中的 'this is the correct format‘。如果你没有准守这个输出格式，测试脚本将无法通过。
* 可以修改 `mr/worker.go`、`mr/master.go` 和 `mr/rpc.go`。除吃之外的其他文件也可以修改，但是最后需要与原文件保持一致。
* Worker 应该将 Map 的结果存储到文件中，方便 Reduce 进行操作和读写。
* 当 MapReduce 任务完成时 `main/mrmaster.go` 会期待 `mr/master.go` 实现一个返回 true 的 `Done()` 方法，`main/mrmaster.go` 将会退出。
* 当项目完成后必需退出所有 Worker 程序。一个比较简单的实现方式是使用 `call()`：
    * 如果这个 Worker 程序无法联系 Master 程序，那它可以假设 Master 退出，即项目完成了 Worker 程序可以退出。
    * 根据这个设计方法，Master 程序还可以创建一个 `please exit` 的伪任务发送给 Worker 程序。

### 实验提示

* 一种开始的思路是通过 `mr/worker.go` 的 `Worker()` 发送一个 RPC 请求给 Master 程序获取任务。接着去修改 Master 程序响应一个唯一的任务和文件名给 Worker 程序。接着去修改 Worker 程序读取该文件并调用 Map 函数，例如 `mrsequential.go`。
* 以 `.so` 结尾的文件中包含 Map 和 Reduce 函数，将作为 Go 插件包在运行是被使用。
* 如果你更改 `mr/` 目录下的内容，则需要重新编译 MapReduce 插件，例如 `go build -buildmode=plugin ../mrapps/wc.go`。
* 项目运行 Worker 程序需要共享文件系统，当只有一台主机时这很容易实现。当 Worker 程序运行在不同主机上时，则需要 GFS 这样的全局文件系统。
* 中间文件的命名规则为 `mr-X-Y`，其中 X 是 Map 的任务编号，Y 是 Reduce 的任务编号。
* Worker 程序在运行 Map 任务时将会需要临时存储 key-value 键值对在文件中，以便 Reduce 任务可以正确的读取数据。可以使用 `encoding/json` 包来进行 `key-value` 键值对来存储：

```
# 写入
enc := json.NewEncoder(file)
for _, kv := ... {
  err := enc.Encode(&kv)
}
```

```
# 读取
dec := json.NewDecoder(file)
for {
  var kv KeyValue
  if err := dec.Decode(&kv); err != nil {
    break
  }
 kva = append(kva, kv)
}
```

* Worker 程序的 Map 任务可以使用 worker.go 中的 `ihash(key)` 函数来为通过 key 获取 Reduce 任务。
* 可以参考 `mrsequential.go` 中关于读取 Map 输入文件、对中间文件的 key-value 键值对排序、存储 Reduce 输出文件。
* Master 程序作为 RPC 服务器将会产生并发请求，不要忘记共享事务锁。
* 使用 Go 的竞态检测器 `go build -race` 和 `go run -race`，在 `test-mr.sh` 有关于如何为测试脚本启动竞态检测器。
* Worker 程序有时会需要等待执行，例如：Reduce 需要等待 Map 执行后产生中间文件才能运行。一种解决方法是通过告诉 Master 程序并调用 `time.Sleep()` 进行睡眠和唤醒。另一种解决方法是在 Master 程序 RPC 服务器中创建一个循环接收器，可以是 `time.Sleep()` 或者 `sync.Cond`。Go 在线程中为 RPC 运行处理文件，所以处理文件等待的过程中不会影响到 Master 程序处理其他 RPC 请求。
* Master 程序无法可靠的区分崩溃的 Worker 程序、某种情况导致工作停滞的 Worker 程序、执行太慢而无法使用的 Worker 程序，你所能做的就是让 Master 程序等待一段时间，如果 10 秒内没有响应则假设 Worker 程序已经死亡，并重新派遣新的任务到其他 Worker 程序中。
* 如果要测试崩溃修复可以使用 `mrapps/crash.go` 插件，它将会随机在 Map 或者 Reduce 程序中运行。
* 为了确保崩溃的 Worker 程序无意间写入文件，MapReduce 论文提到了使用临时文件并在完全写入后进行重命名。可以使用 `ioutil.TempFile` 去创建一个临时文件，然后使用 `os.Rename` 去将它重命名。
* `test-mr.sh` 在 `mr-tmp` 子目录中运行所有进程，如果出现错误需要查看中间文件或者输出可以到 `mr-tmp` 中浏览对应文件。

### 项目提交

* 重要事件: 需要运行 `test-mr.sh` 后才能提交项目。
* 重要事件: 需要登录网站去确认你已经成功提交项目代码。
* 使用 `make lab1` 这个命令去打包你的项目，并将项目上传到 https://6824.scripts.mit.edu/2020/handin.py/ 网站中。
* 使用 MIT 证书或者 API KEY 对项目进行签名：

```
$ cd ~/6.824
$ echo XXX > api.key
$ make lab1
```

注意: 你可能会有多次项目提交，我们将会通过时间戳来获取你最后一次提交的项目。