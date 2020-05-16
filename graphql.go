package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
	"io/ioutil"
	"log"
	"net/http"
)

func InitialiseGraphQl(mdb musicDatabase) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(
			createQueryType(
				createAlbumType(
					createTrackType(),
				),
			),
		),
	})
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	handler := gqlhandler.New(&gqlhandler.Config{
		Schema: &schema,
	})
	http.Handle("/graphql", handler)
	log.Println("Server started at http://localhost:3000/graphql")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func createQueryType(albumType *graphql.Object) graphql.ObjectConfig {
	return graphql.ObjectConfig{Name: "QueryType", Fields: graphql.Fields{
		"album": &graphql.Field{
			Type: albumType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"]
				log.Printf("fetching album with id: %s", id)
				return fetchAlbumByiD(id.(string))
			},
		},
	}}
}

func createAlbumType(trackType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Album",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"artist": &graphql.Field{
				Type: graphql.String,
			},
			"tracks": &graphql.Field{
				Type: graphql.NewList(trackType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					album, _ := p.Source.(*Album)
					log.Printf("fetching tracks of album with id: %s", album.Id)
					return fetchTracksByAlbumID(album.Id)
				},
			},
		},
	})
}

func createTrackType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Track",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"filePath": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

func fetchAlbumByiD(albumId string) (*Album, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/albums/%s", albumId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", "could not fetch data", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not read data")
	}
	result := Album{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, errors.New("could not unmarshal data")
	}
	return &result, nil
}

func fetchTracksByAlbumID(albumId string) ([]Track, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/albums/%s/tracks", albumId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", "could not fetch data", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not read data")
	}
	result := []Track{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, errors.New("could not unmarshal data")
	}
	return result, nil
}
