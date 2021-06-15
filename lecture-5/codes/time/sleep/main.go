package main

import "time"

func main() {
	// 启动函数
	time.Sleep(1 * time.Second)
	println("started")

	// 启动五秒周期
	go periodic()
	time.Sleep(5 * time.Second)
}

func periodic() {
	for {
		// 每一秒输出一次
		println("tick")
		time.Sleep(1 * time.Second)
	}
}

// 该代码不会终止运行
