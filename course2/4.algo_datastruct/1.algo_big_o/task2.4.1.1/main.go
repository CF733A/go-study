package main

import (
	"fmt"
	"runtime"
	"time"
)

func factorialRecursive(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorialRecursive(n-1)
}

func factorialIterative(n int) int {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}

// выдает true, если реализация быстрее и false, если медленнее
func compareWhichFactorialIsFaster() map[string]bool {
	result := make(map[string]bool)
	n := 1000

	startR := time.Now()
	_ = factorialRecursive(n)
	endR := time.Since(startR)

	startI := time.Now()
	_ = factorialIterative(n)
	endI := time.Since(startI)

	result["iterative"] = endI < endR
	result["recursive"] = endR < endI

	return result
}

func main() {
	fmt.Println("Go version:", runtime.Version())
	fmt.Println("Go OS/Arch:", runtime.GOOS, "/", runtime.GOARCH)

	fmt.Println("Which factorial is faster?")
	fmt.Println(compareWhichFactorialIsFaster())
}
