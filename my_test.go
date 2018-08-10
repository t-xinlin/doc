import (
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
	"log"
	"reflect"
	"strings"
	"testing"
	"os"
	"context"
	"math/rand"
	"container/list"
	"encoding/json"
	"net/url"
	"runtime"
	"sync"
	"errors"
	"net/http"
	"go/parser"
	"go/token"
	"go/ast"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func init() {
	//log.Printf("my_test init...")
}

const DefaultExpireTime = 0

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func getICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

type SayHello interface {
	sayHello()
}
type Jp struct {
	name string
}

func (a Jp) sayHello() {
	fmt.Println(a.name)
}

type Chinese struct {
	name string
}

func (a Chinese) sayHello() {
	fmt.Println(a.name)
}
func greet(i interface{}) {
	//r, ok := i.(SayHello)
	switch t := i.(type) {
	default:
		fmt.Printf("unexpected type %T", t) // %T prints whatever type t has
	case bool:
		fmt.Printf("boolean %t\n", t) // t has type bool
	case int:
		fmt.Printf("integer %d\n", t) // t has type int
	case *bool:
		fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
	case *int:
		fmt.Printf("pointer to integer %d\n", *t) // t has type *int
	case SayHello:
		fmt.Println("greet say")
		t.sayHello()
	}

}

type Float float64

type Person struct {
	name string
	age  int
}

type T struct {
	A int
	B string
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) ReflectCallFunc() {
	fmt.Println("reflect learn")
}

func (u User) FuncHasArgs(name string, age int) {
	fmt.Println("FuncHasArgs name: ", name, ", age:", age, "and origal User.Name:", u.Name)
}

func (u User) FuncNoArgs() {
	fmt.Println("FuncNoArgs")
}

func sum(j int) (s int) {
	s = j + 100
	fmt.Printf("\nsum %+v", s)
	return s
}

func sum2(j []int) (s int) {
	//s = j + 100
	fmt.Printf("\n#################sum2 %+v", j)
	j[5] = 99
	return 0
}

func sum1(j []int) (s int) {
	//s = j + 100
	fmt.Printf("\nsum1 %+v", j)
	j[5] = 99
	return 0

}

func context_test() {
	//d := time.Now().Add(50 * time.Millisecond)
	//ctx, cancel := context.WithDeadline(context.Background(), d)
	//// Even though ctx will be expired, it is good practice to call its
	//// cancelation function in any case. Failure to do so may keep the
	//// context and its parent alive longer than necessary.
	//defer cancel()
	//select {
	//case <-time.After(1 * time.Second):
	//	fmt.Printf("overslept")
	//case <-ctx.Done():
	//	fmt.Printf("err: %+v",ctx.Err())
	//	fmt.Printf("\n---")
	//}

	//d := time.Now().Add(50 * time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(),d)
	//// Even though ctx will be expired, it is good practice to call its
	//// cancelation function in any case. Failure to do so may keep the
	//// context and its parent alive longer than necessary.
	//defer cancel()
	//select {
	//case <- time.After(1 * time.Second):
	//	fmt.Printf("overslept")
	//case <- ctx.Done():
	//	fmt.Printf("err: %+v",ctx.Err())
	//}

	//// gen generates integers in a separate goroutine and
	//// sends them to the returned channel.
	//// The callers of gen need to cancel the context once
	//// they are done consuming generated integers not to leak
	//// the internal goroutine started by gen.
	//gen := func(ctx context.Context) <-chan int {
	//	dst := make(chan int)
	//	n := 1
	//	go func() {
	//		for {
	//			select {
	//			case <-ctx.Done():
	//				return // returning not to leak the goroutine
	//			case dst <- n:
	//				n++
	//			}
	//		}
	//	}()
	//
	//	return dst
	//}
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel() // cancel when we are finished consuming integers
	//
	//for n := range gen(ctx) {
	//	fmt.Println(n)
	//	if n == 5 {
	//		break
	//	}
	//}

	//type favContextKey string
	//
	//f := func(ctx context.Context, k favContextKey) {
	//	if v := ctx.Value(k); v != nil {
	//		fmt.Println("found value:", v)
	//		return
	//	}
	//	fmt.Println("key not found:", k)
	//}
	//
	//k := favContextKey("language")
	//ctx := context.WithValue(context.Background(), k, "Go")
	//f(ctx, k)
	//f(ctx, favContextKey("color"))

}

func Utf8Index(str, substr string) int {
	asciiPos := strings.Index(str, substr)
	if asciiPos == -1 || asciiPos == 0 {
		return asciiPos
	}
	pos := 0
	totalSize := 0
	reader := strings.NewReader(str)
	for _, size, err := reader.ReadRune(); err == nil; _, size, err = reader.ReadRune() {
		totalSize += size
		pos++
		// 匹配到
		if totalSize == asciiPos {
			return pos
		}
	}
	return pos
}

var mychan = make(chan int)

func dgucofunc1() {
	fmt.Println("In dguco_func1")
	mychan <- 1
}

func dgucofunc2() {
	a := <-mychan
	fmt.Println("In dguco_func2 a=", a)
}

func dgucoproducer(c chan int, max int) {
	for i := 0; i < max; i++ {
		time.Sleep(time.Second * 2)
		c <- i
	}
	//close(c)
}

func dgucoconsumer(c chan int) {
	ok := true
	value := 0
	for ok {
		time.Sleep(time.Second * 2)
		fmt.Println("Wait receive")
		if value, ok = <-c; ok {
			fmt.Println(value)
		}
		if ok == false {
			fmt.Println("*******Break********")
		}
	}
}

var mychannel = make(chan bool)

func dgucotimertask() {
	fmt.Println("My TimerTask..")
}

func dgucotimer() {
	//返回一个定时器
	mytimer := time.NewTicker(time.Millisecond * 1000)
	select {
	case <-mytimer.C:
		go dgucotimertask()
	}
}

type Work struct {
}

func (w Work) Refuse() {
	fmt.Println("Refuse")
}

func (w Work) Do() {
	fmt.Println("Do")
}

func worker(i int, ch chan Work, quit chan struct{}) {
	for {
		select {
		case w := <-ch:
			if quit == nil {
				w.Refuse();
				fmt.Println("worker", i, "refused", w)
				break
			}
			w.Do();
			fmt.Println("worker", i, "processed", w)
		case <-quit:
			fmt.Println("worker", i, "quitting")
			quit = nil
		}
	}
}

func makeWork(ch chan Work) {

}

func Benchmark_ticker(testB *testing.B) {
	var ticker *time.Ticker
	ticker = time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			log.Printf("time ticker run ")
			time.Sleep(time.Second * 5)
		}
	}()

	for ; ; {
		i := 0
		i++
	}
}

