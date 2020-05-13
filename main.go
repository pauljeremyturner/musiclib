package main

import (
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(4)

	mdr := NewMetaDataReader()
	mdp := NewMetaDataProcessor()

	NewRestRouter(mdr, mdp).RestRouter()
}
