package main

import (
	"fmt"
	"sort"
)

func main() {
	intArr := [8]int{5, 2, 8, 1, 9, 3, 7, 4}
	floatArr := [8]float64{5.5, 2.2, 8.8, 1.1, 9.9, 3.3, 7.7, 4.4}

	sortedIntDesc := sortDescInt(intArr)
	sortedIntAsc := sortAscInt(intArr)
	sortedFloatDesc := sortDescFloat(floatArr)
	sortedFloatAsc := sortAscFloat(floatArr)

	fmt.Println("Sorted Int Array (Descending):", sortedIntDesc)
	fmt.Println("Sorted Int Array (Ascending):", sortedIntAsc)
	fmt.Println("Sorted Float Array (Descending):", sortedFloatDesc)
	fmt.Println("Sorted Float Array (Ascending):", sortedFloatAsc)
}

func sortAscInt(arr [8]int) [8]int {
	newSl := arr[:]
	sort.Ints(newSl)
	return [8]int(arr)
}


func sortDescInt(arr [8]int) [8]int {
	newSl := arr[:]
	sort.Sort(sort.Reverse(sort.IntSlice(newSl)))
	return [8]int(arr)
}

func sortAscFloat(arr [8]float64) [8]float64 {
	newSl := arr[:]
	sort.Float64s(newSl)
	return [8]float64(arr)
}

func sortDescFloat(arr [8]float64) [8]float64 {
	newSl := arr[:]
	sort.Sort(sort.Reverse(sort.Float64Slice(newSl)))
	return [8]float64(arr)
}