type Man struct {
	Name string
	Age  int
}

const str_dic = `The make built-in function allocates and initializes an object of type
	slice, map, or chan (only). Like new, the first argument is a type, not a
value. Unlike new, make's return type is the same as the type of its
argument, not a pointer to it. The specification of the result depends on
the type:
Slice: The size specifies the length. The capacity of the slice is
equal to its length. A second integer argument may be provided to
specify a different capacity; it must be no smaller than the
length. For example, make([]int, 0, 10) allocates an underlying array
of size 10 and returns a slice of length 0 and capacity 10 that is
backed by this underlying array.
Map: An empty map is allocated with enough space to hold the
specified number of elements. The size may be omitted, in which case
a small starting size is allocated.
Channel: The channel's buffer is initialized with the specified
buffer capacity. If zero, or the size is omitted, the channel is
unbuffered.`

func dict_word(w string) {
	t := strings.Replace(w, ",", "", -1)
	d := strings.Split(t, " ")
	m := make(map[string]int)
	for _, v := range d {
		if _, ok := m[v]; ok {
			m[v] = m[v] + 1
		} else {
			m[v] = 1
		}
	}

	for k, v := range m {
		fmt.Printf("%+v >>> tatol: %+v\n", k, v)
	}

}

func Benchmark_map(testB *testing.B) {
	dict_word(str_dic)

	map_man := make(map[string]Man)
	map_man["str1"] = Man{Name: "xiaoming", Age: 16}
	map_man["str2"] = Man{Name: "xiaobai", Age: 19}
	map_man["str3"] = Man{Name: "zhangsan", Age: 20}
	map_man["str4"] = Man{Name: "lilu", Age: 18}
	for _, v := range map_man {
		fmt.Printf("name: %+v  age: %+v\n", v.Name, v.Age)
		v.Name = "test"
	}

	for _, v := range map_man {
		fmt.Printf("aftet modify-- name: %+v  age: %+v\n", v.Name, v.Age)
		v.Name = "test"
	}

	for k, v := range map_man {
		//fmt.Printf("modify-- name: %+v  age: %+v\n", v.Name, v.Age)
		v.Name = "test"
		map_man[k] = v
	}

	for _, v := range map_man {
		fmt.Printf("aftet modify -- name: %+v  age: %+v\n", v.Name, v.Age)
		v.Name = "test"
	}

	map1 := make(map[string]string)
	map1["a"] = "AAA"
	map1["b"] = "BBB"
	map1["c"] = "CCC"

	for k, v := range map1 {
		fmt.Printf("map k=%+v  v=%+v\n", k, v)
	}
	for _, v := range map1 {
		fmt.Printf("map v=%+v\n", v)
	}

	array := [...]int64{1, 2, 3, 4}
	for k, v := range array {
		fmt.Printf("arr k=%+v  v=%+v\n", k, v)
	}
	for _, v := range array {
		fmt.Printf("arr  v=%+v\n", v)
	}

	slice := array[:3:3]
	for k, v := range slice {
		fmt.Printf("arr k=%+v  v=%+v\n", k, v)
	}
	fmt.Println(cap(slice))
	slice[0] = 444
	slice = append(slice, 888)
	slice = append(slice, 999)
	fmt.Println(slice)
	fmt.Println(array)
}

