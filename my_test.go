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
