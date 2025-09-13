package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Need one or more arguments")
		return
	}

	var minNum, maxNum float64
	for i := 1; i < len(arguments); i++ {
		n, err := strconv.ParseFloat(arguments[i], 64)
		if err != nil {
			continue
		}

		if i == 1 {
			minNum = n
			maxNum = n
			continue
		}

		if n < minNum {
			minNum = n
		}

		if n > maxNum {
			maxNum = n
		}
	}

	fmt.Println("Min: ", minNum)
	fmt.Println("Max: ", maxNum)
}
