package main

import (
	"fmt"
	"gotest.tools/assert"
	"strconv"
	"testing"
)

type mockMetaDataState struct {
	tracks []Track
}

var md MetaData

func (m mockMetaDataState) ProcessMetadata() []Track {

	tracks := []Track{}

	for i := 0; i < 6; i++ {
		t := Track{
			Id:          int32(i),
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

func newMockMetaData() MetaData {
	return &mockMetaDataState{
		tracks: []Track{},
	}
}

func setupLibraryTest() {
	md = newMockMetaData()
}

func TestOrganizeIntoAlbums(t *testing.T) {
	setupLibraryTest()

	lib.LoadFromPath(md)

	album := lib.AlbumsByTitle["album"]

	for i, track := range album.Tracks {
		assert.Equal(t, "title"+strconv.Itoa(i), track.Title)
	}
}

func TestOrganizeByTrackTitle(t *testing.T) {
	setupLibraryTest()

	lib.LoadFromPath(md)

	for k, v := range lib.TracksByTitle {
		assert.Equal(t, k, v.Title)
	}
}