func Benchmark_string(testB *testing.B) {

	var version1 = "v0.0.1"
	var version2 = "v0.0.3"

	if version2 > version1 {
		log.Printf("version1< version2 ? %+v", version1 < version2)
	} else if version2 < version1 {
		log.Printf("version1 > version2 ?%+v", version1 > version2)
	} else if version2 == version1 {
		log.Printf("version1 = version2 ?%+v", version1 == version2)
	}

	log.Printf("---------------------")

	data := "A\xfe\x02\xff\x04"
	for _, v := range data {
		fmt.Printf("%#x ", v)
	}
	//prints: 0x41 0xfffd 0x2 0xfffd 0x4 (not ok)

	fmt.Println()
	for _, v := range []byte(data) {
		fmt.Printf("%#x ", v)
	}
	//prints: 0x41 0xfe 0x2 0xff 0x4 (good)

}

func Benchmark_channel(testB *testing.B) {

	ch, quit := make(chan Work), make(chan struct{})
	go makeWork(ch)
	for i := 0; i < 4; i++ {
		go worker(i, ch, quit)
	}
	time.Sleep(5 * time.Second)
	close(quit)
	time.Sleep(2 * time.Second)

	return
	dgucotimer()
	for {
		time.Sleep(time.Second * 10)
	}

	//------------------------

	c00 := make(chan int)
	defer close(c00)
	go dgucoproducer(c00, 10)
	go dgucoconsumer(c00)
	for {
		time.Sleep(time.Second * 10)
	}
	fmt.Println("Done")

	return
	go dgucofunc1()
	go dgucofunc2()

	for {
		time.Sleep(time.Second * 10)
	}
	//

	timeout := make(chan bool, 1)
	t := 0 != 0
	go func() {
		for {
			time.Sleep(time.Second * 15) // sleep 3 seconds
			timeout <- true
		}

	}()
	ch1 := make(chan int, 0)
	ch2 := 0
	i1 := 0
	go func(i int) {
		for {
			time.Sleep(time.Second * 1) // sleep 3 seconds
			i++
			ch1 <- i
		}
	}(i1)
	fmt.Println("start...")
	//for !t {
	select {
	case ch2 = <-ch1:
		log.Println("ch 1 = !", ch2)
	case t = <-timeout:
		log.Println("timeout!", t)

		//default:
		//	fmt.Println("default!")
	}
	//}
	log.Println("end...")

}

