package main

import (
	"fmt"
	"sync"
)

//
// Fetcher
//
type Fetcher interface {
	Fetch(url string) (urls []string, err error)
}

//
// 页面抓取结果
//
type fakeResult struct {
	body string
	urls []string
}

//
// 页面抓取集合
//
type fakeFetcher map[string]*fakeResult

//
// 页面抓取函数
//
func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found: %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("missing: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

var fetcher fakeFetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

//
// Serial Crawler
//
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	// 链接已被抓取直接返回
	if fetched[url] {
		return
	}

	// 标识链接已被抓取防止重复抓取
	fetched[url] = true

	// 获取链接集合
	urls, err := fetcher.Fetch(url)

	if err != nil {
		return
	}

	for _, url := range urls {
		Serial(url, fetcher, fetched)
	}

	return
}

//
// 定义状态管理接口
//
type FetchState struct {
	mutex   sync.Mutex
	fetched map[string]bool
}

//
// 创建状态管理
//
func makeState() *FetchState {
	state := &FetchState{}
	state.fetched = make(map[string]bool)
	return state
}

//
// 共享内存抓取数据
//
func ConcurrentMutex(url string, fetcher Fetcher, state *FetchState) {
	state.mutex.Lock()
	fetched := state.fetched[url]
	state.fetched[url] = true
	state.mutex.Unlock()

	if fetched {
		return
	}

	urls, err := fetcher.Fetch(url)

	if err != nil {
		return
	}

	var done sync.WaitGroup

	for _, url := range urls {
		done.Add(1)
		go func(url string) {
			defer done.Done()
			ConcurrentMutex(url, fetcher, state)
		}(url)
	}

	done.Wait()
	return
}

//
//  Worker
//
func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	if err != nil {
		ch <- []string{}
	} else {
		ch <- urls
	}
}

//
// Master
//
func master(ch chan []string, fetcher Fetcher) {
	n := 1
	fetched := make(map[string]bool)
	for urls := range ch {
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1
				go worker(u, ch, fetcher)
			}
		}
		n -= 1
		if n == 0 {
			break
		}
	}
}

//
// 利用 Channel 抓取数据
//
func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url}
	}()
	master(ch, fetcher)
}

//
// 函数入口
//
func main() {
	// 定义抓取链接
	url := "http://golang.org/"

	fmt.Printf("=== Serial ===\n")
	Serial(url, fetcher, make(map[string]bool))

	fmt.Printf("=== ConcurrentMutex ===\n")
	ConcurrentMutex(url, fetcher, makeState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	ConcurrentChannel(url, fetcher)
}
