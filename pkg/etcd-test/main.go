package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	//endpoints      = []string{"100.101.197.15:3379",}
	endpoints = []string{"127.0.0.1:3379"}
)
var cli1 *clientv3.Client
var cli2 *clientv3.Client
var cli3 *clientv3.Client

func put_v() {
	cli1, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, //[]string{"10.39.0.6:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
	i := 0
	for {
		time.Sleep(time.Second * 1)
		_, err = cli1.Put(context.Background(), key_001, "world "+strconv.Itoa(i))
		if err != nil {
			log.Printf("PUT error: %s\n", err)
		}
		i++
	}
	//ctx, cancel := context.(context.Background(), requestTimeout)
	//_, err = cli1.Put(context.Background(), "/test/hello", "world")
	//cancel()
	//defer cli.Close()
}

func get_v() {
	cli2, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, //[]string{"10.39.0.6:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	for {
		time.Sleep(time.Second * 1)
		start := time.Now()
		resp, err := cli2.Get(context.Background(), key_002)
		if err != nil {
			log.Printf("GET error: %s\n", err)
		}
		for _, ev := range resp.Kvs {
			log.Printf("GET %s : %s, Version: %+v, CreateRevision: %+v, ModRevision: %+v, cost: %+v\n", ev.Key, ev.Value, ev.Version, ev.CreateRevision, ev.ModRevision, time.Now().Sub(start))
		}
	}

	//defer cli.Close()
}

func put_v1() {
	cli3, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, //[]string{"10.39.0.6:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
	i := 0
	for {
		time.Sleep(time.Second)
		_, err = cli3.Put(context.TODO(), key_003, "xyz"+strconv.Itoa(i))
		//ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		_, err = cli3.Txn(context.TODO()).
			If(clientv3.Compare(clientv3.Value(key_003), ">", "abc")).
			Then(clientv3.OpPut(key_003, "XYZ")).
			Else(clientv3.OpPut(key_003, "ABC")).
			Commit()
	}

}

func keep_alive() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, //[]string{"10.39.0.6:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
	//defer cli.Close()
	//var err error
	ctx := context.TODO()               //, cancel := context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	leaseResp, err := cli.Grant(ctx, 5) //租约时间设定为10秒
	//cancel()
	if err != nil {
		//return err
	}

	kvc := clientv3.NewKV(cli)
	var txnResp *clientv3.TxnResponse
	txnResp, err = kvc.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key_002), ">", 0)).
		Then(clientv3.OpPut(key_002, "8888", clientv3.WithLease(clientv3.LeaseID(leaseResp.ID)))).
		Else(clientv3.OpPut(key_002, "9999")).
		Commit()

	//_, err = cli.Put(ctx, key_002, "9999")
	//cancel()

	if err != nil {
		log.Printf("Put error %s\n", err)
	}

	if !txnResp.Succeeded {
		log.Printf("txnResp 失败")
	}

	//ticker := time.NewTicker(time.Second * 10)
	//for range ticker.C {
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(6)))
		//time.Sleep(time.Second * 3)
		//leaseResp, err := cli.Grant(ctx, 10) //租约时间设定为10秒
		//ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(rand.Intn(3)))
		//leaseResp, err = cli.Grant(ctx, int64(rand.Intn(15))) //租约时间设定为10秒
		_, err = cli.KeepAliveOnce(ctx, leaseResp.ID)
		//cancel()
		//v := "888888888"
		if err != nil {
			//return err
			t := int64(rand.Intn(5))
			log.Printf("续租时间: %+vs\n", t)
			//v = strconv.FormatInt(time.Now().UnixNano(),10)

			leaseResp, err = cli.Grant(ctx, t) //租约时间设定为10秒
			if err != nil {
				log.Printf("cli.Grant error %s", err)
			}
			//clientv3.OpPut(key_002, v, clientv3.WithLease(clientv3.LeaseID(leaseResp.ID)))
			cli.Put(ctx, key_002, "from: "+getLocalIp()+"  time:"+time.Now().Format("2006-01-02 15:04:05"), clientv3.WithLease(clientv3.LeaseID(leaseResp.ID)))

			continue
		}
		//cli.Put(ctx, key_002, v)
		log.Printf("续租成功 %s\n", key_002)
	}
}