func Benchmark_panic_defer(testB *testing.B) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover")
		}
	}()
	j := 0
	i := 1 / j
	fmt.Printf("recover i", i)

	//---------------
	for i := 0; i < 10; i++ {
		defer sum(i)
	}

	var arr = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//for i:=0; i<10; i++{
	defer sum1(arr)
	arr[0] = 88
	//}
	fmt.Printf("**********arr: %+v", arr)
	//a10 := defer_test(10)
	//fmt.Printf("\nmain : %+v", a10)

}

func Benchmark_reflect(testB *testing.B) {
	user := User{1, "test", 13}
	var i interface{}
	i = user

	uValue := reflect.ValueOf(i)
	uType := reflect.TypeOf(i)
	fmt.Println("uValue: ", uValue)
	fmt.Println(uValue.Interface()) //转换为interface类型,unpack uValue.Interface().(User)
	fmt.Println(uValue.Type())
	fmt.Println("uValue,string: ", uType.String())
	fmt.Println("uType: ", uType.Name())

	for i := 0; i < uType.NumField(); i++ { //获取field信息
		field := uType.Field(i)
		value := uValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	for i := 0; i < uType.NumMethod(); i++ { // 获取method信息
		method := uType.Method(i)
		fmt.Printf("method[%d] = %s \n", i, method.Name)

		//m:=uValue.MethodByName(method.Name)
		//args := []reflect.Value{reflect.ValueOf("xiong"), reflect.ValueOf(30+i)}
		//m:=reflect.ValueOf(method)
		//m.Call(args)

	}

	fmt.Println(uValue.Kind())
	fmt.Println(uType.Kind())

	//user := User{1, "test", 13}
	//uValue := reflect.ValueOf(user)
	//uType  := reflect.TypeOf(user)

	fmt.Printf("--------------------\n")

	m1 := uValue.MethodByName("FuncHasArgs")
	m2 := uValue.MethodByName("FuncNoArgs")
	m, b1 := uType.MethodByName("FuncNoArgs")
	args := []reflect.Value{reflect.ValueOf("xiong"), reflect.ValueOf(30)}
	m1.Call(args)

	args = make([]reflect.Value, 0)
	m2.Call(args)

	fmt.Println("m1:", m1)
	fmt.Println("m2:", m2)
	fmt.Printf("m:%#v,isfound:%v\n", m, b1)
	fmt.Println(m1)

	var x1 int = 8
	value1 := reflect.ValueOf(x1)
	fmt.Println(value1.Type()) //int
	fmt.Println(value1.Kind()) //int
	var x Float = 3.14
	value2 := reflect.ValueOf(x)
	fmt.Println(value2.Type()) //Float
	fmt.Println(value2.Kind()) //float64
	person := Person{}
	value3 := reflect.ValueOf(person)
	fmt.Println(value3.Type()) //Person
	fmt.Println(value3.Kind()) //struct

}

const (
	b  = 1 << (10 * iota)
	kb
	mb
	gb
	tb
	pb
)

func Benchmark_iota(testB *testing.B) {
	var man SayHello
	man = Jp{name: "JP"}
	man.sayHello()
	greet(man)
	man = Chinese{name: "中国人"}
	man.sayHello()
	greet(man)
	const (
		a = iota //0
		b        //1
		c        //2
		d = "ha" //独立值，iota += 1
		e        //"ha"   iota += 1
		f = 100  //iota +=1
		g        //100  iota +=1
		h = iota //7,恢复计数
		i        //8
	)
	log.Printf("a=%+v, b=%+v, c=%+v, d=%+v, e=%+v, f=%+v, g=%+v, h=%+v, i=%+v", a, b, c, d, e, f, g, h, i)

	const (
		a1 = iota * 8
		b1
		c1
	)

	log.Printf("a1=%+v, b1=%+v, c1=%+v", a1, b1, c1)

	//////---------------------------
	fmt.Println("b=", b)
	fmt.Println("kb=", kb)
	fmt.Println("mb=", mb)
	fmt.Println("gb=", gb)
	fmt.Println("tb=", tb)
	fmt.Println("pb=", pb)
}

func Benchmark_ping(testB *testing.B) {
	//store := make(map[string]string, 0)
	//store["a"] = "1"
	//store["b"] = "2"
	//store["c"] = "3"
	//store["d"] = "4"
	//
	//log.Printf("*****2*****%+v", store["c"])
	//
	//store = make(map[string]string, 0)
	//
	//v, ok := store["c"]
	//
	//log.Printf("****2******v= %+v,  ok=%+v", v, ok)
	//
	////---------------------------------------------------------------
	//
	//h1 := "rest://100.101.197.164:30101/"
	//u, err := url.Parse(h1)
	//if err != nil {
	//	panic(err)
	//}
	//
	//raddr, err := net.ResolveIPAddr("ip", "100.101.197.164")
	//
	//if err != nil {
	//	fmt.Printf("Fail to resolve %s, %s\n", u.Host, err)
	//	return
	//}
	//
	//fmt.Printf("Ping %s (%s):\n\n", raddr.String(), u.Host)
	//
	//for i := 1; i < 6; i++ {
	//	if err = sendICMPRequest(getICMP(uint16(i)), raddr); err != nil {
	//		fmt.Printf("Error: %s\n", err)
	//	}
	//	time.Sleep(2 * time.Second)
	//}

	//var (
	//	icmp   ICMP
	//	laddr  = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
	//	raddr, _ = net.ResolveIPAddr("ip", h1)
	//	)
	//	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//	defer conn.Close()

	//我们将解析这个 URL 示例，它包含了一个 scheme，认证信息，主机名，端口，路径，查询参数和片段。
	////s := "postgres://user:pass@host.com:5432/path?k=v#f"
	//s := "rest://100.101.197.164:30101/"
	////解析这个 URL 并确保解析没有出错。
	//u, err := url.Parse(s)
	//if err != nil {
	//	panic(err)
	//}
	//
	////fmt.Printf("----------%+v, %+v, %+v, %+v", u.Host, u.Port())
	//h := strings.Split(u.Host, ":")
	//fmt.Println(h[0])
	//fmt.Println(h[1])

	//直接访问 scheme。
	//fmt.Println(u.Scheme)
	////User 包含了所有的认证信息，这里调用 Username和 Password 来获取独立值。
	//fmt.Println(u.User)
	//fmt.Println(u.User.Username())
	//p, _ := u.User.Password()
	//fmt.Println(p)
	////Host 同时包括主机名和端口信息，如过端口存在的话，使用 strings.Split() 从 Host 中手动提取端口。
	//fmt.Println(u.Host)
	//h := strings.Split(u.Host, ":")
	//fmt.Println(h[0])
	//fmt.Println(h[1])
	////这里我们提出路径和查询片段信息。
	//fmt.Println(u.Path)
	//fmt.Println(u.Fragment)
	////要得到字符串中的 k=v 这种格式的查询参数，可以使用 RawQuery 函数。你也可以将查询参数解析为一个map。已解析的查询参数 map 以查询字符串为键，对应值字符串切片为值，所以如何只想得到一个键对应的第一个值，将索引位置设置为 [0] 就行了。
	//fmt.Println(u.RawQuery)
	//m, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(m)
	//fmt.Println(m["k"][0])

	//conn, err := net.DialIP("ip4:icmp", nil, destAddr)

	//conn, err := net.Dial("tcp", "100.101.197.164:30101")
	//if err != nil {
	//	fmt.Printf("Fail to connect to remote host: %s\n", err)
	//	return err
	//}
	//defer conn.Close()
	//
	//var buffer bytes.Buffer
	//binary.Write(&buffer, binary.BigEndian, icmp)
	//
	//if _, err := conn.Write(buffer.Bytes()); err != nil {
	//	fmt.Printf("Write")
	//	return err
	//}
	//
	//tStart := time.Now()
	//
	//conn.SetReadDeadline((time.Now().Add(time.Second * 2)))
	//
	//recv := make([]byte, 1024)
	//receiveCnt, err := conn.Read(recv)
	//
	//if err != nil {
	//	fmt.Printf("Read")
	//	return
	//}
	//
	//tEnd := time.Now()
	//duration := tEnd.Sub(tStart).Nanoseconds() / 1e6
	//
	//fmt.Printf("%d bytes from %s: seq=%d time=%dms\n", receiveCnt, destAddr.String(), icmp.SequenceNum, duration)
}

func Benchmark_channel1(testB *testing.B) {
	//var wg sync.WaitGroup
	//var urls = []string{
	//	"http://www.google.com",
	//	"http://www.google.org",
	//	"http://www.baidu.com",
	//	"http://www.baidu.com",
	//	"http://www.baidu.com",
	//}

	cacheChan := make(chan struct{}, 5)

	for i := 0; i < 5; i++ {
		cacheChan <- struct{}{}
	}

	go func() {
		for {
			cacheChan <- struct{}{}
			//defer wg.Add(-1)
		}
	}()

	var temp = make(chan struct{})
	go func() {
		for msg := range cacheChan {
			select {
			case temp <- msg:
				log.Printf("case----")
			default:
				//err := WriteMessageToBackend(&msgBuf, msg, c.backend)
				//if err != nil {
				//	// ... handle errors ...
				//}
				log.Printf("default----")
			}
		}
	}()

	//for i := 0; i < 5; i++ {
	//	cacheChan <- struct{}{}
	//}
	//
	//for _, url := range urls {
	//	<-cacheChan
	//	wg.Add(1)
	//	go func(url string) {
	//		cacheChan <- struct{}{}
	//		time.Sleep(time.Second * 2)
	//		fmt.Printf(">>>>>>>>>  %+v \n", url)
	//
	//		defer wg.Add(-1)
	//	}(url)
	//}

	//defer close(cacheChan)
	//wg.Wait()

	for ; ; {
		i := 0
		i++
	}
	fmt.Printf(">>>>>>>>>> over \n")

}

var logg *log.Logger

func someHandler() {
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	doStuff()

	//10秒后取消doStuff
	//time.Sleep(10 * time.Second)

}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
//func doStuff(ctx context.Context) {
func doStuff() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	//time.Sleep(1 * time.Second)
	//is_out := false
	logg.Printf("start")
	//go func() {
	//	for {
	//		if is_out {
	//			return
	//		}
	//		time.Sleep(10 * time.Second)
	//	}
	//
	//}()

	select {
	case <-ctx.Done():
		logg.Printf("done-timeout")
		//is_out = true
		//return
		//default:
		//	logg.Printf("work")
		//	//time.Sleep(10 * time.Second)
		//	//time.Sleep(20 * time.Second)
		//	time.Sleep(1 * time.Second)
	}

	logg.Printf("end")
}

