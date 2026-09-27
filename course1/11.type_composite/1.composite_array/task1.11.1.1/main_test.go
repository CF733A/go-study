package main

import (
	"testing"
)

func TestMain(t *testing.T){
	main()
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		xs   [8]int
		want int
	}{
		{"All positive numbers", [8]int{1, 2, 3, 4, 5, 6, 7, 8}, 36},
		{"All zeros", [8]int{0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{"Mixed numbers", [8]int{-1, 2, -3, 4, -5, 6, -7, 8}, 4},
		{"All negative numbers", [8]int{-1, -2, -3, -4, -5, -6, -7, -8}, -36},
		{"Single non-zero", [8]int{0, 0, 0, 0, 5, 0, 0, 0}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.xs); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name string
		xs   [8]int
		want float64
	}{
		{"All positive numbers", [8]int{1, 2, 3, 4, 5, 6, 7, 8}, 4.5},
		{"All zeros", [8]int{0, 0, 0, 0, 0, 0, 0, 0}, 0.0},
		{"Mixed numbers", [8]int{-1, 2, -3, 4, -5, 6, -7, 8}, 0.5},
		{"All negative numbers", [8]int{-1, -2, -3, -4, -5, -6, -7, -8}, -4.5},
		{"Single non-zero", [8]int{0, 0, 0, 0, 10, 0, 0, 0}, 1.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := average(tt.xs); got != tt.want {
				t.Errorf("average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAverageFloat(t *testing.T) {
	tests := []struct {
		name string
		xs   [8]float64
		want float64
	}{
		{"All positive numbers", [8]float64{1, 2, 3, 4, 5, 6, 7, 8}, 4.5},
		{"All zeros", [8]float64{0, 0, 0, 0, 0, 0, 0, 0}, 0.0},
		{"Mixed numbers", [8]float64{-1.5, 2.5, -3.5, 4.5, -5.5, 6.5, -7.5, 8.5}, 0.5},
		{"All negative numbers", [8]float64{-1, -2, -3, -4, -5, -6, -7, -8}, -4.5},
		{"Single non-zero", [8]float64{0, 0, 0, 0, 10.5, 0, 0, 0}, 1.3125},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := averageFloat(tt.xs); got != tt.want {
				t.Errorf("averageFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		xs   [8]int
		want [8]int
	}{
		{"All positive numbers", [8]int{1, 2, 3, 4, 5, 6, 7, 8}, [8]int{8, 7, 6, 5, 4, 3, 2, 1}},
		{"All zeros", [8]int{0, 0, 0, 0, 0, 0, 0, 0}, [8]int{0, 0, 0, 0, 0, 0, 0, 0}},
		{"Mixed numbers", [8]int{-1, 2, -3, 4, -5, 6, -7, 8}, [8]int{8, -7, 6, -5, 4, -3, 2, -1}},
		{"Sequential numbers", [8]int{10, 20, 30, 40, 50, 60, 70, 80}, [8]int{80, 70, 60, 50, 40, 30, 20, 10}},
		{"Palindrome array", [8]int{1, 2, 3, 4, 4, 3, 2, 1}, [8]int{1, 2, 3, 4, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverse(tt.xs); got != tt.want {
				t.Errorf("reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
