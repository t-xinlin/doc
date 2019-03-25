package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

//解析
//考点:结构体比较
//进行结构体比较时候，只有相同类型的结构体才可以比较，结构体是否相同不但与属性类型个数有关，还与属性顺序相关。
//
//sn3:= struct {
//name string
//age  int
//}{age:11,name:"qq"}
//sn3与sn1就不是相同的结构体了，不能比较。 还有一点需要注意的是结构体是相同的，但是结构体属性中有不可以比较的类型，如map,slice。 如果该结构属性都是可以比较的，那么就可以使用“==”进行比较操作。
//
//可以使用reflect.DeepEqual进行比较
//
//if reflect.DeepEqual(sn1, sm) {
//fmt.Println("sn1 ==sm")
//}else {
//fmt.Println("sn1 !=sm")
//}
//所以编译不通过： invalid operation: sm1 == sm2

//18.是否可以编译通过？如果通过，输出什么？
func Foo(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}
func Test_interface_struct_equal(t *testing.T) {
	var x *int = nil
	Foo(x)
}

//解析
//考点：interface内部结构

//19.是否可以编译通过？如果通过，输出什么？
//func GetValue1(m map[int]string, id int) (string, bool) {
//	if _, exist := m[id]; exist {
//		return "存在数据", true
//	}
//	return nil, false
//}


//解析
//考点：函数返回值类型
//nil 可以用作 interface、function、pointer、map、slice 和 channel 的“空值”。但是如果不特别指定的话，Go 语言不能识别类型，所以会报错。报:cannot use nil as type string in return argument.

//20.是否可以编译通过？如果通过，输出什么？
const (
	x = iota
	y
	z = "zz"
	k
	p = iota
)

func Test_iota_008888(t *testing.T) {
	fmt.Println(x, y, z, k, p)
}

//解析
//考点：iota
//结果:
//
//0 1 zz zz 4

//21.编译执行下面代码会出现什么?

//var (
//	size := 1024
//	max_size = size * 2
//)
//
//func Test_iota_00888822(t *testing.T) {
//	println(size, max_size)
//}

//解析
//考点:变量简短模式
//变量简短模式限制：
//
//定义变量同时显式初始化
//不能提供数据类型
//只能在函数内部使用
//结果：
//
//syntax error: unexpected :=

//22.下面函数有什么问题？

const cl = 100

var bl = 123

func Test_const_008888(t *testing.T) {
	//println(&bl, bl)
	//println(&cl, cl)
}

//解析
//考点:常量
//常量不同于变量的在运行期分配内存，常量通常会被编译器在预处理阶段直接展开，作为指令数据使用，

//cannot take the address of cl
//23.编译执行下面代码会出现什么?

func Test_goto_008888(t *testing.T) {
	//for i:=0;i<10 ;i++  {
	//loop:
	//	println(i)
	//}
	//goto loop
}

//解析
//考点：goto
//goto不能跳转到其他函数或者内层代码
//
//goto loop jumps into block starting at
//24.编译执行下面代码会出现什么?

func Test_Type_Alias_008888(t *testing.T) {
	//type MyInt1 int
	//type MyInt2 = int
	//var i int = 9
	//var i1 MyInt1 = i
	//var i2 MyInt2 = i
	//fmt.Println(i1, i2)
}

//解析
//考点：**Go 1.9 新特性 Type Alias **
//基于一个类型创建一个新类型，称之为defintion；基于一个类型创建一个别名，称之为alias。 MyInt1为称之为defintion，虽然底层类型为int类型，但是不能直接赋值，需要强转； MyInt2称之为alias，可以直接赋值。
//
//结果:
//cannot use i (type int) as type MyInt1 in assignment
//25.编译执行下面代码会出现什么?

type User1 struct {
}
type MyUser1 User1
type MyUser2 = User1

func (i MyUser1) m1() {
	fmt.Println("MyUser1.m1")
}
func (i User1) m2() {
	fmt.Println("User.m2")
}

func Test_Type_Alias_0088881111(t *testing.T) {
	var i1 MyUser1
	var i2 MyUser2
	i1.m1()
	i2.m2()
}

//解析
//考点：**Go 1.9 新特性 Type Alias **
//因为MyUser2完全等价于User，所以具有其所有的方法，并且其中一个新增了方法，另外一个也会有。 但是
//
//i1.m2()
//是不能执行的，因为MyUser1没有定义该方法。 结果:
//
//MyUser1.m1
//User.m2
//26.编译执行下面代码会出现什么?

type T1 struct {
}

func (t T1) m1() {
	fmt.Println("T1.m1")
}

type T2 = T1
type MyStruct11 struct {
	T1
	T2
}

func Test_Type_Alias(t *testing.T) {
	//my := MyStruct{}
	//my.m1()
}

//解析
//考点：**Go 1.9 新特性 Type Alias **
//是不能正常编译的,异常：
//
//ambiguous selector my.m1
//结果不限于方法，字段也也一样；也不限于type alias，type defintion也是一样的，只要有重复的方法、字段，就会有这种提示，因为不知道该选择哪个。 改为:
//
//my.T1.m1()
//my.T2.m1()
//type alias的定义，本质上是一样的类型，只是起了一个别名，源类型怎么用，别名类型也怎么用，保留源类型的所有方法、字段等。

