package produce

func Nothing()  {

}

//
//import (
//	"fmt"
//	"net/url"
//	"os"
//	"regexp"
//	"time"
//
//	"github.com/gocolly/colly"
//	"github.com/nats-io/go-nats"
//)
//
//var domain2Collector = map[string]*colly.Collector{}
//var nc *nats.Conn
//var maxDepth = 10
//var natsURL = "nats://localhost:4222"
//
//func factory(urlStr string) *colly.Collector {
//	u, _ := url.Parse(urlStr)
//	return domain2Collector[u.Host]
//}
//
//func initABCDECollector() *colly.Collector {
//	c := colly.NewCollector(
//		colly.AllowedDomains("www.baidu.com"),
//		colly.MaxDepth(maxDepth),
//	)
//
//	c.OnResponse(func(resp *colly.Response) {
//		// 做一些爬完之后的善后工作
//		// 比如页面已爬完的确认存进 MySQL
//	})
//
//	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
//		// 我们认为匹配该模式的是该网站的详情页
//		detailRegex, _ := regexp.Compile(`/go/go\?p=\d+$`)
//		// 匹配下面模式的是该网站的列表页
//		listRegex, _ := regexp.Compile(`/t/\d+#\w+`)
//
//		// 基本的反爬虫策略
//		link := e.Attr("href")
//		time.Sleep(time.Second * 2)
//
//		// 正则 match 列表页的话，就 visit
//		if listRegex.Match([]byte(link)) {
//			c.Visit(e.Request.AbsoluteURL(link))
//		}
//		// 正则 match 落地页的话，就发消息队列
//		if detailRegex.Match([]byte(link)) {
//			err := nc.Publish("tasks", []byte(link))
//			if nil != err {
//				fmt.Printf("Pub Error: %s", err.Error())
//			}
//			nc.Flush()
//		}
//	})
//	return c
//}
//
//func initHIJKLCollector() *colly.Collector {
//	c := colly.NewCollector(
//		colly.AllowedDomains("studygolang.com"),
//		colly.MaxDepth(maxDepth),
//	)
//
//	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
//	})
//
//	return c
//}
//
//func init() {
//	domain2Collector["www.baidu.com"] = initABCDECollector()
//	domain2Collector["studygolang.com"] = initHIJKLCollector()
//
//	var err error
//	nc, err = nats.Connect(natsURL)
//	if err != nil {
//		fmt.Printf("nats connect error: %s", err.Error())
//		os.Exit(1)
//	}
//}
//
//func main() {
//	urls := []string{"https://www.baidu.com", "https://studygolang.com"}
//	for _, url := range urls {
//		instance := factory(url)
//		fmt.Printf("\ninstance===%+v", instance)
//		instance.Visit(url)
//	}
//
//	ticker := time.NewTicker(time.Second)
//	for range ticker.C {
//		err := nc.Publish("tasks", []byte("Hello world!"))
//		if nil != err {
//			fmt.Printf("Pub Error: %s", err.Error())
//		}
//		nc.Flush()
//	}
//
//}
