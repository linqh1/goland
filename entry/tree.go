package main

import (
	"fmt"
	"golang.org/x/tour/tree"
	"sort"
)

// Walk 步进 tree t 将所有的值从 tree 发送到 channel ch。
func Walk(t *tree.Tree, ch chan int){
	ch <- t.Value
	if t.Left != nil {
		Walk(t.Left,ch)
	}
	if t.Right != nil {
		Walk(t.Right,ch)
	}
}

// Same 检测树 t1 和 t2 是否含有相同的值。
func Same(t1, t2 *tree.Tree) bool{
	var l1,l2 = Len(t1),Len(t2)
	if l1 != l2 {
		fmt.Printf("Tree len is not same,len(tree1):%v,len(tree2):%v\n",l1,l2)
		return false
	}
	c1,c2 := make(chan int,l1),make(chan int,l1)
	go Walk(t1,c1)
	go Walk(t2,c2)
	s1,s2 := make([]int,l1),make([]int,l1)
	for i:=0;i<l1;i++ {
		s1[i] = <-c1
		s2[i] = <-c2
	}
	sort.Ints(s1)
	sort.Ints(s2)
	for i:=0;i<l1;i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func Len(t *tree.Tree) int {
	var l int
	treeLenCalculate(t,&l)
	return l
}

func treeLenCalculate(t1 *tree.Tree,len *int){
	*len += 1
	if t1.Left != nil {
		treeLenCalculate(t1.Left,len)
	}
	if t1.Right != nil {
		treeLenCalculate(t1.Right,len)
	}
}

func main() {
	tree1 := tree.New(1)
	tree2 := tree.New(1)
	fmt.Println(tree1,tree2)
	fmt.Println(Same(tree1,tree2))
}
