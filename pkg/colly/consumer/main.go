package consumer

func Nothing()  {

}

//
//import (
//	"fmt"
//	"github.com/nats-io/go-nats"
//	"net/url"
//	"os"
//	"time"
//
//	"github.com/gocolly/colly"
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
//		colly.AllowedDomains("www.abcdefg.com"),
//		colly.MaxDepth(maxDepth),
//	)
//	return c
//}
//
//func initV2fxCollector(url string) *colly.Collector {
//	c := colly.NewCollector(
//		colly.AllowedDomains(url),
//		colly.MaxDepth(maxDepth),
//	)
//	return c
//}
//
//func init() {
//	domain2Collector["www.baidu.com"] = initV2fxCollector("www.baidu.com")
//	domain2Collector["studygolang.com"] = initV2fxCollector("studygolang.com")
//
//	var err error
//	nc, err = nats.Connect(natsURL)
//	if err != nil {
//		os.Exit(1)
//	}
//}
//
//func startConsumer() {
//	nc, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		return
//	}
//
//	sub, err := nc.QueueSubscribeSync("tasks", "workers")
//	if err != nil {
//		return
//	}
//
//	var msg *nats.Msg
//	for {
//		msg, err = sub.NextMsg(time.Hour * 10000)
//		if err != nil {
//			break
//		}
//
//		fmt.Printf("\nmsg: %s", string(msg.Data))
//
//		//urlStr := string(msg.Data)
//		//ins := factory(urlStr)
//		//// 因为最下游拿到的一定是对应网站的落地页
//		//// 所以不用进行多余的判断了，直接爬内容即可
//		//ins.Visit(urlStr)
//		// 防止被封杀
//		time.Sleep(time.Second)
//	}
//}
//
//func main() {
//	startConsumer()
//}
