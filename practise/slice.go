package practise

import (
	"fmt"
)

func SliceTest()  {
	s1 := make([]int, 5)
	for i:=0;i<len(s1);i++ {
		s1[i] = i
	}
	printSlice(s1)//[0 1 2 3 4]
	s2 := s1[1:3]
	printSlice(s2)//[1 2]
	s2[0] = 11
	printSlice(s2)//[11 2]
	printSlice(s1)//[0 11 2 3 4]
	s3 := append(s2,5,6,7)//这个时候s2与s3就不再有关系
	printSlice(s3)//[11 2 5 6 7]
	s3[0] = 111
	printSlice(s3)//[111 2 5 6 7]
	printSlice(s2)//[11 2]
	printSlice(s1)//[0 11 2 3 4]
}

func printSlice(s []int)  {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

