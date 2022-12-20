package main

import (
	"fmt"
	"os"
)

func main () {
	data, err := os.ReadFile("/home/mikolaszko/t.txt")	
	fmt.Print(string(data), err)
}
