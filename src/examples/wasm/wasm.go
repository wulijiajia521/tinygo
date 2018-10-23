package main

import "os"

func main() {
}

//go:export add
func add(a, b int) int {
	os.Stdout.Write([]byte("adding two numbers"))
	return a + b
}
