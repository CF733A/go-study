package main

import (
	"fmt"
)

var Operate = func(f func(xs ...interface{}) interface{}, i ...interface{}) interface{} {
	return f(i...)
}

var Concat = func(xs ...interface{}) interface{} {
	var strResult string
	for _, val := range xs {
		val, ok := val.(string)
		if ok {
			strResult += val
		}
	}
	return strResult
}

var Sum = func(xs ...interface{}) interface{} {
	var ires int
	var fres float64
	isFloat := false
	for _, val := range xs {
		switch v := val.(type) {
		case int:
			ires += v
		case float64:
			fres += v
			isFloat = true
		default:
			continue
		}
	}
	if isFloat {
		return fres
	}
	return ires
}

func main() {
	fmt.Println(Operate(Concat, "Hello, ", "World!"))  // Вывод: "Hello, World!"
	fmt.Println(Operate(Sum, 1, 2, 3, 4, 5))           // Вывод: 15
	fmt.Println(Operate(Sum, 1.1, 2.2, 3.3, 4.4, 5.5)) // Вывод: 16.5
}