func Benchmark_context(testB *testing.B) {
	logg = log.New(os.Stdout, "", log.Ltime)
	someHandler()
	logg.Printf("finish")
}

func context_withTimeout(is_first *bool, prev_map *map[string]string) {
	rand.Seed(int64(time.Now().Nanosecond()))
	ch := make(chan bool)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		if *is_first {
			(*prev_map)["a"] = "1"
		} else {
			(*prev_map)["a"] = "2"
		}
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
		ch <- true
	}()

	select {
	case <-ch:
		*is_first = false
		(*prev_map)["b"] = "fast"
		return
	case <-ctx.Done():
		*is_first = true
		(*prev_map)["b"] = "slow"
		return
	}
}

func longRunningCalculation(timeCost int) chan string {
	result := make(chan string)
	go func() {
		time.Sleep(time.Second * (time.Duration(timeCost)))
		result <- "Done"
	}()
	return result
}

func jobWithTimeoutHandler() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	select {
	case <-ctx.Done():
		//log.Println(ctx.Err())
		log.Printf("done-timeout")
		return
	case <-longRunningCalculation(3):
		log.Printf("work finish")
	}
	return
}

func jobWithTimeoutHandler1() {

	for {
		time.Sleep(time.Duration(2 * time.Second))
		ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
		//defer cancel()

		arr := make([]int, 5)

		for i, _ := range arr {
			log.Printf("---%+v", i)
			select {
			//case <-time.After(time.Second * 5):
			//	log.Printf("############# %+v #############", "time after!!!")
			case <-ctx.Done():
				//log.Println(ctx.Err())
				log.Printf("---%+v", "context timeout!!!")
				goto loop
				//return
				//case <-longRunningCalculation(3):
				//	log.Printf("work finish")
			default:
				t := time.Duration(rand.Intn(5))
				log.Printf("---ping %+v", t)
				time.Sleep(time.Duration(t * time.Second))


			}
		}
	loop:
		log.Printf("---%+v ", "finish!!!")
	}

}

