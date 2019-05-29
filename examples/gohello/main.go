package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Printf("hello, I am %s %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("args: %v\n", os.Args)
	fmt.Printf("env: %v\n", os.Environ())
}
