package main

import (
	"fmt"
	"strings"
	"time"
)

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	for len(left) > 0 && len(right) > 0 {
		if left[0] <= right[0] {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}
	result = append(result, left...)
	result = append(result, right...)
	return result
}

func insertionSort(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func selectionSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
}

func quicksort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	left, right := 0, len(arr)-1
	pivot := len(arr) / 2

	arr[pivot], arr[right] = arr[right], arr[pivot]

	for i := range arr {
		if arr[i] < arr[right] {
			arr[left], arr[i] = arr[i], arr[left]
			left++
		}
	}

	arr[left], arr[right] = arr[right], arr[left]

	quicksort(arr[:left])
	quicksort(arr[left+1:])

	return arr
}

func GeneralSort(arr []int) {
	waySort := chekSortTime(arr)
	switch {
	case strings.Contains(waySort, "InsertionSorter"):
		insertionSort(arr)
	case strings.Contains(waySort, "MergedSorter"):
		mergeSort(arr)
	case strings.Contains(waySort, "SelectionSorter"):
		selectionSort(arr)
	case strings.Contains(waySort, "QuickSorter"):
		quicksort(arr)
	}
}

type Sorter interface {
	Sort([]int)
}

type MergedSorter struct{}

func (m *MergedSorter) Sort(data []int) {
	mergeSort(data)
}

type InsertionSorter struct{}

func (i *InsertionSorter) Sort(data []int) {
	insertionSort(data)
}

type SelectionSorter struct{}

func (s *SelectionSorter) Sort(data []int) {
	selectionSort(data)
}

type QuickSorter struct{}

func (q *QuickSorter) Sort(data []int) {
	quicksort(data)
}

func chekSortTime(data []int) string {
	sorters := []Sorter{&MergedSorter{}, &InsertionSorter{}, &SelectionSorter{}, &QuickSorter{}}
	result := make(map[string]time.Duration)
	iterations := 1000000

	for _, sort := range sorters {
		var totalTime time.Duration

		for i := 0; i < iterations; i++ {

			copyData := make([]int, len(data))
			copy(copyData, data)

			start := time.Now()
			sort.Sort(copyData)
			end := time.Since(start)

			totalTime += end
		}
		AVG := totalTime / time.Duration(iterations)
		typeName := fmt.Sprintf("%T", sort)
		result[typeName] = AVG
	}

	var fastestTime time.Duration = 1 * time.Hour
	var fastestName string

	for name, dura := range result {
		if dura < fastestTime {
			fastestName = name
			fastestTime = dura
		}
	}
	return fastestName
}

func main() {
	data := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Println("Original: ", data)

	sortedData := mergeSort(data)
	fmt.Println("Sorted by Merge Sort: ", sortedData)

	data = []int{64, 34, 25, 12, 22, 11, 90}
	insertionSort(data)
	fmt.Println("Sorted by Insertion Sort: ", data)

	data = []int{64, 34, 25, 12, 22, 11, 90}
	selectionSort(data)
	fmt.Println("Sorted by Selection Sort: ", data)

	data = []int{64, 34, 25, 12, 22, 11, 90}
	sortedData = quicksort(data)
	fmt.Println("Sorted by Quicksort: ", sortedData)

	data = []int{64, 34, 25, 12, 22, 11, 90}
	GeneralSort(data)
	fmt.Println("Sorted by GeneralSort: ", data)
}
