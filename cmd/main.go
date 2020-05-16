package main

import (
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(4)


	database := NewDatabase()

	mdr := NewMetaDataReader(database)
	go InitialiseGraphQl(database)
	NewRestRouter(mdr, database).RestRouter()

}
