package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/fx"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main11() {
	fx_uber()
}

//Annotated(位于annotated.go文件) 主要用于采用annotated的方式,提供Provide注入类型
func fx_uber1(){
	type t3 struct {
		Name string
	}

	targets := struct {
		fx.In

		V1 *t3 `name:"n1"`
	}{}

	app := fx.New(
		fx.Provide(fx.Annotated{
			Name:"n1",
			Target: func() *t3{
				return &t3{"hello world"}
			},
		}),
		fx.Populate(&targets),
	)
	app.Start(context.Background())
	defer app.Stop(context.Background())

	fmt.Printf("the result is = '%v'\n", targets.V1.Name)
}

//type App struct {
//	err          error
//	container    *dig.Container        // 容器
//	lifecycle    *lifecycleWrapper    // 生命周期
//	provides     []interface{}            // 注入的类型实现类
//	invokes      []interface{}
//	logger       *fxlog.Logger
//	startTimeout time.Duration
//	stopTimeout  time.Duration
//	errorHooks   []ErrorHandler
//
//	donesMu sync.RWMutex
//	dones   []chan os.Signal
//}

// 新建一个App对象
//func New(opts ...Option) *App {
//	logger := fxlog.New()   // 记录Fx日志
//	lc := &lifecycleWrapper{lifecycle.New(logger)}  // 生命周期
//
//	app := &App{
//		container:    dig.New(dig.DeferAcyclicVerification()),
//		lifecycle:    lc,
//		logger:       logger,
//		startTimeout: DefaultTimeout,
//		stopTimeout:  DefaultTimeout,
//	}
//
//	for _, opt := range opts {  // 提供的Provide和Populate的操作
//		opt.apply(app)
//	}
//
//	// 进行app相关一些操作
//	for _, p := range app.provides {
//		app.provide(p)
//	}
//	app.provide(func() Lifecycle { return app.lifecycle })
//	app.provide(app.shutdowner)
//	app.provide(app.dotGraph)
//
//	if app.err != nil {  // 记录app初始化过程是否正常
//		app.logger.Printf("Error after options were applied: %v", app.err)
//		return app
//	}
//
//	// 执行invoke
//	if err := app.executeInvokes(); err != nil {
//		app.err = err
//
//		if dig.CanVisualizeError(err) {
//			var b bytes.Buffer
//			dig.Visualize(app.container, &b, dig.VisualizeError(err))
//			err = errorWithGraph{
//				graph: b.String(),
//				err:   err,
//			}
//		}
//		errorHandlerList(app.errorHooks).HandleError(err)
//	}
//	return app
//}

func fx_uber() {
	type T3 struct {
		Name string
	}
	type T4 struct {
		Age int
	}

	type Result struct {
		fx.Out
		V1 *T3 `name:"n1"`
		V2 *T3 `name:"n2"`
	}
	type ResultG struct {
		fx.Out
		V1 *T3 `group: "g"`
		V2 *T3 `group: "g"`
	}
	targets := struct {
		fx.In
		V1 *T3 `name:"n1"`
		V2 *T3 `name:"n2"`
	}{}
	targetsG := struct {
		fx.In
		Group []*T3 `group: "g"`
	}{}
	var t3 *T3
	var t4 *T4
	var reader io.Reader

	app := fx.New(
		// io.reader的应用
		// 提供构造函数
		fx.Provide(func() io.Reader {
			return strings.NewReader("hello world")
		}),

		fx.Provide(func() *T3{
			return &T3{"your name"}
		}),
		fx.Provide(func() *T4{
			return &T4{19}
		}),

		fx.Provide(func() Result{
			return Result{
				V1:&T3{"v1 name"},
				V2:&T3{"v2 name"},
			}
		}),
		fx.Provide(func() ResultG{
			return ResultG{
				V1:&T3{"v1 name G"},
				V2:&T3{"v2 name G"},
			}
		}),
		fx.Populate(&reader), // 通过依赖注入完成变量与具体类的映射
		fx.Populate(&t3),
		fx.Populate(&t4),
		fx.Populate(&targets),
		fx.Populate(&targetsG),
	)
	app.Start(context.Background())
	defer app.Stop(context.Background())

	// 使用
	// reader变量已与fx.Provide注入的实现类关联了
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Panic("read occur error, ", err)
	}
	fmt.Printf("the result is reader:'%s' t3: %v  t4: %v  targets.v1.name: %v  targets.v2.name: %v\n", string(bs), t3, t4, targets.V1.Name, targets.V2.Name)
}

type openWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}
