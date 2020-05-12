package main

import (
	"bytes"
	"container/list"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func main() {
	fmt.Println("run")
	http.HandleFunc("/spans", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("ReadAll error: %+v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		if r != nil {
			defer r.Body.Close()
		}
		log.Printf(">>>>>>>>>>>>>>>>>>>>>%+v", string(bytes))
		var status int = 200
		w.WriteHeader(status)
		log.Printf("Rec: : %+v  -> Back[  version: 1.0.0   httpStatusCode: %+v]", string(bytes), status)
		w.Write([]byte("Svc version: 1.0.0"))
	})
	http.ListenAndServe("127.0.0.1:9898", nil)
}

// ICMP is use
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

// DefaultExpireTime is time out
const DefaultExpireTime = 0

// CheckSum is time out
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
				w.Refuse()
				fmt.Println("worker", i, "refused", w)
				break
			}
			w.Do()
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

	for {
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
	b = 1 << (10 * iota)
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

	for {
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

		for i := range arr {
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

	for k := range arr {
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
		case ' ': // 空格符会直接 break，返回 false // 和其他语言不一样
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

func Test_channel_3(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	fmt.Printf("go1\n") // 1
	ch <- 2
	fmt.Printf("go2\n") // 2
	ch <- 3
	fmt.Printf("go3\n")                       // 3
	fmt.Printf("-------------------------\n") // 3
	fmt.Printf("go%d\n", <-ch)                // 1
	fmt.Printf("go%d\n", <-ch)                // 2
	fmt.Printf("go%d\n", <-ch)                // 3

}

//在类型断言语句中，断言失败则会返回目标类型的“零值”，断言变量与原来变量混用可能出现异常情况：
func Test_asset1(t *testing.T) {
	var data interface{} = "great"

	// data 混用
	if data, ok := data.(int); ok {
		fmt.Println("[is an int], data: ", data)
	} else {
		fmt.Println("[not an int], data: ", data) // [isn't a int], data:  0
	}

}

func Test_asset2(t *testing.T) {
	var data interface{} = "great"

	if res, ok := data.(int); ok {
		fmt.Println("[is an int], data: ", res)
	} else {
		fmt.Println("[not an int], data: ", data) // [not an int], data:  great
	}

}

type data struct {
	name string
}

type printer interface {
	print()
}

func (p *data) print() {
	fmt.Println("name: ", p.name)
}

//使用指针作为方法的 receiver
//只要值是可寻址的，就可以在值上直接调用指针方法。即是对一个方法，它的 receiver 是指针就足矣。

//但不是所有值都是可寻址的，比如 map 类型的元素、通过 interface 引用的变量：
func Test_point(t *testing.T) {
	d1 := data{"one"}
	d1.print() // d1 变量可寻址，可直接调用指针 receiver 的方法

	//var in printer = data{"two"}
	//in.print() // 类型不匹配

	//m := map[string]data{
	//	"x": data{"three"},
	//}
	//m["x"].print() // m["x"] 是不可寻址的    // 变动频繁
}

//如果 map 一个字段的值是 struct 类型，则无法直接更新该 struct 的单个字段：
// 无法直接更新 struct 的字段值
//因为 map 中的元素是不可寻址的。需区分开的是，slice 的元素可寻址
func Test_uptate_map(t *testing.T) {
	//m := map[string]data{
	//	"x": {"Tom"},
	//}
	//m["x"].name = "Jerry"

	s := []data{{"Tom"}, {name: "Tonny"}}
	s[0].name = "Jerry"
	fmt.Println(s) // [{Jerry}]

	m := map[string]*data{
		"x": {"Tom"},
	}

	m["x"].name = "Jerry" // 直接修改 m["x"] 中的字段
	fmt.Println(m["x"])   // &{Jerry}

	//m = map[string]*data{
	//	"x": {"Tom"},
	//}
	//m["z"].name = "what???"//报错
	//fmt.Println(m["x"])
}

//53. nil interface 和 nil interface 值
//虽然 interface 看起来像指针类型，但它不是。interface 类型的变量只有在类型和值均为 nil 时才为 nil
//
//如果你的 interface 变量的值是跟随其他变量变化的（雾），与 nil 比较相等时小心：

func Test_interface_nil(t *testing.T) {
	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) // <nil> true
	fmt.Println(in, in == nil)     // <nil> true

	in = data
	fmt.Println(in, in == nil) // <nil> false    // data 值为 nil，但 in 值不为 nil
	fmt.Printf("=========%+v", in)
}

//如果你的函数返回值类型是 interface，更要小心这个坑：
// 错误示例
func Test_interface_nil_1(t *testing.T) {
	doIt := func(arg int) interface{} {
		var result *struct{} = nil
		if arg > 0 {
			result = &struct{}{}
		}
		return result
	}

	if res := doIt(-1); res != nil {
		fmt.Println("Good result: ", res) // Good result:  <nil>
		fmt.Printf("%T\n", res)           // *struct {}    // res 不是 nil，它的值为 nil
		fmt.Printf("%v\n", res)           // <nil>
	}
}

// 正确示例
func Test_interface_nil_2(t *testing.T) {
	doIt := func(arg int) interface{} {
		var result *struct{} = nil
		if arg > 0 {
			result = &struct{}{}
		} else {
			return nil // 明确指明返回 nil
		}
		return result
	}

	if res := doIt(-1); res != nil {
		fmt.Println("Good result: ", res)
	} else {
		fmt.Println("Bad result: ", res) // Bad result:  <nil>
	}
}

//type data1 struct {
//	Name string
//}

func chtest() {

}

// channel
func Test_channel_02(t *testing.T) {
	ch := make(chan int, 5)
	//ch1 := make(chan int, 5)
	go func() {
		tick := time.NewTicker(time.Second * 2)
		i := 0
		for {
			i++
			//	ch <- i
			fmt.Printf("Tick\n")
			<-tick.C
		}

	}()
	j := 0
	tick := time.NewTicker(time.Second * 2)
	for {
		j++
		fmt.Printf("j: %+v\n", j)
		select {
		case <-tick.C:
			fmt.Printf("tick.C: %+v\n", j)

		case v := <-ch:
			fmt.Printf("ch: %+v\n", v)
			//default:
			//	fmt.Printf("default\n")

			//default:
			//	fmt.Printf("default\n")
		}
	}

}

//golang读取关闭channel遇到的问题/如何优雅关闭channel
func TestReadFromClosedChan(t *testing.T) {
	asChan := func(vs ...int) <-chan int {
		c := make(chan int)
		go func() {
			for _, v := range vs {
				c <- v
				time.Sleep(time.Second)
			}
			close(c)
		}()
		return c
	}
	merge := func(a, b <-chan int) <-chan int {
		c := make(chan int)
		go func() {
			for {
				select {
				case v := <-a:
					c <- v
				case v := <-b:
					c <- v
				}
			}
		}()
		return c
	}

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}

//   chan<- //只写
func counter(out chan<- int) {
	defer close(out)
	for i := 0; i < 100; i++ {
		fmt.Printf("生产%+v\n", i)
		out <- i //如果对方不读 会阻塞
		time.Sleep(time.Second * 1)
	}
}

//   <-chan //只读
func printer1(in <-chan int) {
	for {
		for num := range in {
			fmt.Printf("消费%+v\n", num)
		}
	}

}

//chan
func TestChan(t *testing.T) {
	//ticker := time.NewTicker(time.Second * 50)

	ch := make(chan int) //   chan   //读写

	go counter(ch)  //生产者
	go printer1(ch) //消费者

	//<-ticker.C
	<-time.After(time.Second * 10)
	fmt.Println("done")
}

func TestDefer(t *testing.T) {
	a1()
	b1()
	c := c1()
	fmt.Printf("c= %+v\n", c)

}

//规则一 当defer被声明时，其参数就会被实时解析
//我们通过以下代码来解释这条规则:
func a1() {
	i := 0
	defer fmt.Printf("a1() run i=%+v\n", i)
	i++
	return
}

func b1() {
	for i := 0; i < 4; i++ {
		defer fmt.Printf("b1 i= %+v\n", i)
	}
}

//输出结果是12. 在开头的时候，我们说过defer是在return调用之后才执行的。 这里需要明确的是defer代码块的作用域仍然在函数之内，结合上面的函数也就是说，defer的作用域仍然在c函数之内。因此defer仍然可以读取c函数内的变量(如果无法读取函数内变量，那又如何进行变量清除呢....)。
//当执行return 1 之后，i的值就是1. 此时此刻，defer代码块开始执行，对i进行自增操作。 因此输出2.
//掌握了defer以上三条使用规则，那么当我们遇到defer代码块时，就可以明确得知defer的预期结果。
func c1() (i int) {
	defer func() {
		i++
	}()
	i = 1
	return 8
}

func Test_IP_PORT(t *testing.T) {
	url1 := fmt.Sprintf("http://%s%s/v1/mesher/health", "127.0.0.1", ":30102")

	fmt.Printf("kkkkkkkkk-======%+v\n", url1)

	//我们将解析这个 URL 示例，它包含了一个 scheme，认证信息，主机名，端口，路径，查询参数和片段。
	//s := "postgres://user:pass@host.com:5432/path?k=v#f"

	s := "http://197.1.0.1:9999"

	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(u.Host)
	h1 := strings.Split(u.Host, ":")
	fmt.Println(h1[0])
	fmt.Println(h1[1])

	return

	//直接访问 scheme。
	fmt.Println(u.Scheme)
	//User 包含了所有的认证信息，这里调用 Username和 Password 来获取独立值。
	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	p, _ := u.User.Password()
	fmt.Println(p)
	//Host 同时包括主机名和端口信息，如过端口存在的话，使用 strings.Split() 从 Host 中手动提取端口。
	fmt.Println(u.Host)
	h := strings.Split(u.Host, ":")
	fmt.Println(h[0])
	fmt.Println(h[1])
	//这里我们提出路径和查询片段信息。
	fmt.Println(u.Path)
	fmt.Println(u.Fragment)
	//要得到字符串中的 k=v 这种格式的查询参数，可以使用 RawQuery 函数。你也可以将查询参数解析为一个map。已解析的查询参数 map 以查询字符串为键，对应值字符串切片为值，所以如何只想得到一个键对应的第一个值，将索引位置设置为 [0] 就行了。
	fmt.Println(u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["k"][0])
}

func Lookup(meta []int32, target int32) {
	left := 0
	right := len(meta) - 1
	for i := 0; i < len(meta); i++ {
		if meta[left]+meta[right] > target {
			right--
		} else if meta[left]+meta[right] < target {
			left++
		} else {
			fmt.Println(fmt.Sprintf("%d, %d", meta[left], meta[right]))
			return
		}
	}
	fmt.Println("未找到匹配数据")
}

//1.如何在一个给定有序数组中找两个和为某个定值的数，要求时间复杂度为O(n), 比如给｛1，2，4，5，8，11，15｝和15？
func Test_Lookup(t *testing.T) {
	fmt.Printf("1<<31-1 = %d", 1<<31-1)
	arr := []int32{1, 2, 3, 6, 8, 9, 11, 12, 14, 15, 16, 18, 20, 23, 35, 27, 29, 30, 31, 32, 39, 45, 46, 58, 59, 60, 62}
	Lookup(arr, 58)
}

//2.给定一个数组代表股票每天的价格，请问只能买卖一次的情况下，最大化利润是多少？
// 日期不重叠的情况下，可以买卖多次呢？ 输入：{100,80,120,130,70,60,100,125}，
// 只能买一次：65(60买进，125卖出)；
// 可以买卖多次：115(80买进，130卖出；60买进，125卖出)？

func Test_Buy_sale(t *testing.T) {
	price := []int{100, 80, 120, 130, 70, 60, 100, 125}
	var buyDay = -1
	var saleDay = -1

	var buyPrice = 1000
	var salePrice = -1000

	type Op struct {
		BuyDay, SaleDay     int
		BuyPrice, SalePrice int
		Ear                 int
	}

	var opList = []Op{}
	for k, todayPrice := range price {
		//find buy price
		if buyDay == -1 {
			if todayPrice < buyPrice {
				buyPrice = todayPrice
				continue
			}
			buyDay = k - 1
			continue
		}

		//find sale price
		if todayPrice > salePrice {
			salePrice = todayPrice
			if k < len(price)-1 {
				continue
			}
		}

		if k < len(price)-1 {
			saleDay = k - 1
		} else {
			saleDay = k
		}

		opList = append(opList, Op{
			BuyDay:    buyDay,
			BuyPrice:  buyPrice,
			SaleDay:   saleDay,
			SalePrice: salePrice,
			Ear:       salePrice - buyPrice,
		})

		buyDay = -1
		saleDay = -1

		buyPrice = 1000
		salePrice = -1000

	}

	for _, v := range opList {
		fmt.Printf("==%+v\n", v)
	}
}

//1. 写出下面代码输出内容。
func Test_defer_call(t *testing.T) {
	defer_call()
}

func defer_call() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	panic("触发异常")
}

//考点：defer执行顺序
//解答：
//defer 是后进先出。
//panic 需要等defer 结束后才会向上传递。 出现panic恐慌时候，会先按照defer的后入先出的顺序执行，最后才会执行panic。

//2. 以下代码有什么问题，说明原因。

func Test_foreachl(t *testing.T) {
	type student struct {
		Name string
		Age  int
	}
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	// 错误写法
	for _, stu := range stus {
		m[stu.Name] = &stu
	}

	for k, v := range m {
		println(k, "错误写法=>", v.Name)
	}

	// 正确
	for i := 0; i < len(stus); i++ {
		m[stus[i].Name] = &stus[i]
	}
	for k, v := range m {
		println(k, "正确写法=>", v.Name)
	}
}

//知识：foreach
//解答：
//这样的写法初学者经常会遇到的，很危险！ 与Java的foreach一样，都是使用副本的方式。所以m[stu.Name]=&stu实际上一致指向同一个指针， 最终该指针的值为遍历的最后一个struct的值拷贝。 就像想修改切片元素的属性：

//3. 下面的代码会输出什么，并说明原因
func Test_goroutine_001(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("A: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("B: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

//考点：go执行的随机性和闭包
//解答：
//谁也不知道执行后打印的顺序是什么样的，所以只能说是随机数字。 但是A:均为输出10，B:从0~9输出(顺序不定)。 第一个go func中i是外部for的一个变量，地址不变化。遍历完成后，最终i=10。 故go func执行时，i的值始终是10。
//
//第二个go func中i是函数参数，与外部for中的i完全是两个变量。 尾部(i)将发生值拷贝，go func内部指向值拷贝地址。

//4. 下面代码会输出什么？

type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("People showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func oop_00() {
	//var t People
	t := Teacher{}
	t.ShowA()
	t.ShowB()
}
func Test_oop(t *testing.T) {
	oop_00()
}

//知识：go的组合继承
//解答：
//这是Golang的组合模式，可以实现OOP的继承。 被组合的类型People所包含的方法虽然升级成了外部类型Teacher这个组合类型的方法（一定要是匿名字段），但它们的方法(ShowA())调用时接受者并没有发生变化。 此时People类型并不知道自己会被什么类型组合，当然也就无法调用方法时去使用未知的组合者Teacher类型的功能。

//5. 下面代码会触发异常吗？请详细说明
func Test_select_rand(t *testing.T) {
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}
}

//考点：select随机性
//解答：
//select会随机选择一个可用通用做收发操作。 所以代码是有肯触发异常，也有可能不会。 单个chan如果无缓冲时，将会阻塞。但结合 select可以在多个chan间等待执行。有三点原则：
//
//select 中只要有一个case能return，则立刻执行。
//当如果同一时间有多个case均能return则伪随机方式抽取任意一个执行。
//如果没有一个case能return则可以执行”default”块。

//6.下面代码输出什么？
func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

//考点：defer执行顺序
//解答：
//这道题类似第1题 需要注意到defer执行顺序和值传递 index:1肯定是最后执行的，但是index:1的第三个参数是一个函数，所以最先被调用calc("10",1,2)==>10,1,2,3 执行index:2时,与之前一样，需要先调用calc("20",0,2)==>20,0,2,2 执行到b=1时候开始调用，index:2==>calc("2",0,2)==>2,0,2,2 最后执行index:1==>calc("1",1,3)==>1,1,3,4

func Test_defer_000(t *testing.T) {
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1
}

//7.请写出以下輸出内容
func Test_append_000(t *testing.T) {
	s := make([]int, 5)
	fmt.Println(s)
	s = append(s, 1, 2, 3)
	fmt.Println(s)
}

//考点：make默认值和append
//解答：
//make初始化是由默认值的哦，此处默认值为0

//8.下面的代码有什么问题?
type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}

//func (ua *UserAges) Get(name string) int {
//	if age, ok := ua.ages[name]; ok {
//		return age
//	}
//	return -1
//}

//修改后
func (ua *UserAges) Get(name string) int {
	ua.Lock()
	defer ua.Unlock()
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

func Test_map_routine(t *testing.T) {
	u := UserAges{}
	u.ages = make(map[string]int, 0)

	tiker := time.NewTicker(time.Second * 30)

	go func() {
		i := 0
		for {
			i++
			key := fmt.Sprintf("key_%d", i)
			u.Add(key, i)
			fmt.Printf("ADD  key=%+v, value=%+v\n", key, i)

		}
	}()

	go func() {
		time.Sleep(time.Second * 3)
		i := 0
		for {
			i++
			key := fmt.Sprintf("key_%d", i)
			fmt.Printf("GET  key=%+v, value=%+v\n", key, u.Get(key))
		}
	}()

	<-tiker.C
}

//考点：map线程安全
//解答：
//可能会出现fatal error: concurrent map read and map write. 修改一下看看效果

//9. 下面的迭代会有什么问题？
type threadSafeSet struct {
	sync.RWMutex
	s []interface{}
}

func (set *threadSafeSet) Iter() <-chan interface{} {
	// ch := make(chan interface{}) // 解除注释看看！
	ch := make(chan interface{}, len(set.s))
	go func() {
		set.RLock()

		for elem, value := range set.s {
			ch <- elem
			fmt.Printf("Iter:%+v  %+v\n", elem, value)
		}

		close(ch)
		set.RUnlock()

	}()
	return ch
}

func Test_Iter(t *testing.T) {
	th := threadSafeSet{
		s: []interface{}{"1", "2"},
	}
	//v := <-th.Iter()

	for v := range th.Iter() {
		fmt.Printf("%+v %+v\n", "ch", v)
	}

}

//考点：chan缓存池
//解答：
//看到这道题，我也在猜想出题者的意图在哪里。 chan?sync.RWMutex?go?chan缓存池?迭代? 所以只能再读一次题目，就从迭代入手看看。 既然是迭代就会要求set.s全部可以遍历一次。但是chan是为缓存的，那就代表这写入一次就会阻塞。 我们把代码恢复为可以运行的方式，看看效果

//10. 以下代码能编译过去吗？为什么？

type People1 interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func Test_oop_interface(t *testing.T) {
	var peo People1 = &Stduent{}
	think := "bitch"
	fmt.Println(peo.Speak(think))
}

//考点：golang的方法集
//解答：
//编译不通过！ 做错了！？说明你对golang的方法集还有一些疑问。 一句话：golang的方法集仅仅影响接口实现和方法表达式转化，与通过实例或者指针调用方法无关。

//11. 以下代码打印出来什么内容，说出为什么。

type People2 interface {
	Show()
}

type Student struct{}

func (stu *Student) Show() {

}

func live() People2 {
	var stu *Student
	return stu
}

func Test_oop_interface22(t *testing.T) {
	if live() == nil {
		fmt.Println("AAAAAAA")
	} else {
		fmt.Println("BBBBBBB")
	}
}

//考点：interface内部结构
//解答：
//很经典的题！ 这个考点是很多人忽略的interface内部结构。 go中的接口分为两种一种是空的接口类似这样：
//
//var in interface{}
//另一种如题目：
//
//type People interface {
//	Show()
//}
//他们的底层结构如下：
//
//type eface struct {      //空接口
//	_type *_type         //类型信息
//	data  unsafe.Pointer //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
//}
//type iface struct {      //带有方法的接口
//	tab  *itab           //存储type信息还有结构实现方法的集合
//	data unsafe.Pointer  //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
//}
//type _type struct {
//	size       uintptr  //类型大小
//	ptrdata    uintptr  //前缀持有所有指针的内存大小
//	hash       uint32   //数据hash值
//	tflag      tflag
//	align      uint8    //对齐
//	fieldalign uint8    //嵌入结构体时的对齐
//	kind       uint8    //kind 有些枚举值kind等于0是无效的
//	alg        *typeAlg //函数指针数组，类型实现的所有方法
//	gcdata    *byte
//	str       nameOff
//	ptrToThis typeOff
//}
//type itab struct {
//	inter  *interfacetype  //接口类型
//	_type  *_type          //结构类型
//	link   *itab
//	bad    int32
//	inhash int32
//	fun    [1]uintptr      //可变大小 方法集合
//}
//可以看出iface比eface 中间多了一层itab结构。 itab 存储_type信息和[]fun方法集，从上面的结构我们就可得出，因为data指向了nil 并不代表interface 是nil， 所以返回值并不为空，这里的fun(方法集)定义了接口的接收规则，在编译的过程中需要验证是否实现接口 结果：
//
//BBBBBBB

//12.是否可以编译通过？如果通过，输出什么？
func Test_type_009(t *testing.T) {
	//i := GetValue()
	//switch i.(type) {
	//case int:
	//	println("int")
	//case string:
	//	println("string")
	//case interface{}:
	//	println("interface")
	//default:
	//	println("unknown")
	//}

}

func GetValue() int {
	return 1
}

//解析
//考点：type
//编译失败，因为type只能使用在interface

//13.下面函数有什么问题？
//func funcMui(x,y int)(sum int,error){
//	return x+y,nil
//}
//解析
//考点：函数返回值命名
//在函数有多个返回值时，只要有一个返回值有指定命名，其他的也必须有命名。 如果返回值有有多个返回值必须加上括号； 如果只有一个返回值并且有命名也需要加上括号； 此处函数第一个返回值有sum名称，第二个未命名，所以错误。

//14.是否可以编译通过？如果通过，输出什么？

func Test_defer_00888(t *testing.T) {
	println(DeferFunc1(1))
	println(DeferFunc2(1))
	println(DeferFunc3(1))
}

func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc2(i int) int {
	t := i
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc3(i int) (t int) {
	defer func() {
		t += i
	}()
	return 2
}

//解析
//考点:defer和函数返回值
//需要明确一点是defer需要在函数结束前执行。 函数返回值名字会在函数起始处被初始化为对应类型的零值并且作用域为整个函数 DeferFunc1有函数返回值t作用域为整个函数，在return之前defer会被执行，所以t会被修改，返回4; DeferFunc2函数中t的作用域为函数，返回1; DeferFunc3返回3

//15.是否可以编译通过？如果通过，输出什么？
func Test_new_008888(t *testing.T) {
	//list := new([]int)
	//list = append(list, 1)
	//fmt.Println(list)
}

//解析
//考点：new
//list:=make([]int,0)

//16.是否可以编译通过？如果通过，输出什么？

func Test_append_008888(t *testing.T) {
	//s1 := []int{1, 2, 3}
	//s2 := []int{4, 5}
	//s1 = append(s1, s2)
	//fmt.Println(s1)
}

//解析
//考点：append
//append切片时候别漏了'...'

//17.是否可以编译通过？如果通过，输出什么？
func Test_struct_compare_008888(t *testing.T) {
	//sn1 := struct {
	//	age  int
	//	name string
	//}{age: 11, name: "qq"}
	//sn2 := struct {
	//	age  int
	//	name string
	//}{age: 11, name: "qq"}
	//
	//if sn1 == sn2 {
	//	fmt.Println("sn1 == sn2")
	//}
	//
	//sm1 := struct {
	//	age int
	//	m   map[string]string
	//}{age: 11, m: map[string]string{"a": "1"}}
	//sm2 := struct {
	//	age int
	//	m   map[string]string
	//}{age: 11, m: map[string]string{"a": "1"}}
	//
	//if sm1 == sm2 {
	//	fmt.Println("sm1 == sm2")
	//}
}
