///////////////////////////黑魔法？
package main

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"
)

type MyStruct struct {
	A int
	B int
}

var sizeOfMyStruct = int(unsafe.Sizeof(MyStruct{}))

func MyStructToBytes(s *MyStruct) []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfMyStruct
	x.Cap = sizeOfMyStruct
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func BytesToMyStruct(b []byte) *MyStruct {
	return (*MyStruct)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}

// returns &s[0], which is not allowed in go
func stringPointer(s string) unsafe.Pointer {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(p.Data)
}

// returns &b[0], which is not allowed in go
func bytePointer(b []byte) unsafe.Pointer {
	p := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(p.Data)
}

// convert b to string without copy
func byteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type MuxType string

func (m *MuxType) String() (s string) {
	pMT := (*reflect.StringHeader)(unsafe.Pointer(m))
	pStr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pStr.Data = pMT.Data
	pStr.Len = pMT.Len
	return
}

type Option struct {
	Key   string
	Value interface{}
}

type PluginOpOption func(*Option)

func OptionsToOp(opts ...PluginOpOption) (op Option) {
	for _, opt := range opts {
		opt(&op)
	}
	//if op.Limit == 0 {
	//	op.Offset = -1
	//	op.Limit = DEFAULT_PAGE_COUNT
	//}
	return
}

func Do(ctx context.Context, opts ...PluginOpOption) {
	op := OptionsToOp(opts...)
	fmt.Printf("==opt==%+v", op)
}

var Key PluginOpOption = func(op *Option) { op.Key = "key1111" }
var Value PluginOpOption = func(op *Option) { op.Value = "value222" }

//自重写程序
var q = `/* Go quine */
package main

import "fmt"

func main() {
    fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60)
}

var q = `

func main() {
	/* Go quine */

	fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60)

	//Do(nil, Value, Key)
	//
	//interval, _ := time.ParseDuration("5s")
	//fmt.Printf("interval %d", interval)
	//return
	//
	//var s MuxType = "hello world"
	//fmt.Printf("s==%p==\n", &s)
	//s1 := s.String()
	//fmt.Printf("s1==%p==\n", &s1)
	//fmt.Printf("s1 string==%s==\n", s1)
	//
	//fmt.Printf("==%s==\n", byteString([]byte("hello world!")))
	//return
}

//自重写程序
func main1() { print(c + "\x60" + c + "\x60") }

var c = `package main;func main(){print(c+"\x60"+c+"\x60")};var c=`

//aa6ro2g@icloud.com
//ATt778877
//
//jb2h4r7
//
//ir6z9sl