func Benchmark_context_withTimeout(testB *testing.B) {
	jobWithTimeoutHandler1()
	return
	is_first := true
	prev_map := make(map[string]string)
	for {
		context_withTimeout(&is_first, &prev_map)
		time.Sleep(time.Second)
		log.Println(prev_map, is_first)
	}
}

func delte_array(arr []string, i int) (result []string) {
	index := 0
	endIndex := len(arr) - 1
	result = make([]string, 0)

	for k, _ := range arr {
		if k == i {
			result = append(result, arr[index:k]...)
			index = k + 1
		} else if k == endIndex {
			result = append(result, arr[index:endIndex+1]...)
		}
	}

	return result
}
func Benchmark_array(testB *testing.B) {
	var s = []string{"ca", "ab", "ec", "ca", "ab", "ec", "ca", "ab", "ab"}
	log.Printf("befor %+v", s)
	t := delte_array(s, 9)
	log.Printf("after %+v", t)
}

func Contains(l *list.List, value interface{}) (bool, *list.Element) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return true, e
		}
	}
	return false, nil
}
func Benchmark_list(testB *testing.B) {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)
	l.PushBack(5)
	l.PushBack(6)
	l.PushBack(7)
	l.PushBack(8)
	l.PushBack(9)
	log.Printf("befor list : %+v", l)
	if contain, e := Contains(l, 9); contain {
		l.Remove(e)
	}
	log.Printf("after list : %+v", l)
}

