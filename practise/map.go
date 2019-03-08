package practise

import "fmt"

func MapTest(){
	//max := int(^uint(0) >> 1)
	//fmt.Printf("%b %d %v %x\n",max,max,max,max)
	//judgeIfExist()
}

func judgeIfExist(){
	map1 := make(map[string]bool)
	map1["shindou"] = true
	map1["sai"] = true
	if _,isexist := map1["akira"];!isexist {
		map1["akira"] = true
		fmt.Printf("Type of isexist is : %T\n",isexist)
	}
	map1["kouyou"] = false
	fmt.Println(map1)
	if _,isexist := map1["kouyou"];!isexist {
		map1["kouyou"] = true
	}
	fmt.Println(map1)
	delete(map1,"sai")
	fmt.Println(map1)
}