type EventType string

const (
	key_001              = "/test/v1/hello"
	key_002              = "/test/v2/hello"
	key_003              = "/test/v3/hello"
	EVT_INIT   EventType = "INIT"
	EVT_CREATE EventType = "CREATE"
	EVT_UPDATE EventType = "UPDATE"
	EVT_DELETE EventType = "DELETE"
	EVT_EXPIRE EventType = "EXPIRE"
	EVT_ERROR  EventType = "ERROR"
)

var del_flag int = 0

func watch_00() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, //[]string{"10.39.0.6:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), key_002, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.DELETE:
				log.Printf("WATCH === 删除%s %q : %q,  , Version: %+v, CreateRevision: %+v, ModRevision: %+v, \n", ev.Type, ev.Kv.Key, ev.Kv.Value, ev.Kv.Version, ev.Kv.CreateRevision, ev.Kv.ModRevision)
			case mvccpb.PUT:
				log.Printf("WATCH === 存放%s %q : %q,  , Version: %+v, CreateRevision: %+v, ModRevision: %+v, \n", ev.Type, ev.Kv.Key, ev.Kv.Value, ev.Kv.Version, ev.Kv.CreateRevision, ev.Kv.ModRevision)
			}
		}
	}

	if err != nil {
		println(err)
	}
}

func getLocalIp() (IpAddr string) {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Printf("Get local IP addr failed!!!")
		IpAddr = "localhost"
		return
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				IpAddr = ipnet.IP.String()
				return
			}
		}
	}
	IpAddr = "localhost"
	return
}

type Pele struct {
	Name string
	Age  int
}

type Pels struct {
	rwMux sync.RWMutex
	Mp    map[string]*Pele
}

var p Pels

func main() {
	p = Pels{}
	m := make(map[string]*Pele, 0)
	m["key01"] = &Pele{Name: "li", Age: 19}
	m["key02"] = &Pele{Name: "lu", Age: 20}
	p.Mp = m

	go func() {
		for {
			p.rwMux.Lock()
			for k, v := range p.Mp {
				fmt.Println("init", k, v)
			}

			p.Mp["key02"].Name = "louuu"
			for k, v := range p.Mp {
				fmt.Println("update", k, v)
			}
			delete(p.Mp, "key01")

			for k, v := range p.Mp {
				fmt.Println("del", k, v)
			}
			p.rwMux.Unlock()
		}

	}()

	go func() {

		for {
			p.rwMux.Lock()
			for k, v := range p.Mp {
				fmt.Println("init", k, v)
			}

			p.Mp["key02"].Name = fmt.Sprintf("louuu")
			for k, v := range p.Mp {
				fmt.Println("update", k, v)
			}
			delete(p.Mp, "key01")

			for k, v := range p.Mp {
				fmt.Println("del", k, v)
			}
			p.rwMux.Unlock()
		}

	}()

	time.Sleep(time.Second * 30)

	return

	go put_v()
	go get_v()
	go put_v1()
	go keep_alive()
	go watch_00()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, //[]string{"10.39.0.6:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
	defer cli.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//_, err = cli.Put(ctx, "/test/hello", "world")
	//cancel()

	//ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//resp, err := cli.Get(ctx, "/test/hello")
	//cancel()
	//
	//for _, ev := range resp.Kvs {
	//	fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	//}

	//_, err = cli.Put(context.TODO(), "key", "xyz")
	//ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//_, err = cli.Txn(ctx).
	//	If(clientv3.Compare(clientv3.Value("key"), ">", "abc")).
	//	Then(clientv3.OpPut("key", "XYZ")).
	//	Else(clientv3.OpPut("key", "ABC")).
	//	Commit()
	//cancel()

	rch := cli.Watch(context.Background(), key_001, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			log.Printf("WATCH === %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

	if err != nil {
		println(err)
	}

}