func get() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) //prints: 10000 10000 <byte_addr_x>
	return raw[:3]
}

func Benchmark_point(testB *testing.B) {

	data := []*struct{ num int }{{1}, {2}, {3}}

	for _, v := range data {
		v.num *= 10
	}

	fmt.Println(data[0], data[1], data[2])            //prints &{10} &{20} &{30}
	fmt.Printf("\n-------------------------------\n") //prints &{10} &{20} &{30}
	data1 := get()
	fmt.Println(len(data1), cap(data1), &data1[0]) //prints: 3 10000 <byte_addr_x>
}

func Benchmark_point1(testB *testing.B) {
	s1 := []int{1, 2, 3}
	fmt.Println(len(s1), cap(s1), s1) //prints 3 3 [1 2 3]

	s2 := s1[1:]
	fmt.Println(len(s2), cap(s2), s2) //prints 2 2 [2 3]

	for i := range s2 {
		s2[i] += 20
	}

	//still referencing the same array
	fmt.Println(s1) //prints [1 22 23]
	fmt.Println(s2) //prints [22 23]

	s2 = append(s2, 4)

	for i := range s2 {
		s2[i] += 10
	}

	//s1 is now "stale"
	fmt.Println(s1) //prints [1 22 23]
	fmt.Println(s2) //prints [32 33 14]
}

