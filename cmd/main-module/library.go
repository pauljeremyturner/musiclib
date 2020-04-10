package main

import (
	"log"
)

func (r *Library) LoadFromPath(md MetaData) {

	tracks := md.ProcessMetadata()

	albumTrackMap := make(map[string][]Track)

	for _, t := range tracks {
		album := t.Album
		albumTrackMap[album] = append(albumTrackMap[album], t)

		r.TracksByTitle[t.Title] = t
		r.TracksById[t.Id] = t
	}

	var id int32 = 0

	for k, v := range albumTrackMap {
		artistSet := make(map[string]bool)
		for _, t := range v {
			artistSet[t.Artist] = true
		}
		var combinedArtist string
		for artist := range artistSet {
			combinedArtist = combinedArtist + artist + " "
		}
		album := Album{Id: id, Title: k, Tracks: v, Artist: combinedArtist}
		id++
		r.AlbumsByTitle[k] = album
		r.AlbumsById[album.Id] = album

	}

	log.Printf("Loaded %3d albums", len(r.AlbumsById))

}
