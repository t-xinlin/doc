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

