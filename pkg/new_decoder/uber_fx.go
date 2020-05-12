package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"
)

//在Fx中提供的构造函数都是惰性调用，可以通过invocations在application启动来完成一些必要的初始化工作：fx.Invoke(function); 通过也可以按需自定义实现LiftCycle的Hook对应的OnStart和OnStop用来完成手动启动容器和关闭，来满足一些自己实际的业务需求。
// Logger构造函数
func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	logger.Print("Executing NewLogger.")
	return logger
}

// http.Handler构造函数
func NewHandler(logger *log.Logger) (http.Handler, error) {
	logger.Print("Executing NewHandler.")
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		logger.Print("Got a request.")
	}), nil
}

// http.ServeMux构造函数
func NewMux(lc fx.Lifecycle, logger *log.Logger) *http.ServeMux {
	logger.Print("Executing NewMux.")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{ // 自定义生命周期过程对应的启动和关闭的行为
		OnStart: func(context.Context) error {
			logger.Print("Starting HTTP server.")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Print("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})

	return mux
}

// 注册http.Handler
func Register(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}

func main() {
	app := fx.New(
		fx.Provide(
			NewLogger,
			NewHandler,
			NewMux,
		),
		fx.Invoke(Register), // 通过invoke来完成Logger、Handler、ServeMux的创建
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil { // 手动调用Start
		log.Fatal(err)
	}

	http.Get("http://localhost:8080/") // 具体操作

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil { // 手动调用Stop
		log.Fatal(err)
	}

}
