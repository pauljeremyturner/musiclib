package main

import (
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(4)

	NewRestRouter().RestRouter()
}
