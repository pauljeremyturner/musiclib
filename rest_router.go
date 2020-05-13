package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type restRouterState struct {
	mdr MetaDataReader
	mdp MetaDataProcessor
	lib Library
}

func NewRestRouter(inMd MetaDataReader, inMdp MetaDataProcessor) RestRouter {
	return &restRouterState{
		mdr: inMd,
		mdp: inMdp,
	}
}

type RestRouter interface {
	RestRouter()
	LoadLibrary(writer http.ResponseWriter, request *http.Request)
	GetLibrary(writer http.ResponseWriter, request *http.Request)
	GetAlbum(writer http.ResponseWriter, request *http.Request)
}

func (r *restRouterState) RestRouter() {

	router := mux.NewRouter()

	albumRouter := router.PathPrefix("/albums").Subrouter()
	trackRouter := router.PathPrefix("/tracks").Subrouter()
	libraryRouter := router.PathPrefix("/librarys").Subrouter()

	albumRouter.
		HandleFunc("/{albumid}", r.GetAlbum).
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
	bytes, _ := json.Marshal(r.lib)

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

		tracks := r.mdr.ReadMetaData(requestBody.Path)

		library := r.mdp.TransformMetaData(tracks)

		r.lib = library

		writer.WriteHeader(http.StatusAccepted)
	}
}

func (r *restRouterState) GetTrack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	i, paramErr := strconv.Atoi(params["trackid"])

	if paramErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	albumId := int(i)
	if album, ok := r.lib.AlbumsById[albumId]; ok {
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

func (r *restRouterState) GetAlbum(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	i, paramErr := strconv.Atoi(params["albumid"])

	if paramErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	trackId := int(i)
	if track, ok := r.lib.TracksById[trackId]; ok {
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