type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func Benchmark_point2(testB *testing.B) {
	data := []field{{"one"}, {"two"}, {"three"}}

	for _, v := range data {
		go v.print()
	}

	time.Sleep(3 * time.Second)
	//goroutines print: three, three, three

	//data := []field{{"one"},{"two"},{"three"}}
	//for _,v := range data {
	//	v := v
	//	go v.print()
	//}
	//
	//time.Sleep(3 * time.Second)
	//goroutines print: one, two, three

	//data := []*field{{"one"}, {"two"}, {"three"}}
	//for _, v := range data {
	//	go v.print()
	//}
	//time.Sleep(3 * time.Second)

}

func Benchmark_defer1(testB *testing.B) {
	var i int = 1

	defer fmt.Println("result =>", func() int { return i * 2 }())
	i++
	//prints: result => 2 (not ok if you expected 4)
}

func Benchmark_append(testB *testing.B) {
	var s []int
	s = append(s, 1)
	log.Printf("------------>>>>>>>>>>>>>>> %+v", s)
}

func Test_fallthrough(t *testing.T) {
	isSpace := func(char byte) bool {
		switch char {
		case ' ':       // 空格符会直接 break，返回 false // 和其他语言不一样
			fallthrough // 返回 true
			//case '\n':

		case '\t':
			return true
		}
		return false
	}
	fmt.Println(isSpace('\t')) // true
	fmt.Println(isSpace(' '))  // false

}

func Test_XOR(t *testing.T) {
	var a uint8 = 0x82
	var b uint8 = 0x02
	fmt.Printf("%08b [A]\n", a)
	fmt.Printf("%08b [B]\n", b)

	fmt.Printf("%08b (NOT B)\n", ^b)
	fmt.Printf("%08b ^ %08b = %08b [B XOR 0xff]\n", b, 0xff, b^0xff)

	fmt.Printf("%08b ^ %08b = %08b [A XOR B]\n", a, b, a^b)
	fmt.Printf("%08b & %08b = %08b [A AND B]\n", a, b, a&b)
	fmt.Printf("%08b &^%08b = %08b [A 'AND NOT' B]\n", a, b, a&^b)
	fmt.Printf("%08b&(^%08b)= %08b [A AND (NOT B)]\n", a, b, a&(^b))

}

// 状态名称可能是 int 也可能是 string，指定为 json.RawMessage 类型
func Test_json(t *testing.T) {
	records := [][]byte{
		[]byte(`{"status":200, "tag":"one"}`),
		[]byte(`{"status":"ok", "tag":"two"}`),
	}

	for idx, record := range records {
		var result struct {
			StatusCode uint64
			StatusName string
			Status     json.RawMessage `json:"status"`
			Tag        string          `json:"tag"`
		}

		err := json.NewDecoder(bytes.NewReader(record)).Decode(&result)
		//checkError(err)

		var name string
		err = json.Unmarshal(result.Status, &name)
		if err == nil {
			result.StatusName = name
		}

		var code uint64
		err = json.Unmarshal(result.Status, &code)
		if err == nil {
			result.StatusCode = code
		}

		fmt.Printf("[%v] result => %+v\n", idx, result)
	}
}

//
//var channel chan int = make(chan int)
//// 或
//channel := make(chan int)

func Test_channel(t *testing.T) {
	var messages chan string = make(chan string)
	go func(message string) {
		messages <- message // 存消息
	}("Ping!")

	fmt.Println(<-messages) // 取消息
}

var complete chan int = make(chan int)

func loop1() {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d\n", i)
		time.Sleep(time.Second)
	}

	complete <- 0 // 执行完毕了，发个消息
}

func Test_channel_1(t *testing.T) {
	go loop1()
	<-complete // 直到线程跑完, 取到消息. 在此阻塞住

}

var ch chan int = make(chan int)

func foo(id int) { //id: 这个routine的标号
	ch <- id
}

func Test_channel_2(t *testing.T) {
	// 开启5个routine
	for i := 0; i < 5; i++ {
		go foo(i)
	}

	// 取出信道中的数据
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Printf("%d\n", <-ch)
	}
}
