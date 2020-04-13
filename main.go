package main

import (
	"github.com/pauljeremyturner/musiclib/processor"
	"github.com/pauljeremyturner/musiclib/reader"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(4)

	mdr := reader.NewMetaDataReader()
	mdp := processor.NewMetaDataProcessor()

	NewRestRouter(mdr, mdp).RestRouter()
}
