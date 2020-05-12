package main

//import (
//	"fmt"
//	"sync"
//)
//
//var total struct {
//	sync.Mutex
//	value int
//}
//
//func worker(wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	for i := 0; i <= 100; i++ {
//		total.Lock()
//		total.value += i
//		total.Unlock()
//	}
//}
//
//func main() {
//	var wg sync.WaitGroup
//	wg.Add(2)
//	go worker(&wg)
//	go worker(&wg)
//	wg.Wait()
//
//	fmt.Println(total.value)
//}

import (
	"fmt"
	"github.com/t-xinlin/doc/pkg/mutex/pubsub"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&total, i)
	}
}

// 生产者: 生成 factor 整数倍的序列
func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

// 消费者
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main1() {
	ch := make(chan int, 64) // 成果队列

	go Producer(3, ch) // 生成 3 的倍数的序列
	go Producer(5, ch) // 生成 5 的倍数的序列
	go Consumer(ch)    // 消费 生成的队列

	// Ctrl+C 退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Printf("quit (%v)\n", <-c)

	var wg sync.WaitGroup
	wg.Add(2)

	go worker(&wg)
	go worker(&wg)
	wg.Wait()
	fmt.Println(total)

	//var config atomic.Value // 保存当前配置信息
	//
	//// 初始化配置信息
	//config.Store(loadConfig())
	//
	//// 启动一个后台线程, 加载更新后的配置信息
	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		config.Store(loadConfig())
	//	}
	//}()
	//
	//// 用于处理请求的工作者线程始终采用最新的配置信息
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		for r := range requests() {
	//			c := config.Load()
	//			// ...
	//		}
	//	}()
	//}
}
func main() {
	p := pubsub.NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello,  world!")
	p.Publish("hello, golang!")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}
