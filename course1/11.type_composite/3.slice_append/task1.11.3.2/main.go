package main

import "fmt"


func main(){
	mySlice := []int{1,2,3}
	appendInt(&mySlice, 1,2,3)
	fmt.Println(mySlice)
}

func appendInt(xs *[]int, x... int) {
    *xs = append(*xs, x...)
}