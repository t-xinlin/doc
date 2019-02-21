package main

import (
	"fmt"
	beegoyaml "github.com/astaxie/beego/config/yaml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Myconf struct {
	Ipport             string
	StartSendTime      string
	SendMaxCountPerDay int
	Devices            []Device
	WarnFrequency      int
	SendFrequency      int
}
type Device struct {
	DevId string
	Nodes []Node
}
type Node struct {
	PkId     string
	BkId     string
	Index    string
	MinValue float32
	MaxValue float32
	DataType string
}

type Cse struct {
	ConfName   string   `yaml:confname`
	Setting    Setting  `yaml:set`
	FriendShip []string `yaml:friendship`
}

type Setting struct {
	Name string `yaml:name`
	Age  int    `yaml:age`
}

func main() {
	data1, _ := ioutil.ReadFile("test/yaml/conf1.yaml")
	fmt.Println(string(data1))
	t1 := Cse{}
	//把yaml形式的字符串解析成struct类型
	yaml.Unmarshal(data1, &t1)
	fmt.Println("初始数据", t1)
	//if (t1.Ipport == "") {
	//	fmt.Println("配置文件设置错误")
	//	return;
	//}
	d1, _ := yaml.Marshal(&t1)
	fmt.Printf("看看\n%+v", string(d1))

	for k := range t1.FriendShip {
		fmt.Printf("FriendShip: %+v\n", t1.FriendShip[k])
	}

	return

	//data, _ := ioutil.ReadFile("test/yaml/conf.yaml")
	//fmt.Println(string(data))
	//t := Myconf{}
	////把yaml形式的字符串解析成struct类型
	//yaml.Unmarshal(data, &t)
	//fmt.Println("初始数据", t)
	//if (t.Ipport == "") {
	//	fmt.Println("配置文件设置错误")
	//	return;
	//}
	//d, _ := yaml.Marshal(&t)
	//fmt.Println("看看 :", string(d))
	//return
	conf, err := beegoyaml.ReadYmlReader("test/yaml/conf.yaml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("init: %+v", conf)
	//name := conf["name"]
	//favourite := conf["favourite"]
	//fmt.Println(name)
	//fmt.Println(favourite)
}
