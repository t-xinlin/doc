package main

import "fmt"

func main() {
	s1 := []int{1, 2, 3}
	fmt.Println(len(s1), cap(s1), s1) //prints 3 3 [1 2 3]

	s2 := s1[1:]
	fmt.Println(len(s2), cap(s2), s2)         //prints 2 2 [2 3]for i := range s2 { s2[i]+=20}//still referencing the same array
	fmt.Printf("pointer: %p   %#v\n", s1, s1) //prints [1 22 23]
	fmt.Printf("pointer: %p   %#v\n", s2, s2) //prints [22 23]

	s2 = append(s2, 4)
	fmt.Println(len(s2), cap(s2), s2)
	for i := range s2 {
		s2[i] += 10
	} //s1 is now "stale"
	fmt.Printf("pointer: %p   %#v\n", s1, s1) //prints [1 22 23]
	fmt.Printf("pointer: %p   %#v\n", s2, s2) //prints [32 33 14]}

	data := make([]byte, 10000)
	data[0] = 88

	fmt.Println(len(data), cap(data), &data[0], data[0]) //prints: 3 10000 <byte_addr_x>}
	data1 := data[:3]

	//data2 := make([]byte, 3)
	data2 := data[:3]
	//copy(data2, data)
	data2 = append(data2, 6, 7, 6, 7, 3, 4, 5, 6)

	for i := range data2 {
		data2[i] *= 19
	}

	data[0] = 78
	data1[0] = 99

	fmt.Println(len(data1), cap(data1), &data1[0], data[0], data1[0], data2[0]) //prints: 3 10000 <byte_addr_x>}
}
