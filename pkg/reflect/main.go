package main

import (
	"fmt"
	"reflect"
)

func main() {
	user := User{1, 19, "Allen.Wu", "add"}
	DofileAndMethod(user)

}

type User struct {
	Id   int
	Age  int
	Name string
	Add  string
}

func (u User) ReflecFunc(name string, age int) {
	fmt.Printf("Method Call>>>>>>>>>>>>>>User: %+v Name: %s  Age: %d\n", u, name, age)
}

func DofileAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)
	fmt.Println("Get type is ", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("Get  all filed is ", getValue)

	pointer := reflect.ValueOf(&input)
	newv := pointer.Elem()
	newv.CanSet()
	newv.Set(reflect.ValueOf(User{111, 23, "zhangsan", "shenzhen"}))
	fmt.Println("User input: ", input)

	mValue := getValue.MethodByName("ReflecFunc")
	args := []reflect.Value{reflect.ValueOf("lss"), reflect.ValueOf(20)}
	mValue.Call(args)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
}
