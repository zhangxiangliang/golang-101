package main

import (
	"sync"
	"time"
)

var done bool
var mu sync.Mutex

func main() {
	// 开始任务
	println("started")
	time.Sleep(1 * time.Second)

	// 启动五秒周期
	go periodic()
	time.Sleep(5 * time.Second)

	// 利用 Done 来表示运行截止至
	mu.Lock()
	done = true
	mu.Unlock()

	// 关闭任务
	println("cancelled")
	time.Sleep(3 * time.Second)
}

func periodic() {
	for {

		// 每秒输出一次
		println("tick")
		time.Sleep(1 * time.Second)

		// 检查是否 Done 已经停止运行
		mu.Lock()
		if done {
			return
		}
		mu.Unlock()
	}
}
