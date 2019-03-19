package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uber-go/fx"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MyInterface interface {
	run()
}

type Conf struct {
	Name string
	Age  int
}

func (conf *Conf) run() {
	println("run go_________")
}

func fx_uber() {
	var reader io.Reader

	app := fx.New(
		// io.reader的应用
		// 提供构造函数
		fx.Provide(func() io.Reader {
			return strings.NewReader("hello world")
		}),
		fx.Populate(&reader), // 通过依赖注入完成变量与具体类的映射
	)
	app.Start(context.Background())
	defer app.Stop(context.Background())

	// 使用
	// reader变量已与fx.Provide注入的实现类关联了
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Panic("read occur error, ", err)
	}
	fmt.Printf("the result is '%s' \n", string(bs))
}

func main() {
	fx_uber()
	return
	var myInterface MyInterface
	myInterface = new(Conf)
	myInterface.run()

	total := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				if i != j && j != k && i != k {
					total++
					println(i, " ", j, " ", k)
				}
			}
		}
	}
	println("total： ", total)
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
