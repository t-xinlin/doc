package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := make(chan int)
	ch := UseTicker(intChan)
	arr:=make([]int,0)
	i :=0
	for e := range intChan {
		if len(arr)>=50{
			break
		}
		//arr[i] = e
		arr = append(arr, e)
		i++
		fmt.Println("Recv: ", e)
	}
	ch <- true
	time.Sleep(time.Second * 2)
	fmt.Println("End Recv sum", arr)
	close(ch)

}

func UseTicker(intchan chan int) chan bool {
	stop := make(chan bool)
	ticker := time.NewTicker(time.Millisecond * 1)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//fmt.Println("ticker")
				setValue(intchan)
			case <-stop:
				fmt.Println("stop")
				return

			}
		}
	}()
	return stop
}

func setValue(intChan chan int) {
	select {
	case intChan <- 0:
	case intChan <- 1:
	case intChan <- 2:
	case intChan <- 3:
	case intChan <- 4:
	case intChan <- 5:
	case intChan <- 6:
	case intChan <- 7:
	case intChan <- 8:
	case intChan <- 9:
	//case intChan <- 'A':
	//case intChan <- 'B':
	//case intChan <- 'C':


	}
}
