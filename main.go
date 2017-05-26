package main

import (
	"fmt"
)

func main() {
	a:=[]int{1,2,3,4,5,6,7}
	b:=a[2:5]
	fmt.Println(b)
	c:=b[-1:2]
	fmt.Println(c)

	dispatcher := NewDispatcher()
	dispatcher.SetNetService(&NetService{})
	dispatcher.Dispatch(NewTestSpider())
}
