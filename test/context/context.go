package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

var logg *log.Logger

func someHandler() {
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	go doStuff(ctx)
	go doStuff1(ctx)
	//10秒后取消doStuff
	time.Sleep(10 * time.Second)
	cancel()

}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff1(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logg.Printf("done1")
			return
		default:
			logg.Printf("work1")
		}
	}
}

func main() {
	logg = log.New(os.Stdout, "", log.Ltime)
	//someHandler()
	timeoutHandler()
	logg.Printf("main down")
	time.Sleep(10 * time.Second)
}

func doTimeOutStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)

		if deadline, ok := ctx.Deadline(); ok { //设置了deadl
			logg.Printf("deadline set")
			if time.Now().After(deadline) {
				logg.Printf(ctx.Err().Error())
				return
			}

		}

		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

func timeoutHandler() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	go doTimeOutStuff(ctx)
	// go doStuff(ctx)

	time.Sleep(10 * time.Second)

	cancel()

}

func timeoutHandler1() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// go doTimeOutStuff(ctx)
	go doStuff(ctx)

	time.Sleep(10 * time.Second)

	cancel()

}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
