package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	//scaner := bufio.NewScanner(os.Stdin)
	//scaner.Split(bufio.ScanLines)
	//for scaner.Scan(){
	//	fmt.Println(scaner.Text())
	//}
	input := "abcdefghijkl"
	scanner := bufio.NewScanner(strings.NewReader(input))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
		return 0, nil, nil
	}
	scanner.Split(split)
	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}
