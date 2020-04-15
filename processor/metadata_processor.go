package processor

import (
	"github.com/pauljeremyturner/musiclib/model"
	"log"
)

type MetaDataProcessor interface {
	TransformMetaData(inTracks []model.Track) model.Library
}

type metaDataProcessorState struct {}

func NewMetaDataProcessor() MetaDataProcessor {
	return metaDataProcessorState{}
}

func (r metaDataProcessorState) TransformMetaData(tracks []model.Track) model.Library {

	lib := model.NewLibrary()

	albumTrackMap := make(map[string][]model.Track)

	for _, t := range tracks {
		album := t.Album
		albumTrackMap[album] = append(albumTrackMap[album], t)

		lib.TracksByTitle[t.Title] = t
		lib.TracksById[t.Id] = t
	}

	var id int = 0

	for k, v := range albumTrackMap {
		artistSet := make(map[string]bool)
		for _, t := range v {
			artistSet[t.Artist] = true
		}
		var combinedArtist string
		for artist := range artistSet {
			combinedArtist = combinedArtist + artist + " "
		}
		album := model.Album{Id: id, Title: k, Tracks: v, Artist: combinedArtist}
		id++
		lib.AlbumsByTitle[k] = album
		lib.AlbumsById[album.Id] = album

	}

	log.Printf("Loaded %3d albums", len(lib.AlbumsById))

	return lib

}
