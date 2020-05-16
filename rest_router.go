package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type restRouterState struct {
	mdr MetaDataReader
	mdb musicDatabase
}

func NewRestRouter(mdr MetaDataReader, mdb musicDatabase) RestRouter {
	return &restRouterState{
		mdr: mdr,
		mdb: mdb,
	}
}

type RestRouter interface {
	RestRouter()
	LoadLibrary(writer http.ResponseWriter, request *http.Request)
	GetLibrary(writer http.ResponseWriter, request *http.Request)
	GetAlbumById(writer http.ResponseWriter, request *http.Request)
}

func (r *restRouterState) RestRouter() {

	router := mux.NewRouter()

	albumRouter := router.PathPrefix("/albums").Subrouter()
	trackRouter := router.PathPrefix("/tracks").Subrouter()
	libraryRouter := router.PathPrefix("/librarys").Subrouter()

	albumRouter.
		HandleFunc("/{albumid}", r.GetAlbumById).
		Methods("GET").Schemes("http")

	albumRouter.
		HandleFunc("/{albumid}/tracks", r.GetTracksByAlbum).
		Methods("GET").Schemes("http")

	trackRouter.
		HandleFunc("/{trackid}", r.GetTrack).
		Methods("GET").
		Schemes("http")

	libraryRouter.
		HandleFunc("/", r.LoadLibrary).
		Methods("POST").
		Schemes("http")

	libraryRouter.
		HandleFunc("/", r.GetLibrary).
		Methods("GET").
		Schemes("http")

	http.ListenAndServe(":8080", router)
}

func (r *restRouterState) GetLibrary(writer http.ResponseWriter, request *http.Request) {

	tracks, _ := r.mdb.LoadTracks()
	bytes, _ := json.Marshal(tracks)

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func (r *restRouterState) LoadLibrary(writer http.ResponseWriter, request *http.Request) {

	var requestBody LibraryRequest

	err := func() error {
		bytes, err1 := ioutil.ReadAll(request.Body)
		if err1 != nil {
			log.Printf("Failed read post body err: %s", err1)
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

		err := r.mdr.ReadMetaData(requestBody.Path)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.WriteHeader(http.StatusAccepted)
		}
	}
}

func (r *restRouterState) GetTracksByAlbum(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	albumid, ok := params["albumid"]; if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}



	tracks, ok := r.mdb.LoadTracksByAlbumId(albumid); if ok {
		bytes, err := json.Marshal(tracks)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			writer.Write(bytes)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func (r *restRouterState) GetTrack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	trackid := params["trackid"]


	if album, ok := r.mdb.LoadTracksById(trackid); ok {
		bytes, err := json.Marshal(album)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			writer.Write(bytes)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func (r *restRouterState) GetAlbumById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	albumid := params["albumid"]



	if track, ok := r.mdb.LoadAlbumById(albumid); ok {
		bytes, err := json.Marshal(track)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			writer.Write(bytes)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}
