package main

//
//import (
//	"context"
//	"errors"
//	"github.com/coreos/etcd/clientv3"
//	"github.com/coreos/etcd/clientv3/concurrency"
//	"log"
//	"math/rand"
//	"sync"
//	"time"
//)
//
//var globalMux sync.Mutex
//
//type MyLock struct {
//	Key       string
//	Endpoints []string
//	Session   *concurrency.Session
//	Cli       *clientv3.Client
//	Mutex     *concurrency.Mutex
//}
//
//func Lock(ctx context.Context, key string, endpoints []string) (lock *MyLock, err error) {
//	globalMux.Lock()
//	defer globalMux.Unlock()
//	lock = &MyLock{Key: key, Endpoints: endpoints}
//	err = lock.lock(ctx)
//
//	return
//}
//
//func (mylocak *MyLock) lock(ctx context.Context) (err error) {
//	endpoints := mylocak.Endpoints //[]string{"http://100.101.197.164:3379"} //strings.Split(config.Conf.EtcdClusterNodes, ",")
//	if endpoints == nil || len(endpoints) == 0 {
//		log.Printf("Endpoints is empty.")
//		return errors.New("Endpoints is empty.")
//	} //如果采用分布式锁
//	mylocak.Cli, err = clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Printf("ETCD NewEtcd cli failed , error is %v ", err)
//		return err
//	}
//	mylocak.Session, err = concurrency.NewSession(mylocak.Cli)
//	if err != nil {
//		log.Printf("ETCD NewSession failed , error is %v ", err)
//		return err
//	}
//
//	mylocak.Mutex = concurrency.NewMutex(mylocak.Session, "/my-lock-mbcp-bc-als/")
//	log.Printf("Wait unlock")
//	if err := mylocak.Mutex.Lock(ctx); err != nil {
//		log.Printf("ETCD NewMutex failed , error is %v ", err)
//		return err
//	}
//
//	log.Printf("Get the lock")
//	//谁先抢到这个分布式锁谁就可以得到执行机会，抢不到则阻塞在这里
//	return nil
//
//}
//
//func (mylocak *MyLock) unLock() (err error) {
//	//endpoints := mylocak.Endpoints //[]string{"http://100.101.197.164:3379"} //strings.Split(config.Conf.EtcdClusterNodes, ",")
//	//if endpoints != nil && len(endpoints) != 0 { //如果采用分布式锁
//	//	log.Printf("unLock !!!")
//	//	//
//	//	defer mylocak.Cli.Close()
//	//	defer mylocak.Session.Close()
//	//	defer mylocak.Mutex.Unlock(context.TODO())
//	//	//谁先抢到这个分布式锁谁就可以得到执行机会，抢不到则阻塞在这里
//	//}
//
//	log.Printf("UnLock")
//	if mylocak.Cli != nil {
//		defer mylocak.Cli.Close()
//	}
//
//	if mylocak.Session != nil {
//		defer mylocak.Session.Close()
//	}
//	if mylocak.Mutex != nil {
//		defer mylocak.Mutex.Unlock(context.TODO())
//	}
//
//	return nil
//}
//
//func test_lock() {
//	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
//	lock, err := Lock(ctx, "/my-lock-mbcp-bc-als/", []string{"http://100.101.197.168:3379"})
//	if err != nil {
//		log.Printf("Get the lock error: %s", err)
//	}
//	//
//	//for i := range random(1) {
//	//	//i = time.Duration(random(1))//random(1)
//	//	time.Sleep(time.Second * time.Duration(i))
//	//}
//
//	time.Sleep(time.Second * time.Duration(rand.Intn(3)))
//
//	defer lock.unLock()
//}
//
////func getRandom(){
////	for i := range random(1) {
////		//i = time.Duration(random(1))//random(1)
////		time.Sleep(time.Second * time.Duration(i))
////	}
////}
//
//func random(n int) <-chan int {
//	c := make(chan int)
//	go func() {
//		defer close(c)
//		for i := 0; i < n; i++ {
//			select {
//			case c <- 0:
//			case c <- 1:
//			case c <- 2:
//			case c <- 3:
//			case c <- 4:
//			case c <- 5:
//			case c <- 6:
//			case c <- 7:
//			case c <- 8:
//			case c <- 9:
//			}
//		}
//	}()
//	return c
//}
//
//func main() {
//	ticker := time.NewTicker(time.Second * 1)
//	for range ticker.C {
//		test_lock()
//	}
//
//	//etcdsync.SetDebug(true)
//	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
//	//m := etcdsync.New("/etcdsync", 30, []string{"http://100.101.197.15:3379"})
//	//if m == nil {
//	//	log.Printf("etcdsync.NewMutex failed")
//	//}
//	//err := m.Lock()
//	//if err != nil {
//	//	log.Printf("etcdsync.Lock failed")
//	//} else {
//	//	log.Printf("etcdsync.Lock OK")
//	//}
//	//
//	//log.Printf("Get the lock. Do something here.")
//	//
//	//err = m.Unlock()
//	//if err != nil {
//	//	log.Printf("etcdsync.Unlock failed")
//	//} else {
//	//	log.Printf("etcdsync.Unlock OK")
//	//}
//
//	//endpoints := []string{"http://100.101.197.164:3379"} //strings.Split(config.Conf.EtcdClusterNodes, ",")
//	//if endpoints != nil && len(endpoints) != 0 { //如果采用分布式锁
//	//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	//	if err != nil {
//	//		//lib.Log.Panicf("ETCD NewEtcd cli failed , error is %v " , err)
//	//		log.Printf("ETCD NewEtcd cli failed , error is %v ", err)
//	//	}
//	//	defer cli.Close()
//	//	s, err := concurrency.NewSession(cli)
//	//	if err != nil {
//	//		log.Printf("ETCD NewSession failed , error is %v ", err)
//	//	}
//	//	defer s.Close()
//	//	log.Printf("wait unlock")
//	//	m := concurrency.NewMutex(s, "/my-lock-mbcp-bc-als/")
//	//	if err := m.Lock(context.TODO()); err != nil {
//	//		log.Printf("ETCD NewMutex failed , error is %v ", err)
//	//	}
//	//	log.Printf("Get the lock. Do something here.")
//	//	time.Sleep(time.Second * 10)
//	//	defer m.Unlock(context.TODO())
//	//	//谁先抢到这个分布式锁谁就可以得到执行机会，抢不到则阻塞在这里
//	//}
//
//	////m.pfx是前缀，比如"service/lock/"
//	////s.Lease()是一个64位的整数值，etcd-test v3引入了lease（租约）的概念，concurrency包基于lease封装了session，每一个客户端都有自己的lease，也就是说每个客户端都有一个唯一的64位整形值
//	////m.myKey类似于"service/lock/12345"
//	//m.myKey = fmt.Sprintf("%s%x", m.pfx, s.Lease())
//	//
//	//
//	////etcdv3新引入的多键条件事务，替代了v2中Compare-And-put操作。etcdv3的多键条件事务的语意是先做一个比较（compare）操作，如果比较成立则执行一系列操作，如果比较不成立则执行另外一系列操作。有类似于C语言中的条件表达式。
//	////接下来的这部分实现了如果不存在这个key，则将这个key写入到etcd，如果存在则读取这个key的值这样的功能。
//	////下面这一句，是构建了一个compare的条件，比较的是key的createRevision，如果revision是0，则存入一个key，如果revision不为0，则读取这个key。
//	////revision是etcd一个全局的序列号，每一个对etcd存储进行改动都会分配一个这个序号，在v2中叫index，createRevision是表示这个key创建时被分配的这个序号。当key不存在时，createRivision是0。
//	//cmp := clientv3.Compare(clientv3.CreateRevision(m.myKey), "=", 0)
//	//put := clientv3.OpPut(m.myKey, "", clientv3.WithLease(s.Lease()))
//	//get := clientv3.OpGet(m.myKey)
//	//resp, err := clientv3.Txn(ctx).If(cmp).Then(put).Else(get).Commit()
//	//if err != nil {
//	//	return err
//	//}
//	//m.myRev = resp.Header.Revision
//	//if !resp.Succeeded {
//	//	m.myRev = resp.Responses[0].GetResponseRange().Kvs[0].CreateRevision
//	//}
//	//
//	////如果上面的code操作成功了，则myRev是当前客户端创建的key的revision值。
//	////waitDeletes等待匹配m.pfx （"service/lock/"）这个前缀（可类比在这个目录下的）并且createRivision小于m.myRev-1所有key被删除
//	////如果没有比当前客户端创建的key的revision小的key，则当前客户端者获得锁
//	////如果有比它小的key则等待，比它小的被删除
//	//err = waitDeletes(ctx, client, m.pfx, m.myRev-1)
//
//}
