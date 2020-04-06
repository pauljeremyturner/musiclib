package main

import (
	rest "github.com/pauljeremyturner/musiclib/cmd/main-module/rest"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(4)

	rest.RestRouter()

}
