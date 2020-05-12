package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/t-xinlin/doc/pkg/proto/pb"
	"io/ioutil"
	"os"
)

func write() {
	p1 := &pb.Person{
		Id:   1,
		Name: "小张",
		Phones: []*pb.Phone{
			{Type: pb.PhoneType_HOME, Number: "111111111"},
			{Type: pb.PhoneType_WORK, Number: "222222222"},
		},
	}
	p2 := &pb.Person{
		Id:   2,
		Name: "小王",
		Phones: []*pb.Phone{
			{Type: pb.PhoneType_HOME, Number: "333333333"},
			{Type: pb.PhoneType_WORK, Number: "444444444"},
		},
	}

	//创建地址簿
	book := &pb.ContactBook{}
	book.Persons = append(book.Persons, p1)
	book.Persons = append(book.Persons, p2)

	//编码数据
	data, _ := proto.Marshal(book)
	//把数据写入文件
	ioutil.WriteFile("./proto/test.txt", data, os.ModePerm)
}

func read() {
	//读取文件数据
	data, _ := ioutil.ReadFile("./proto/test.txt")
	book := &pb.ContactBook{}
	//解码数据
	proto.Unmarshal(data, book)
	for _, v := range book.Persons {
		fmt.Println(v.Id, v.Name)
		for _, vv := range v.Phones {
			fmt.Println(vv.Type, vv.Number)
		}
	}
}

func main() {
	write()
	read()
}
