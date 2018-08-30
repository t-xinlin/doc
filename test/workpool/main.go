func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("------------------------start------------------------\n")

	//test_big()

	//go func() {
	//	http.ListenAndServe("127.0.0.1:8888", nil)
	//}()

	//test_fabonacci()
	//test_once()

	//test_proxy()
	//test_HasSuffix_HasPrefix()
	//test_fibonacci()
	//m := WordCount(STR1)

	//fmt.Printf("result1 :  %+v\n", m)

	//工作池
	TestWorkPool()
}
