package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"library"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/pauljeremyturner/musiclib/cmd/main-module/model"
)

var lib model.Library

func RestRouter() {

	r := mux.NewRouter()

	albumRouter := r.PathPrefix("/albums").Subrouter()
	trackRouter := r.PathPrefix("/tracks").Subrouter()
	libraryRouter := r.PathPrefix("/librarys").Subrouter()

	albumRouter.
		HandleFunc("/", AllAlbums).
		Methods("GET").
		Schemes("http")

	albumRouter.
		HandleFunc("/{album}", GetAlbum).
		Methods("GET").Schemes("http")

	trackRouter.HandleFunc("/", AllTracks).
		Methods("GET").
		Schemes("http")

	trackRouter.
		HandleFunc("/{track}", GetTrack).
		Methods("GET").
		Schemes("http")

	libraryRouter.
		HandleFunc("/", LoadLibrary).
		Methods("POST").
		Schemes("http")

	http.ListenAndServe(":8080", r)
}

func LoadLibrary(writer http.ResponseWriter, request *http.Request) {

	var requestBody model.LibraryRequest

	bytes, _ := ioutil.ReadAll(request.Body)
	err := json.Unmarshal(bytes, &requestBody)

	if err != nil {

	}

	log.Printf("Loading music from %s", requestBody.Path)

	go library.LoadFromPath(requestBody.Path)

	writer.WriteHeader(http.StatusAccepted)
}

func GetTrack(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Content-Type", "application/json")

	fmt.Fprintf(writer, "hello world")
}

func AllTracks(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "hello world")
}

func GetAlbum(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "hello world")
}

func AllAlbums(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprintf(writer, "hello world")
}
