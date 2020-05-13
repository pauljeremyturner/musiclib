package main

import (
	"fmt"
	"gotest.tools/assert"
	"strconv"
	"testing"
)

var testTracks []Track
var mdp MetaDataProcessor
var library Library

func getTestTracks() []Track {

	var tracks []Track

	for i := 0; i < 6; i++ {
		t := Track{
			Id:          int(i),
			Title:       "title" + strconv.Itoa(i),
			Artist:      "artist",
			TrackNumber: TrackNumber{i, 6},
			Album:       "album",
			AlbumArtist: "album-artist",
			Composer:    "composer",
			FilePath:    fmt.Sprintf("/home/paul/Music/album/track-%2d.mp3", i),
		}
		tracks = append(tracks, t)
	}

	return tracks
}

func setupLibraryTest() {
	testTracks = getTestTracks()
	mdp = NewMetaDataProcessor()
	library = mdp.TransformMetaData(testTracks)
}

func TestMain(m *testing.M) {
	setupLibraryTest()

	m.Run()
}

func TestOrganizeIntoAlbums(t *testing.T) {

	album := library.AlbumsByTitle["album"]

	for i, track := range album.Tracks {
		assert.Equal(t, "title"+strconv.Itoa(i), track.Title)
	}
}

func TestOrganizeByTrackTitle(t *testing.T) {

	for k, v := range library.TracksByTitle {
		assert.Equal(t, k, v.Title)
	}
}
