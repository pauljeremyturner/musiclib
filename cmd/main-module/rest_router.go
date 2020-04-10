package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var lib Library

func RestRouter() {

	router := mux.NewRouter()

	albumRouter := router.PathPrefix("/albums").Subrouter()
	trackRouter := router.PathPrefix("/tracks").Subrouter()
	libraryRouter := router.PathPrefix("/librarys").Subrouter()

	albumRouter.
		HandleFunc("/", allAlbums).
		Methods("GET").
		Schemes("http")

	albumRouter.
		HandleFunc("/{albumid}", getAlbum).
		Methods("GET").Schemes("http")

	trackRouter.HandleFunc("/", allTracks).
		Methods("GET").
		Schemes("http")

	trackRouter.
		HandleFunc("/{trackid}", getTrack).
		Methods("GET").
		Schemes("http")

	libraryRouter.
		HandleFunc("/", loadLibrary).
		Methods("POST").
		Schemes("http")

	libraryRouter.
		HandleFunc("/", getLibrary).
		Methods("GET").
		Schemes("http")

	http.ListenAndServe(":8080", router)
}

func getLibrary(writer http.ResponseWriter, request *http.Request) {
	bytes, _ := json.Marshal(lib)

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func loadLibrary(writer http.ResponseWriter, request *http.Request) {

	var requestBody LibraryRequest

	err := func() error {
		bytes, err1 := ioutil.ReadAll(request.Body)
		if err1 != nil {
			log.Printf("Failed read post body", err1)
			return err1
		}
		err2 := json.Unmarshal(bytes, &requestBody)
		if err2 != nil {
			log.Printf("Failed unmarshal json post body, error: %s", err2)
			return err2
		}
		return nil
	}()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Printf("Loading music from %s", requestBody.Path)

		md := NewMetaData(requestBody.Path)
		lib = NewLibrary()
		lib.LoadFromPath(md)
		writer.WriteHeader(http.StatusAccepted)
	}
}

func getTrack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	track := lib.TracksByTitle[params["trackid"]]

	bytes, err := json.Marshal(track)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)
	}
}

func allTracks(writer http.ResponseWriter, request *http.Request) {
}

func getAlbum(writer http.ResponseWriter, request *http.Request) {
}

func allAlbums(writer http.ResponseWriter, request *http.Request) {

}
