package library

import (
	"metadata"

	model "github.com/pauljeremyturner/musiclib/cmd/main-module/model"
)

func LoadFromPath(path string) model.Library {

	tracks := metadata.ProcessMetadata(path)

	albumTrackMap := make(map[string][]model.Track)

	for _, t := range tracks {
		album := t.Album
		albumTrackMap[album] = append(albumTrackMap[album], t)
	}

	library := model.Library{}

	for k, v := range albumTrackMap {
		artistSet := make(map[string]bool)
		for _, t := range v {
			artistSet[t.Artist] = true
		}
		var combinedArtist string
		for artist := range artistSet {
			combinedArtist = combinedArtist + artist + " "
		}
		album := model.Album{Title: k, Tracks: v, Artist: combinedArtist}

		library.Albums = append(library.Albums, album)

	}

	return library
}
