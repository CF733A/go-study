package main

import "fmt"

func main(){
	fmt.Println(appendInt([]int{1,2,3}, 1,2,3))
}

func appendInt(xs []int, x... int) []int {
    return append(xs, x...)
}