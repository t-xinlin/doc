package main

import (
	"github.com/t-xinlin/doc/test/rpc/interfaces"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	//client, err := interfaces.DialHelloService("tcp", "127.0.0.1:7777")
	client, err := interfaces.DialHelloServiceJSON("tcp", "127.0.0.1:7777")
	if nil != err {
		log.Fatal("Dial error ", err)
	}

	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		var reply string
		//err = client.Call("HelloService.SayHello", "Hi", &reply)
		client.SayHello("Hi", &reply)
		if nil != err {
			log.Printf("call error %s", err.Error())
		}
		log.Printf("Reply: %s", reply)

		var d string =	`{"method":"com.xl.interface.SayHello","params":["hello"],"id":0}`
		resp, err :=http.Post("http://127.0.0.1:6666/jsonrpc","application/json", strings.NewReader(d))
		if err != nil{
			log.Printf("Error: %s", err.Error())
		} else {
			bytes, _ := ioutil.ReadAll(resp.Body)
			log.Printf("httpRpc: %s", string(bytes))
		}

	}
}