//27.编译执行下面代码会出现什么?

var ErrDidNotWork = errors.New("did not work")

func DoTheThing(reallyDoIt bool) (err error) {
	if reallyDoIt {
		result, err := tryTheThing()
		if err != nil || result != "it worked" {
			err = ErrDidNotWork
		}
	}
	return err
}

func tryTheThing() (string, error) {
	return "", ErrDidNotWork
}

func Test_filed_00888822222(t *testing.T) {
	fmt.Println(DoTheThing(true))
	fmt.Println(DoTheThing(false))
}

//解析
//考点：变量作用域
//因为 if 语句块内的 err 变量会遮罩函数作用域内的 err 变量，结果：
//
//<nil>
//<nil>
//改为：

//func DoTheThing(reallyDoIt bool) (err error) {
//	var result string
//	if reallyDoIt {
//		result, err = tryTheThing()
//		if err != nil || result != "it worked" {
//			err = ErrDidNotWork
//		}
//	}
//	return err
//}
//28.编译执行下面代码会出现什么?

func test() []func() {
	var funs []func()
	for i := 0; i < 2; i++ {
		funs = append(funs, func() {
			println(&i, i)
		})
	}
	return funs
}

func Test_filed2222_00888822222(t *testing.T) {
	funs := test()
	for _, f := range funs {
		f()
	}
}

//解析
//考点：闭包延迟求值
//for循环复用局部变量i，每一次放入匿名函数的应用都是想一个变量。 结果：
//
//0xc042046000 2
//0xc042046000 2
//如果想不一样可以改为：
//
//func test() []func()  {
//	var funs []func()
//	for i:=0;i<2 ;i++  {
//		x:=i
//		funs = append(funs, func() {
//			println(&x,x)
//		})
//	}
//	return funs
//}
//29.编译执行下面代码会出现什么?

func test222(x int) (func(), func()) {
	return func() {
			println(x)
			x += 10
		}, func() {
			println(x)
		}
}

func Test_filed222222_00888822222(t *testing.T) {
	a, b := test222(100)
	a()
	b()
}

//解析
//考点：闭包引用相同变量*
//结果：
//
//100
//110
//30.编译执行下面代码会出现什么?

func Test_panic_revover_00888822222(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("1  ++++")
			fmt.Println(err)
		} else {
			fmt.Println("fatal")
		}
	}()

	defer func() {
		fmt.Println("2  ++++")
		panic("defer panic")
	}()
	panic("panic")
}

func Test_panic_revover_0088882222211(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("1  ++++")
			f := err.(func() string)
			fmt.Println(err, f(), reflect.TypeOf(err).Kind().String())
		} else {
			fmt.Println("fatal")
		}
	}()

	defer func() {
		panic(func() string {
			fmt.Println("2  ++++")
			return "defer panic"
		})
	}()
	//panic("panic")
	panic(func() string {
		fmt.Println("3  ++++")
		return "panic"
	})
}

//解析
//考点：panic仅有最后一个可以被revover捕获
//触发panic("panic")后顺序执行defer，但是defer中还有一个panic，所以覆盖了之前的panic("panic")
//
//defer panic

func Test_0009999(t *testing.T) {
	t.Log("========")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", "package main; var a = 0", parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	var v visitor
	ast.Walk(v, f)

}

type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	var s string
	switch node := n.(type) {
	case *ast.Ident:
		s = node.Name
	case *ast.BasicLit:
		s = node.Value
	}
	fmt.Printf("%s%T: %s\n", strings.Repeat("\t", int(v)), n, s)
	return v + 1
}

const (
	// InterfacePrio signifies the address is on a local interface
	InterfacePrio int = iota

	// BoundPrio signifies the address has been explicitly bounded to.
	BoundPrio

	// UpnpPrio signifies the address was obtained from UPnP.
	UpnpPrio

	// HTTPPrio signifies the address was obtained from an external HTTP service.
	HTTPPrio

	// ManualPrio signifies the address was provided by --externalip.
	ManualPrio
)

type ByteSize float64

const (
	_           = iota             // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota) // 1 << (10*1)
	MB                             // 1 << (10*2)
	GB                             // 1 << (10*3)
	TB                             // 1 << (10*4)
	PB                             // 1 << (10*5)
	EB                             // 1 << (10*6)
	ZB                             // 1 << (10*7)
	YB                             // 1 << (10*8)
)

func Test_iota11111(t *testing.T) {
	log.Printf("iota InterfacePrio=%d, BoundPrio=%d, UpnpPrio=%d, HTTPPrio=%d, ManualPrio=%d", InterfacePrio, BoundPrio, UpnpPrio, HTTPPrio, ManualPrio)
	log.Printf("iota KB=%+v, MB=%+v, GB=%+v, TB=%+v, PB=%+v,  EB=%+v,  ZB=%+v,  YB=%+v", KB, MB, GB, TB, PB, EB, ZB, YB)
}
